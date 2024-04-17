package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Ip address of server in local router
var LocalIp = "192.168.1.100"

type Stats struct {
	StartTime    time.Time
	Redirections uint
	Pings        uint
	LastPing     time.Time
	mu           sync.Mutex
}

var stats = &Stats{
	StartTime:    time.Now(),
	LastPing:     time.Time{}, // use IsZero
	Redirections: 0,
	Pings:        0,
}

func (s *Stats) Ping() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.LastPing = time.Now()
	s.Pings += 1
}

type Db struct {
	Ip string `json:"ip"`
}

func getIpToRedirect(r *http.Request) string {
	clientIp, err := getClientIP(r)

	if err != nil {
		log.Printf("redirect err - %s", err)
		return db.Ip
	}

	if clientIp == db.Ip {
		return LocalIp
	}

	return db.Ip
}

var redirectBaseUrl = "/"

func (d *Db) RedirectHandler(r *http.Request) string {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.Redirections += 1

	return "http://" + getIpToRedirect(r) + ":9876/" + r.URL.Path[len(homeBaseUrl):]

}

func (d *Db) RedirectHome(r *http.Request) string {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.Redirections += 1

	return "http://" + getIpToRedirect(r) + ":9876/" + r.URL.Path[len(homeBaseUrl):]
}

func (d *Db) RedirectJellyfin(r *http.Request) string {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.Redirections += 1

	return "http://" + getIpToRedirect(r) + ":8096/" + r.URL.Path[len(jellyfinBaseUrl):]
}

var db Db

var homeBaseUrl = "/hp/"

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if len(db.Ip) == 0 {
		noIpResponse(w)
	} else {

		redirectUrl := db.RedirectHome(r)
		ip, _ := getClientIP(r)

		logg(r, http.StatusTemporaryRedirect, fmt.Sprintf("client: %s ; redirect: %s", ip, redirectUrl))

		http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)

	}
}

var jellyfinBaseUrl = "/tv/"

func jellyfinHandler(w http.ResponseWriter, r *http.Request) {
	if len(db.Ip) == 0 {
		noIpResponse(w)
	} else {
		redirectUrl := db.RedirectJellyfin(r)
		ip, _ := getClientIP(r)

		logg(r, http.StatusTemporaryRedirect, fmt.Sprintf("client: %s ; redirect: %s", ip, redirectUrl))

		http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)

	}
}

var statsBaseUrl = "/stats/"

func statsHandler(w http.ResponseWriter, r *http.Request) {

	jsonData, err := json.Marshal(stats)
	if err != nil {

		logg(r, http.StatusInternalServerError, err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	logg(r, http.StatusOK, fmt.Sprintf("stats sent: %s", jsonData))

}

var updateIpBaseUrl = "/update/"

func updateIpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		var data Db
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {

			logg(r, http.StatusBadRequest, err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		stats.Ping()
		db.Ip = data.Ip

		logg(r, http.StatusOK, fmt.Sprintf("received update: %s", data.Ip))

		w.WriteHeader(http.StatusOK)
		return
	}

	logg(r, http.StatusForbidden, "BAD METHOD")

	w.WriteHeader(http.StatusForbidden)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(homeBaseUrl, homeHandler)
	mux.HandleFunc(jellyfinBaseUrl, jellyfinHandler)
	mux.HandleFunc(statsBaseUrl, statsHandler)
	mux.HandleFunc(updateIpBaseUrl, updateIpHandler)

	server := &http.Server{
		Addr:    ":9876",
		Handler: mux,
	}

	log.Println("Serving on - http://0.0.0.0:9876")

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}
