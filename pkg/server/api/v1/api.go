package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/heathcliff26/go-wol/pkg/wol"
)

type Response struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /wake/{macAddr}", WakeHandler)
	return router
}

// /wake/{macAddr}
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
