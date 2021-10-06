package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//启动http server
func startHttpServer() {
	err := http.ListenAndServe("127.0.0.1:8081", nil)
	if err != nil {
		fmt.Printf("http server failed, err:%v\n", err)
		return
	}
}

// 将用户请求的request header写到response header
func index(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("VERSION")
	requestHeader := r.Header
	for key, value := range requestHeader {
		s := strings.Join(value, "")
		w.Header().Set(key, s)
	}
	// 根据环境变量获取version.如:V1.0.0
	w.Header().Set("Version", version)
	_, err := w.Write([]byte("ok"))
	logInfo(r, err)
}

//心跳检测
func healthz(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(strconv.Itoa(http.StatusOK)))
	logInfo(r, err)

}

//记录请求信息
func logInfo(request *http.Request, err error) {
	// httpCode:200
	statusCode := http.StatusOK
	if err != nil {
		// httpCode:500
		statusCode = http.StatusInternalServerError
	}
	ip := ip(request)
	fmt.Printf("客服端ip：%s，http状态码:%d", ip, statusCode)
}

// 获取IP地址,可获取反向代理的IP
func ip(req *http.Request) string {
	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = strings.Split(req.RemoteAddr, ":")[0]
		}
	}
	return addr
}
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/healthz", healthz)
	// 启动http服务器,listen:8081
	startHttpServer()
}
