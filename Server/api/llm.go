package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/genai"
)

var formPrompt string
var guidelines string
var calculatePrompt string
var analyzePrompt string
var BASE_URL string
var API_KEY string

func readFileElsePanic(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func init() {
	formPrompt = readFileElsePanic("./docs/识别文件.md")
	BASE_URL = readFileElsePanic("./secret/BASE_URL")
	API_KEY = readFileElsePanic("./secret/API_KEY")
	guidelines = readFileElsePanic("./docs/保研条例.md")
	calculatePrompt = readFileElsePanic("./docs/审核信息.md")
	analyzePrompt = readFileElsePanic("./docs/信息分析.md")
}

type FormResult struct {
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Files []string `json:"files"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	genai.SetDefaultBaseURLs(genai.BaseURLParameters{
		GeminiURL: BASE_URL,
		VertexURL: BASE_URL,
	})
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: API_KEY,
	})
	if err != nil {
		http.Error(w, "Failed to create LLM client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if client.ClientConfig().Backend == genai.BackendVertexAI {
		http.Error(w, "Vertex AI backend is not supported for this operation", http.StatusInternalServerError)
		return
	}

	parts := []*genai.Part{
		{Text: formPrompt},
		{Text: "保研条例如下，请严格按照条例要求进行材料填写，否则不予通过。\n" + guidelines},
	}
	for _, fname := range req.Files {
		filename := "./upload/" + fname
		fileBytes, err := os.ReadFile(filename)
		if err != nil {
			http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		MIMEType := http.DetectContentType(fileBytes)

		fmt.Println("检测文件", filename, "格式为", MIMEType)

		parts = append(parts, &genai.Part{
			InlineData: &genai.Blob{
				Data:     fileBytes,
				MIMEType: MIMEType,
			},
		})
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		http.Error(w, "LLM error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	ret := result.Text()
	ret = strings.ReplaceAll(ret, "```json", "")
	ret = strings.ReplaceAll(ret, "```", "")
	ret = strings.ReplaceAll(ret, " ", "")
	ret = strings.ReplaceAll(ret, "\n", "")

	var retJSON FormResult
	if err := json.Unmarshal([]byte(ret), &retJSON); err != nil {
		http.Error(w, "Failed to parse LLM response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    retJSON,
	})
}

func RegisterLLMRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/llm/form", FormHandler)
}

type LLMCalculateResult struct {
	AiScore       float64 `json:"aiScore"`
	AiConfidence  float64 `json:"aiConfidence"`
	AiSuggestions string  `json:"aiSuggestions"`
	AiRiskLevel   string  `json:"aiRiskLevel"`
}

func CalculateScore(res *MaterialUploadRequest) (*LLMCalculateResult, error) {
	ctx := context.Background()
	genai.SetDefaultBaseURLs(genai.BaseURLParameters{
		GeminiURL: BASE_URL,
		VertexURL: BASE_URL,
	})
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: API_KEY,
	})
	if err != nil {
		return nil, err
	}

	if client.ClientConfig().Backend == genai.BackendVertexAI {
		return nil, fmt.Errorf("Vertex AI backend is not supported for this operation")
	}

	parts := []*genai.Part{
		{Text: calculatePrompt},
		{Text: fmt.Sprintf("材料标题：%s\n材料类别：%s\n材料标签：%v\n材料描述：%s\n", res.Title, res.Category, res.Tags, res.Description)},
		{Text: "保研条例如下，请严格按照条例要求进行材料填写，否则不予通过。\n" + guidelines},
	}
	for _, fname := range res.Files {
		filename := "./upload/" + fname
		fileBytes, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		MIMEType := http.DetectContentType(fileBytes)

		fmt.Println("检测文件", filename, "格式为", MIMEType)

		parts = append(parts, &genai.Part{
			InlineData: &genai.Blob{
				Data:     fileBytes,
				MIMEType: MIMEType,
			},
		})
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		return nil, err
	}
	ret := result.Text()
	ret = strings.ReplaceAll(ret, "```json", "")
	ret = strings.ReplaceAll(ret, "```", "")
	ret = strings.ReplaceAll(ret, " ", "")
	ret = strings.ReplaceAll(ret, "\n", "")

	var score LLMCalculateResult

	if err := json.Unmarshal([]byte(ret), &score); err != nil {
		return nil, err
	}

	return &score, nil
}

type MaterialRecord struct {
	MaterialId   string  `json:"materialId"`
	AccountId    string  `json:"accountId"`
	Type         string  `json:"type"`
	Category     string  `json:"category"`
	Id           string  `json:"id"`
	Project      string  `json:"project"`
	AwardDate    string  `json:"awardDate"`
	AwardType    string  `json:"awardType"`
	TeamRank     string  `json:"teamRank"`
	SelfScore    float64 `json:"selfScore"`
	ScoreBasis   string  `json:"scoreBasis"`
	CollegeScore float64 `json:"collegeScore"`
}

type MaterialFile struct {
	FileURL  string `json:"fileUrl"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
}

type MaterialDetail struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Category      string         `json:"category"`
	Tags          []string       `json:"tags"`
	Files         []MaterialFile `json:"files"`
	Status        string         `json:"status"`
	Uploader      string         `json:"uploader"`
	UploadTime    string         `json:"uploadTime"`
	Reviewer      string         `json:"reviewer"`
	ReviewTime    string         `json:"reviewTime"`
	ReviewComment string         `json:"reviewComment"`
	AiScore       float64        `json:"aiScore"`
	AiConfidence  float64        `json:"aiConfidence"`
	AiSuggestions []string       `json:"aiSuggestions"`
	AiRiskLevel   string         `json:"aiRiskLevel"`
}

func materialRecordExists(materialId string) (bool, error) {
	materialId = strings.TrimSpace(materialId)
	if materialId == "" {
		return false, errors.New("materialId is required")
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS material_records (
	materialId TEXT PRIMARY KEY,
	accountId TEXT,
	type TEXT,
	category TEXT,
	id TEXT,
	project TEXT,
	awardDate TEXT,
	awardType TEXT,
	teamRank TEXT,
	selfScore REAL,
	scoreBasis TEXT,
	collegeScore REAL
)`); err != nil {
		return false, err
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(1) FROM material_records WHERE materialId = ?`, materialId).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func saveMaterialRecord(record *MaterialRecord) error {
	if record == nil {
		return errors.New("record is nil")
	}
	if strings.TrimSpace(record.MaterialId) == "" {
		return errors.New("record.MaterialId is required")
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		return err
	}
	defer db.Close()

	const createTableSQL = `CREATE TABLE IF NOT EXISTS material_records (
	materialId TEXT PRIMARY KEY,
	accountId TEXT,
	type TEXT,
	category TEXT,
	id TEXT,
	project TEXT,
	awardDate TEXT,
	awardType TEXT,
	teamRank TEXT,
	selfScore REAL,
	scoreBasis TEXT,
	collegeScore REAL
)`

	if _, err := db.Exec(createTableSQL); err != nil {
		return err
	}

	_, err = db.Exec(`INSERT OR REPLACE INTO material_records (
	materialId, accountId, type, category, id, project, awardDate, awardType, teamRank, selfScore, scoreBasis, collegeScore
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,

		record.MaterialId,
		record.AccountId,
		record.Type,
		record.Category,
		record.Id,
		record.Project,
		record.AwardDate,
		record.AwardType,
		record.TeamRank,
		record.SelfScore,
		record.ScoreBasis,
		record.CollegeScore,
	)
	return err
}

func getMaterialDetail(materialID string) (*MaterialDetail, error) {
	materialID = strings.TrimSpace(materialID)
	if materialID == "" {
		return nil, errors.New("materialID is required")
	}

	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	const createTableSQL = `CREATE TABLE IF NOT EXISTS materials (
		id TEXT PRIMARY KEY,
		title TEXT,
		description TEXT,
		category TEXT,
		tags TEXT,
		files TEXT,
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
	)`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	var (
		id, title, description, category, tagsStr, filesJSON              string
		status, uploader, uploadTime, reviewer, reviewTime, reviewComment string
		aiSuggestionsStr, aiRiskLevel                                     string
		aiScore, aiConfidence                                             float64
	)

	row := db.QueryRow(`SELECT id, title, description, category, tags, files, status, uploader, uploadTime, reviewer, reviewTime, reviewComment, aiScore, aiConfidence, aiSuggestions, aiRiskLevel FROM materials WHERE id = ?`, materialID)
	if err := row.Scan(
		&id, &title, &description, &category, &tagsStr, &filesJSON,
		&status, &uploader, &uploadTime, &reviewer, &reviewTime, &reviewComment,
		&aiScore, &aiConfidence, &aiSuggestionsStr, &aiRiskLevel,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("material %s not found", materialID)
		}
		return nil, err
	}

	tags := make([]string, 0)
	for _, tag := range strings.Split(tagsStr, ",") {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}

	aiSuggestions := make([]string, 0)
	for _, suggestion := range strings.Split(aiSuggestionsStr, ",") {
		trimmed := strings.TrimSpace(suggestion)
		if trimmed != "" {
			aiSuggestions = append(aiSuggestions, trimmed)
		}
	}

	files := make([]MaterialFile, 0)
	if strings.TrimSpace(filesJSON) != "" {
		if err := json.Unmarshal([]byte(filesJSON), &files); err != nil {
			return nil, fmt.Errorf("parse files json failed: %w", err)
		}
	}

	return &MaterialDetail{
		ID:            id,
		Title:         title,
		Description:   description,
		Category:      category,
		Tags:          tags,
		Files:         files,
		Status:        status,
		Uploader:      uploader,
		UploadTime:    uploadTime,
		Reviewer:      reviewer,
		ReviewTime:    reviewTime,
		ReviewComment: reviewComment,
		AiScore:       aiScore,
		AiConfidence:  aiConfidence,
		AiSuggestions: aiSuggestions,
		AiRiskLevel:   aiRiskLevel,
	}, nil
}

func DealMaterialToRecord(materialId string) error {
	ok, err := materialRecordExists(materialId)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	ctx := context.Background()

	detail, err := getMaterialDetail(materialId)
	if err != nil {
		return err
	}

	genai.SetDefaultBaseURLs(genai.BaseURLParameters{
		GeminiURL: BASE_URL,
		VertexURL: BASE_URL,
	})
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: API_KEY,
	})
	if err != nil {
		return err
	}

	if client.ClientConfig().Backend == genai.BackendVertexAI {
		return fmt.Errorf("Vertex AI backend is not supported for this operation")
	}

	parts := []*genai.Part{
		{Text: analyzePrompt},
		{Text: fmt.Sprintf("材料标题：%s\n材料类别：%s\n材料标签：%v\n材料描述：%s\n审核意见：%s\n", detail.Title, detail.Category, detail.Tags, detail.Description, detail.ReviewComment)},
		{Text: "保研条例如下，请严格按照条例要求进行材料填写，否则不予通过。\n" + guidelines},
	}
	for _, fname := range detail.Files {
		filename := "./upload/" + fname.FileName
		fileBytes, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		MIMEType := http.DetectContentType(fileBytes)

		fmt.Println("检测文件", filename, "格式为", MIMEType)

		parts = append(parts, &genai.Part{
			InlineData: &genai.Blob{
				Data:     fileBytes,
				MIMEType: MIMEType,
			},
		})
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		return err
	}
	ret := result.Text()
	ret = strings.ReplaceAll(ret, "```json", "")
	ret = strings.ReplaceAll(ret, "```", "")
	ret = strings.ReplaceAll(ret, " ", "")
	ret = strings.ReplaceAll(ret, "\n", "")

	var score MaterialRecord

	if err := json.Unmarshal([]byte(ret), &score); err != nil {
		return err
	}

	err = saveMaterialRecord(&score)
	if err != nil {
		return err
	}

	return nil
}
