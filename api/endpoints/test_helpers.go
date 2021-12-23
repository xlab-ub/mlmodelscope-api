package endpoints

import (
	"api/api_db"
	"api/db"
	"bytes"
	"encoding/json"
	"github.com/c3sr/mq/interfaces"
	"net/http"
	"os"
)

var testDb db.Db

func openDatabase() {
	os.Setenv("DB_DRIVER", "sqlite")
	os.Setenv("DB_HOST", "models_test.sqlite")
	testDb, _ = api_db.GetDatabase()
	testDb.Migrate()
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

func validPredictRequestBody(framework string) (body *predictRequestBody) {
	return &predictRequestBody{
		Architecture: "amd64",
		Framework:    framework,
		Inputs:       []string{"input_url"},
		Model:        "AlexNet-v1.0",
	}
}

type messageQueueSpy struct {
	channel        *channelSpy
	publishChannel string
}

type channelSpy struct {
	message string
}

func (c *channelSpy) SendMessage(message string) (string, error) {
	c.message = message

	return "x", nil
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

func (m *messageQueueSpy) Shutdown() {
	panic("implement me")
}

func (m *messageQueueSpy) GetPublishChannel(name string) (interfaces.Channel, error) {
	m.channel = &channelSpy{}
	m.publishChannel = name

	return m.channel, nil
}

func (m *messageQueueSpy) SubscribeToChannel(name string) (<-chan interfaces.Message, error) {
	panic("implement me")
}
