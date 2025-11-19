package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	bonusTypeAcademic          = "academic"
	bonusTypeComprehensive     = "comprehensive"
	bonusCategoryAcademic      = "学术专长"
	bonusCategoryComprehensive = "综合素质"
	academicScoreCap           = 15.0
	comprehensiveScoreCap      = 5.0
)

type BonusRecord struct {
	ID           string  `json:"id"`
	Project      string  `json:"project"`
	AwardDate    string  `json:"awardDate"`
	AwardLevel   string  `json:"awardLevel"`
	AwardType    string  `json:"awardType"`
	TeamRank     string  `json:"teamRank"`
	SelfScore    float64 `json:"selfScore"`
	ScoreBasis   string  `json:"scoreBasis"`
	CollegeScore float64 `json:"collegeScore"`
}

type BonusSummaryItem struct {
	Category   string  `json:"category"`
	TotalScore float64 `json:"totalScore"`
	ItemCount  int     `json:"itemCount"`
}

type BonusSummaryResponse struct {
	TotalScore float64            `json:"totalScore"`
	Items      []BonusSummaryItem `json:"items"`
}

func RegisterBonusRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/bonus/academic/list", bonusListHandler(bonusTypeAcademic))
	mux.HandleFunc("/api/bonus/comprehensive/list", bonusListHandler(bonusTypeComprehensive))
	mux.HandleFunc("/api/bonus/summary", bonusSummaryHandler)
}

func bonusListHandler(targetType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		accountID := strings.TrimSpace(r.URL.Query().Get("accountId"))
		if accountID == "" {
			http.Error(w, "accountId is required", http.StatusBadRequest)
			return
		}

		records, err := queryBonusRecords(accountID)
		if err != nil {
			http.Error(w, fmt.Sprintf("query records failed: %v", err), http.StatusInternalServerError)
			return
		}

		normalizedTarget := normalizeBonusType(targetType)
		filtered := make([]BonusRecord, 0)
		for _, rec := range records {
			if normalizeBonusType(rec.Category) == normalizedTarget {
				filtered = append(filtered, convertToBonusRecord(rec))
			}
		}

		writeJSON(w, map[string]interface{}{
			"code":   0,
			"msg":    "",
			"status": "ok",
			"data":   filtered,
		})
	}
}

func bonusSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	accountID := strings.TrimSpace(r.URL.Query().Get("accountId"))
	if accountID == "" {
		http.Error(w, "accountId is required", http.StatusBadRequest)
		return
	}

	records, err := queryBonusRecords(accountID)
	if err != nil {
		http.Error(w, fmt.Sprintf("query records failed: %v", err), http.StatusInternalServerError)
		return
	}

	academicScore := 0.0
	comprehensiveScore := 0.0
	academicCount := 0
	comprehensiveCount := 0

	for _, rec := range records {
		normalized := normalizeBonusType(rec.Type)
		switch normalized {
		case bonusTypeAcademic:
			academicScore += rec.CollegeScore
			academicCount++
		case bonusTypeComprehensive:
			comprehensiveScore += rec.CollegeScore
			comprehensiveCount++
		}
	}

	cappedAcademic := math.Min(academicScore, academicScoreCap)
	cappedComprehensive := math.Min(comprehensiveScore, comprehensiveScoreCap)

	summary := BonusSummaryResponse{
		TotalScore: cappedAcademic + cappedComprehensive,
		Items: []BonusSummaryItem{
			{
				Category:   bonusCategoryAcademic,
				TotalScore: cappedAcademic,
				ItemCount:  academicCount,
			},
			{
				Category:   bonusCategoryComprehensive,
				TotalScore: cappedComprehensive,
				ItemCount:  comprehensiveCount,
			},
		},
	}

	writeJSON(w, map[string]interface{}{
		"code":   0,
		"msg":    "",
		"status": "ok",
		"data":   summary,
	})
}

func queryBonusRecords(accountID string) ([]MaterialRecord, error) {
	db, err := sql.Open("sqlite3", "./user_info.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err := ensureMaterialRecordsTable(db); err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT materialId, accountId, IFNULL(type, ''), IFNULL(category, ''), IFNULL(id, ''), IFNULL(project, ''), IFNULL(awardDate, ''), IFNULL(awardType, ''), IFNULL(teamRank, ''), IFNULL(selfScore, 0), IFNULL(scoreBasis, ''), IFNULL(collegeScore, 0) FROM material_records WHERE accountId = ?`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]MaterialRecord, 0)
	for rows.Next() {
		var rec MaterialRecord
		if err := rows.Scan(
			&rec.MaterialId,
			&rec.AccountId,
			&rec.Type,
			&rec.Category,
			&rec.Id,
			&rec.Project,
			&rec.AwardDate,
			&rec.AwardType,
			&rec.TeamRank,
			&rec.SelfScore,
			&rec.ScoreBasis,
			&rec.CollegeScore,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func ensureMaterialRecordsTable(db *sql.DB) error {
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
	_, err := db.Exec(createTableSQL)
	return err
}

func convertToBonusRecord(rec MaterialRecord) BonusRecord {
	return BonusRecord{
		ID:           rec.Id,
		Project:      rec.Project,
		AwardDate:    rec.AwardDate,
		AwardLevel:   rec.Category,
		AwardType:    rec.AwardType,
		TeamRank:     rec.TeamRank,
		SelfScore:    rec.SelfScore,
		ScoreBasis:   rec.ScoreBasis,
		CollegeScore: rec.CollegeScore,
	}
}

func normalizeBonusType(raw string) string {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	switch normalized {
	case "academic", "academics", "学术专长", "学术", "speciality", "specialty":
		return bonusTypeAcademic
	case "comprehensive", "comprehensiveness", "综合", "综合素质", "performance":
		return bonusTypeComprehensive
	default:
		return normalized
	}
}

func writeJSON(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}
