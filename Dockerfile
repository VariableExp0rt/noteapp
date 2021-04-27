FROM golang:1.16.3-alpine3.13 as build

RUN mkdir -p $GOPATH/src/build
RUN adduser -u 10001 usr -D
COPY . $GOPATH/src/build
WORKDIR $GOPATH/src/build
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/notes

FROM alpine
COPY --from=build /bin/notes /bin/notes
COPY --from=build /etc/passwd /etc/passwd
USER usr
CMD [ "/bin/notes", "--port \":8081\"" ]