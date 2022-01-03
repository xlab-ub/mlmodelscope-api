package endpoints

import (
	"api/api_mq"
	"fmt"
	"github.com/gin-gonic/gin"
)

type predictRequestBody struct {
	Architecture          string   `json:"architecture,omitempty" binding:"required"`
	Framework             string   `json:"framework,omitempty" binding:"required"`
	Inputs                []string `json:"inputs,omitempty" binding:"required"`
	Model                 string   `json:"model,omitempty" binding:"required"`
	DesiredResultModality string   `json:"desiredResultModality,omitempty" binding:"required"`
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

	channelName := fmt.Sprintf("agent-%s-%s", requestBody.Framework, requestBody.Architecture)
	messageQueue := api_mq.GetMessageQueue()
	channel, _ := messageQueue.GetPublishChannel(channelName)

	_, _ = channel.SendMessage("do some work")

	c.JSON(200, &predictResponseBody{ExperimentId: "1"})
}
