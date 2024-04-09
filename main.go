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
	log.Printf("%s", clientIp)
	if err != nil {
		log.Printf("redirect err - %s", err)
		return db.Ip
	}

	if clientIp == db.Ip {
		return LocalIp
	}

	return db.Ip
}

func (d *Db) RedirectHome(r *http.Request) string {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.Redirections += 1

	return "http://" + getIpToRedirect(r) + ":9876/"
}

func (d *Db) RedirectJellyfin(r *http.Request) string {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.Redirections += 1

	return "http://" + getIpToRedirect(r) + ":8096/"
}

var db Db

var redirectBaseUrl = "/"

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] - %s", r.RemoteAddr, r.URL.Path)
	if len(db.Ip) == 0 {
		if stats.LastPing.IsZero() {
			fmt.Fprintf(w, "Didn't receive any pings since start %s", stats.StartTime.String())
		} else {
			fmt.Fprintf(w, "Waiting for ping, last ping received: %.0f seconds ago", time.Since(stats.LastPing).Seconds())
		}
	} else {
		http.Redirect(w, r, db.RedirectJellyfin(r), http.StatusTemporaryRedirect)

	}
}

var statsBaseUrl = "/stats"

func statsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	jsonData, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

var updateIpBaseUrl = "/update"

func updateIpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data Db
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stats.Ping()
		db.Ip = data.Ip
		log.Printf("received update: %s", data.Ip)
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusForbidden)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(redirectBaseUrl, redirectHandler)
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
