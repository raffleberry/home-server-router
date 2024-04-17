package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func logg(r *http.Request, returnCode int, str string) {
	log.Printf("%s : %s - %d - %s", r.Method, r.RequestURI, returnCode, str)
}

func getClientIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")
	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs,
		// meaning the last one is the actual client IP `splitIps[len(splitIps)-1]`
		// i want the user's ip which is the first 0
		netIP := net.ParseIP(splitIps[0])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}

func noIpResponse(w http.ResponseWriter) {
	if stats.LastPing.IsZero() {
		fmt.Fprintf(w, "Didn't receive any pings since start %s", stats.StartTime.String())
	} else {
		fmt.Fprintf(w, "Waiting for ping, last ping received: %.0f seconds ago", time.Since(stats.LastPing).Seconds())
	}
}
