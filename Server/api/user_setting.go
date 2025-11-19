package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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

	var setting struct {
		AccountId     string `json:"accountId"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		Nickname      string `json:"nickname"`
		CountryRegion string `json:"countryRegion"`
		Area          string `json:"area"`
		Address       string `json:"address"`
		Profile       string `json:"profile"`
		Avatar        string `json:"avatar,omitempty"`
	}
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

	// 创建 users 表（如果不存在）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		name TEXT,
		avatar TEXT,
		job TEXT,
		organization TEXT,
		location TEXT,
		email TEXT,
		introduction TEXT,
		personalWebsite TEXT,
		jobName TEXT,
		organizationName TEXT,
		locationName TEXT,
		phone TEXT,
		registrationDate TEXT,
		accountId TEXT UNIQUE,
		certification INTEGER,
		role TEXT,
		updateTime DATETIME DEFAULT (datetime('now'))
	)`)
	if err != nil {
		http.Error(w, "Failed to create users table", http.StatusInternalServerError)
		return
	}

	// 根据 accountId 更新用户信息，updateTime 自动由 SQLite datetime('now') 生成
	// 修正字段映射顺序以匹配 users 表结构
	// 检查 accountId 是否存在
	var exists int
	_ = db.QueryRow(`SELECT COUNT(1) FROM users WHERE accountId = ?`, setting.AccountId).Scan(&exists)
	if exists == 0 {
		// 如果不存在则插入新记录
		_, err = db.Exec(`INSERT INTO users (accountId, name, email, avatar, jobName, locationName, organizationName, introduction, updateTime) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
			setting.AccountId, setting.Name, setting.Email, setting.Avatar, setting.Nickname, setting.Area, setting.CountryRegion, setting.Profile)
	} else {
		// 如果存在则更新记录
		_, err = db.Exec(`UPDATE users 
			SET name=?, email=?, avatar=?, jobName=?, locationName=?, organizationName=?, introduction=?, updateTime=datetime('now') 
			WHERE accountId=?`,
			setting.Name, setting.Email, setting.Avatar, setting.Nickname, setting.Area, setting.CountryRegion, setting.Profile, setting.AccountId)
	}
	if err != nil {
		http.Error(w, "Failed to insert or update record: "+err.Error(), http.StatusInternalServerError)
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
func thirdDealer(username, password string) error {
	wsURL := "ws://localhost:8081/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 发送登录指令
	if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("login_lnt_password %s %s", username, password))); err != nil {
		return err
	}

	var session string
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		parts := strings.SplitN(string(message), " ", 2)
		prefix := parts[0]
		content := ""
		if len(parts) > 1 {
			content = parts[1]
		}
		if prefix == "Session" {
			session = content
			break
		} else if prefix == "Error" {
			return fmt.Errorf("websocket error: %s", content)
		}
	}

	if !strings.HasPrefix(session, "V2-1") {
		return fmt.Errorf("获取错误")
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte("profile "+session)); err != nil {
		return err
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		parts := strings.SplitN(string(message), " ", 2)
		prefix := parts[0]
		content := ""
		if len(parts) > 1 {
			content = parts[1]
		}
		if prefix == "Profile" {
			var profile struct {
				Username         string `json:"username"`
				Password         string `json:"password"`
				Name             string `json:"name"`
				Avatar           string `json:"avatar"`
				Job              string `json:"job"`
				Organization     string `json:"organization"`
				Location         string `json:"location"`
				Email            string `json:"email"`
				Introduction     string `json:"introduction"`
				PersonalWebsite  string `json:"personalWebsite"`
				JobName          string `json:"jobName"`
				OrganizationName string `json:"organizationName"`
				LocationName     string `json:"locationName"`
				Phone            string `json:"phone"`
				RegistrationDate string `json:"registrationDate"`
				AccountId        string `json:"accountId"`
				Certification    int    `json:"certification"`
				Role             string `json:"role"`
				UpdateTime       string `json:"updateTime"`
			}
			if err := json.Unmarshal([]byte(content), &profile); err != nil {
				return err
			}

			db, err := sql.Open("sqlite3", "./user_info.db")
			if err != nil {
				return err
			}
			defer db.Close()

			_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT UNIQUE,
				password TEXT,
				name TEXT,
				avatar TEXT,
				job TEXT,
				organization TEXT,
				location TEXT,
				email TEXT,
				introduction TEXT,
				personalWebsite TEXT,
				jobName TEXT,
				organizationName TEXT,
				locationName TEXT,
				phone TEXT,
				registrationDate TEXT,
				accountId TEXT UNIQUE,
				certification INTEGER,
				role TEXT,
				updateTime DATETIME DEFAULT (datetime('now'))
			)`)
			if err != nil {
				return err
			}

			_, err = db.Exec(`INSERT OR REPLACE INTO users (
				username, password, name, avatar, job, organization, location, email, introduction, personalWebsite,
				jobName, organizationName, locationName, phone, registrationDate, accountId, certification, role, updateTime
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				username, password, profile.Name, profile.Avatar, profile.Job, profile.Organization, profile.Location, profile.Email,
				profile.Introduction, profile.PersonalWebsite, profile.JobName, profile.OrganizationName, profile.LocationName,
				profile.Phone, profile.RegistrationDate, profile.AccountId, profile.Certification, profile.Role, profile.UpdateTime)
			return err
		} else if prefix == "Error" {
			return fmt.Errorf("profile error: %s", content)
		}
	}
}

func UserAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
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

	// 验证用户是否存在且密码正确并初始化示例数据
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		name TEXT,
		avatar TEXT,
		job TEXT,
		organization TEXT,
		location TEXT,
		email TEXT,
		introduction TEXT,
		personalWebsite TEXT,
		jobName TEXT,
		organizationName TEXT,
		locationName TEXT,
		phone TEXT,
		registrationDate TEXT,
		accountId TEXT UNIQUE,
		certification INTEGER,
		role TEXT,
		updateTime DATETIME DEFAULT (datetime('now'))
	)`)
	if err != nil {
		http.Error(w, "Failed to create users table", http.StatusInternalServerError)
		return
	}

	var accountId string

	// 插入示例数据：admin、reviewer、user，密码为用户名，包含 updateTime
	db.Exec(`INSERT OR IGNORE INTO users (
		username, password, name, avatar, job, organization, location, email, introduction, personalWebsite,
		jobName, organizationName, locationName, phone, registrationDate, accountId, certification, role, updateTime
	) VALUES 
		('admin', 'admin', '管理员', 'https://example.com/admin.png', '系统管理员', 'HCIBGA', '上海',
		'admin@example.com', '系统管理员账户', 'https://hcibga.example.com', '管理员', 'HCIBGA公司', '上海市',
		'13800138000', '2025-01-01', '1', 1, 'admin', datetime('now')),

		('reviewer', 'reviewer', '审核员', 'https://example.com/reviewer.png', '资料审核员', 'HCIBGA', '北京',
		'reviewer@example.com', '负责审核材料的账户', 'https://hcibga.example.com', '审核员', 'HCIBGA公司', '北京市',
		'13900139000', '2025-01-02', '2', 1, 'reviewer', datetime('now')),

		('user', 'user', '普通用户', 'https://example.com/user.png', '工程师', 'HCIBGA', '广州',
		'user@example.com', '普通用户账户', 'https://hcibga.example.com', '前端工程师', 'HCIBGA公司', '广州市',
		'13700137000', '2025-01-03', '3', 0, 'user', datetime('now'))`)

	// 插入后再次查询 accountId 和 role
	var role string
	err = db.QueryRow(`SELECT accountId, role FROM users WHERE username = ?`, req.Username).Scan(&accountId, &role)
	if err != nil {
		if tdErr := thirdDealer(req.Username, req.Password); tdErr != nil {
			http.Error(w, "thirdDealer failed: "+tdErr.Error(), http.StatusInternalServerError)
			return
		}
		_ = db.QueryRow(`SELECT accountId, role FROM users WHERE username = ?`, req.Username).Scan(&accountId, &role)
	}

	// 验证密码
	var dbPassword string
	if err := db.QueryRow(`SELECT password FROM users WHERE username = ?`, req.Username).Scan(&dbPassword); err != nil {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}
	if dbPassword != req.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// 创建表（如果不存在）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS login_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		ip TEXT,
		login_time DATETIME,
		role TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to create table", http.StatusInternalServerError)
		return
	}

	ip := r.RemoteAddr
	_, err = db.Exec(`INSERT INTO login_records (username, ip, login_time, role) VALUES (?, ?, datetime('now'), ?)`, req.Username, ip, role)
	if err != nil {
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		return
	}

	// 查询刚插入的记录返回
	row := db.QueryRow(`SELECT id, username, ip, login_time, role FROM login_records ORDER BY id DESC LIMIT 1`)
	var record struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		IP        string `json:"ip"`
		LoginTime string `json:"login_time"`
		Role      string `json:"role"`
	}
	// 将返回的 id 替换为 accountId
	if err := row.Scan(&record.ID, &record.Username, &record.IP, &record.LoginTime, &record.Role); err != nil {
		http.Error(w, "Failed to fetch record", http.StatusInternalServerError)
		return
	}

	record.ID = accountId

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
func UserLatestActivityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountId string `json:"accountId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.AccountId == "" {
		http.Error(w, "accountId is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 创建或确保存在活动表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_activities (
		id TEXT PRIMARY KEY,
		userId INTEGER,
		type TEXT,
		title TEXT,
		timestamp TEXT,
		description TEXT,
		status TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to create activities table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 初始化示例活动数据（仅用于测试） - 已按要求注释
	/*
		 db.Exec(`INSERT OR IGNORE INTO user_activities (
			id, userId, type, title, timestamp, description, status
		) VALUES
		('act_001', ?, 'material', '提交了材料《项目计划书》', '2025-11-14T10:00:00Z', '提交了项目计划书用于评审', 'completed'),
		('act_002', ?, 'review', '审核了材料《设计文档》', '2025-11-13T15:30:00Z', '设计文档审核通过', 'approved')`, req.AccountId, req.AccountId)
	*/

	rows, err := db.Query(`SELECT id, type, title, timestamp, description, status 
		FROM user_activities WHERE userId = ? ORDER BY timestamp DESC`, req.AccountId)
	if err != nil {
		http.Error(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Activity struct {
		ID          string `json:"id"`
		Type        string `json:"type"`
		Title       string `json:"title"`
		Timestamp   string `json:"timestamp"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	activities := make([]Activity, 0)
	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.ID, &a.Type, &a.Title, &a.Timestamp, &a.Description, &a.Status); err != nil {
			http.Error(w, "Row scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		activities = append(activities, a)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    activities,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type MyProjectRecord struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	PeopleNumber int           `json:"peopleNumber"`
	Contributors []UserSetting `json:"contributors"`
}

func MyProjectListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountId string `json:"accountId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.AccountId == "" {
		http.Error(w, "accountId is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 创建或确保存在项目表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		description TEXT,
		peopleNumber INTEGER,
		ownerAccountId TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to create projects table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 创建或确保存在贡献者表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS project_contributors (
		projectId INTEGER,
		name TEXT,
		email TEXT,
		avatar TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to create contributors table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 在查询前先写入一些示例数据 - 已按要求注释
	/*
		db.Exec(`INSERT OR REPLACE INTO projects (id, name, description, peopleNumber, ownerAccountId) VALUES
			(1, '项目A', '这是项目A的描述', 5, ?),
			(2, '项目B', '这是项目B的描述', 3, ?)`, req.AccountId, req.AccountId)
	*/

	/*
		db.Exec(`INSERT OR REPLACE INTO project_contributors (projectId, name, email, avatar) VALUES
			(1, '张三', 'zhangsan@example.com', 'https://example.com/avatar1.png'),
			(1, '李四', 'lisi@example.com', 'https://example.com/avatar2.png')`)
	*/

	// 查询项目
	rows, err := db.Query(`SELECT id, name, description, peopleNumber FROM projects WHERE ownerAccountId = ?`, req.AccountId)
	if err != nil {
		http.Error(w, "Query projects failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	projects := make([]MyProjectRecord, 0)
	for rows.Next() {
		var p MyProjectRecord
		p.Contributors = make([]UserSetting, 0)
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PeopleNumber); err != nil {
			http.Error(w, "Row scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 查询贡献者
		cRows, err := db.Query(`SELECT name, email, avatar FROM project_contributors WHERE projectId = ?`, p.ID)
		if err != nil {
			http.Error(w, "Query contributors failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		for cRows.Next() {
			var c UserSetting
			if err := cRows.Scan(&c.Name, &c.Email, &c.Avatar); err != nil {
				http.Error(w, "Contributor scan failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
			p.Contributors = append(p.Contributors, c)
		}
		cRows.Close()

		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    projects,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type MyTeamRecord struct {
	ID           int    `json:"id"`
	Avatar       string `json:"avatar"`
	Name         string `json:"name"`
	PeopleNumber int    `json:"peopleNumber"`
}

func MyTeamListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountId string `json:"accountId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.AccountId == "" {
		http.Error(w, "accountId is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 创建或确保存在团队表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS teams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		avatar TEXT,
		name TEXT,
		peopleNumber INTEGER,
		ownerAccountId TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to create teams table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 在查询前先写入一些示例数据 - 已按要求注释
	/*
		db.Exec(`INSERT OR REPLACE INTO teams (id, avatar, name, peopleNumber, ownerAccountId) VALUES
			(101, 'https://example.com/team1.png', '研发一组', 8, ?),
			(102, 'https://example.com/team2.png', '设计团队', 5, ?)`, req.AccountId, req.AccountId)
	*/

	// 查询团队
	rows, err := db.Query(`SELECT id, avatar, name, peopleNumber FROM teams WHERE ownerAccountId = ?`, req.AccountId)
	if err != nil {
		http.Error(w, "Query teams failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	teams := make([]MyTeamRecord, 0)
	for rows.Next() {
		var t MyTeamRecord
		if err := rows.Scan(&t.ID, &t.Avatar, &t.Name, &t.PeopleNumber); err != nil {
			http.Error(w, "Row scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		teams = append(teams, t)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    teams,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type EnterpriseCertificationModel struct {
	AccountType          int    `json:"accountType"`
	Status               int    `json:"status"`
	Time                 string `json:"time"`
	LegalPerson          string `json:"legalPerson"`
	CertificateType      string `json:"certificateType"`
	AuthenticationNumber string `json:"authenticationNumber"`
	EnterpriseName       string `json:"enterpriseName"`
}

type CertificationRecord struct {
	CertificationType    int    `json:"certificationType"`
	CertificationContent string `json:"certificationContent"`
	Status               int    `json:"status"`
	Time                 string `json:"time"`
}

type UnitCertification struct {
	EnterpriseInfo EnterpriseCertificationModel `json:"enterpriseInfo"`
	Record         []CertificationRecord        `json:"record"`
}

func UserCertificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountId string `json:"accountId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.AccountId == "" {
		http.Error(w, "accountId is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 创建或确保存在企业认证表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS enterprise_certification (
		accountId TEXT PRIMARY KEY,
		accountType INTEGER,
		status INTEGER,
		time TEXT,
		legalPerson TEXT,
		certificateType TEXT,
		authenticationNumber TEXT,
		enterpriseName TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to create enterprise_certification table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 创建或确保存在认证记录表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS certification_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		accountId TEXT,
		certificationType INTEGER,
		certificationContent TEXT,
		status INTEGER,
		time TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to create certification_records table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 插入示例数据 - 已按要求注释
	/*
		db.Exec(`INSERT OR REPLACE INTO enterprise_certification (accountId, accountType, status, time, legalPerson, certificateType, authenticationNumber, enterpriseName) VALUES
			(?, 1, 2, '2025-11-14 10:00:00', '王五', '营业执照', '91320100MA1X0XXXX', '示例科技有限公司')`, req.AccountId)
	*/

	/*
		db.Exec(`INSERT OR REPLACE INTO certification_records (accountId, certificationType, certificationContent, status, time) VALUES
			(?, 1, '企业营业执照', 2, '2025-11-14 10:00:00'),
			(?, 2, '税务登记证', 1, '2025-11-10 09:00:00')`, req.AccountId, req.AccountId)
	*/

	// 查询企业认证信息
	var enterpriseInfo EnterpriseCertificationModel
	row := db.QueryRow(`SELECT accountType, status, time, legalPerson, certificateType, authenticationNumber, enterpriseName FROM enterprise_certification WHERE accountId = ?`, req.AccountId)
	if err := row.Scan(&enterpriseInfo.AccountType, &enterpriseInfo.Status, &enterpriseInfo.Time, &enterpriseInfo.LegalPerson, &enterpriseInfo.CertificateType, &enterpriseInfo.AuthenticationNumber, &enterpriseInfo.EnterpriseName); err != nil {
		http.Error(w, "Failed to fetch enterprise certification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 查询认证记录
	rows, err := db.Query(`SELECT certificationType, certificationContent, status, time FROM certification_records WHERE accountId = ?`, req.AccountId)
	if err != nil {
		http.Error(w, "Failed to fetch certification records: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	records := make([]CertificationRecord, 0)
	for rows.Next() {
		var rec CertificationRecord
		if err := rows.Scan(&rec.CertificationType, &rec.CertificationContent, &rec.Status, &rec.Time); err != nil {
			http.Error(w, "Record scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": UnitCertification{
			EnterpriseInfo: enterpriseInfo,
			Record:         records,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

var qrStore sync.Map

type qrData struct {
	Status  string
	Session string
}

func setQRData(qrId string, data *qrData) {
	qrStore.Store(qrId, data)
}

func getQRData(qrId string) (*qrData, bool) {
	val, ok := qrStore.Load(qrId)
	if !ok {
		return nil, false
	}
	return val.(*qrData), ok
}

const aesChars = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"

func randomString(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	result := make([]byte, length)
	for i := range result {
		idx := r.Intn(len(aesChars))
		result[i] = aesChars[idx]
	}

	return string(result)
}

func QRCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	wsURL := "ws://localhost:8081/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		http.Error(w, "Failed to connect websocket: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte("login_lnt_qr")); err != nil {
		http.Error(w, "Failed to send get_qrcode: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var qrUrl, qrId string
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, "Failed to read message: "+err.Error(), http.StatusInternalServerError)
			return
		}
		parts := strings.Split(string(message), " ")
		if parts[0] == "QrCodeId" {
			uuid := parts[1]
			qrUrl = "https://ids.xmu.edu.cn/authserver/qrCode/getCode?uuid=" + uuid
			qrId = randomString(16)
			break
		}
	}

	go func() {
		defer conn.Close()
		data := &qrData{Status: "pending"}
		setQRData(qrId, data)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				data.Status = "expired"
				return
			}
			parts := strings.Split(string(message), " ")
			if parts[0] == "Error" && parts[1] == "QrCodeExpired" {
				data.Status = "expired"
				return
			} else if parts[0] == "Session" {
				data.Status = "done"
				data.Session = parts[1]
				return
			}
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]string{
			"qrUrl": qrUrl,
			"qrId":  qrId,
		},
	})
}

func QRStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	qrId := r.URL.Query().Get("qrId")
	if qrId == "" {
		http.Error(w, "qrId is required", http.StatusBadRequest)
		return
	}

	data, ok := getQRData(qrId)
	if !ok {
		http.Error(w, "QR code not found", http.StatusNotFound)
		return
	}

	if data.Status == "expired" {
		defer qrStore.Delete(qrId)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]string{
			"status": data.Status,
		},
	})
}

func QRAuthResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		QrId string `json:"qrId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.QrId == "" {
		http.Error(w, "qrId is required", http.StatusBadRequest)
		return
	}

	data, ok := getQRData(req.QrId)
	if !ok {
		http.Error(w, "QR code not found", http.StatusNotFound)
		return
	}
	defer qrStore.Delete(req.QrId)

	if data.Status != "done" {
		http.Error(w, "QR code not done", http.StatusBadRequest)
		return
	}

	wsURL := "ws://localhost:8081/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		http.Error(w, "Failed to connect websocket: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, []byte("profile_lnt "+data.Session)); err != nil {
		http.Error(w, "Failed to send profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, "Failed to read message: "+err.Error(), http.StatusInternalServerError)
			return
		}
		parts := strings.SplitN(string(message), " ", 2)
		prefix := parts[0]
		content := ""
		if len(parts) > 1 {
			content = parts[1]
		}
		if prefix == "Profile" {
			var profile struct {
				Username         string `json:"username"`
				Password         string `json:"password"`
				Name             string `json:"name"`
				Avatar           string `json:"avatar"`
				Job              string `json:"job"`
				Organization     string `json:"organization"`
				Location         string `json:"location"`
				Email            string `json:"email"`
				Introduction     string `json:"introduction"`
				PersonalWebsite  string `json:"personalWebsite"`
				JobName          string `json:"jobName"`
				OrganizationName string `json:"organizationName"`
				LocationName     string `json:"locationName"`
				Phone            string `json:"phone"`
				RegistrationDate string `json:"registrationDate"`
				AccountId        string `json:"accountId"`
				Certification    int    `json:"certification"`
				Role             string `json:"role"`
				UpdateTime       string `json:"updateTime"`
			}
			if err := json.Unmarshal([]byte(content), &profile); err != nil {
				http.Error(w, "Invalid profile data: "+err.Error(), http.StatusBadRequest)
				return
			}

			db, err := sql.Open("sqlite3", "./user_info.db")
			if err != nil {
				http.Error(w, "Failed to open database: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer db.Close()

			_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT UNIQUE,
				password TEXT,
				name TEXT,
				avatar TEXT,
				job TEXT,
				organization TEXT,
				location TEXT,
				email TEXT,
				introduction TEXT,
				personalWebsite TEXT,
				jobName TEXT,
				organizationName TEXT,
				locationName TEXT,
				phone TEXT,
				registrationDate TEXT,
				accountId TEXT UNIQUE,
				certification INTEGER,
				role TEXT,
				updateTime DATETIME DEFAULT (datetime('now'))
			)`)
			if err != nil {
				http.Error(w, "Failed to create table: "+err.Error(), http.StatusInternalServerError)
				return
			}

			var existingId int
			err = db.QueryRow(`SELECT id FROM users WHERE username = ?`, profile.Username).Scan(&existingId)
			if err != nil && err != sql.ErrNoRows {
				_, err = db.Exec(`INSERT OR REPLACE INTO users (
				username, password, name, avatar, job, organization, location, email, introduction, personalWebsite,
				jobName, organizationName, locationName, phone, registrationDate, accountId, certification, role, updateTime
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					profile.Username, randomString(16), profile.Name, profile.Avatar, profile.Job, profile.Organization, profile.Location, profile.Email,
					profile.Introduction, profile.PersonalWebsite, profile.JobName, profile.OrganizationName, profile.LocationName,
					profile.Phone, profile.RegistrationDate, profile.AccountId, profile.Certification, profile.Role, profile.UpdateTime)
				if err != nil {
					http.Error(w, "Failed to insert profile: "+err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 0,
				"msg":  "QR auth success",
				"data": map[string]string{
					"token": profile.AccountId,
					"role":  profile.Role,
				},
			})
			return
		} else if prefix == "Error" {
			http.Error(w, "Profile error: "+content, http.StatusBadRequest)
			return
		}
	}
}

func RegisterUserSettingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/user/save-info", SaveUserInfoHandler)
	mux.HandleFunc("/api/user/auth", UserAuthHandler)
	mux.HandleFunc("/api/user/upload", UploadUserFileHandler)
	mux.HandleFunc("/api/user/latest-activity", UserLatestActivityHandler)
	mux.HandleFunc("/api/user/my-project/list", MyProjectListHandler)
	mux.HandleFunc("/api/user/my-team/list", MyTeamListHandler)
	mux.HandleFunc("/api/user/certification", UserCertificationHandler)
	mux.HandleFunc("/api/user/qr-code", QRCodeHandler)
	mux.HandleFunc("/api/user/qr-status", QRStatusHandler)
	mux.HandleFunc("/api/user/qr-auth-result", QRAuthResultHandler)
}
