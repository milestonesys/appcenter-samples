package main

import (
	"log"
	"net/http"
	"strconv"

	"apigateway-webserver/src/pkg/handlers"
)

// Handlers
var homeHandler *handlers.HomeHandler
var loginHandler *handlers.LoginHandler
var viewHandler *handlers.ViewHandler
var eventHandler *handlers.EventHandler

func main() {
	// Initialize handlers
	homeHandler = handlers.NewHomeHandler()
	loginHandler = handlers.NewLoginHandler()
	http.HandleFunc("/", homeHandler.Handle)
	http.HandleFunc("/_login/", loginHandler.Handle)

	viewHandler = handlers.NewViewHandler()
	eventHandler = handlers.NewEventHandler()
	http.HandleFunc("/view_events/", viewHandler.Handle)
	http.HandleFunc("/view_events/_events_start/", eventHandler.StartSubscriptionHandle)
	http.HandleFunc("/view_events/_events_request/", eventHandler.RequestEventsHandle)

	err := http.ListenAndServe(":"+strconv.Itoa(8080), nil)
	if err != nil {
		log.Fatal("Error while starting the webserver: ", err)
		return
	}
}
