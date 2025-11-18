package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"google.golang.org/genai"
)

var formPrompt string
var guidelines string
var calculatePrompt string
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
