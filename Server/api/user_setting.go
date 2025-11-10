package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// UserSetting 表示用户设置信息
type UserSetting struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar,omitempty"`
}

// SaveUserInfoHandler 保存用户信息（真实数据逻辑）
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

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 创建表（如果不存在）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_settings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		avatar TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to create table", http.StatusInternalServerError)
		return
	}

	// 插入记录（如果已存在则更新）
	_, err = db.Exec(`INSERT INTO user_settings (name, email, avatar) VALUES (?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET email=excluded.email, avatar=excluded.avatar`,
		setting.Name, setting.Email, setting.Avatar)
	if err != nil {
		http.Error(w, "Failed to insert or update record", http.StatusInternalServerError)
		return
	}

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

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 使用 database/sql 记录登录信息
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 创建表（如果不存在）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS login_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		ip TEXT,
		login_time DATETIME
	)`)
	if err != nil {
		http.Error(w, "Failed to create table", http.StatusInternalServerError)
		return
	}

	ip := r.RemoteAddr
	_, err = db.Exec(`INSERT INTO login_records (username, ip, login_time) VALUES (?, ?, datetime('now'))`, req.Username, ip)
	if err != nil {
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		return
	}

	// 查询刚插入的记录返回
	row := db.QueryRow(`SELECT id, username, ip, login_time FROM login_records ORDER BY id DESC LIMIT 1`)
	var record struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		IP        string `json:"ip"`
		LoginTime string `json:"login_time"`
	}
	if err := row.Scan(&record.ID, &record.Username, &record.IP, &record.LoginTime); err != nil {
		http.Error(w, "Failed to fetch record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "auth success",
		"data": record,
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

	// 将上传记录保存到数据库
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS upload_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		path TEXT,
		upload_time DATETIME
	)`)
	if err != nil {
		http.Error(w, "Failed to create table", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`INSERT INTO upload_records (filename, path, upload_time) VALUES (?, ?, datetime('now'))`, handler.Filename, savePath)
	if err != nil {
		http.Error(w, "Failed to insert upload record", http.StatusInternalServerError)
		return
	}

	// 查询刚插入的记录返回
	row := db.QueryRow(`SELECT id, filename, path, upload_time FROM upload_records ORDER BY id DESC LIMIT 1`)
	var record struct {
		ID         int    `json:"id"`
		Filename   string `json:"filename"`
		Path       string `json:"path"`
		UploadTime string `json:"upload_time"`
	}
	if err := row.Scan(&record.ID, &record.Filename, &record.Path, &record.UploadTime); err != nil {
		http.Error(w, "Failed to fetch upload record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "file uploaded",
		"data": record,
	})
}

// RegisterUserSettingRoutes 注册用户设置相关路由
func RegisterUserSettingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/user/save", SaveUserInfoHandler)
	mux.HandleFunc("/api/user/auth", UserAuthHandler)
	mux.HandleFunc("/api/user/upload", UploadUserFileHandler)
}
