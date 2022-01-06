package endpoints

import (
	"api/api_db"
	"api/api_mq"
	"api/db/models"
	"encoding/json"
	"fmt"
	"github.com/c3sr/mq/messages"
	"github.com/gin-gonic/gin"
	"strings"
)

type predictRequestBody struct {
	Architecture          string   `json:"architecture,omitempty" binding:"required"`
	BatchSize             uint     `json:"batchSize,omitempty" binding:"required"`
	DesiredResultModality string   `json:"desiredResultModality,omitempty" binding:"required"`
	Inputs                []string `json:"inputs,omitempty" binding:"required"`
	Model                 uint     `json:"model,omitempty" binding:"required"`
	TraceLevel            string   `json:"traceLevel,omitempty" binding:"required"`
}

type predictResponseBody struct {
	ExperimentId string `json:"experimentId,omitempty"`
}

func Predict(c *gin.Context) {
	requestBody := &predictRequestBody{}
	err := c.Bind(requestBody)
	if err != nil {
		return
	}

	db, _ := api_db.GetDatabase()
	model, err := db.GetModelById(requestBody.Model)
	if err != nil {
		c.JSON(400, err)

		return
	}

	channelName := fmt.Sprintf("agent-%s-%s", strings.ToLower(model.Framework.Name), requestBody.Architecture)
	messageQueue := api_mq.GetMessageQueue()
	channel, _ := messageQueue.GetPublishChannel(channelName)
	message := makePredictMessage(requestBody, model)
	messageBytes, _ := json.Marshal(message)
	correlationId, _ := channel.SendMessage(string(messageBytes))

	var inputs []models.TrialInput
	for _, input := range requestBody.Inputs {
		inputs = append(inputs, models.TrialInput{
			URL: input,
		})
	}

	db.CreateTrial(&models.Trial{
		ID:      correlationId,
		ModelID: model.ID,
		Inputs:  inputs,
	})

	c.JSON(200, &predictResponseBody{ExperimentId: correlationId})
}

func makePredictMessage(request *predictRequestBody, model *models.Model) *messages.PredictByModelName {
	return &messages.PredictByModelName{
		BatchSize:             1,
		DesiredResultModality: request.DesiredResultModality,
		Inputs:                request.Inputs,
		ModelName:             fmt.Sprintf("%s_%s", strings.ToLower(model.Name), model.Version),
		Warmups:               1,
		TraceLevel:            request.TraceLevel,
		TracerAddress:         "trace.mlmodelscope.org",
		UseGpu:                false,
	}
}
