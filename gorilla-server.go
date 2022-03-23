package main

import (
	"context"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Define basic blog post structure
type BlogPost struct {
	Title      string `json:"title"`
	Last       int
	Number     int `json:"number"`
	Next       int
	Timestamp  string   `json:"timestamp"`
	Main       []string `json:"main"`
	ParsedMain template.HTML
}

// require my blog HTML template for a template parser
var blogTemplate = template.Must(template.ParseFiles("./assets/documents/blogtemplate.html"))

// blog handler
func blog_handler(w http.ResponseWriter, r *http.Request) {
	// locate JSON file
	blogstr := r.URL.Path[len("/blog/"):] + ".json"

	// open json file
	f, err := os.Open("db/" + blogstr)
	// if no file found
	if err != nil {
		// return error status
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// defer the closing of this json file until the page is done loading content
	defer f.Close()

	// define a blog post object
	var post BlogPost
	// decode the JSON data into the object
	if err := json.NewDecoder(f).Decode(&post); err != nil {
		// if an error occurs return HTTP status
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse the post's object data into the template
	post.ParsedMain = template.HTML(strings.Join(post.Main, " "))
	// get value for last post
	post.Last = post.Number - 1
	// get value for next post
	post.Next = post.Number + 1

	// execute and serve the template
	if err := blogTemplate.Execute(w, post); err != nil {
		// if an error occurs return status
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// print serve status
	log.Printf("Blog %d served.", post.Number)
}

// favicon handler
func favicon_handler(w http.ResponseWriter, r *http.Request) {
	// Serve the favicon
	http.ServeFile(w, r, "./assets/art/favicon.ico")
}

// teapot handler
func teapot_handler(w http.ResponseWriter, r *http.Request) {
	// return teapot state
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusTeapot)
	// serve a teapot image with linkage
	log.Print("Tea Served.")
	io.WriteString(w, "<html><h1><a href='https://datatracker.ietf.org/doc/html/rfc2324/'>HTCPTP</h1><img src='https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Ftaooftea.com%2Fwp-content%2Fuploads%2F2015%2F12%2Fyixing-dark-brown-small.jpg&f=1&nofb=1' alt='Im a teapot'></a><html>")
}

// Main server program
func main() {
	// initialize a time.Duration variable to hold a wait time-period
	var wait time.Duration
	// define a graceful termination period
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Initialize a new router
	router := mux.NewRouter()

	srv := &http.Server{
		// address to listen on
		Addr: "127.0.0.1:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		// router to serve as my handler
		Handler: router, // Pass our instance of gorilla/mux in.
	}

	// blog post handler
	router.HandleFunc("/blog/post{number}", blog_handler).Methods("GET")
	// site icon server
	router.HandleFunc("/favicon.ico", favicon_handler).Methods("GET")
	// teapot handler
	router.HandleFunc("/teapot", teapot_handler).Methods("GET")
	// define the fileserver root dir
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./"))))
	// pass all requests to my router
	http.Handle("/", router)
	// print listener status
	log.Print("Listening at http://localhost:8080")

	// Run the server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// define channel to recieve server shutdown signal
	shutdown_chan := make(chan os.Signal, 1)
	// I'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(shutdown_chan, os.Interrupt)

	// Block until I receive the signal in channel c
	<-shutdown_chan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, I could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if my application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	// exit the program succesfully
	os.Exit(0)
}
