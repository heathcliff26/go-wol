package server

import (
	"io"
	"net/http"
	"testing"

	"github.com/heathcliff26/go-wol/pkg/server/config"
	"github.com/heathcliff26/go-wol/static"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	assert := assert.New(t)

	c := config.Config{
		Port: 8080,
		SSL: config.SSLConfig{
			Enabled: true,
			Cert:    "test.crt",
			Key:     "test.key",
		},
		Hosts: []config.Host{
			{
				Name: "testName",
				MAC:  "testMAC",
			},
		},
	}

	s, err := NewServer(c)

	if !assert.NoError(err, "Should create new server") || !assert.NotNil(s, "Server should not be empty") {
		t.FailNow()
	}

	assert.Equal(":8080", s.addr, "Server should have address set")
	assert.Equal(c.SSL, s.ssl, "Server should have SSL Config")

	assert.Contains(s.indexHTML, "testName", "Should add hostname to index.html")
	assert.Contains(s.indexHTML, "testMAC", "Should add MAC to index.html")
	assert.Contains(s.indexHTML, "onclick=\"wake('testMAC', 'testName');\"", "Should add onClick javascript function call with MAC")
	assert.NotEmpty(s.indexChecksum, "Server should have checksum of index.html")
}

func TestServer(t *testing.T) {
	t.Run("SSL", func(t *testing.T) {
		assert := assert.New(t)

		c := config.Config{
			Port: 8080,
			SSL: config.SSLConfig{
				Enabled: true,
				Cert:    "test.crt",
				Key:     "test.key",
			},
			Hosts: []config.Host{
				{
					Name: "testName",
					MAC:  "testMAC",
				},
			},
		}

		s, err := NewServer(c)
		if !assert.NoError(err, "Should create server without error") {
			t.FailNow()
		}

		assert.Error(s.Run(), "Server should fail to run, as the ssl certificate and key do not exist")
	})
	c := config.Config{
		Port: 8080,
		Hosts: []config.Host{
			{
				Name: "testName",
				MAC:  "testMAC",
			},
		},
	}
	s, err := NewServer(c)
	if !assert.NoError(t, err, "Should create server without error") {
		t.FailNow()
	}

	go func() {
		err := s.Run()
		if err != nil {
			t.Logf("Failed to run server: %v", err)
			t.Fail()
		}
	}()
	address := "http://localhost:8080"

	t.Run("IndexHandler", func(t *testing.T) {
		assert := assert.New(t)

		for _, path := range []string{"/", "/index.html"} {
			res, err := http.Get(address + path)
			t.Cleanup(func() {
				res.Body.Close()
			})

			assert.NoErrorf(err, "Should not return error for request to %s", path)
			assert.Equal(http.StatusOK, res.StatusCode, "Should return 200 for request to %s", path)
			assert.Equalf(s.indexChecksum, res.Header.Get("ETag"), "Should have ETag set on request to %s", path)
			assert.NotEmptyf(res.Header.Get("Cache-Control"), "Should have Cache-Control header set on request to %s", path)

			body, err := io.ReadAll(res.Body)
			assert.NoErrorf(err, "Should read body without error for request to %s", path)
			assert.Equalf(s.indexHTML, string(body), "Body should match index.html on path %s", path)
		}

		resNotFound, err := http.Get(address + "/something")
		t.Cleanup(func() {
			resNotFound.Body.Close()
		})
		assert.NoError(err, "Should receive valid response for random path")
		assert.Equal(http.StatusNotFound, resNotFound.StatusCode, "Should return 404 for random path")

		req, _ := http.NewRequest(http.MethodGet, address+"/", nil)
		req.Header.Add("If-None-Match", s.indexChecksum)

		resCache, err := (&http.Client{}).Do(req)
		t.Cleanup(func() {
			resCache.Body.Close()
		})
		assert.NoError(err, "Should receive valid response to call with cache")
		assert.Equal(http.StatusNotModified, resCache.StatusCode, "Should receive cache hit")
	})
	t.Run("API", func(t *testing.T) {
		assert := assert.New(t)

		res, err := http.Get(address + "/api/not-a-mac")
		t.Cleanup(func() {
			res.Body.Close()
		})

		assert.NoError(err)
		assert.Equal(http.StatusBadRequest, res.StatusCode, "Should receive a bad request response when using a malformed mac address")
	})
	t.Run("CSS", func(t *testing.T) {
		assert := assert.New(t)

		file, err := static.CSS.ReadFile("css/bootstrap.css")
		assert.NoError(err, "Should read file from static")

		res, err := http.Get(address + "/css/bootstrap.css")
		t.Cleanup(func() {
			res.Body.Close()
		})

		assert.NoError(err, "Request should not return an error")
		assert.Equal(http.StatusOK, res.StatusCode, "Request should not return an error")
		body, _ := io.ReadAll(res.Body)
		assert.Equal(string(file), string(body), "Response should match file")
	})
	t.Run("Icons", func(t *testing.T) {
		assert := assert.New(t)

		file, err := static.Icons.ReadFile("icons/favicon.svg")
		assert.NoError(err, "Should read file from static")

		res, err := http.Get(address + "/icons/favicon.svg")
		t.Cleanup(func() {
			res.Body.Close()
		})

		assert.NoError(err, "Request should not return an error")
		assert.Equal(http.StatusOK, res.StatusCode, "Request should not return an error")
		body, _ := io.ReadAll(res.Body)
		assert.Equal(string(file), string(body), "Response should match file")
	})
	t.Run("JS", func(t *testing.T) {
		assert := assert.New(t)

		file, err := static.JS.ReadFile("js/wake.js")
		assert.NoError(err, "Should read file from static")

		res, err := http.Get(address + "/js/wake.js")
		t.Cleanup(func() {
			res.Body.Close()
		})

		assert.NoError(err, "Request should not return an error")
		assert.Equal(http.StatusOK, res.StatusCode, "Request should not return an error")
		body, _ := io.ReadAll(res.Body)
		assert.Equal(string(file), string(body), "Response should match file")
	})
}
