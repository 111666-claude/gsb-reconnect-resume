package session

import (
	"encoding/json"
	"net/http"
)

// Handler 暴露开局与续局两个接口。
func (m *Manager) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/session/open", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Player string `json:"player"`
			AtMs   int64  `json:"atMs"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		_ = json.NewEncoder(w).Encode(m.Open(req.Player, req.AtMs))
	})

	mux.HandleFunc("/session/resume", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RequestID string `json:"requestId"`
			Token     string `json:"token"`
			Version   int    `json:"version"`
			AtMs      int64  `json:"atMs"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		snap, reason := m.Resume(req.RequestID, req.Token, req.Version, req.AtMs)
		if reason != "" {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": reason})
			return
		}
		_ = json.NewEncoder(w).Encode(snap)
	})

	return mux
}
