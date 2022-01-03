package endpoints

import (
	"api/api_db"
	"api/api_mq"
	"fmt"
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

	_, _ = channel.SendMessage("do some work")

	c.JSON(200, &predictResponseBody{ExperimentId: "1"})
}
