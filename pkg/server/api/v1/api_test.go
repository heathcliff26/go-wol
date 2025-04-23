package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heathcliff26/go-wol/pkg/server/storage"
	"github.com/heathcliff26/go-wol/pkg/server/storage/file"
	"github.com/stretchr/testify/assert"
)

func TestWakeHandler(t *testing.T) {
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
			Name:   "InvalidMAC",
			MAC:    "Not-a-mac-address",
			Status: http.StatusBadRequest,
			Response: Response{
				Status: "error",
				Reason: "Invalid MAC address",
			},
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /api/{macAddr}", WakeHandler)

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

func TestAddHostHandler(t *testing.T) {
	tMatrix := []struct {
		Name, MAC, Host string
		Status          int
		Response        Response
	}{
		{
			Name:   "Success",
			MAC:    "00:11:22:33:44:55",
			Host:   "TestHost",
			Status: http.StatusOK,
			Response: Response{
				Status: "ok",
			},
		},
		{
			Name:   "InvalidMAC",
			MAC:    "Invalid-MAC",
			Host:   "TestHost",
			Status: http.StatusBadRequest,
			Response: Response{
				Status: "error",
				Reason: "Invalid MAC address",
			},
		},
	}

	tmpDir := t.TempDir()

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			cfg := storage.StorageConfig{
				Type: "file",
				File: file.FileBackendConfig{
					Path: tmpDir + "/" + tCase.Name + "-hosts.yaml",
				},
			}
			fileBackend, err := storage.NewStorage(cfg)
			if !assert.NoError(err, "Should create file backend without error") {
				t.FailNow()
			}

			handler := &apiHandler{storage: fileBackend}
			mux := http.NewServeMux()
			mux.HandleFunc("PUT /hosts/{macAddr}/{name}", handler.AddHostHandler)

			req := httptest.NewRequest(http.MethodPut, "/hosts/"+tCase.MAC+"/"+tCase.Host, nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			assert.Equal(tCase.Status, rr.Result().StatusCode, "Should return correct status code")

			var res Response
			err = json.Unmarshal(rr.Body.Bytes(), &res)
			assert.NoError(err, "Response should be json")

			assert.Equal(tCase.Response, res, "Response should match")
		})
	}
}

func TestRemoveHostHandler(t *testing.T) {
	tMatrix := []struct {
		Name, MAC string
		Status    int
		Response  Response
	}{
		{
			Name:   "Success",
			MAC:    "00:11:22:33:44:55",
			Status: http.StatusOK,
			Response: Response{
				Status: "ok",
			},
		},
		{
			Name:   "InvalidMAC",
			MAC:    "Invalid-MAC",
			Status: http.StatusBadRequest,
			Response: Response{
				Status: "error",
				Reason: "Invalid MAC address",
			},
		},
	}

	tmpDir := t.TempDir()

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			cfg := storage.StorageConfig{
				Type: "file",
				File: file.FileBackendConfig{
					Path: tmpDir + "/" + tCase.Name + "-hosts.yaml",
				},
			}
			fileBackend, err := storage.NewStorage(cfg)
			if !assert.NoError(err, "Should create file backend without error") {
				t.FailNow()
			}

			handler := &apiHandler{storage: fileBackend}
			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /hosts/{macAddr}", handler.RemoveHostHandler)

			req := httptest.NewRequest(http.MethodDelete, "/hosts/"+tCase.MAC, nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			assert.Equal(tCase.Status, rr.Result().StatusCode, "Should return correct status code")

			var res Response
			err = json.Unmarshal(rr.Body.Bytes(), &res)
			assert.NoError(err, "Response should be json")

			assert.Equal(tCase.Response, res, "Response should match")
		})
	}
}
