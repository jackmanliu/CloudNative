package main

import (
	"fmt"
	"github.com/jackmanliu/CloudNative/module10/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	os.Setenv("VERSION", "v1.0")
}

func rootHandler(rw http.ResponseWriter, r *http.Request) {
	// 设置访问延时
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	delay := rand.Intn(3000)
	time.Sleep(time.Millisecond * time.Duration(delay))

	// 接收客户端 request，并将 request 中带的 header 写入 response header
	for k, v := range r.Header {
		for _, value := range v {
			log.Printf("set key:%s and value:%s in response header\n", k, value)
			rw.Header().Add(k, value)
		}
	}

	// 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	version := os.Getenv("VERSION")
	rw.Header().Add("VERSION", version)
	log.Printf("set VERTION: %s in response header\n", version)

	// Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	clientIP := getClientIP(r)
	log.Printf("Acess successfully. client ip:%s Response code: %d\n", clientIP, http.StatusOK)

	rw.Write([]byte(fmt.Sprintf("<h1>Respond in %d ms<h1>", delay)))
}

// 当访问 localhost/healthz 时，应返回 200
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
	log.Print("Starting http server")
	metrics.Register()
	server := http.NewServeMux()
	server.HandleFunc("/", rootHandler)
	server.HandleFunc("/healthz", healthCheck)
	server.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("Http server startup failed: %s\n", err.Error())
	}
}
