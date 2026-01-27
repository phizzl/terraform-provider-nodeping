package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

type MockNodePingServer struct {
	Server   *httptest.Server
	mu       sync.RWMutex
	contacts map[string]map[string]interface{}
	checks   map[string]map[string]interface{}
}

func NewMockNodePingServer() *MockNodePingServer {
	m := &MockNodePingServer{
		contacts: make(map[string]map[string]interface{}),
		checks:   make(map[string]map[string]interface{}),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/contacts", m.handleContacts)
	mux.HandleFunc("/contacts/", m.handleContact)
	mux.HandleFunc("/checks", m.handleChecks)
	mux.HandleFunc("/checks/", m.handleCheck)

	m.Server = httptest.NewServer(mux)
	return m
}

func (m *MockNodePingServer) Close() {
	m.Server.Close()
}

func (m *MockNodePingServer) URL() string {
	return m.Server.URL
}

func (m *MockNodePingServer) handleContacts(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m.contacts)

	case http.MethodPost:
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
			return
		}

		id := "MOCK-CONTACT-" + generateID()
		contact := map[string]interface{}{
			"_id":         id,
			"type":        "contact",
			"customer_id": "MOCK-CUSTOMER",
			"name":        req["name"],
			"custrole":    req["custrole"],
			"addresses":   make(map[string]interface{}),
		}

		if newAddrs, ok := req["newaddresses"].([]interface{}); ok {
			addresses := make(map[string]interface{})
			for _, addr := range newAddrs {
				addrMap := addr.(map[string]interface{})
				addrID := generateID()
				addresses[addrID] = addrMap
			}
			contact["addresses"] = addresses
		}

		m.contacts[id] = contact
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contact)

	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (m *MockNodePingServer) handleContact(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := strings.TrimPrefix(r.URL.Path, "/contacts/")

	switch r.Method {
	case http.MethodGet:
		contact, ok := m.contacts[id]
		if !ok {
			http.Error(w, `{"error": "contact not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contact)

	case http.MethodPut:
		contact, ok := m.contacts[id]
		if !ok {
			http.Error(w, `{"error": "contact not found"}`, http.StatusNotFound)
			return
		}

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if name, ok := req["name"]; ok {
			contact["name"] = name
		}
		if custrole, ok := req["custrole"]; ok {
			contact["custrole"] = custrole
		}

		m.contacts[id] = contact
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contact)

	case http.MethodDelete:
		if _, ok := m.contacts[id]; !ok {
			http.Error(w, `{"error": "contact not found"}`, http.StatusNotFound)
			return
		}
		delete(m.contacts, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "id": id})

	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (m *MockNodePingServer) handleChecks(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m.checks)

	case http.MethodPost:
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
			return
		}

		id := "MOCK-CHECK-" + generateID()
		check := map[string]interface{}{
			"_id":         id,
			"customer_id": "MOCK-CUSTOMER",
			"type":        req["type"],
			"label":       req["label"],
			"enable":      req["enabled"],
			"interval":    req["interval"],
			"state":       1,
			"created":     1609459200000,
			"modified":    1609459200000,
			"parameters": map[string]interface{}{
				"target":    req["target"],
				"threshold": req["threshold"],
				"sens":      req["sens"],
			},
		}

		m.checks[id] = check
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(check)

	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (m *MockNodePingServer) handleCheck(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := strings.TrimPrefix(r.URL.Path, "/checks/")

	switch r.Method {
	case http.MethodGet:
		check, ok := m.checks[id]
		if !ok {
			http.Error(w, `{"error": "check not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(check)

	case http.MethodPut:
		check, ok := m.checks[id]
		if !ok {
			http.Error(w, `{"error": "check not found"}`, http.StatusNotFound)
			return
		}

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if label, ok := req["label"]; ok {
			check["label"] = label
		}
		if enabled, ok := req["enabled"]; ok {
			check["enable"] = enabled
		}

		m.checks[id] = check
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(check)

	case http.MethodDelete:
		if _, ok := m.checks[id]; !ok {
			http.Error(w, `{"error": "check not found"}`, http.StatusNotFound)
			return
		}
		delete(m.checks, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "id": id})

	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

var idCounter int
var idMu sync.Mutex

func generateID() string {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return string(rune('A'+idCounter%26)) + string(rune('A'+(idCounter/26)%26)) + string(rune('0'+idCounter%10))
}

func (m *MockNodePingServer) AddContact(id string, contact map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.contacts[id] = contact
}

func (m *MockNodePingServer) AddCheck(id string, check map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checks[id] = check
}

func (m *MockNodePingServer) GetContact(id string) (map[string]interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.contacts[id]
	return c, ok
}

func (m *MockNodePingServer) GetCheck(id string) (map[string]interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.checks[id]
	return c, ok
}
