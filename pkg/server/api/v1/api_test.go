package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/heathcliff26/go-wol/pkg/server/storage"
	"github.com/heathcliff26/go-wol/pkg/server/storage/file"
	"github.com/heathcliff26/go-wol/pkg/server/storage/types"
	"github.com/heathcliff26/go-wol/pkg/server/storage/valkey"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGetHostsHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)

		cfg := storage.StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: "testdata/hosts.yaml",
			},
		}

		storageBackend, err := storage.NewStorage(cfg)
		require.NoError(t, err, "Should create file backend without error")

		handler := &apiHandler{storage: storageBackend}
		mux := http.NewServeMux()
		mux.HandleFunc("GET /hosts", handler.GetHostsHandler)

		req := httptest.NewRequest(http.MethodGet, "/hosts", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		assert.Equal(http.StatusOK, rr.Result().StatusCode, "Should return status code 200")
		var res []types.Host
		err = json.Unmarshal(rr.Body.Bytes(), &res)
		assert.NoError(err, "Response should be json")
		assert.Len(res, 2, "Should return 2 hosts")
	})
	t.Run("StorageError", func(t *testing.T) {
		assert := assert.New(t)

		mr := miniredis.RunT(t)

		cfg := storage.StorageConfig{
			Type: "valkey",
			Valkey: valkey.ValkeyConfig{
				Addrs: []string{mr.Addr()},
			},
		}

		storageBackend, err := storage.NewStorage(cfg)
		require.NoError(t, err, "Should create valkey backend without error")

		// Close miniredis to simulate storage error
		mr.Close()

		handler := &apiHandler{storage: storageBackend}
		mux := http.NewServeMux()
		mux.HandleFunc("GET /hosts", handler.GetHostsHandler)

		req := httptest.NewRequest(http.MethodGet, "/hosts", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Result().StatusCode, "Should return status code 200")
	})
}

func TestAddHostHandler(t *testing.T) {
	tMatrix := []struct {
		Name     string
		Host     types.Host
		Readonly bool
		Status   int
		Response Response
	}{
		{
			Name: "Success",
			Host: types.Host{
				MAC:  "00:11:22:33:44:55",
				Name: "TestHost",
			},
			Status: http.StatusOK,
			Response: Response{
				Status: "ok",
			},
		},
		{
			Name: "InvalidMAC",
			Host: types.Host{
				MAC:  "Invalid-MAC",
				Name: "TestHost",
			},
			Status: http.StatusBadRequest,
			Response: Response{
				Status: "error",
				Reason: "Invalid MAC address",
			},
		},
		{
			Name: "InvalidHost",
			Host: types.Host{
				MAC:  "00:11:22:33:44:55",
				Name: "Invalid-Host@not_a_domain",
			},
			Status: http.StatusBadRequest,
			Response: Response{
				Status: "error",
				Reason: "Invalid hostname",
			},
		},
		{
			Name: "ReadonlyStorage",
			Host: types.Host{
				MAC:  "00:11:22:33:44:55",
				Name: "TestHost",
			},
			Readonly: true,
			Status:   http.StatusForbidden,
			Response: Response{
				Status: "error",
				Reason: "Storage is readonly",
			},
		},
	}

	tmpDir := t.TempDir()

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			cfg := storage.StorageConfig{
				Type: "file",
				File: file.FileBackendConfig{
					Path: tmpDir + "/" + tCase.Name + "-hosts.yaml",
				},
				Readonly: tCase.Readonly,
			}
			storageBackend, err := storage.NewStorage(cfg)
			require.NoError(err, "Should create file backend without error")

			router := NewRouter(storageBackend)

			body, err := json.Marshal(tCase.Host)
			require.NoError(err, "Should encode host to JSON")

			req := httptest.NewRequest(http.MethodPut, "/hosts", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(tCase.Status, rr.Result().StatusCode, "Should return correct status code")

			var res Response
			err = json.Unmarshal(rr.Body.Bytes(), &res)
			assert.NoError(err, "Response should be json")

			assert.Equal(tCase.Response, res, "Response should match")
		})
	}

	t.Run("InvalidRequestBody", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		cfg := storage.StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: tmpDir + "/InvalidRequestBody-hosts.yaml",
			},
		}
		storageBackend, err := storage.NewStorage(cfg)
		require.NoError(err, "Should create file backend without error")

		router := NewRouter(storageBackend)

		req := httptest.NewRequest(http.MethodPut, "/hosts", bytes.NewReader([]byte("This is a text, not a JSON object")))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusBadRequest, rr.Result().StatusCode, "Should return status code 400")

		var res Response
		err = json.Unmarshal(rr.Body.Bytes(), &res)
		assert.NoError(err, "Response should be json")

		assert.Equal(Response{
			Status: "error",
			Reason: "Request body must be a valid host JSON object",
		}, res, "Response should match")
	})
}

func TestRemoveHostHandler(t *testing.T) {
	tMatrix := []struct {
		Name, MAC string
		Readonly  bool
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
		{
			Name:     "ReadonlyStorage",
			MAC:      "00:11:22:33:44:55",
			Readonly: true,
			Status:   http.StatusForbidden,
			Response: Response{
				Status: "error",
				Reason: "Storage is readonly",
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
				Readonly: tCase.Readonly,
			}
			storageBackend, err := storage.NewStorage(cfg)
			require.NoError(t, err, "Should create file backend without error")

			handler := &apiHandler{storage: storageBackend}
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

func TestStorageErrors(t *testing.T) {
	hostsFile := t.TempDir() + "/hosts.yaml"

	cfg := storage.StorageConfig{
		Type: "file",
		File: file.FileBackendConfig{
			Path: hostsFile,
		},
	}

	storageBackend, err := storage.NewStorage(cfg)
	require.NoError(t, err, "Should create file backend without error")

	require.NoError(t, os.Chmod(hostsFile, 0444), "Should set file permissions without error")

	router := NewRouter(storageBackend)

	t.Run("AddHost", func(t *testing.T) {
		assert := assert.New(t)

		body, err := json.Marshal(types.Host{
			MAC:  "FF:FF:FF:FF:FF:FF",
			Name: "testhost",
		})
		require.NoError(t, err, "Should encode host to JSON")

		req := httptest.NewRequest(http.MethodPut, "/hosts", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Result().StatusCode, "Should return correct status code")

		var res Response
		err = json.Unmarshal(rr.Body.Bytes(), &res)
		assert.NoError(err, "Response should be json")

		expectedResponse := Response{
			Status: "error",
			Reason: "Failed to add host",
		}

		assert.Equal(expectedResponse, res, "Response should match")
	})
	t.Run("RemoveHost", func(t *testing.T) {
		assert := assert.New(t)

		req := httptest.NewRequest(http.MethodDelete, "/hosts/FF:FF:FF:FF:FF:FF", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Result().StatusCode, "Should return correct status code")

		var res Response
		err = json.Unmarshal(rr.Body.Bytes(), &res)
		assert.NoError(err, "Response should be json")

		expectedResponse := Response{
			Status: "error",
			Reason: "Failed to remove host",
		}

		assert.Equal(expectedResponse, res, "Response should match")
	})
}
