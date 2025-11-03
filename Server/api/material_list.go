// Package api 实现材料列表相关接口
package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// 获取材料列表
func GetMaterialListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 模拟返回数据
	data := []map[string]interface{}{
		{
			"id":         "uuid-001",
			"title":      "项目申报材料",
			"category":   "document",
			"status":     "approved",
			"fileSize":   102400,
			"uploader":   "张三",
			"uploadTime": "2025-11-03T10:00:00Z",
			"reviewer":   "李四",
			"reviewTime": "2025-11-03T12:00:00Z",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// 删除材料
func DeleteMaterialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/material/")
	id := strings.TrimSpace(path)
	if id == "" {
		http.Error(w, "Missing material ID", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"code":    200,
		"message": "删除成功",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// 审核材料
func ReviewMaterialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type ReviewRequest struct {
		MaterialId string `json:"materialId"`
		Status     string `json:"status"`
		Comment    string `json:"comment"`
	}

	var req ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"code":    200,
		"message": "审核成功",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RegisterMaterialListRoutes 注册材料列表路由
func RegisterMaterialListRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/material/list", GetMaterialListHandler)
	mux.HandleFunc("/api/material/review", ReviewMaterialHandler)
	mux.HandleFunc("/api/material/", DeleteMaterialHandler)
}
