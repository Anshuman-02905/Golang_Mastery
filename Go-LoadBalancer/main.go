package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// Server interface defines the basic methods for a backend server
// Each server should provide its address, a health check, and serve requests

type Server interface {
	GetAddr() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

// SimpleServer represents a single backend server
// It contains an address and a reverse proxy for forwarding requests

type SimpleServer struct {
	Addr  string
	Proxy *httputil.ReverseProxy
}

// LoadBalancer holds a list of backend servers and manages request distribution

type LoadBalancer struct {
	Port            string
	Servers         []Server
	RoundRobincount int
}

// handle_err logs and terminates the program on error

func handle_err(err error) {
	log.Fatalf("Error occurred: %v \n", err)
	os.Exit(1)
}

// newLoadBalancer initializes a load balancer with the provided servers

func newLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		Port:            port,
		Servers:         servers,
		RoundRobincount: 0,
	}
}

// newSimpleServer creates a SimpleServer with a reverse proxy

func newSimpleServer(addr string) *SimpleServer {
	serverURL, err := url.Parse(addr)
	if err != nil {
		handle_err(err)
		return nil
	}

	return &SimpleServer{
		Addr:  addr,
		Proxy: httputil.NewSingleHostReverseProxy(serverURL),
	}
}

// GetAddr returns the server address

func (s *SimpleServer) GetAddr() string { return s.Addr }

// IsAlive returns true indicating that the server is always considered alive

func (s *SimpleServer) IsAlive() bool { return true }

// Serve forwards the request to the reverse proxy

func (s *SimpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(rw, r)
}

// GetNextAvailableServer selects the next available server using Round-Robin scheduling

func (lb *LoadBalancer) GetNextAvailableServer() Server {
	server := lb.Servers[lb.RoundRobincount%len(lb.Servers)]
	if !server.IsAlive() {
		lb.RoundRobincount++
		server = lb.Servers[lb.RoundRobincount%len(lb.Servers)]
	}
	lb.RoundRobincount++
	return server
}

// ServeProxy forwards the request to the selected backend server

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.GetNextAvailableServer()
	fmt.Printf("Forwarding request to: %v \n", targetServer.GetAddr())
	targetServer.Serve(rw, r)
}

func main() {
	// Define the list of backend servers
	servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.bing.com"),
		newSimpleServer("https://www.duckduckgo.com"),
	}

	// Initialize the load balancer
	lb := newLoadBalancer("8000", servers)

	// Define request handler to forward traffic via the load balancer
	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(rw, r)
	}

	http.HandleFunc("/", handleRedirect)
	fmt.Printf("Load balancer running on port %s \n", lb.Port)
	http.ListenAndServe(":"+lb.Port, nil)
}
