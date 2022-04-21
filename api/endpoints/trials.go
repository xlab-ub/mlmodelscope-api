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
	ID               string                       `json:"id"`
	Probability      float64                      `json:"probability"`
	Type             string                       `json:"type"`
	BoundingBox      *featureBoundingBox          `json:"bounding_box,omitempty"`
	Classification   *featureClassification       `json:"classification,omitempty"`
	ImageEnhancement *featureImageEnhancement     `json:"image,omitempty"`
	InstanceSegment  *featureInstanceSegmentation `json:"instance_segment,omitempty"`
	SemanticSegment  *featureSemanticSegment      `json:"semantic_segment,omitempty"`
}

type featureBoundingBox struct {
	Index int64   `json:"index"`
	Label string  `json:"label"`
	XMax  float64 `json:"xmax"`
	XMin  float64 `json:"xmin"`
	YMax  float64 `json:"ymax"`
	YMin  float64 `json:"ymin"`
}

type featureClassification struct {
	Index uint   `json:"index"`
	Label string `json:"label"`
}

type featureImageEnhancement struct {
	Data []byte `json:"data"`
}

type featureInstanceSegmentation struct {
	FloatMask []float32 `json:"float_mask"`
	Height    int32     `json:"height,omitempty"`
	Index     int32     `json:"index,omitempty"`
	IntMask   []int32   `json:"int_mask"`
	Label     string    `json:"label,omitempty"`
	MaskType  string    `json:"mask_type,omitempty"`
	Width     int32     `json:"width,omitempty"`
	Xmax      float32   `json:"xmax,omitempty"`
	Xmin      float32   `json:"xmin,omitempty"`
	Ymax      float32   `json:"ymax,omitempty"`
	Ymin      float32   `json:"ymmin,omitempty"`
}

type featureSemanticSegment struct {
	Height uint     `json:"height"`
	Labels []string `json:"labels"`
	Mask   []uint   `json:"int_mask"`
	Width  uint     `json:"width"`
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
