package api_mq

import "github.com/c3sr/mq/interfaces"

type nullMessageQueue struct{}

func (n *nullMessageQueue) Acknowledge(message interfaces.Message) error {
	return nil
}

func (n *nullMessageQueue) Nack(message interfaces.Message) error {
	return nil
}

func (n *nullMessageQueue) Shutdown() {
}

func (n *nullMessageQueue) GetPublishChannel(name string) (interfaces.Channel, error) {
	return &nullChannel{}, nil
}

func (n *nullMessageQueue) SubscribeToChannel(name string) (<-chan interfaces.Message, error) {
	return nil, nil
}

func NullMessageQueue() interfaces.MessageQueue {
	return &nullMessageQueue{}
}

type nullChannel struct{}

func (n *nullChannel) SendMessage(message string) (string, error) {
	return "", nil
}

func (n *nullChannel) SendResponse(message string, correlationId string) error {
	panic("implement me")
}
