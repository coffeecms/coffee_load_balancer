package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Server structure to hold server information
type Server struct {
	Address string // e.g., "10.220.3.1"
	Port    string // e.g., "23252"
	Weight  int    // e.g., 10
}

// LoadBalancer structure
type LoadBalancer struct {
	servers         []Server
	algorithm       string
	mutex           sync.Mutex
	current         int
	connections     map[string]int // Number of current connections for each server
	rateLimit       int             // Requests per second
	shutdownChannel chan struct{}
}

// Load server list from a file
func (lb *LoadBalancer) LoadServers(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var servers []Server
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue // Skip malformed lines
		}
		address := parts[0]
		port := parts[1]
		weight := 1
		if len(parts) > 2 {
			weight, _ = strconv.Atoi(parts[2])
		}
		servers = append(servers, Server{Address: address, Port: port, Weight: weight})
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	lb.servers = servers
	lb.connections = make(map[string]int)
	return nil
}

// Automatically reload server list periodically
func (lb *LoadBalancer) ReloadServersPeriodically(filename string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Reloading server list...")
			err := lb.LoadServers(filename)
			if err != nil {
				log.Printf("Error reloading servers: %v", err)
			}
		case <-lb.shutdownChannel:
			return
		}
	}
}

// Rate limiting middleware
func (lb *LoadBalancer) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(time.Second / time.Duration(lb.rateLimit)):
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	})
}

// Select the next server based on the algorithm
func (lb *LoadBalancer) NextServer() *Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	var server *Server

	switch lb.algorithm {
	case "round_robin":
		server = &lb.servers[lb.current]
		lb.current = (lb.current + 1) % len(lb.servers)
	case "least_connections":
		minConnections := int(^uint(0) >> 1) // Max int value
		for _, s := range lb.servers {
			if lb.connections[s.Address+":"+s.Port] < minConnections {
				minConnections = lb.connections[s.Address+":"+s.Port]
				server = &s
			}
		}
	case "weighted_round_robin":
		totalWeight := 0
		for _, s := range lb.servers {
			totalWeight += s.Weight
		}
		randWeight := rand.Intn(totalWeight)
		for _, s := range lb.servers {
			if randWeight < s.Weight {
				server = &s
				break
			}
			randWeight -= s.Weight
		}
	case "ip_hash":
		// Implement IP Hashing if needed
	}

	return server
}

// Handle incoming requests and forward them to backend servers
func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	server := lb.NextServer()
	if server == nil {
		http.Error(w, "No servers available", http.StatusServiceUnavailable)
		return
	}

	serverURL := fmt.Sprintf("http://%s:%s", server.Address, server.Port)

	// Create a new request to forward to the backend server
	req, err := http.NewRequest(r.Method, fmt.Sprintf("%s%s", serverURL, r.RequestURI), r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error contacting server", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Update the number of connections for the server
	lb.connections[server.Address+":"+server.Port]++

	// Copy headers from the backend server response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	// Decrease the number of connections after processing
	lb.connections[server.Address+":"+server.Port]--
}

// Handle requests using goroutines
func (lb *LoadBalancer) requestHandler(w http.ResponseWriter, r *http.Request) {
	go lb.HandleRequest(w, r)
}

func main() {
	lb := &LoadBalancer{
		algorithm:       "round_robin", // Choose algorithm as needed
		rateLimit:       1000,          // Example rate limit: 1000 requests per second
		shutdownChannel: make(chan struct{}),
	}

	err := lb.LoadServers("servers.conf")
	if err != nil {
		log.Fatalf("Error loading servers: %v", err)
	}

	go lb.ReloadServersPeriodically("servers.conf", 5*time.Second)

	// HTTP server
	//go func() {
	//	http.Handle("/", lb.RateLimit(http.HandlerFunc(lb.requestHandler)))
	//	log.Println("Starting HTTP server on port 8080...")
	//	err := http.ListenAndServe(":8080", nil)
	//	if err != nil {
	//		log.Fatalf("HTTP server error: %v", err)
	//	}
	//}()

	// HTTPS server
	go func() {
		http.Handle("/", lb.RateLimit(http.HandlerFunc(lb.requestHandler)))
		log.Println("Starting HTTPS server on port 443...")
		err := http.ListenAndServeTLS(":443", "lowlevelforest.com.crt", "lowlevelforest.com.key", nil)
		if err != nil {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	close(lb.shutdownChannel)
	log.Println("Shutting down gracefully...")
}
