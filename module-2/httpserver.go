package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func init() {
	os.Setenv("VERSION", "v1.0")
}

func rootHandler(rw http.ResponseWriter, r *http.Request) {
	// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
	for k, v := range r.Header {
		for _, value := range v {
			log.Printf("set key:%s and value:%s in response header\n", k, value)
			rw.Header().Add(k, value)
		}
	}

	// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	version := os.Getenv("VERSION")
	rw.Header().Add("VERSION", version)
	log.Printf("set VERTION: %s in response header\n", version)

	//3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	clientIP := getClientIP(r)
	log.Printf("Acess successfully. client ip:%s Response code: %d\n", clientIP, http.StatusOK)
}

// 4. 当访问 localhost/healthz 时，应返回 200
func healthCheck(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "server is sound and responsive with return code:", 200)
}

// 获取客户端IP地址
func getClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func main() {
	server := http.NewServeMux()
	server.HandleFunc("/", rootHandler)
	server.HandleFunc("/healthz", healthCheck)
	err := http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("Http server startup failure: %s\n", err.Error())
	}
}
