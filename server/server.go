package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer(port string) {
	router := mux.NewRouter()

	router.HandleFunc("/login", LoginHandlerGet).Methods("GET")
	router.HandleFunc("/api/login", LoginHandlerPost).Methods("POST")
	router.HandleFunc("/profile/{username}", ProfileHandler).Methods("GET")
	router.HandleFunc("/api/logout", LogoutHandler).Methods("GET")
	router.HandleFunc("/unauthorized", UnauthorizedLoginHandler).Methods("GET")
	router.HandleFunc("/profile/{username}/add", ProfileAddQuoteGet).Methods("GET")
	router.HandleFunc("/profile/add", ProfileAddQuote).Methods("POST")
	router.HandleFunc("/profile/get/quotes", ProfileGetQuotes).Methods("GET")
	router.HandleFunc("/profile/delete", ProfileDeleteQuote).Methods("DELETE")
	router.HandleFunc("/profile/put", ProfileUpdateQuote).Methods("PUT")

	// Serving the static folder
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	fmt.Println("Server starting on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
