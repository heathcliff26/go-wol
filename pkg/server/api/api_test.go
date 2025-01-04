package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	tMatrix := []struct {
		Name, MAC string
		Status    int
		Response  Response
	}{
		{
			Name:   "Success",
			MAC:    "ff:ff:ff:ff:ff:ff",
			Status: http.StatusOK,
			Response: Response{
				Status: "ok",
			},
		},
		{
			Name:   "WrongMAC",
			MAC:    "Not-a-mac-address",
			Status: http.StatusBadRequest,
			Response: Response{
				Status: "error",
				Reason: "Failed to create magic packet from MAC address",
			},
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /api/{macAddr}", Api)

			req := httptest.NewRequest(http.MethodGet, "/api/"+tCase.MAC, nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			assert.Equal(tCase.Status, rr.Result().StatusCode, "Should return correct status code")

			var res Response
			err := json.Unmarshal(rr.Body.Bytes(), &res)
			assert.NoError(err, "Response should be json")

			assert.Equal(tCase.Response, res, "Response should match")
		})
	}
}
