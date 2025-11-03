package api

import (
	"encoding/json"
	"net/http"
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

	projects := []Project{
		{ID: "P001", Name: "项目一"},
		{ID: "P002", Name: "项目二"},
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

	activities := []Activity{
		{ID: "A001", Content: "提交了材料", Time: "2025-11-01 10:00"},
		{ID: "A002", Content: "审核通过", Time: "2025-11-02 15:30"},
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

	teams := []Team{
		{ID: "T001", Name: "团队一"},
		{ID: "T002", Name: "团队二"},
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
