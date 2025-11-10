package api

import (
	"database/sql"
	"encoding/json"

	//"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// UploadCheckHandler - 检查上传状态
func UploadCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":  "ready",
		"message": "Upload service is available.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// UploadFileHandler - 文件上传接口
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File upload error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempDir := filepath.Join(os.TempDir(), "hci_upload_material")
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		http.Error(w, "Cannot create temp dir", http.StatusInternalServerError)
		return
	}

	dstPath := filepath.Join(tempDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Cannot create destination file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"filename": handler.Filename,
		"size":     handler.Size,
		"path":     dstPath,
		"result":   "success",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// MaterialLLMFillHandler - 材料自动填充接口（改为数据库操作）
func MaterialLLMFillHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	materialId, _ := req["materialId"].(string)
	if materialId == "" {
		http.Error(w, "Missing materialId", http.StatusBadRequest)
		return
	}

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

	// 更新材料状态为已填充
	_, err = db.Exec(`UPDATE materials SET status = 'filled' WHERE id = ?`, materialId)
	if err != nil {
		http.Error(w, "Update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"materialId": materialId,
		"status":     "filled",
		"message":    "Material auto-filled successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// MaterialUploadHandler - 材料上传完成接口（改为数据库操作）
func MaterialUploadHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	id, _ := req["id"].(string)
	title, _ := req["title"].(string)
	category, _ := req["category"].(string)
	fileSize, _ := req["fileSize"].(float64)
	uploader, _ := req["uploader"].(string)

	if id == "" || title == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

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

	// 插入材料记录
	_, err = db.Exec(`INSERT INTO materials (id, title, category, status, fileSize, uploader, uploadTime) VALUES (?, ?, ?, 'uploaded', ?, ?, datetime('now'))`,
		id, title, category, int64(fileSize), uploader)
	if err != nil {
		http.Error(w, "Insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"uploadId": id,
		"status":   "uploaded",
		"title":    title,
		"category": category,
		"fileSize": fileSize,
		"uploader": uploader,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RegisterMaterialUploadRoutes 注册所有上传相关接口
func RegisterMaterialUploadRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/upload/check", UploadCheckHandler)
	mux.HandleFunc("/api/upload/file", UploadFileHandler)
	mux.HandleFunc("/api/material/llm-fill", MaterialLLMFillHandler)
	mux.HandleFunc("/api/material/upload", MaterialUploadHandler)
}
