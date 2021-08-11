package mocks

import (
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
	"github.com/stretchr/testify/mock"
)

type MockAmqpService struct {
	mock.Mock
}

func (m *MockAmqpService) Send(queueConfig model.QueueConfig, publishConfig model.PublishConfig) error {
	ret := m.Called(queueConfig, publishConfig)

	var r1 error

	if ret.Get(0) != nil {
		r1 = ret.Get(1).(error)
	}

	return r1
}

func (m *MockAmqpService) Consume(queueConfig model.QueueConfig, consumeConfig model.ConsumeConfig) error {
	ret := m.Called(queueConfig, consumeConfig)

	var r1 error

	if ret.Get(0) != nil {
		r1 = ret.Get(1).(error)
	}

	return r1
}