package api

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

// 模拟学生数据
var students = []Student{
	{1, "张三", "2021001", "计算机科学", "计科1班", 5.0},
	{2, "李四", "2021002", "软件工程", "软工1班", 3.5},
	{3, "王五", "2021003", "人工智能", "AI1班", 4.2},
}

// 导出学生信息为 Excel（CSV）
func ExportStudentsExcel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.ms-excel")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"students_%d.csv\"", time.Now().Unix()))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// 写入表头
	writer.Write([]string{"id", "name", "studentId", "major", "class", "score"})

	// 写入内容
	for _, s := range students {
		record := []string{
			strconv.Itoa(s.ID),
			s.Name,
			s.StudentID,
			s.Major,
			s.Class,
			fmt.Sprintf("%.2f", s.Score),
		}
		writer.Write(record)
	}
}

// 导出学生信息为 TXT
func ExportStudentsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"students_%d.txt\"", time.Now().Unix()))

	for _, s := range students {
		line := fmt.Sprintf("ID: %d, 姓名: %s, 学号: %s, 专业: %s, 班级: %s, 加分: %.2f\n",
			s.ID, s.Name, s.StudentID, s.Major, s.Class, s.Score)
		w.Write([]byte(line))
	}
}

// 注册路由函数
func RegisterExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/export/students/excel", ExportStudentsExcel)
	mux.HandleFunc("/api/export/students/txt", ExportStudentsTxt)
}
