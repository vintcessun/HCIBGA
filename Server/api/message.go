package api

import (
	"encoding/json"
	"net/http"
)

// MessageListRequest 定义消息列表请求体
// 当前仅校验 accountId 是否存在，未来可接入数据库或其他数据源
// 暂时返回空数据占位
func MessageListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountID string `json:"accountId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.AccountID == "" {
		http.Error(w, "accountId is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "请求成功",
		"data": []interface{}{},
	})
}

// MessageReadHandler 处理消息已读请求
// 当前阶段不做任何状态修改，仅校验参数并返回固定成功
func MessageReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountID string  `json:"accountId"`
		IDs       []int64 `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.AccountID == "" {
		http.Error(w, "accountId is required", http.StatusBadRequest)
		return
	}
	if len(req.IDs) == 0 {
		http.Error(w, "ids is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "请求成功",
		"data": true,
	})
}

// RegisterMessageRoutes 注册消息中心路由
func RegisterMessageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/message/list", MessageListHandler)
	mux.HandleFunc("/api/message/read", MessageReadHandler)
}
