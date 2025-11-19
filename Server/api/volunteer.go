package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

const volunteerWSURL = "ws://localhost:8081/ws"

// VolunteerHoursRequest 请求体
// username/password 来源于前端表单，accountId 由当前登录用户 token 附带
type VolunteerHoursRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	AccountID string `json:"accountId"`
}

// VolunteerHoursResponse 响应体（透传志愿汇字段）
type VolunteerHoursResponse struct {
	CreditHours float64 `json:"credit_hours"`
	HonorHours  float64 `json:"honor_hours"`
	TotalHours  float64 `json:"total_hours"`
}

// VolunteerCreditHandler 处理 /api/volunteer/credit
func VolunteerCreditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VolunteerHoursRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.AccountID = strings.TrimSpace(req.AccountID)
	if req.Username == "" || req.Password == "" || req.AccountID == "" {
		http.Error(w, "用户名、密码和 accountId 均不能为空", http.StatusBadRequest)
		return
	}

	displayName, err := getDisplayNameByAccount(req.AccountID)
	if err != nil {
		http.Error(w, "获取用户姓名失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hours, err := fetchVolunteerHours(req.Username, req.Password, displayName)
	if err != nil {
		http.Error(w, "查询志愿时数失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":   0,
		"status": "ok",
		"msg":    "请求成功",
		"data":   hours,
	})
}

func fetchVolunteerHours(username, password, displayName string) (*VolunteerHoursResponse, error) {
	token, err := loginVolunteer(username, password)
	if err != nil {
		return nil, err
	}
	return queryVolunteerHours(token, displayName)
}

func loginVolunteer(username, password string) (string, error) {
	conn, _, err := websocket.DefaultDialer.Dial(volunteerWSURL, nil)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	payload := fmt.Sprintf("login_zyh365 %s %s", username, password)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
		return "", err
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return "", err
		}
		prefix, content := splitVolunteerMessage(string(message))
		switch prefix {
		case "Token":
			return strings.TrimSpace(content), nil
		case "Error":
			return "", errors.New(strings.TrimSpace(content))
		}
	}
}

func queryVolunteerHours(token, displayName string) (*VolunteerHoursResponse, error) {
	conn, _, err := websocket.DefaultDialer.Dial(volunteerWSURL, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	payload := fmt.Sprintf("hours_zyh365 %s %s", token, displayName)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
		return nil, err
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return nil, err
		}
		prefix, content := splitVolunteerMessage(string(message))
		switch prefix {
		case "Hours":
			var resp VolunteerHoursResponse
			if err := json.Unmarshal([]byte(content), &resp); err != nil {
				return nil, err
			}
			return &resp, nil
		case "Error":
			return nil, errors.New(strings.TrimSpace(content))
		}
	}
}

func splitVolunteerMessage(message string) (string, string) {
	trimmed := strings.TrimSpace(message)
	parts := strings.SplitN(trimmed, " ", 2)
	prefix := parts[0]
	content := ""
	if len(parts) > 1 {
		content = parts[1]
	}
	return prefix, content
}

func getDisplayNameByAccount(accountID string) (string, error) {
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		return "", fmt.Errorf("数据库连接失败: %w", err)
	}
	defer db.Close()

	var name string
	err = db.QueryRow(`SELECT name FROM users WHERE accountId = ?`, accountID).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("未找到 accountId=%s 对应的用户", accountID)
		}
		return "", fmt.Errorf("查询用户姓名失败: %w", err)
	}

	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("当前账号未设置姓名，请先完善用户信息")
	}

	return strings.TrimSpace(name), nil
}

// RegisterVolunteerRoutes 注册志愿接口
func RegisterVolunteerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/volunteer/credit", VolunteerCreditHandler)
}
