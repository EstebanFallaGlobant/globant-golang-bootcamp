package structs

import (
	"github.com/stretchr/testify/mock"
)

type MockWordCounter struct {
	mock.Mock
}

func (mock *MockWordCounter) CountWords(text string) []WordCount {

	args := mock.Called(text)

	return args.Get(0).([]WordCount)
}
