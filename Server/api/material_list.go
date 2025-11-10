// Package api 实现材料列表相关接口
package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// 获取材料列表
func GetMaterialListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 连接 SQLite 数据库
	db, err := sql.Open("sqlite3", "./hci.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 确保 materials 表存在
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS materials (
		id TEXT PRIMARY KEY,
		title TEXT,
		category TEXT,
		status TEXT,
		fileSize INTEGER,
		uploader TEXT,
		uploadTime TEXT,
		reviewer TEXT,
		reviewTime TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to ensure table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 查询所有材料记录
	rows, err := db.Query(`SELECT id, title, category, status, fileSize, uploader, uploadTime, reviewer, reviewTime FROM materials`)
	if err != nil {
		http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var id, title, category, status, uploader, uploadTime, reviewer, reviewTime string
		var fileSize int64
		if err := rows.Scan(&id, &title, &category, &status, &fileSize, &uploader, &uploadTime, &reviewer, &reviewTime); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, map[string]interface{}{
			"id":         id,
			"title":      title,
			"category":   category,
			"status":     status,
			"fileSize":   fileSize,
			"uploader":   uploader,
			"uploadTime": uploadTime,
			"reviewer":   reviewer,
			"reviewTime": reviewTime,
		})
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

	// 连接 SQLite 数据库
	db, err := sql.Open("sqlite3", "./hci.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 删除指定材料记录
	result, err := db.Exec(`DELETE FROM materials WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Material not found", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"code":    200,
		"message": "删除成功",
		"id":      id,
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

	// 连接 SQLite 数据库
	db, err := sql.Open("sqlite3", "./hci.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 更新指定材料的审核状态和意见
	_, err = db.Exec(`UPDATE materials SET status = ?, reviewer = ?, reviewTime = datetime('now') WHERE id = ?`,
		req.Status, req.Comment, req.MaterialId)
	if err != nil {
		http.Error(w, "Update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"code":    200,
		"message": "审核成功",
		"id":      req.MaterialId,
		"status":  req.Status,
		"comment": req.Comment,
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
