package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestVersionCompare_Compare(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	type Payload struct {
		Data struct {
			CompareFrom string `json:"compare_from,omitempty"`
			CompareTo   string `json:"compare_to,omitempty"`
		} `json:"data"`
	}

	type Response struct {
		Code string `json:"code"`
		Data struct {
			Result string `json:"result,omitempty"`
		} `json:"data"`
	}

	t.Run("status 200", func(t *testing.T) {
		tests := []struct {
			compareFrom      string
			compareTo        string
			resultStatusCode int
			resultAppCode    string
			result           string
		}{
			{
				"1.0.0",
				"1.0.0",
				http.StatusOK,
				"app.compare_success",
				"1.0.0 is \"equal\" to 1.0.0",
			},
			{
				"1.1.0",
				"1.0.0",
				http.StatusOK,
				"app.compare_success",
				"1.1.0 is \"after\" 1.0.0",
			},
			{
				"0.1.0",
				"1.0.0",
				http.StatusOK,
				"app.compare_success",
				"0.1.0 is \"before\" 1.0.0",
			},
		}

		for i := range tests {
			test := tests[i]

			payload := Payload{}
			payload.Data.CompareFrom = test.compareFrom
			payload.Data.CompareTo = test.compareTo

			b, err := json.Marshal(payload)
			require.NoError(t, err)

			r, err := http.NewRequest("POST", "/foo", bytes.NewBuffer(b))
			require.NoError(t, err)

			base := Base{}
			base.Logger = logger
			h := &VersionCompare{}
			base.H = h.Compare

			recorder := httptest.NewRecorder()
			base.ServeHTTP(recorder, r)
			require.Equal(t, test.resultStatusCode, recorder.Code)
			resp := Response{}
			err = json.Unmarshal(recorder.Body.Bytes(), &resp)
			require.Equal(t, test.resultAppCode, resp.Code)
			require.Equal(t, test.result, resp.Data.Result)
		}
	})

	// t.Run("status 422", func(t *testing.T) {
	//
	// })
}
