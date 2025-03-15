package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/heathcliff26/go-wol/pkg/server/api"
	"github.com/heathcliff26/go-wol/pkg/server/config"
	"github.com/heathcliff26/go-wol/static"
	"github.com/heathcliff26/simple-fileserver/pkg/middleware"
)

type Server struct {
	addr          string
	ssl           config.SSLConfig
	indexHTML     string
	indexChecksum string
}

func NewServer(c config.Config) (*Server, error) {
	tmpl, err := template.New("index.html").Parse(string(static.IndexTemplate))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, c.Hosts)
	if err != nil {
		return nil, err
	}
	indexHTML := buf.String()

	checksum := sha256.Sum256([]byte(indexHTML))

	return &Server{
		addr:          ":" + strconv.Itoa(c.Port),
		ssl:           c.SSL,
		indexHTML:     indexHTML,
		indexChecksum: hex.EncodeToString(checksum[:]),
	}, nil
}

func (s *Server) indexHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("ETag", s.indexChecksum)
	res.Header().Set("Cache-Control", "public, max-age=300")

	if match := req.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, s.indexChecksum) {
			res.WriteHeader(http.StatusNotModified)
			return
		}
	}

	count, err := res.Write([]byte(s.indexHTML))
	if err != nil {
		slog.Error("Failed to write index.html to client", "err", err, slog.Int("written", count))
	}
}

// Starts the server and exits with error if that fails
func (s *Server) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", s.indexHandler)
	router.HandleFunc("GET /index.html", s.indexHandler)
	router.HandleFunc("GET /api/{macAddr}", api.API)
	router.Handle("GET /css/", StaticFileServer(static.CSS))
	router.Handle("GET /js/", StaticFileServer(static.JS))

	server := http.Server{
		Addr:    s.addr,
		Handler: middleware.Logging(router),
	}

	var err error
	if s.ssl.Enabled {
		slog.Info("Starting server", slog.String("addr", s.addr), slog.String("sslKey", s.ssl.Key), slog.String("sslCert", s.ssl.Cert))
		err = server.ListenAndServeTLS(s.ssl.Cert, s.ssl.Key)
	} else {
		slog.Info("Starting server", slog.String("addr", s.addr))
		err = server.ListenAndServe()
	}

	// This just means the server was closed after running
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("Server closed, exiting")
		return nil
	}
	return err
}
