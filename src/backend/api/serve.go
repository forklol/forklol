package api

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"backend/config"
	"time"
	"backend/bitcoin"
)

const JSON_PATH = "/tmp/forklol.json"

func Serve(port string, c chan bool) {
	r := mux.NewRouter()
	registerRoutes(r)

	serve := fmt.Sprintf("0.0.0.0:%s", port)

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
	r.HandleFunc("/{pepper}.fake.csv", handleRequestCached).Methods("GET")
	r.HandleFunc("/poll/{seq}", handlePollRequest).Methods("GET")
	r.HandleFunc("/last", handleLastRequest).Methods("GET")
	r.HandleFunc("/exchangerate", handleExchangerateRequest).Methods("GET")
}

func handleExchangerateRequest(w http.ResponseWriter, r *http.Request) {
	bitcoin.ExchangeRateLock.RLock()
	j := make([]byte, len(*bitcoin.ExchangeRateJson))
	copy(j, *bitcoin.ExchangeRateJson)
	bitcoin.ExchangeRateLock.RUnlock()

	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=15")
	w.Write(j)
}

func handleLastRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	http.ServeFile(w, r, JSON_PATH+".last")
}

func handleRequestCached(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=15")
	http.ServeFile(w, r, JSON_PATH)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	http.ServeFile(w, r, JSON_PATH)
}

var newblock uint64 = 0

func handlePollRequest(w http.ResponseWriter, r *http.Request) {

	oldblock := newblock

	tcr := time.NewTicker(5 * time.Second)
	tmr := time.NewTimer(85 * time.Second)

	for {
		select {
		case <-tcr.C:
			if oldblock < newblock && r.Context().Err() == nil {
				handleLastRequest(w, r)
				return
			}

			if r.Context().Err() != nil {
				return
			}
		case <-tmr.C:
			handleLastRequest(w, r)
			return
		}
	}
}

func checkBlocks(c chan bool) {
	for {
		select {
		case <-c:
			newblock++
		}
	}
}
