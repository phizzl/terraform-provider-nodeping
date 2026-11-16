package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListChecks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/checks" {
			t.Errorf("expected path /checks, got %s", r.URL.Path)
		}

		checks := map[string]Check{
			"201205050153W2Q4C-0J2HSIRF": {
				ID:         "201205050153W2Q4C-0J2HSIRF",
				CustomerID: "201205050153W2Q4C",
				Label:      "Test Check",
				Type:       "HTTP",
				Enabled:    "active",
				Parameters: CheckParameters{
					Target:    "https://example.com",
					Threshold: 5,
				},
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(checks)
	}))
	defer server.Close()

	c := NewClient(ClientConfig{
		APIToken: "test-token",
		BaseURL:  server.URL,
	})

	checks, err := c.ListChecks(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(checks) != 1 {
		t.Errorf("expected 1 check, got %d", len(checks))
	}

	check, ok := checks["201205050153W2Q4C-0J2HSIRF"]
	if !ok {
		t.Fatal("expected check not found")
	}

	if check.Label != "Test Check" {
		t.Errorf("expected label 'Test Check', got %q", check.Label)
	}
}

func TestGetCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		check := Check{
			ID:         "201205050153W2Q4C-0J2HSIRF",
			CustomerID: "201205050153W2Q4C",
			Label:      "Test Check",
			Type:       "HTTP",
			Enabled:    "active",
			State:      1,
			Parameters: CheckParameters{
				Target:    "https://example.com",
				Threshold: 5,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(check)
	}))
	defer server.Close()

	c := NewClient(ClientConfig{
		APIToken: "test-token",
		BaseURL:  server.URL,
	})

	check, err := c.GetCheck(context.Background(), "201205050153W2Q4C-0J2HSIRF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if check.ID != "201205050153W2Q4C-0J2HSIRF" {
		t.Errorf("expected ID '201205050153W2Q4C-0J2HSIRF', got %q", check.ID)
	}

	if check.State != 1 {
		t.Errorf("expected state 1, got %d", check.State)
	}
}

func TestGetCheckNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Check not found"})
	}))
	defer server.Close()

	c := NewClient(ClientConfig{
		APIToken: "test-token",
		BaseURL:  server.URL,
	})

	_, err := c.GetCheck(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	_, ok := err.(*NotFoundError)
	if !ok {
		t.Errorf("expected *NotFoundError, got %T", err)
	}
}

func TestCreateCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req CheckCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Type != "HTTP" {
			t.Errorf("expected type 'HTTP', got %q", req.Type)
		}

		check := Check{
			ID:         "201205050153W2Q4C-NEWCHK",
			CustomerID: "201205050153W2Q4C",
			Label:      req.Label,
			Type:       req.Type,
			Enabled:    req.Enabled,
			Parameters: CheckParameters{
				Target: req.Target,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(check)
	}))
	defer server.Close()

	c := NewClient(ClientConfig{
		APIToken: "test-token",
		BaseURL:  server.URL,
	})

	check, err := c.CreateCheck(context.Background(), CheckCreateRequest{
		Type:    "HTTP",
		Target:  "https://example.com",
		Label:   "New Check",
		Enabled: "active",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if check.ID != "201205050153W2Q4C-NEWCHK" {
		t.Errorf("expected ID '201205050153W2Q4C-NEWCHK', got %q", check.ID)
	}
}

func TestUpdateCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		check := Check{
			ID:         "201205050153W2Q4C-0J2HSIRF",
			CustomerID: "201205050153W2Q4C",
			Label:      "Updated Check",
			Type:       "HTTP",
			Enabled:    "active",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(check)
	}))
	defer server.Close()

	c := NewClient(ClientConfig{
		APIToken: "test-token",
		BaseURL:  server.URL,
	})

	check, err := c.UpdateCheck(context.Background(), "201205050153W2Q4C-0J2HSIRF", CheckUpdateRequest{
		CheckCreateRequest: CheckCreateRequest{
			Label: "Updated Check",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if check.Label != "Updated Check" {
		t.Errorf("expected label 'Updated Check', got %q", check.Label)
	}
}

func TestDeleteCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DeleteResponse{OK: true, ID: "201205050153W2Q4C-0J2HSIRF"})
	}))
	defer server.Close()

	c := NewClient(ClientConfig{
		APIToken: "test-token",
		BaseURL:  server.URL,
	})

	err := c.DeleteCheck(context.Background(), "201205050153W2Q4C-0J2HSIRF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// inactiveCheckJSON is the shape NodePing returns for a disabled check. For an
// active check "queue" is a string; for an inactive one it is the boolean
// false. Reported against 0.4.1, where Check.Queue was a string and the
// unmarshal failed, so every plan with an inactive check in it broke:
//
//	failed to unmarshal response: json: cannot unmarshal bool into Go
//	struct field Check.queue of type string
//
// The raw JSON is written out by hand on purpose. Encoding a Check struct, as
// the other tests do, can only ever produce the types the struct already has
// and so could never reproduce this.
const inactiveCheckJSON = `{
  "_id": "201205050153W2Q4C-0J2HSIRF",
  "customer_id": "201205050153W2Q4C",
  "label": "inactive check",
  "type": "HTTPADV",
  "enable": "inactive",
  "queue": false,
  "status": false,
  "uuid": false,
  "_rev": false,
  "state": 0,
  "parameters": {
    "target": "https://example.com",
    "sendheaders": {"x-test": "1"}
  }
}`

func TestGetCheckInactive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(inactiveCheckJSON))
	}))
	defer server.Close()

	c := NewClient(ClientConfig{APIToken: "test-token", BaseURL: server.URL})

	check, err := c.GetCheck(context.Background(), "201205050153W2Q4C-0J2HSIRF")
	if err != nil {
		t.Fatalf("reading an inactive check must not fail: %v", err)
	}
	if check.Enabled != "inactive" {
		t.Errorf("Enabled = %q, want %q", check.Enabled, "inactive")
	}
	if check.Parameters.Target != "https://example.com" {
		t.Errorf("Target = %q, want the configured target", check.Parameters.Target)
	}
}

// The list endpoint decodes the same struct, so one inactive check used to
// break nodeping_checks for the whole account, not just the one check.
func TestListChecksWithInactiveCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
  "ACTIVE-1": {"_id": "ACTIVE-1", "type": "HTTP", "enable": "active", "queue": "bQeue1"},
  "INACTIVE-1": ` + inactiveCheckJSON + `
}`))
	}))
	defer server.Close()

	c := NewClient(ClientConfig{APIToken: "test-token", BaseURL: server.URL})

	checks, err := c.ListChecks(context.Background())
	if err != nil {
		t.Fatalf("listing checks that include an inactive one must not fail: %v", err)
	}
	if len(checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(checks))
	}
}

// NodePing answers an update with the check as it stood before the change was
// fully applied: a check switched from disabled to enabled can still report
// "queue": false in the PUT response. Decoding that response failed the same
// way the read did, after NodePing had already stored the change -- so the
// apply errored while the check was in fact updated, and the state was not.
func TestUpdateCheckResponseWithInactiveQueue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(inactiveCheckJSON))
	}))
	defer server.Close()

	c := NewClient(ClientConfig{APIToken: "test-token", BaseURL: server.URL})

	req := CheckUpdateRequest{CheckCreateRequest: CheckCreateRequest{Type: "HTTPADV", Enabled: "active"}}
	if _, err := c.UpdateCheck(context.Background(), "201205050153W2Q4C-0J2HSIRF", req); err != nil {
		t.Fatalf("an update answered with queue=false must not fail: %v", err)
	}
}
