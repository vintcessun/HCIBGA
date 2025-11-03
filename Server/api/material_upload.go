package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

// MaterialLLMFillHandler - 材料自动填充接口（模拟 LLM）
func MaterialLLMFillHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	resp := map[string]interface{}{
		"filledData": fmt.Sprintf("Auto-filled data for material: %v", req["materialId"]),
		"status":     "completed",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// MaterialUploadHandler - 材料上传完成接口
func MaterialUploadHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"uploadId": fmt.Sprintf("mat-%d", 1000+os.Getpid()),
		"status":   "uploaded",
		"details":  req,
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
