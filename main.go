package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Cấu trúc để lưu thông tin máy chủ
type Server struct {
	Address string
	Weight  int
}

// Load balancer cấu trúc
type LoadBalancer struct {
	servers       []Server
	algorithm     string
	mutex         sync.Mutex
	current       int
	connections   map[string]int // Số lượng kết nối hiện tại cho mỗi máy chủ
}

// Đọc danh sách máy chủ từ file
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
		address := parts[0]
		weight := 1
		if len(parts) > 1 {
			weight, _ = strconv.Atoi(parts[1])
		}
		servers = append(servers, Server{Address: address, Weight: weight})
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	lb.servers = servers
	lb.connections = make(map[string]int)
	return nil
}

// Tự động tải lại danh sách máy chủ sau mỗi khoảng thời gian
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
		}
	}
}

// Chọn máy chủ tiếp theo dựa trên thuật toán
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
			if lb.connections[s.Address] < minConnections {
				minConnections = lb.connections[s.Address]
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

// Xử lý yêu cầu và chuyển tiếp đến máy chủ backend
func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	server := lb.NextServer()
	if server == nil {
		http.Error(w, "No servers available", http.StatusServiceUnavailable)
		return
	}

	// Tạo một yêu cầu mới để gửi đến máy chủ backend
	req, err := http.NewRequest(r.Method, fmt.Sprintf("http://%s%s", server.Address, r.RequestURI), r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Sao chép các header từ yêu cầu gốc
	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error contacting server", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Cập nhật số lượng kết nối cho máy chủ
	lb.connections[server.Address]++

	// Sao chép các header từ phản hồi của máy chủ backend
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	// Giảm số lượng kết nối sau khi xử lý xong
	lb.connections[server.Address]--
}

func main() {
	lb := &LoadBalancer{
		algorithm: "round_robin", // Hoặc "least_connections", "weighted_round_robin", "ip_hash"
	}

	err := lb.LoadServers("servers.conf")
	if err != nil {
		log.Fatalf("Error loading servers: %v", err)
	}

	go lb.ReloadServersPeriodically("servers.conf", 5*time.Second)

	// HTTP server
	go func() {
		log.Println("Starting HTTP server on port 8080...")
		err := http.ListenAndServe(":8080", http.HandlerFunc(lb.HandleRequest))
		if err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// HTTPS server
	go func() {
		log.Println("Starting HTTPS server on port 8443...")
		err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", http.HandlerFunc(lb.HandleRequest))
		if err != nil {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()

	// Keep the application running
	select {}
}
