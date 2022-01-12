package endpoints

import (
	"api/api_db"
	"api/db"
	"api/db/models"
	"bytes"
	"encoding/json"
	"github.com/c3sr/mq/interfaces"
	"net/http"
	"os"
)

var testDb *db.Db

func openDatabase() {
	os.Setenv("DB_DRIVER", "sqlite")
	os.Setenv("DB_HOST", "models_test.sqlite")
	testDb, _ = api_db.GetDatabase()
	testDb.Migrate()
}

func createTestModelAndFramework() {
	testDb.CreateModel(&models.Model{
		Attributes:  models.ModelAttributes{},
		Description: "Test Model",
		Details:     models.ModelDetails{},
		Framework: &models.Framework{
			Name:    "PyTorch",
			Version: "1.0.0",
		},
		Input:   models.ModelOutput{},
		License: "",
		Name:    "Test_Model",
		Output:  models.ModelOutput{},
		Version: "1.0",
	})
}

func cleanupTestDatabase() {
	api_db.CloseDatabase()
	os.Remove("models_test.sqlite")
}

type nullWriter struct{}

func (w *nullWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func NewJsonRequest(method string, url string, body interface{}) (request *http.Request) {
	jsonBody, _ := json.Marshal(body)
	request, _ = http.NewRequest(method, url, bytes.NewReader(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	return
}

func validPredictRequestBody() (body *predictRequestBody) {
	return &predictRequestBody{
		Architecture:          "amd64",
		BatchSize:             1,
		DesiredResultModality: "image_classification",
		Inputs:                []string{"input_url"},
		Model:                 1,
		TraceLevel:            "NO_TRACE",
	}
}

type messageQueueSpy struct {
	channel        *channelSpy
	correlationId  string
	publishChannel string
}

type channelSpy struct {
	mq      *messageQueueSpy
	message string
}

func (c *channelSpy) SendMessage(message string) (string, error) {
	c.message = message

	return c.mq.correlationId, nil
}

func (c *channelSpy) SendResponse(message string, correlationId string) error {
	panic("implement me")
}

func (m *messageQueueSpy) Acknowledge(message interfaces.Message) error {
	panic("implement me")
}

func (m *messageQueueSpy) Nack(message interfaces.Message) error {
	panic("implement me")
}

func (m *messageQueueSpy) NotifyClose(chan error) {
	panic("implement me")
}

func (m *messageQueueSpy) Shutdown() {
	panic("implement me")
}

func (m *messageQueueSpy) GetPublishChannel(name string) (interfaces.Channel, error) {
	m.channel = &channelSpy{
		mq: m,
	}
	m.publishChannel = name

	return m.channel, nil
}

func (m *messageQueueSpy) SubscribeToChannel(name string) (<-chan interfaces.Message, error) {
	panic("implement me")
}
