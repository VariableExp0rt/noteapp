package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/VariableExp0rt/dddexample/adding"
	"github.com/VariableExp0rt/dddexample/auth"
	"github.com/VariableExp0rt/dddexample/deleting"
	"github.com/VariableExp0rt/dddexample/listing"
	"github.com/VariableExp0rt/dddexample/middleware"
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
	log.Fatal(http.ListenAndServe(":8080", s.router))
	log.Print("Server listening on http://localhost:8080.")
}

func (s *Server) Shutdown(ctx context.Context) {
	s.Shutdown(ctx)
	log.Print("Server shutting down.")
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

func (s *Server) RegisterRoutes(adder adding.Service, lister listing.Service, deleter deleting.Service, updater updating.Service, authep auth.Service) {

	//wrap all routes in AuthMiddleware, except login
	s.router.HandleFunc("/login", auth.MakeUserLoginEndpoint(authep)).Methods("POST")
	s.router.HandleFunc("/signup", auth.MakeUserSignUpEndpoint(authep)).Methods("POST")
	s.router.HandleFunc("/notes/{id}", middleware.AuthMiddleware(listing.MakeGetNoteEndpoint(lister))).Methods("GET")
	s.router.HandleFunc("/notes", middleware.AuthMiddleware(listing.MakeGetNotesEndpoint(lister))).Methods("GET")
	s.router.HandleFunc("/notes/{id}/delete", middleware.AuthMiddleware(deleting.MakeDeleteNoteEndpoint(deleter))).Methods("POST")
	s.router.HandleFunc("/notes/{id}/update", middleware.AuthMiddleware(updating.MakeUpdateNoteEndpoint(updater))).Methods("POST")
	s.router.HandleFunc("/notes", middleware.AuthMiddleware(adding.MakeAddNoteEndpoint(adder))).Methods("POST")
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
	var userStorage auth.Repository

	noteStorage = &storage.BoltStorage{DB: store}
	userStorage = &storage.BoltStorage{DB: store}

	adder := adding.NewService(noteStorage)
	lister := listing.NewService(noteStorage)
	updater := updating.NewService(noteStorage)
	deleter := deleting.NewService(noteStorage)
	ath := auth.NewService(userStorage)

	r := mux.NewRouter()

	srv := Server{
		config: &Config{},
		logger: NewLogger(),
		router: r,
	}

	srv.RegisterRoutes(
		adder,
		lister,
		deleter,
		updater,
		ath,
	)

	go func() {
		srv.Run()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	os.Exit(0)
}
