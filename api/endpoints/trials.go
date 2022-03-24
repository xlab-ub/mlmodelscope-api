package endpoints

import (
	"api/api_db"
	"api/db/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type trialResponse struct {
	ID          string        `json:"id"`
	Inputs      []string      `json:"inputs"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	Results     trialResults  `json:"results,omitempty"`
	Model       *models.Model `json:"model,omitempty"`
}

type trialResults struct {
	Responses []responseFeatures `json:"responses,omitempty"`
	TraceId   traceId            `json:"trace_id,omitempty"`
}

type responseFeatures struct {
	Features []responseFeature `json:"features,omitempty"`
}

type responseFeature struct {
	ID             string                `json:"id"`
	Probability    float64               `json:"probability"`
	Type           string                `json:"type"`
	Classification featureClassification `json:"classification"`
}

type featureClassification struct {
	Index uint   `json:"index"`
	Label string `json:"label"`
}

type traceId struct {
	Id string `json:"id,omitempty"`
}

func GetTrial(c *gin.Context) {
	if id := c.Param("trial_id"); id == "" {
		c.JSON(404, NotFoundResponse)
		return
	} else {
		db, err := api_db.GetDatabase()
		if err != nil {
			c.JSON(500, NewErrorResponse(err))
			return
		}

		trial, err := db.GetTrialById(id)
		if err != nil {
			log.Printf("[WARN] %s", err.Error())
			c.JSON(404, NotFoundResponse)
			return
		}

		framework, err := db.QueryFrameworks(&models.Framework{ID: uint(trial.Model.FrameworkID)})
		if err != nil {
			log.Printf("[WARN] %s", err.Error())
			c.JSON(500, NewErrorResponse(err))
			return
		}

		trial.Model.Framework = framework
		response := trialToResponse(trial)
		c.JSON(200, response)
	}
}

func trialToResponse(t *models.Trial) (r *trialResponse) {
	var inputs []string
	for _, input := range t.Inputs {
		inputs = append(inputs, input.URL)
	}

	var results trialResults
	err := json.Unmarshal([]byte(t.Result), &results)
	if err != nil {
		log.Printf("[WARN] %s", err.Error())
	}
	r = &trialResponse{
		ID:          t.ID,
		Inputs:      inputs,
		CompletedAt: t.CompletedAt,
		Model:       t.Model,
		Results:     results,
	}

	return
}
