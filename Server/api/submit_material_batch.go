package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// BatchItem 表示批量材料项
type BatchItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Reviewer string `json:"reviewer,omitempty"`
}

// SubmitMaterialBatchHandler 处理批量材料列表获取
func SubmitMaterialBatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "./hci.db")
	if err != nil {
		http.Error(w, "数据库连接失败", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 创建表（如果不存在）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS submit_material_batch (
		id TEXT PRIMARY KEY,
		name TEXT,
		status TEXT,
		reviewer TEXT
	)`)
	if err != nil {
		http.Error(w, "创建表失败", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, name, status, reviewer FROM submit_material_batch")
	if err != nil {
		http.Error(w, "查询数据失败", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []BatchItem
	for rows.Next() {
		var item BatchItem
		err := rows.Scan(&item.ID, &item.Name, &item.Status, &item.Reviewer)
		if err != nil {
			http.Error(w, "数据解析错误", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": items,
	})
}

// RegisterSubmitMaterialBatchRoutes 注册路由
func RegisterSubmitMaterialBatchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/submit-material/batch-list", SubmitMaterialBatchHandler)
}
