package log

import "github.com/stretchr/testify/mock"

type LoggerMock struct{ mock.Mock }

func (l *LoggerMock) Error(msg string, err error, args ...interface{}) {
	l.Mock.Called(msg, err, args)
}

func (l *LoggerMock) Info(msg string, args ...interface{}) {
	l.Mock.Called(msg, args)
}

func (l *LoggerMock) Debug(msg string, args ...interface{}) {
	l.Mock.Called(msg, args)
}

func (l *LoggerMock) With(args ...interface{}) Logger {
	res := l.Mock.Called(args)
	return res.Get(0).(Logger)
}
