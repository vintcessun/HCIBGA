package api

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// 学生结构体
type Student struct {
	ID        int
	Name      string
	StudentID string
	Major     string
	Class     string
	Score     float64
}

// 模拟学生数据（暂时保留，TXT 导出未改造前使用）
var students = []Student{
	{1, "张三", "2021001", "计算机科学", "计科1班", 5.0},
	{2, "李四", "2021002", "软件工程", "软工1班", 3.5},
	{3, "王五", "2021003", "人工智能", "AI1班", 4.2},
}

// 导出学生信息为 Excel（CSV）
func ExportStudentsExcel(w http.ResponseWriter, r *http.Request) {
	// 连接数据库
	db, err := sql.Open("sqlite3", "./students.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("数据库连接失败: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 确保表存在
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		studentId TEXT,
		major TEXT,
		class TEXT,
		score REAL
	)`)
	if err != nil {
		http.Error(w, fmt.Sprintf("初始化数据表失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 查询学生信息
	rows, err := db.Query(`SELECT id, name, studentId, major, class, score FROM students`)
	if err != nil {
		http.Error(w, fmt.Sprintf("查询数据失败: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "application/vnd.ms-excel")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"students_%d.csv\"", time.Now().Unix()))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write([]string{"id", "name", "studentId", "major", "class", "score"}); err != nil {
		http.Error(w, fmt.Sprintf("写入表头失败: %v", err), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.StudentID, &s.Major, &s.Class, &s.Score); err != nil {
			http.Error(w, fmt.Sprintf("读取数据失败: %v", err), http.StatusInternalServerError)
			return
		}
		record := []string{
			strconv.Itoa(s.ID),
			s.Name,
			s.StudentID,
			s.Major,
			s.Class,
			fmt.Sprintf("%.2f", s.Score),
		}
		if err := writer.Write(record); err != nil {
			http.Error(w, fmt.Sprintf("写入数据失败: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// 导出学生信息为 TXT
func ExportStudentsTxt(w http.ResponseWriter, r *http.Request) {
	// 连接数据库
	db, err := sql.Open("sqlite3", "./students.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("数据库连接失败: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 确保表存在
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		studentId TEXT,
		major TEXT,
		class TEXT,
		score REAL
	)`)
	if err != nil {
		http.Error(w, fmt.Sprintf("初始化数据表失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 查询学生信息
	rows, err := db.Query(`SELECT id, name, studentId, major, class, score FROM students`)
	if err != nil {
		http.Error(w, fmt.Sprintf("查询数据失败: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"students_%d.txt\"", time.Now().Unix()))

	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.StudentID, &s.Major, &s.Class, &s.Score); err != nil {
			http.Error(w, fmt.Sprintf("读取数据失败: %v", err), http.StatusInternalServerError)
			return
		}
		line := fmt.Sprintf("ID: %d, 姓名: %s, 学号: %s, 专业: %s, 班级: %s, 成绩: %.2f",
			s.ID, s.Name, s.StudentID, s.Major, s.Class, s.Score)
		if _, err := w.Write([]byte(line + "\n")); err != nil {
			http.Error(w, fmt.Sprintf("写入数据失败: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// 注册路由函数
func RegisterExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/export/students/excel", ExportStudentsExcel)
	mux.HandleFunc("/api/export/students/txt", ExportStudentsTxt)
}
