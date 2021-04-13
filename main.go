package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/VariableExp0rt/dddexample/adding"
	"github.com/VariableExp0rt/dddexample/deleting"
	"github.com/VariableExp0rt/dddexample/listing"
	"github.com/VariableExp0rt/dddexample/notes"
	"github.com/VariableExp0rt/dddexample/storage"
	"github.com/VariableExp0rt/dddexample/updating"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	buf = new(bytes.Buffer)
)

type Config struct {
	//some other configuration
}

type Server struct {
	config *Config
	logger *zap.SugaredLogger
	router *mux.Router
}

func (s *Server) Run() {
	s.logger.Info("Server listening on 127.0.0.1:8080.")
	s.logger.Fatal(http.ListenAndServe(":8080", s.router))
}

func NewLogger() *zap.SugaredLogger {

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", "app.log"}
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Unable to create new logger, exiting with error: %v", err.Error())
		return nil
	}
	defer logger.Sync()

	suggar := logger.Sugar()

	return suggar
}

func NewRouter() *mux.Router {
	return mux.NewRouter()
}

func (s *Server) RegisterRoutes(adder adding.Service, lister listing.Service, deleter deleting.Service, updater updating.Service) {
	s.router.HandleFunc("/notes/{id}", listing.MakeGetNoteEndpoint(lister)).Methods("GET")
	s.router.HandleFunc("/notes", listing.MakeGetNotesEndpoint(lister)).Methods("GET")
	s.router.HandleFunc("/notes/{id}/delete", deleting.MakeDeleteNoteEndpoint(deleter)).Methods("POST")
	s.router.HandleFunc("/notes/{id}/update", updating.MakeUpdateNoteEndpoint(updater))
	s.router.HandleFunc("/notes", adding.MakeAddNoteEndpoint(adder)).Methods("POST")
}

func main() {

	//Other flags
	pflag.String("db", "/tmp/my.db", "Supply a path for Bolt to open the database.")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	store, err := bolt.Open(viper.GetString("db"), 600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("Unable to open database at %v. Error: %v", viper.GetString("db"), err)
	}

	defer store.Close()

	var noteStorage notes.Repository
	noteStorage = &storage.BoltStorage{DB: store}

	adder := adding.NewService(noteStorage)
	lister := listing.NewService(noteStorage)
	updater := updating.NewService(noteStorage)
	deleter := deleting.NewService(noteStorage)

	r := NewRouter()

	srv := Server{
		config: &Config{},
		logger: NewLogger(),
		router: r,
	}
	srv.logger.Info("Registering handler routes with server.")
	srv.RegisterRoutes(adder, lister, deleter, updater)
	srv.logger.Info("Successfully registered handler routes with server.")

	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		defer cancel()
		srv.logger.Info("Running server.")
		srv.Run()
	}()

	<-ctx.Done()
	srv.logger.Info("Received stop signal. Exiting.")
}
