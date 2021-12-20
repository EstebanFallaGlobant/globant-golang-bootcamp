package wordcounterapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	wcStructs "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/wordcounterapi/structs"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
	"github.com/stretchr/testify/assert"
)

var app App

func Test_WordCounterAPI(t *testing.T) {
	testCases := []struct {
		name               string
		input              string
		method             string
		handler            string
		expectedStatus     int
		expectedCollection []wcStructs.WordCount
	}{
		{
			name:               "Test with valid string",
			input:              "This is a test string",
			expectedStatus:     200,
			expectedCollection: getEmptyWordCountCollection(),
			method:             "GET",
			handler:            "count",
		},
		{
			name:               "Test with empty string",
			input:              "   ",
			expectedStatus:     400,
			expectedCollection: getEmptyWordCountCollection(),
			method:             "GET",
			handler:            "count",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(wcStructs.MockWordCounter)
			mock.On("CountWords", tt.input).Return(tt.expectedCollection)

			app.Initialize(mock)

			reqResult := new(wcStructs.WordCounterResponse)

			reqPath := "/"

			if !util.IsEmptyString(tt.handler) {
				reqPath += tt.handler + "/"
			}

			reqPath += tt.input

			t.Logf("Request path: %s", reqPath)

			req, _ := http.NewRequest(tt.method, reqPath, nil)
			res := executeRequest(req)

			json.Unmarshal(res.Body.Bytes(), &reqResult)

			assert.Equal(t, tt.expectedStatus, reqResult.Status)
		})
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func getEmptyWordCountCollection() []wcStructs.WordCount {
	return make([]wcStructs.WordCount, 0)
}
