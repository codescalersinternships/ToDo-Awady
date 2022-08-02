package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	_ "restapi/docs"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

var DBFILE = "todo.db"
var LISTENURL = ":8080"

func main() {
	c := cors.AllowAll()

	if f, ok := os.LookupEnv("DBFILE"); ok {
		DBFILE = f
	}
	if u, ok := os.LookupEnv("LISTENURL"); ok {
		LISTENURL = u
	}

	r := mux.NewRouter()
	handler := c.Handler(r)
	a, err := NewApp(DBFILE, LISTENURL, handler)
	if err != nil {
		fmt.Printf("error: %q, DBFILE: %q, LISTENURL: %q", err.Error(), DBFILE, LISTENURL)
		panic("Couldn't initialize app")
	}

	r.Use(Middleware)
	r.HandleFunc("/todo", a.GetAllToDosHandler).Methods("GET")
	r.HandleFunc("/todo", a.AddToDoHandler).Methods("POST")
	r.HandleFunc("/todo/{id}", a.GetToDoHandler).Methods("GET")
	r.HandleFunc("/todo/{id}", a.UpdateToDoHandler).Methods("PUT")
	r.HandleFunc("/todo/{id}", a.DeleteToDoHandler).Methods("DELETE")
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := a.srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
