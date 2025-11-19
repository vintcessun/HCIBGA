// Package api 实现材料列表相关接口
package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// 获取材料列表
func GetMaterialListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 连接 SQLite 数据库
	// 使用 user_info.db 存储材料数据以与用户信息保持一致
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 在查询之前插入示例数据以验证查询逻辑
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS materials (
		id TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		category TEXT,
		tags TEXT,
		files TEXT, -- 存储多文件信息的JSON字符串
		status TEXT,
		uploader TEXT,
		uploadTime TEXT,
		reviewer TEXT,
		reviewTime TEXT,
		reviewComment TEXT,
		aiScore REAL,
		aiConfidence REAL,
		aiSuggestions TEXT,
		aiRiskLevel TEXT
	)`)

	/*
		// 从 auth 中获取当前用户 ID 以插入示例数据
		authDB, _ := sql.Open("sqlite3", "./user_info.db")
		defer authDB.Close()
		var sampleUploader string
		_ = authDB.QueryRow(`SELECT accountId FROM users WHERE username = (SELECT username FROM login_records ORDER BY id DESC LIMIT 1)`).Scan(&sampleUploader)

		fmt.Printf("调试: sampleUploader = %s\n", sampleUploader)
		if _, err := db.Exec(`INSERT OR REPLACE INTO materials
			(id, title, description, category, tags, files, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			"approved_sample_1", "示例材料", "这是一个通过审核的示例材料", "教材", "教育,数学",
			`[{"fileUrl":"/uploads/files/sample.pdf","fileName":"sample.pdf","fileSize":102400},{"fileUrl":"/uploads/files/sample_extra.pdf","fileName":"sample_extra.pdf","fileSize":2048}]`,
			"approved", sampleUploader, "2025-11-15T08:00:00Z", "reviewer01", "2025-11-15T09:00:00Z", "资料齐全", 0.95, 0.9, "增加封面,优化目录", "low"); err != nil {
			fmt.Println("插入 approved_sample_1 出错:", err)
		}

		if _, err := db.Exec(`INSERT OR REPLACE INTO materials
			(id, title, description, category, tags, files, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			"pending_sample_1", "另一个示例材料", "这是另一个正在审核的示例材料", "论文", "科研,计算机",
			`[{"fileUrl":"/uploads/files/sample2.pdf","fileName":"sample2.pdf","fileSize":204800}]`,
			"pending", sampleUploader, "2025-11-16T10:00:00Z", "reviewer02", "2025-11-16T11:00:00Z", "内容详实", 0.92, 0.85, "补充参考文献,优化格式", "medium"); err != nil {
			fmt.Println("插入 pending_sample_1 出错:", err)
		}

		if _, err := db.Exec(`INSERT OR REPLACE INTO materials
			(id, title, description, category, tags, files, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			"rejected_sample_1", "被拒绝的示例材料", "这是一个被拒绝的示例材料", "报告", "工作总结",
			`[{"fileUrl":"/uploads/files/sample3.pdf","fileName":"sample3.pdf","fileSize":51200}]`,
			"rejected", sampleUploader, "2025-11-17T12:00:00Z", "reviewer03", "2025-11-17T13:00:00Z", "内容不完整", 0.4, 0.3, "补充数据分析,完善结论", "high"); err != nil {
			fmt.Println("插入 rejected_sample_1 出错:", err)
		}
	*/

	// 确保 materials 表存在，并包含完整字段（多文件JSON存储）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS materials (
		id TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		category TEXT,
		tags TEXT,
		files TEXT, -- 存储多文件信息的JSON字符串
		status TEXT,
		uploader TEXT,
		uploadTime TEXT,
		reviewer TEXT,
		reviewTime TEXT,
		reviewComment TEXT,
		aiScore REAL,
		aiConfidence REAL,
		aiSuggestions TEXT,
		aiRiskLevel TEXT
	)`)
	if err != nil {
		http.Error(w, "Failed to ensure table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 从 auth 中获取当前用户 ID
	authDB, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Auth database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer authDB.Close()
	var accountId string
	row := authDB.QueryRow(`SELECT accountId FROM users WHERE username = (SELECT username FROM login_records ORDER BY id DESC LIMIT 1)`)
	if err := row.Scan(&accountId); err != nil {
		http.Error(w, "Failed to get accountId from auth: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// 查询role
	var role string
	if err := authDB.QueryRow(`SELECT role FROM users WHERE accountId = ?`, accountId).Scan(&role); err != nil {
		http.Error(w, "Failed to get user role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var rows *sql.Rows
	if role == "user" {
		rows, err = db.Query(`SELECT id, title, description, category, tags, files, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel FROM materials WHERE uploader = ?`, accountId)
	} else {
		rows, err = db.Query(`SELECT id, title, description, category, tags, files, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel FROM materials`)
	}
	if err != nil {
		http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	data := make([]map[string]interface{}, 0)

	for rows.Next() {
		var id, title, description, category, tagsStr, filesJSON, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiSuggestionsStr, aiRiskLevel string
		var aiScore, aiConfidence float64
		if err := rows.Scan(
			&id,
			&title,
			&description,
			&category,
			&tagsStr,
			&filesJSON, // 多文件 JSON 字符串
			&status,
			&uploader,
			&uploadTime,
			&reviewer,
			&reviewTime,
			&reviewComment,
			&aiScore,
			&aiConfidence,
			&aiSuggestionsStr,
			&aiRiskLevel,
		); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// 将tags和aiSuggestions从字符串转为数组
		tags := strings.Split(tagsStr, ",")
		aiSuggestions := strings.Split(aiSuggestionsStr, ",")

		// 解析files JSON
		var files []map[string]interface{}
		_ = json.Unmarshal([]byte(filesJSON), &files)

		data = append(data, map[string]interface{}{
			"id":            id,
			"title":         title,
			"description":   description,
			"category":      category,
			"tags":          tags,
			"files":         files,
			"status":        status,
			"uploader":      uploader,
			"uploadTime":    uploadTime,
			"reviewer":      reviewer,
			"reviewTime":    reviewTime,
			"reviewComment": reviewComment,
			"aiReviewResult": map[string]interface{}{
				"score":       aiScore,
				"confidence":  aiConfidence,
				"suggestions": aiSuggestions,
				"riskLevel":   aiRiskLevel,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"code": 0,
		"data": data,
		"msg":  "success",
	}
	json.NewEncoder(w).Encode(resp)
}

// 删除材料
func DeleteMaterialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/material/")
	id := strings.TrimSpace(path)
	if id == "" {
		http.Error(w, "Missing material ID", http.StatusBadRequest)
		return
	}

	// 连接 SQLite 数据库
	// 使用 user_info.db 存储材料数据以与用户信息保持一致
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 删除指定材料记录
	result, err := db.Exec(`DELETE FROM materials WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Material not found", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"code":    200,
		"message": "删除成功",
		"id":      id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// 审核材料
func ReviewMaterialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type ReviewRequest struct {
		MaterialId string `json:"materialId"`
		Status     string `json:"status"`
		Comment    string `json:"comment"`
	}

	var req ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 连接 SQLite 数据库
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 获取当前记录的reviewer和status
	var currentReviewer, currentStatus string
	err = db.QueryRow(`SELECT reviewer, status FROM materials WHERE id = ?`, req.MaterialId).Scan(&currentReviewer, &currentStatus)
	if err != nil {
		http.Error(w, "Query current reviewer/status error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 将reviewer字段转为切片
	var reviewers []string
	if strings.TrimSpace(currentReviewer) != "" {
		reviewers = strings.Split(currentReviewer, ",")
	}

	// 获取当前用户accountId作为review人
	authDB2, _ := sql.Open("sqlite3", "./user_info.db")
	defer authDB2.Close()
	var reviewerName string
	_ = authDB2.QueryRow(`SELECT accountId FROM users WHERE username = (SELECT username FROM login_records ORDER BY id DESC LIMIT 1)`).Scan(&reviewerName)

	// 如果status是rejected，直接更新为rejected
	if strings.ToLower(req.Status) == "rejected" {
		if !containsReviewer(reviewers, reviewerName) {
			reviewers = append(reviewers, reviewerName)
		}
		_, err = db.Exec(`UPDATE materials SET status = 'rejected', reviewer = ?, reviewTime = datetime('now'), reviewComment = ? WHERE id = ?`,
			strings.Join(reviewers, ","), req.Comment, req.MaterialId)
		if err != nil {
			http.Error(w, "Update rejected error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if strings.ToLower(req.Status) == "approved" {
		// 如果status是approved，先添加reviewer但不立即改状态，直到3人
		if !containsReviewer(reviewers, reviewerName) {
			go func() {
				err := DealMaterialToRecord(req.MaterialId)
				if err != nil {
					fmt.Println("模型写入记录错误", err)
				}
			}()
			reviewers = append(reviewers, reviewerName)
		}
		newStatus := currentStatus
		if len(reviewers) >= 3 {
			err := DealMaterialToRecord(req.MaterialId)
			if err != nil {
				http.Error(w, "Update Record error"+err.Error(), http.StatusInternalServerError)
				return
			}
			newStatus = "approved"
		}
		_, err = db.Exec(`UPDATE materials SET status = ?, reviewer = ?, reviewTime = datetime('now'), reviewComment = ? WHERE id = ?`,
			newStatus, strings.Join(reviewers, ","), req.Comment, req.MaterialId)
		if err != nil {
			http.Error(w, "Update approved/pending error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// 其他状态只更新reviewer和评论
		if !containsReviewer(reviewers, reviewerName) {
			reviewers = append(reviewers, reviewerName)
		}
		_, err = db.Exec(`UPDATE materials SET reviewer = ?, reviewTime = datetime('now'), reviewComment = ? WHERE id = ?`,
			strings.Join(reviewers, ","), req.Comment, req.MaterialId)
		if err != nil {
			http.Error(w, "Update other status error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp := map[string]interface{}{
		"code":    200,
		"message": "审核成功",
		"id":      req.MaterialId,
		"status":  req.Status,
		"comment": req.Comment,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func containsReviewer(list []string, name string) bool {
	for _, v := range list {
		if strings.TrimSpace(v) == strings.TrimSpace(name) {
			return true
		}
	}
	return false
}

func GetPendingMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountId string `json:"accountId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.AccountId == "" {
		http.Error(w, "Missing accountId", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 确保表存在
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS materials (
		id TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		category TEXT,
		tags TEXT,
		files TEXT, -- 存储多文件信息的JSON字符串
		status TEXT,
		uploader TEXT,
		uploadTime TEXT,
		reviewer TEXT,
		reviewTime TEXT,
		reviewComment TEXT,
		aiScore REAL,
		aiConfidence REAL,
		aiSuggestions TEXT,
		aiRiskLevel TEXT
	)`)

	/*
		// 在查询前插入示例数据（模仿 user_setting.go 中的逻辑，改为多文件JSON存储）
		if _, err := db.Exec(`INSERT OR REPLACE INTO materials
			(id, title, description, category, tags, files, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			"pending_sample_for_user", "示例材料", "这是一个待审核的示例材料", "教材", "教育,数学",
			`[{"fileUrl":"/uploads/files/sample.pdf","fileName":"sample.pdf","fileSize":102400},{"fileUrl":"/uploads/files/sample_extra.pdf","fileName":"sample_extra.pdf","fileSize":2048}]`,
			"pending", req.AccountId, "2025-11-15T08:00:00Z", "teacher01", "2025-11-15T09:00:00Z", "资料齐全", 0.95, 0.9, "增加封面,优化目录", "low"); err != nil {
			http.Error(w, "Failed to insert sample data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	// 查询role
	var role string
	if err := db.QueryRow(`SELECT role FROM users WHERE accountId = ?`, req.AccountId).Scan(&role); err != nil {
		http.Error(w, "Failed to get user role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var rows *sql.Rows
	if role == "user" {
		rows, err = db.Query(`SELECT id, title, description, category, tags, files, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel 
			FROM materials WHERE status = 'pending' AND uploader = ?`, req.AccountId)
	} else {
		rows, err = db.Query(`SELECT id, title, description, category, tags, files, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel 
			FROM materials WHERE status = 'pending'`)
	}
	if err != nil {
		http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, title, description, category, tagsStr, filesJSON, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiSuggestionsStr, aiRiskLevel string
		var aiScore, aiConfidence float64
		if err := rows.Scan(
			&id,
			&title,
			&description,
			&category,
			&tagsStr,
			&filesJSON, // 多文件 JSON 字符串
			&status,
			&uploader,
			&uploadTime,
			&reviewer,
			&reviewTime,
			&reviewComment,
			&aiScore,
			&aiConfidence,
			&aiSuggestionsStr,
			&aiRiskLevel,
		); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 如果pending的reviewer里包含当前accountId，则跳过
		if status == "pending" && strings.Contains(reviewer, req.AccountId) {
			continue
		}

		tags := strings.Split(tagsStr, ",")
		aiSuggestions := strings.Split(aiSuggestionsStr, ",")

		var files []map[string]interface{}
		_ = json.Unmarshal([]byte(filesJSON), &files)

		data = append(data, map[string]interface{}{
			"id":            id,
			"title":         title,
			"description":   description,
			"category":      category,
			"tags":          tags,
			"files":         files,
			"status":        status,
			"uploader":      uploader,
			"uploadTime":    uploadTime,
			"reviewer":      reviewer,
			"reviewTime":    reviewTime,
			"reviewComment": reviewComment,
			"aiReviewResult": map[string]interface{}{
				"score":       aiScore,
				"confidence":  aiConfidence,
				"suggestions": aiSuggestions,
				"riskLevel":   aiRiskLevel,
			},
			"accountId": req.AccountId,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// RegisterMaterialListRoutes 注册材料列表路由
func GetMaterialStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountId string `json:"accountId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.AccountId == "" {
		http.Error(w, "Missing accountId", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var role string
	if err := db.QueryRow(`SELECT role FROM users WHERE accountId = ?`, req.AccountId).Scan(&role); err != nil {
		http.Error(w, "Failed to get user role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var total, pending, approved, rejected int
	byCategory := make(map[string]int)

	if strings.EqualFold(role, "user") {
		accountID := req.AccountId

		if err := db.QueryRow(`SELECT COUNT(*) FROM materials WHERE uploader = ?`, accountID).Scan(&total); err != nil {
			http.Error(w, "Total query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := db.QueryRow(`SELECT COUNT(*) FROM materials WHERE status = 'pending' AND uploader = ?`, accountID).Scan(&pending); err != nil {
			http.Error(w, "Pending query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := db.QueryRow(`SELECT COUNT(*) FROM materials WHERE status = 'approved' AND uploader = ?`, accountID).Scan(&approved); err != nil {
			http.Error(w, "Approved query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := db.QueryRow(`SELECT COUNT(*) FROM materials WHERE status = 'rejected' AND uploader = ?`, accountID).Scan(&rejected); err != nil {
			http.Error(w, "Rejected query error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := db.Query(`SELECT category, COUNT(*) FROM materials WHERE uploader = ? GROUP BY category`, accountID)
		if err != nil {
			http.Error(w, "Category query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var cat string
			var count int
			if err := rows.Scan(&cat, &count); err == nil {
				byCategory[cat] = count
			}
		}
	} else {
		if err := db.QueryRow(`SELECT COUNT(*) FROM materials`).Scan(&total); err != nil {
			http.Error(w, "Total query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := db.QueryRow(`SELECT COUNT(*) FROM materials WHERE status = 'pending'`).Scan(&pending); err != nil {
			http.Error(w, "Pending query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := db.QueryRow(`SELECT COUNT(*) FROM materials WHERE status = 'approved'`).Scan(&approved); err != nil {
			http.Error(w, "Approved query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := db.QueryRow(`SELECT COUNT(*) FROM materials WHERE status = 'rejected'`).Scan(&rejected); err != nil {
			http.Error(w, "Rejected query error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := db.Query(`SELECT category, COUNT(*) FROM materials GROUP BY category`)
		if err != nil {
			http.Error(w, "Category query error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var cat string
			var count int
			if err := rows.Scan(&cat, &count); err == nil {
				byCategory[cat] = count
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"total":      total,
			"pending":    pending,
			"approved":   approved,
			"rejected":   rejected,
			"byCategory": byCategory,
		},
		"message": "success",
	})
}

func RegisterMaterialListRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/material/list", GetMaterialListHandler)
	mux.HandleFunc("/api/material/pending", GetPendingMaterialsHandler)
	mux.HandleFunc("/api/material/review", ReviewMaterialHandler)
	mux.HandleFunc("/api/material/", DeleteMaterialHandler)
	mux.HandleFunc("/api/material/statistics", GetMaterialStatisticsHandler)
}
