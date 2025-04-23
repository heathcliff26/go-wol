package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/heathcliff26/go-wol/pkg/server/storage"
	"github.com/heathcliff26/go-wol/pkg/utils"
	"github.com/heathcliff26/go-wol/pkg/wol"
)

type Response struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type apiHandler struct {
	storage *storage.Storage
}

func NewRouter(storage *storage.Storage) *http.ServeMux {
	handler := &apiHandler{
		storage: storage,
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /wake/{macAddr}", WakeHandler)
	router.HandleFunc("PUT /hosts/{macAddr}/{name}", handler.AddHostHandler)
	router.HandleFunc("DELETE /hosts/{macAddr}", handler.RemoveHostHandler)
	return router
}

// GET /wake/{macAddr}
// Send a magic packet to the specified MAC address
func WakeHandler(res http.ResponseWriter, req *http.Request) {
	macAddr := req.PathValue("macAddr")

	packet, err := wol.CreatePacket(macAddr)
	if err != nil {
		slog.Info("Client send invalid MAC address", slog.String("mac", macAddr), slog.Any("error", err))
		res.WriteHeader(http.StatusBadRequest)
		sendResponse(res, "Invalid MAC address")
		return
	}

	err = packet.Send("")
	if err != nil {
		slog.Info("Failed to send magic packet", slog.String("mac", macAddr), slog.Any("error", err))
		res.WriteHeader(http.StatusBadRequest)
		sendResponse(res, "Failed to send magic packet")
		return
	}

	slog.Info("Send magic packet", slog.String("mac", macAddr))
	sendResponse(res, "")
}

// PUT /hosts/{macAddr}/{name}
// Add a host to the storage
func (h *apiHandler) AddHostHandler(res http.ResponseWriter, req *http.Request) {
	macAddr := req.PathValue("macAddr")
	name := req.PathValue("name")

	if !utils.ValidateMACAddress(macAddr) {
		slog.Debug("Client send invalid MAC address", slog.String("mac", macAddr))
		res.WriteHeader(http.StatusBadRequest)
		sendResponse(res, "Invalid MAC address")
		return
	}

	err := h.storage.AddHost(macAddr, name)
	if err != nil {
		slog.Error("Failed to add host", "mac", macAddr, "name", name, "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		sendResponse(res, "Failed to add host")
		return
	}

	slog.Info("Added host", slog.String("mac", macAddr), slog.String("name", name))
	sendResponse(res, "")
}

// DELETE /hosts/{macAddr}
// Remove a host from the storage
func (h *apiHandler) RemoveHostHandler(res http.ResponseWriter, req *http.Request) {
	macAddr := req.PathValue("macAddr")

	if !utils.ValidateMACAddress(macAddr) {
		slog.Debug("Client send invalid MAC address", slog.String("mac", macAddr))
		res.WriteHeader(http.StatusBadRequest)
		sendResponse(res, "Invalid MAC address")
		return
	}

	err := h.storage.RemoveHost(macAddr)
	if err != nil {
		slog.Error("Failed to remove host", "mac", macAddr, "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		sendResponse(res, "Failed to remove host")
		return
	}

	slog.Info("Removed host", slog.String("mac", macAddr))
	sendResponse(res, "")
}

func sendResponse(rw http.ResponseWriter, reason string) {
	response := Response{
		Status: "error",
		Reason: reason,
	}
	if reason == "" {
		response.Status = "ok"
	}

	b, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		slog.Error("Failed to create Response", "err", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	_, err = rw.Write(b)
	if err != nil {
		slog.Error("Failed to send response to client", "err", err)
	}
}
