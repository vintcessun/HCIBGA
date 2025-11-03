package main

import (
	"log"
	"net/http"

	"github.com/vintcessun/HCIBGA/Server/api"
)

func main() {
	mux := http.NewServeMux()

	// 注册各模块路由
	api.RegisterSubmitMaterialBatchRoutes(mux)
	api.RegisterUserInfoRoutes(mux)
	api.RegisterUserSettingRoutes(mux)

	// 注册其他 API 路由
	api.RegisterExportRoutes(mux)
	api.RegisterInfoImportRoutes(mux)
	api.RegisterMaterialListRoutes(mux)
	api.RegisterMaterialUploadRoutes(mux)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
