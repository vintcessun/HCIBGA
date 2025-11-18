package api

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	//"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type FileCheckRequest struct {
	Md5       string `json:"md5"`
	Filename  string `json:"filename"`
	AccountId string `json:"accountId"`
}

type FileCheckResponse struct {
	Exists bool   `json:"exists"`
	FileID string `json:"file_id,omitempty"`
	URL    string `json:"url,omitempty"`
}

// UploadCheckHandler - 检查文件是否已存在
func UploadCheckHandler(w http.ResponseWriter, r *http.Request) {
	var req FileCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 打开数据库
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 确保 file_map 表存在
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS file_map (
		md5 TEXT PRIMARY KEY,
		filename TEXT,
		file_id TEXT,
		url TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to ensure table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 查询是否存在
	var fileID, url string
	err = db.QueryRow(`SELECT file_id, url FROM file_map WHERE md5 = ? AND filename = ?`, req.Md5, req.Filename).Scan(&fileID, &url)
	if err == sql.ErrNoRows {
		resp := FileCheckResponse{Exists: false}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    resp,
		})
		return
	} else if err != nil {
		http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := FileCheckResponse{
		Exists: true,
		FileID: fileID,
		URL:    url,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    0,
		"data":    resp,
		"message": "success",
	})
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

	// 直接保存到当前目录下的 upload 文件夹
	uploadDir := filepath.Join(".", "upload")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Cannot create upload dir", http.StatusInternalServerError)
		return
	}

	// 计算文件 md5（保持与前端SparkMD5.ArrayBuffer一致，读取原始字节流）
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file for md5 calculation", http.StatusInternalServerError)
		return
	}
	hash := md5.Sum(fileBytes)
	md5Str := fmt.Sprintf("%x", hash[:])
	finalName := md5Str + "-" + handler.Filename
	finalPath := filepath.Join(uploadDir, finalName)

	dst, err := os.Create(finalPath)
	if err != nil {
		http.Error(w, "Cannot create destination file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := dst.Write(fileBytes); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// 写入映射关系到数据库
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	_, err = db.Exec(`INSERT OR REPLACE INTO file_map (md5, filename, file_id, url) VALUES (?, ?, ?, ?)`,
		md5Str, handler.Filename, finalName, "/upload/"+finalName)
	if err != nil {
		http.Error(w, "Failed to insert file map: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"file_id": finalName,
		"url":     "/upload/" + finalName,
		"md5":     md5Str,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    resp,
	})
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

type MaterialUploadRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Files       []string `json:"files"`
	AccountId   string   `json:"accountId"`
}

// MaterialUploadHandler - 材料上传完成接口（改为数据库操作）
func MaterialUploadHandler(w http.ResponseWriter, r *http.Request) {
	var req MaterialUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	id := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v-%v", time.Now().UnixNano(), req.Title))))
	title := req.Title
	description := req.Description
	category := req.Category
	uploader := req.AccountId

	if title == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 使用 user_info.db 以保持与材料列表数据一致
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 确保 materials 表存在，并包含完整的前端需求字段
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS materials (
		id TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		category TEXT,
		tags TEXT,
		files TEXT, -- 存储多文件信息的JSON字符串
		status TEXT,
		uploader TEXT,
		uploadTime TEXT,
		reviewer TEXT,
		reviewTime TEXT,
		reviewComment TEXT,
		aiScore REAL,
		aiConfidence REAL,
		aiSuggestions TEXT,
		aiRiskLevel TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to ensure table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// tags数组转为字符串
	var tagsStr string
	var tagsArr []string
	for _, ts := range req.Tags {
		tagsArr = append(tagsArr, ts)
	}
	tagsStr = strings.Join(tagsArr, ",")

	// files数组转为JSON字符串存储
	var filesJSON string
	var filesArr []map[string]interface{}
	for _, fname := range req.Files {
		info, err := os.Stat(filepath.Join("upload", fname))
		var fSize int64
		if err == nil {
			fSize = info.Size()
		}
		filesArr = append(filesArr, map[string]interface{}{
			"fileUrl":  "/upload/" + fname,
			"fileName": fname,
			"fileSize": fSize,
		})
	}
	b, _ := json.Marshal(filesArr)
	filesJSON = string(b)

	score, err := CalculateScore(&req)
	if err != nil {
		http.Error(w, "Score calculation error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Calculated LLM score:", score)

	// 插入材料记录到数据库, 状态设为pending，初始化 reviewer 为 "" 表示未审核
	_, err = db.Exec(`INSERT INTO materials (
		id, title, description, category, tags, files, status, uploader, uploadTime,
		reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel
	) VALUES (?, ?, ?, ?, ?, ?, 'pending', ?, datetime('now'),
		'', '', '', ?, ?, ?, ?)`,
		id, title, description, category, tagsStr, filesJSON, uploader,
		score.AiScore, score.AiConfidence, score.AiSuggestions, score.AiRiskLevel)
	if err != nil {
		http.Error(w, "Insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"id":          id,
		"title":       title,
		"description": description,
		"category":    category,
		"tags":        strings.Split(tagsStr, ","),
		"files":       json.RawMessage(filesJSON),
		"status":      "pending",
		"uploader":    uploader,
		"uploadTime":  time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": resp,
		"code": 0,
		"msg":  "success",
	})
}

// RegisterMaterialUploadRoutes 注册所有上传相关接口
func ServeUploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 获取文件名
	filename := strings.TrimPrefix(r.URL.Path, "/upload/")
	filename = strings.TrimSpace(filename)
	fmt.Println("Serving file:", filename)
	if filename == "" {
		http.Error(w, "Missing filename", http.StatusBadRequest)
		return
	}

	// 拼接路径并检查文件是否存在
	filePath := filepath.Join("upload", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// 返回文件内容
	http.ServeFile(w, r, filePath)
}

func RegisterMaterialUploadRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/upload/check", UploadCheckHandler)
	mux.HandleFunc("/api/upload/file", UploadFileHandler)
	mux.HandleFunc("/api/material/llm-fill", MaterialLLMFillHandler)
	mux.HandleFunc("/api/material/upload", MaterialUploadHandler)
	mux.HandleFunc("/upload/", ServeUploadFileHandler)
}
