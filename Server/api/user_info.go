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

	var projects []Project
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

	var activities []Activity
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

	var teams []Team
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

// RegisterUserInfoRoutes 注册用户信息相关路由
func RegisterUserInfoRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/user/projects", UserProjectsHandler)
	mux.HandleFunc("/api/user/activities", UserActivitiesHandler)
	mux.HandleFunc("/api/user/teams", UserTeamsHandler)
}
