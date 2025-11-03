package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

// UserSetting 表示用户设置信息
type UserSetting struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar,omitempty"`
}

// SaveUserInfoHandler 保存用户信息
func SaveUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var setting UserSetting
	if err := json.NewDecoder(r.Body).Decode(&setting); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 模拟保存用户信息
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "user info saved",
		"data": setting,
	})
}

// UserAuthHandler 用户认证
func UserAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 模拟认证成功
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "auth success",
	})
}

// UploadUserFileHandler 上传用户文件
func UploadUserFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	saveDir := filepath.Join(os.TempDir(), "hci_user_upload")
	os.MkdirAll(saveDir, os.ModePerm)
	savePath := filepath.Join(saveDir, handler.Filename)

	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "file uploaded",
		"path": savePath,
	})
}

// RegisterUserSettingRoutes 注册用户设置相关路由
func RegisterUserSettingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/user/save", SaveUserInfoHandler)
	mux.HandleFunc("/api/user/auth", UserAuthHandler)
	mux.HandleFunc("/api/user/upload", UploadUserFileHandler)
}
