package api

import (
	"encoding/json"
	"net/http"
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

	// 模拟返回批量材料列表
	items := []BatchItem{
		{ID: "B001", Name: "材料批次1", Status: "pending"},
		{ID: "B002", Name: "材料批次2", Status: "approved", Reviewer: "admin"},
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
