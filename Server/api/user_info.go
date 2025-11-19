package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Project 表示用户项目
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Activity 表示用户动态
type Activity struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

// Team 表示用户团队
type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserProjectsHandler 获取用户项目列表
func UserProjectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM user_projects")
	if err != nil {
		http.Error(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	projects := make([]Project, 0)
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			http.Error(w, "Row scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": projects,
	})
}

// UserActivitiesHandler 获取用户动态
func UserActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	type Activity struct {
		ID      string `json:"id"`
		Content string `json:"content"`
		Time    string `json:"time"`
	}

	// 连接 SQLite 数据库
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, content, time FROM user_activities")
	if err != nil {
		http.Error(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	activities := make([]Activity, 0)
	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.ID, &a.Content, &a.Time); err != nil {
			http.Error(w, "Row scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		activities = append(activities, a)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": activities,
	})
}

// UserTeamsHandler 获取用户团队
func UserTeamsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM user_teams")
	if err != nil {
		http.Error(w, "Query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	teams := make([]Team, 0)
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			http.Error(w, "Row scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		teams = append(teams, t)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": teams,
	})
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
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

	// 确保用户表存在并包含role字段
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
		accountId TEXT,
		certification INTEGER,
		role TEXT,
		updateTime DATETIME DEFAULT (datetime('now'))
	)`)
	if err != nil {
		http.Error(w, "Failed to create users table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 初始化记录
	db.Exec(`INSERT OR IGNORE INTO users (
			username, password, name, avatar, job, organization, location, email, introduction, personalWebsite,
			jobName, organizationName, locationName, phone, registrationDate, accountId, certification, role
		) VALUES
			('admin', 'admin', '管理员', 'https://example.com/admin.png', '系统管理员', 'HCIBGA', '上海',
			'admin@example.com', '系统管理员账户', 'https://hcibga.example.com', '管理员', 'HCIBGA公司', '上海市',
			'13800138000', '2025-01-01', '1', 1, 'admin'),

			('reviewer', 'reviewer', '审核员', 'https://example.com/reviewer.png', '资料审核员', 'HCIBGA', '北京',
			'reviewer@example.com', '负责审核材料的账户', 'https://hcibga.example.com', '审核员', 'HCIBGA公司', '北京市',
			'13900139000', '2025-01-02', '2', 1, 'reviewer'),

			('reviewer1', 'reviewer', '审核员', 'https://example.com/reviewer.png', '资料审核员', 'HCIBGA', '北京',
			'reviewer@example.com', '负责审核材料的账户', 'https://hcibga.example.com', '审核员', 'HCIBGA公司', '北京市',
			'13900139000', '2025-01-02', '21', 1, 'reviewer'),

			('reviewer2', 'reviewer', '审核员', 'https://example.com/reviewer.png', '资料审核员', 'HCIBGA', '北京',
			'reviewer@example.com', '负责审核材料的账户', 'https://hcibga.example.com', '审核员', 'HCIBGA公司', '北京市',
			'13900139000', '2025-01-02', '22', 1, 'reviewer'),

			('reviewer3', 'reviewer', '审核员', 'https://example.com/reviewer.png', '资料审核员', 'HCIBGA', '北京',
			'reviewer@example.com', '负责审核材料的账户', 'https://hcibga.example.com', '审核员', 'HCIBGA公司', '北京市',
			'13900139000', '2025-01-02', '23', 1, 'reviewer'),

			('user', 'user', '普通用户', 'https://example.com/user.png', '工程师', 'HCIBGA', '广州',
			'user@example.com', '普通用户账户', 'https://hcibga.example.com', '前端工程师', 'HCIBGA公司', '广州市',
			'13700137000', '2025-01-03', '3', 0, 'user')`)

	// 根据 accountId 查询完整用户记录
	row := db.QueryRow(`SELECT 
		name, avatar, job, organization, location, email, introduction, personalWebsite,
		jobName, organizationName, locationName, phone, registrationDate, accountId, certification, role, updateTime
		FROM users WHERE accountId = ?`, req.AccountId)

	var name, avatar, job, organization, location, email, introduction, personalWebsite string
	var jobName, organizationName, locationName, phone, registrationDate, accountId, role, updateTime string
	var certification int
	if err := row.Scan(&name, &avatar, &job, &organization, &location, &email, &introduction, &personalWebsite,
		&jobName, &organizationName, &locationName, &phone, &registrationDate, &accountId, &certification, &role, &updateTime); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userState := map[string]interface{}{
		"name":             name,
		"avatar":           avatar,
		"job":              job,
		"organization":     organization,
		"location":         location,
		"email":            email,
		"introduction":     introduction,
		"personalWebsite":  personalWebsite,
		"jobName":          jobName,
		"organizationName": organizationName,
		"locationName":     locationName,
		"phone":            phone,
		"registrationDate": registrationDate,
		"accountId":        accountId,
		"certification":    certification,
		"role":             role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    userState,
	})
}

// RegisterUserInfoRoutes 注册用户信息相关路由
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": nil,
		"msg":  "退出登录成功",
	})
}

func RegisterUserInfoRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/user/logout", LogoutHandler)
	mux.HandleFunc("/api/user/projects", UserProjectsHandler)
	mux.HandleFunc("/api/user/activities", UserActivitiesHandler)
	mux.HandleFunc("/api/user/teams", UserTeamsHandler)
	mux.HandleFunc("/api/user/info", UserInfoHandler)
}
