package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	api "github.com/heathcliff26/go-wol/pkg/server/api/v1"
	"github.com/heathcliff26/go-wol/pkg/server/config"
	"github.com/heathcliff26/go-wol/pkg/server/storage"
	"github.com/heathcliff26/go-wol/static"
	"github.com/heathcliff26/simple-fileserver/pkg/middleware"
)

type Server struct {
	addr    string
	ssl     config.SSLConfig
	storage *storage.Storage
}

func NewServer(cfgServer config.ServerConfig, cfgStorage storage.StorageConfig) (*Server, error) {
	storage, err := storage.NewStorage(cfgStorage)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	return &Server{
		addr:    ":" + strconv.Itoa(cfgServer.Port),
		ssl:     cfgServer.SSL,
		storage: storage,
	}, nil
}

func (s *Server) indexHandler(res http.ResponseWriter, req *http.Request) {
	indexHTML, indexChecksum := s.storage.GetIndexHTML()
	res.Header().Set("ETag", indexChecksum)
	res.Header().Set("Cache-Control", "public, max-age=300")

	if match := req.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, indexChecksum) {
			res.WriteHeader(http.StatusNotModified)
			return
		}
	}

	count, err := res.Write([]byte(indexHTML))
	if err != nil {
		slog.Error("Failed to write index.html to client", "err", err, slog.Int("written", count))
	}
}

// Starts the server and exits with error if that fails
func (s *Server) Run() error {
	assetFS := StaticFileServer(static.Assets)

	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", s.indexHandler)
	router.HandleFunc("GET /index.html", s.indexHandler)
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", api.NewRouter(s.storage)))
	router.Handle("GET /css/", assetFS)
	router.Handle("GET /icons/", assetFS)
	router.Handle("GET /js/", assetFS)

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
	return fmt.Errorf("failed to start server: %w", err)
}
