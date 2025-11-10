// Package api 实现信息导入模块的接口
package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
)

// UploadDir 临时文件保存目录
var UploadDir = filepath.Join(os.TempDir(), "hci_info_import")

// 导入 Excel 文件
func ImportExcelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	if ext != ".xls" && ext != ".xlsx" {
		resp := map[string]interface{}{
			"code":    400,
			"message": "文件格式错误，仅支持 .xls/.xlsx",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	os.MkdirAll(UploadDir, 0755)
	dstPath := filepath.Join(UploadDir, fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename))
	dst, _ := os.Create(dstPath)
	defer dst.Close()
	io.Copy(dst, file)

	// 解析 Excel 并写入数据库
	importedRows := 0
	var errors []string
	func() {
		// 打开 Excel 文件
		xlFile, err := excelize.OpenFile(dstPath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Excel 解析错误: %v", err))
			return
		}
		rows, err := xlFile.GetRows("Sheet1")
		if err != nil {
			errors = append(errors, fmt.Sprintf("读取表格错误: %v", err))
			return
		}

		// 连接数据库
		db, err := sql.Open("sqlite3", "./app.db")
		if err != nil {
			errors = append(errors, fmt.Sprintf("数据库连接失败: %v", err))
			return
		}
		defer db.Close()

		// 确保 students 表存在
		createTableSQL := `CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			studentId TEXT,
			major TEXT,
			class TEXT,
			score REAL
		);`
		if _, err := db.Exec(createTableSQL); err != nil {
			errors = append(errors, fmt.Sprintf("创建表失败: %v", err))
			return
		}

		// 遍历 Excel 行
		for i, row := range rows {
			if i == 0 {
				continue // 跳过表头
			}
			if len(row) < 6 {
				errors = append(errors, fmt.Sprintf("第 %d 行数据不足6列", i+1))
				continue
			}
			scoreVal, _ := strconv.ParseFloat(row[5], 64)
			_, err := db.Exec(`INSERT INTO students (name, studentId, major, class, score) VALUES (?, ?, ?, ?, ?)`,
				row[1], row[2], row[3], row[4], scoreVal)
			if err != nil {
				errors = append(errors, fmt.Sprintf("插入第 %d 行失败: %v", i+1, err))
				continue
			}
			importedRows++
		}
	}()

	resp := map[string]interface{}{
		"code":    200,
		"message": "Excel 文件已成功导入",
		"data": map[string]interface{}{
			"rows":   importedRows,
			"errors": errors,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// 导入 TXT 文件或直接文本
func ImportTxtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	var lines int
	var errors []string

	// splitLines 按行分割字符串
	splitLines := func(s string) []string {
		var result []string
		current := ""
		for _, r := range s {
			if r == '\n' || r == '\r' {
				if current != "" {
					result = append(result, current)
					current = ""
				}
			} else {
				current += string(r)
			}
		}
		if current != "" {
			result = append(result, current)
		}
		return result
	}

	// parseFields 按逗号或制表符分割字段
	parseFields := func(line string) []string {
		if len(line) == 0 {
			return []string{}
		}
		// 优先按逗号分割
		if fields := filepath.SplitList(line); len(fields) > 1 {
			return fields
		}
		// 按制表符分割
		return filepath.SplitList(line)
	}

	// 连接数据库
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("数据库连接失败: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 确保 students 表存在
	createTableSQL := `CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		studentId TEXT,
		major TEXT,
		class TEXT,
		score REAL
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		http.Error(w, fmt.Sprintf("创建表失败: %v", err), http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		ext := filepath.Ext(handler.Filename)
		if ext != ".txt" {
			resp := map[string]interface{}{
				"code":    400,
				"message": "TXT 文件格式错误或内容为空",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		os.MkdirAll(UploadDir, 0755)
		dstPath := filepath.Join(UploadDir, fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename))
		dst, _ := os.Create(dstPath)
		defer dst.Close()
		io.Copy(dst, file)

		// 读取 TXT 文件并解析
		content, err := os.ReadFile(dstPath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("读取文件失败: %v", err))
		} else {
			linesArr := splitLines(string(content))
			for i, line := range linesArr {
				fields := parseFields(line)
				if len(fields) < 6 {
					errors = append(errors, fmt.Sprintf("第 %d 行数据不足6列", i+1))
					continue
				}
				scoreVal, _ := strconv.ParseFloat(fields[5], 64)
				_, err := db.Exec(`INSERT INTO students (name, studentId, major, class, score) VALUES (?, ?, ?, ?, ?)`,
					fields[1], fields[2], fields[3], fields[4], scoreVal)
				if err != nil {
					errors = append(errors, fmt.Sprintf("插入第 %d 行失败: %v", i+1, err))
					continue
				}
				lines++
			}
		}
	} else {
		// 如果 file 为空，尝试读取 text 字段内容
		text := r.FormValue("text")
		if len(text) == 0 {
			resp := map[string]interface{}{
				"code":    400,
				"message": "TXT 文件格式错误或内容为空",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		linesArr := splitLines(text)
		for i, line := range linesArr {
			fields := parseFields(line)
			if len(fields) < 6 {
				errors = append(errors, fmt.Sprintf("第 %d 行数据不足6列", i+1))
				continue
			}
			scoreVal, _ := strconv.ParseFloat(fields[5], 64)
			_, err := db.Exec(`INSERT INTO students (name, studentId, major, class, score) VALUES (?, ?, ?, ?, ?)`,
				fields[1], fields[2], fields[3], fields[4], scoreVal)
			if err != nil {
				errors = append(errors, fmt.Sprintf("插入第 %d 行失败: %v", i+1, err))
				continue
			}
			lines++
		}
	}

	resp := map[string]interface{}{
		"code":    200,
		"message": "导入成功",
		"data": map[string]interface{}{
			"lines":  lines,
			"errors": errors,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RegisterInfoImportRoutes 注册信息导入相关路由
func RegisterInfoImportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/info/import/excel", ImportExcelHandler)
	mux.HandleFunc("/api/info/import/txt", ImportTxtHandler)
}
