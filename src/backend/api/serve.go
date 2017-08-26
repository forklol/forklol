package api

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"backend/config"
)

const JSON_PATH = "/tmp/forklol.json"

var blockchan chan bool
var clientchannels []chan bool

var lock sync.Mutex
var countLock sync.Mutex

var clientCount = 0

func Serve(port string, c chan bool) {
	r := mux.NewRouter()
	registerRoutes(r)

	serve := fmt.Sprintf("0.0.0.0:%s", port)

	blockchan = c
	clientchannels = make([]chan bool, 0, 8096)

	go checkBlocks(c)

	log.Printf("API started on %s\n", serve)
	log.Fatal(http.ListenAndServe(serve, getConfiguredHandler(r)))
}

func getConfiguredHandler(h http.Handler) http.Handler {
	if config.OPTIONS.Debug == true {
		return getCorsHandler(h)
	}

	return h
}

func getCorsHandler(handler http.Handler) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "OPTIONS"})
	return handlers.CORS(originsOk, headersOk, methodsOk)(handler)
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", handleRequest).Methods("GET")
	r.HandleFunc("/poll/{seq}", handlePollRequest).Methods("GET")
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	http.ServeFile(w, r, JSON_PATH)
}

func handlePollRequest(w http.ResponseWriter, r *http.Request) {
	c := make(chan bool)

	lock.Lock()
	clientchannels = append(clientchannels, c)
	lock.Unlock()

	select {
	case <-c:
		if r.Context().Err() == nil {
			countLock.Lock()
			clientCount++
			countLock.Unlock()
			handleRequest(w, r)
		}
	}
}

func checkBlocks(c chan bool) {
	for {
		select {
		case <-c:

			lock.Lock()
			for _, client := range clientchannels {
				client <- true
			}
			lock.Unlock()

			clientchannels = make([]chan bool, 0, 8096)

			go func() {
				time.Sleep(time.Second * 3)
				log.Printf("Sent new blockdata to %d clients.\n", clientCount)
				clientCount = 0
			}()
		}
	}
}
