package dashboard

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// SyncStore persists one user profile per manually issued token. Token values
// are never persisted: the approval list indexes their SHA-256 digests.
type SyncStore struct {
	mu  sync.Mutex
	dir string
}
type approvedToken struct {
	User   string `json:"user"`
	Digest string `json:"digest"`
	Label  string `json:"label,omitempty"`
	Role   string `json:"role,omitempty"`
}
type storedPrefs struct {
	Revision               uint64             `json:"revision"`
	Subscriptions          []syncSubscription `json:"subscriptions"`
	ExcludedDomains        []string           `json:"excludedDomains"`
	ExcludedArticleDomains []string           `json:"excludedArticleDomains"`
 MaxArticleAgeDays int `json:"maxArticleAgeDays,omitempty"`
 LastReadingSelection string `json:"lastReadingSelection,omitempty"`
	UpdatedAt              string             `json:"updatedAt"`
}
type syncSubscription struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	URL      string `json:"url,omitempty"`
	Category string `json:"category,omitempty"`
}
type userPrefs struct {
	User string `json:"user"`
	Role string `json:"role,omitempty"`
	storedPrefs
}

func NewSyncStore(dir string) (*SyncStore, error) {
	if dir == "" {
		return nil, errors.New("sync directory is empty")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &SyncStore{dir: dir}, nil
}
func digest(token string) string             { d := sha256.Sum256([]byte(token)); return hex.EncodeToString(d[:]) }
func (s *SyncStore) file(name string) string { return filepath.Join(s.dir, name) }
func (s *SyncStore) loadTokens() ([]approvedToken, error) {
	b, e := os.ReadFile(s.file("approved-tokens.json"))
	if errors.Is(e, os.ErrNotExist) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	var v []approvedToken
	if e = json.Unmarshal(b, &v); e != nil {
		return nil, e
	}
	return v, nil
}
func (s *SyncStore) loadPrefs() (map[string]storedPrefs, error) {
	b, e := os.ReadFile(s.file("user-preferences.json"))
	if errors.Is(e, os.ErrNotExist) {
		return map[string]storedPrefs{}, nil
	}
	if e != nil {
		return nil, e
	}
	v := map[string]storedPrefs{}
	if e = json.Unmarshal(b, &v); e != nil {
		return nil, e
	}
	if v == nil {
		v = map[string]storedPrefs{}
	}
	return v, nil
}
func atomicJSON(path string, v any) error {
	b, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		return e
	}
	b = append(b, '\n')
	f, e := os.CreateTemp(filepath.Dir(path), ".newsgo-*.tmp")
	if e != nil {
		return e
	}
	defer os.Remove(f.Name())
	if e = f.Chmod(0600); e != nil {
		f.Close()
		return e
	}
	if _, e = f.Write(b); e != nil {
		f.Close()
		return e
	}
	if e = f.Sync(); e != nil {
		f.Close()
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	return os.Rename(f.Name(), path)
}

// IssueToken prints the token only to its caller. No plaintext token is saved.
func (s *SyncStore) IssueToken(user, label string) (string, error) {
	if len(user) < 1 || len(user) > 64 || strings.Trim(user, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_- ") != "" {
		return "", errors.New("user must be 1-64 letters, numbers, spaces, hyphens or underscores")
	}
	var raw [32]byte
	if _, e := rand.Read(raw[:]); e != nil {
		return "", e
	}
	token := hex.EncodeToString(raw[:])
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, e := s.loadTokens()
	if e != nil {
		return "", e
	}
	for _, entry := range entries {
		if entry.User == user {
			return "", fmt.Errorf("user %q already has a token; revoke before reissuing", user)
		}
	}
	entries = append(entries, approvedToken{User: user, Digest: digest(token), Label: label})
	if e = atomicJSON(s.file("approved-tokens.json"), entries); e != nil {
		return "", e
	}
	return token, nil
}
func (s *SyncStore) Revoke(user string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, e := s.loadTokens()
	if e != nil {
		return e
	}
	out := make([]approvedToken, 0, len(entries))
	found := false
	for _, v := range entries {
		if v.User == user {
			found = true
		} else {
			out = append(out, v)
		}
	}
	if !found {
		return fmt.Errorf("unknown user %q", user)
	}
	return atomicJSON(s.file("approved-tokens.json"), out)
}
func validPrefs(p storedPrefs) bool {
 if p.LastReadingSelection != "" && p.LastReadingSelection != "all" {
  found := false
  for _, sub := range p.Subscriptions { if sub.Kind == "feed" && sub.URL == p.LastReadingSelection { found = true; break } }
  if !found { return false }
 }
	if p.MaxArticleAgeDays != 0 && p.MaxArticleAgeDays != 1 && p.MaxArticleAgeDays != 3 && p.MaxArticleAgeDays != 7 && p.MaxArticleAgeDays != 14 && p.MaxArticleAgeDays != 30 && p.MaxArticleAgeDays != 60 { return false }
 if len(p.Subscriptions) > 250 || len(p.ExcludedDomains) > 500 || len(p.ExcludedArticleDomains) > 500 {
		return false
	}
	ids := map[string]bool{}
	for _, r := range p.Subscriptions {
		if len(r.ID) < 1 || len(r.ID) > 2200 || len(r.Name) < 1 || len(r.Name) > 140 || ids[r.ID] {
			return false
		}
		ids[r.ID] = true
		if r.Kind == "feed" {
			if !strings.HasPrefix(r.ID, "feed:") || !strings.HasPrefix(r.URL, "https://") && !strings.HasPrefix(r.URL, "http://") || len(r.URL) > 2048 {
				return false
			}
		}
		if r.Kind == "shared" {
			if !strings.HasPrefix(r.ID, "shared:") {
				return false
			}
		} else if r.Kind != "feed" {
			return false
		}
	}
	for _, v := range append(append([]string{}, p.ExcludedDomains...), p.ExcludedArticleDomains...) {
		if v == "" || len(v) > 253 || strings.ContainsAny(v, " /:@\n\t") || strings.ToLower(v) != v || !strings.Contains(v, ".") {
			return false
		}
	}
	return true
}
func (s *SyncStore) ServeSync(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method != http.MethodGet && r.Method != http.MethodPut {
		w.Header().Set("Allow", "GET, PUT")
		http.Error(w, "Method not allowed", 405)
		return
	}
	// A token is a bearer secret; accept it only in Authorization, never URLs.
	parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || len(parts[1]) != 64 {
		http.Error(w, "Approved token required", 401)
		return
	}
	if r.Method == http.MethodPut {
		if o := r.Header.Get("Origin"); o != "" { // Reject cross-origin writes.
			host := r.Host
			if o != "https://"+host && o != "http://"+host {
				http.Error(w, "Cross-origin save rejected", 403)
				return
			}
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tokens, e := s.loadTokens()
	if e != nil {
		http.Error(w, "Token storage unavailable", 500)
		return
	}
	user := ""
	requested := digest(parts[1])
	for _, t := range tokens {
		if t.Digest == requested {
			user = t.User
			break
		}
	}
	if user == "" {
		http.Error(w, "Token not approved", 401)
		return
	}
	records, e := s.loadPrefs()
	if e != nil {
		http.Error(w, "Preference storage unavailable", 500)
		return
	}
	current, exists := records[user]
	if r.Method == http.MethodGet {
		if !exists {
			w.WriteHeader(404)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "No server preferences saved for this user", "role": tokenRole(tokens, user)})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(userPrefs{User: user, Role: tokenRole(tokens, user), storedPrefs: current})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 128<<10)
	var input storedPrefs
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if e = dec.Decode(&input); e != nil {
		http.Error(w, "Invalid preferences JSON", 400)
		return
	}
	var extra any
	if e = dec.Decode(&extra); e == nil {
		http.Error(w, "Unexpected extra JSON", 400)
		return
	}
	if !validPrefs(input) {
		http.Error(w, "Invalid subscriptions or excluded domains", 400)
		return
	}
	if exists && input.Revision != current.Revision || !exists && input.Revision != 0 {
		http.Error(w, "Conflict: newer server preferences exist; use Load from Server", 409)
		return
	}
	input.Revision++
	input.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	records[user] = input
	if e = atomicJSON(s.file("user-preferences.json"), records); e != nil {
		http.Error(w, "Could not save preferences", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userPrefs{User: user, storedPrefs: input})
}

func tokenRole(tokens []approvedToken, user string) string {
	for _, tok := range tokens {
		if tok.User == user && tok.Role == "admin" {
			return "admin"
		}
	}
	return "user"
}

// authenticate must be called while s.mu is held.
func (s *SyncStore) authenticate(r *http.Request) (approvedToken, int) {
	parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || len(parts[1]) != 64 {
		return approvedToken{}, http.StatusUnauthorized
	}
	tokens, err := s.loadTokens()
	if err != nil {
		return approvedToken{}, http.StatusInternalServerError
	}
	requested := digest(parts[1])
	for _, entry := range tokens {
		if entry.Digest == requested {
			return entry, http.StatusOK
		}
	}
	return approvedToken{}, http.StatusUnauthorized
}

// SetRole promotes or demotes a previously issued token without rotating it.
func (s *SyncStore) SetRole(user, role string) error {
	if role != "admin" && role != "user" {
		return errors.New("role must be admin or user")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, err := s.loadTokens()
	if err != nil {
		return err
	}
	found := false
	for i := range entries {
		if entries[i].User == user {
			found = true
			if role == "user" {
				entries[i].Role = ""
			} else {
				entries[i].Role = "admin"
			}
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown user %q", user)
	}
	return atomicJSON(s.file("approved-tokens.json"), entries)
}
