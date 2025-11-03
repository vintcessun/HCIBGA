// Package api 实现信息导入模块的接口
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadDir 临时文件保存目录
var UploadDir = filepath.Join(os.TempDir(), "hci_info_import")

// 导入 Excel 文件
func ImportExcelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	if ext != ".xls" && ext != ".xlsx" {
		resp := map[string]interface{}{
			"code":    400,
			"message": "文件格式错误，仅支持 .xls/.xlsx",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	os.MkdirAll(UploadDir, 0755)
	dstPath := filepath.Join(UploadDir, fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename))
	dst, _ := os.Create(dstPath)
	defer dst.Close()
	io.Copy(dst, file)

	resp := map[string]interface{}{
		"code":    200,
		"message": "Excel 文件已成功接收",
		"data": map[string]string{
			"taskId": fmt.Sprintf("import_%d", time.Now().Unix()),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// 导入 TXT 文件或直接文本
func ImportTxtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	var lines int
	var errors []string

	file, handler, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		ext := filepath.Ext(handler.Filename)
		if ext != ".txt" {
			resp := map[string]interface{}{
				"code":    400,
				"message": "TXT 文件格式错误或内容为空",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		os.MkdirAll(UploadDir, 0755)
		dstPath := filepath.Join(UploadDir, fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename))
		dst, _ := os.Create(dstPath)
		defer dst.Close()
		io.Copy(dst, file)
		lines = 123 // 模拟行数
	} else {
		// 如果 file 为空，尝试读取 text 字段内容
		text := r.FormValue("text")
		if len(text) == 0 {
			resp := map[string]interface{}{
				"code":    400,
				"message": "TXT 文件格式错误或内容为空",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		lines = len([]rune(text))
	}

	resp := map[string]interface{}{
		"code":    200,
		"message": "导入成功",
		"data": map[string]interface{}{
			"lines":  lines,
			"errors": errors,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RegisterInfoImportRoutes 注册信息导入相关路由
func RegisterInfoImportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/info/import/excel", ImportExcelHandler)
	mux.HandleFunc("/api/info/import/txt", ImportTxtHandler)
}
