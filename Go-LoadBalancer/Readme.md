# Building a Simple Load Balancer in Go Using Round Robin

## Introduction
Load balancers are crucial in modern software development. If you've ever wondered how requests are distributed across multiple servers or why certain websites remain fast during heavy traffic, the answer often lies in efficient load balancing.

In this project, we'll build a simple application load balancer using the Round Robin algorithm in Go. The goal is to understand how a load balancer works under the hood, step by step.

## What is a Load Balancer?
A load balancer is a system that distributes incoming network traffic across multiple servers. It ensures that no single server bears too much load, preventing bottlenecks and improving the overall user experience. Load balancing also ensures that if one server fails, traffic can be automatically re-routed to another available server, reducing downtime and increasing availability.

## Why Use Load Balancers?
- **High Availability:** Ensures traffic is routed to healthy servers even if one server fails.
- **Scalability:** Allows horizontal scaling by adding more servers as traffic increases.
- **Efficiency:** Maximizes resource utilization by distributing workload equally.

## Load Balancing Algorithms
There are different strategies to distribute traffic:
- **Round Robin:** Sequentially distributes requests among servers in order.
- **Weighted Round Robin:** Assigns a weight to each server to control traffic distribution.
- **Least Connections:** Routes traffic to the server with the fewest active connections.
- **IP Hashing:** Selects a server based on the client's IP address.

For this project, we'll focus on implementing a **Round Robin load balancer**.

## What is the Round Robin Algorithm?
A Round Robin algorithm sends each incoming request to the next available server in a circular manner. If server A handles the first request, server B will handle the second, and server C will handle the third. Once all servers have received a request, it starts again from server A.

## Implementation in Go

### Features
- Uses Round Robin scheduling to distribute requests.
- Supports multiple backend servers.
- Implements a reverse proxy to forward traffic.
- Skips unavailable servers to ensure reliability.

### How It Works
1. We define backend servers (`facebook.com`, `bing.com`, and `duckduckgo.com`).
2. The load balancer cycles through these servers in a Round Robin manner.
3. Each request is forwarded to the next server in the list.
4. If a server is unavailable, the load balancer skips it and moves to the next.
5. The process repeats indefinitely, ensuring equal distribution of traffic.

### Code Example
Below is the core implementation of our load balancer:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	GetAddr() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type SimpleServer struct {
	Addr  string
	Proxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	Port            string
	Servers         []Server
	RoundRobincount int
}

func handle_err(err error) {
	log.Fatalf("Error occurred: %v \n", err)
	os.Exit(1)
}

func newLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		Port:            port,
		Servers:         servers,
		RoundRobincount: 0,
	}
}

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

func (s *SimpleServer) GetAddr() string { return s.Addr }
func (s *SimpleServer) IsAlive() bool { return true }
func (s *SimpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(rw, r)
}

func (lb *LoadBalancer) GetNextAvailableServer() Server {
	server := lb.Servers[lb.RoundRobincount%len(lb.Servers)]
	if !server.IsAlive() {
		lb.RoundRobincount++
		server = lb.Servers[lb.RoundRobincount%len(lb.Servers)]
	}
	lb.RoundRobincount++
	return server
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.GetNextAvailableServer()
	fmt.Printf("Forwarding request to: %v \n", targetServer.GetAddr())
	targetServer.Serve(rw, r)
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.bing.com"),
		newSimpleServer("https://www.duckduckgo.com"),
	}
	lb := newLoadBalancer("8000", servers)
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(rw, r)
	})
	fmt.Printf("Load balancer running on port %s \n", lb.Port)
	http.ListenAndServe(":"+lb.Port, nil)
}
```

## Running the Load Balancer
1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/load-balancer-go.git
   cd load-balancer-go
   ```
2. Run the Go application:
   ```sh
   go run main.go
   ```
3. The load balancer will start on port `8000` and forward traffic to backend servers.

## Enhancements
Want to improve this implementation? Try adding:
- **Health Checks:** Dynamically remove unhealthy servers.
- **Weighted Round Robin:** Distribute traffic based on server capacity.
- **Least Connections Strategy:** Optimize request distribution for efficiency.

## Conclusion
This project provides a foundational understanding of how load balancing works. By mastering these concepts, you can design more robust and scalable distributed systems. Happy coding!

## Credits
This project was inspired by [this YouTube tutorial](https://www.youtube.com/watch?v=ZSDYx9eOiqo).

---

## License
This project is open-source and available under the MIT License.

## Author
Developed by [Your Name](https://github.com/yourusername)

