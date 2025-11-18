package main

import (
	"log"
	"net/http"

	"github.com/vintcessun/HCIBGA/Server/api"
)

func main() {
	mux := http.NewServeMux()

	// 添加全局日志中间件
	muxWithLogging := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[访问日志] %s %s 来自 %s User-Agent: %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		// 处理所有OPTIONS请求以支持CORS预检
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}
		// 设置CORS响应头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		mux.ServeHTTP(w, r)
	})

	// 注册各模块路由
	api.RegisterSubmitMaterialBatchRoutes(mux)
	api.RegisterUserInfoRoutes(mux)
	api.RegisterUserSettingRoutes(mux)

	// 注册其他 API 路由
	api.RegisterExportRoutes(mux)
	api.RegisterInfoImportRoutes(mux)
	api.RegisterMaterialListRoutes(mux)
	api.RegisterMaterialUploadRoutes(mux)
	api.RegisterLLMRoutes(mux)

	log.Println("Server started at :8000")
	if err := http.ListenAndServe("127.0.0.1:8000", muxWithLogging); err != nil {
		log.Fatal(err)
	}
}
