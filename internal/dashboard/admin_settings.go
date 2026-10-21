package dashboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

// ImageSettings is global to the NewsGo server, not tied to individual subscriptions.
type ImageSettings struct {
	MetadataTimeoutMs         int `json:"metadataTimeoutMs"`
	MaxHTMLKiB                int `json:"maxHTMLKiB"`
	ParallelRequests          int `json:"parallelRequests"`
	EnrichmentDeadlineSeconds int `json:"enrichmentDeadlineSeconds"`
	MissingImageCacheSeconds  int `json:"missingImageCacheSeconds"`
}
type serverSettings struct {
	Revision  uint64        `json:"revision"`
	Image     ImageSettings `json:"image"`
	UpdatedAt string        `json:"updatedAt,omitempty"`
}

var activeImageSettings atomic.Value

func currentImageSettings() ImageSettings {
	if v := activeImageSettings.Load(); v != nil {
		return v.(ImageSettings)
	}
	return defaultServerSettings().Image
}
func (s *SyncStore) ActivateServerSettings() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, e := s.loadServerSettings()
	if e != nil {
		return e
	}
	activeImageSettings.Store(c.Image)
	return nil
}
func defaultServerSettings() serverSettings {
	return serverSettings{Image: ImageSettings{2500, 512, 12, 25, 90}}
}
func validServerSettings(c serverSettings) bool {
	i := c.Image
	return i.MetadataTimeoutMs >= 500 && i.MetadataTimeoutMs <= 10000 && i.MaxHTMLKiB >= 128 && i.MaxHTMLKiB <= 2048 && i.ParallelRequests >= 1 && i.ParallelRequests <= 16 && i.EnrichmentDeadlineSeconds >= 10 && i.EnrichmentDeadlineSeconds <= 60 && i.MissingImageCacheSeconds >= 30 && i.MissingImageCacheSeconds <= 600
}
func (s *SyncStore) loadServerSettings() (serverSettings, error) {
	b, e := os.ReadFile(s.file("server-settings.json"))
	if errors.Is(e, os.ErrNotExist) {
		return defaultServerSettings(), nil
	}
	if e != nil {
		return serverSettings{}, e
	}
	v := serverSettings{}
	if e = json.Unmarshal(b, &v); e != nil {
		return v, e
	}
	if !validServerSettings(v) {
		return v, errors.New("invalid stored server settings")
	}
	return v, nil
}
func (s *SyncStore) CurrentServerSettings() serverSettings {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, e := s.loadServerSettings()
	if e != nil {
		return defaultServerSettings()
	}
	return v
}
func (s *SyncStore) ServeAdminSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method != http.MethodGet && r.Method != http.MethodPut {
		w.Header().Set("Allow", "GET, PUT")
		http.Error(w, "Method not allowed", 405)
		return
	}
	if r.Method == http.MethodPut {
		if o := r.Header.Get("Origin"); o != "" && o != "https://"+r.Host && o != "http://"+r.Host {
			http.Error(w, "Cross-origin save rejected", 403)
			return
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	token, status := s.authenticate(r)
	if status != 200 {
		http.Error(w, "Approved token required", status)
		return
	}
	if token.Role != "admin" {
		http.Error(w, "Administrator required", 403)
		return
	}
	current, e := s.loadServerSettings()
	if e != nil {
		http.Error(w, "Server settings unavailable", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		_ = json.NewEncoder(w).Encode(current)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var candidate serverSettings
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if e = decoder.Decode(&candidate); e != nil {
		http.Error(w, "Invalid settings JSON", 400)
		return
	}
	var extra any
	if e = decoder.Decode(&extra); e == nil {
		http.Error(w, "Unexpected extra JSON", 400)
		return
	}
	if !validServerSettings(candidate) {
		http.Error(w, "Out-of-range server settings", 400)
		return
	}
	if candidate.Revision != current.Revision {
		http.Error(w, "Settings conflict; reload", 409)
		return
	}
	candidate.Revision++
	candidate.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if e = atomicJSON(s.file("server-settings.json"), candidate); e != nil {
		http.Error(w, fmt.Sprintf("Could not save settings: %v", e), 500)
		return
	}
	activeImageSettings.Store(candidate.Image)
	_ = json.NewEncoder(w).Encode(candidate)
}
