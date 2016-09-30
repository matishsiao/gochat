package main

import (
	_ "crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func WebServer() {
	//r := mux.NewRouter()
	//r.HandleFunc("/", HomeHandler)
	http.HandleFunc("/status", StatusHandler)
	http.Handle("/", http.FileServer(http.Dir("html")))
	//fs := http.FileServer(http.Dir("html/static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
	//http.Handle("/", r)
	WsServer()
	log.Println("Web Service starting.")
	go http.ListenAndServe(CONFIGS.HTTP, nil)
	if CONFIGS.HTTPS != "" && CONFIGS.SSL.Key != "" && CONFIGS.SSL.Crt != "" {
		err := http.ListenAndServeTLS(CONFIGS.HTTPS, CONFIGS.SSL.Crt, CONFIGS.SSL.Key, nil)
		if err != nil {
			fmt.Println("ListenAndServeTLS:", err)
			log.Fatal("Listen SSL Web Server Failed:", err)
		}
	}
}

func jsonParser(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if data != nil {
		json, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(500)
			log.Println("Error generating json", err)
			fmt.Fprintln(w, "Could not generate JSON")
			return
		}
		fmt.Fprint(w, string(json))
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 no data can be find.")
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(403)
	fmt.Fprint(w, "403 Forbidden")
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	jsonParser(data, w)
}
