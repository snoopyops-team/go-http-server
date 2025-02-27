package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func getClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For获取代理IP链
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For可能包含多个IP（代理链），取第一个客户端IP
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}
		}
	}

	// 尝试从X-Real-IP获取
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" && net.ParseIP(realIP) != nil {
		return realIP
	}

	// 直接从RemoteAddr获取（最后手段）
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func handler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	url := r.URL.String()

	// 打印日志
	fmt.Printf("Client IP: %-15s | Method: %-6s | URL: %s\n", clientIP, r.Method, url)

	w.Write([]byte("Request logged"))
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
