package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:5380", "my-secret-token")
	if c.BaseURL == nil || c.BaseURL.String() != "http://localhost:5380" {
		t.Errorf("expected BaseURL http://localhost:5380, got %v", c.BaseURL)
	}
	if c.Token != "my-secret-token" {
		t.Errorf("expected Token my-secret-token, got %s", c.Token)
	}
	if c.HTTPClient == nil {
		t.Error("expected non-nil HTTPClient")
	}
}

func TestAddReservation(t *testing.T) {
	t.Run("validation error - empty IP", func(t *testing.T) {
		c := NewClient("http://localhost:5380", "token")
		err := c.AddReservation("", "00:11:22:33:44:55", "host", "comment")
		if err == nil || err.Error() != "IP address cannot be empty" {
			t.Errorf("expected IP address cannot be empty error, got: %v", err)
		}
	})

	t.Run("validation error - empty MAC", func(t *testing.T) {
		c := NewClient("http://localhost:5380", "token")
		err := c.AddReservation("192.168.1.10", "", "host", "comment")
		if err == nil || err.Error() != "MAC address cannot be empty" {
			t.Errorf("expected MAC address cannot be empty error, got: %v", err)
		}
	})

	t.Run("validation error - empty hostname", func(t *testing.T) {
		c := NewClient("http://localhost:5380", "token")
		err := c.AddReservation("192.168.1.10", "00:11:22:33:44:55", "", "comment")
		if err == nil || err.Error() != "hostname cannot be empty" {
			t.Errorf("expected hostname cannot be empty error, got: %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Path != "/api/dhcp/scopes/addReservedLease" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if r.Header.Get("Authorization") != "Bearer token123" {
				w.Write([]byte(`{"status": "error", "errorMessage": "unauthorized"}`))
				return
			}
			q := r.URL.Query()
			if q.Get("address") != "192.168.1.10" ||
				q.Get("hardwareAddress") != "00:11:22:33:44:55" ||
				q.Get("hostName") != "myhost" ||
				q.Get("comments") != "mycomment" {
				w.Write([]byte(`{"status": "error", "errorMessage": "invalid params"}`))
				return
			}
			w.Write([]byte(`{"status": "ok"}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		err := c.AddReservation("192.168.1.10", "00:11:22:33:44:55", "myhost", "mycomment")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("failure - API error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status": "error", "errorMessage": "DHCP reservation already exists"}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		err := c.AddReservation("192.168.1.10", "00:11:22:33:44:55", "myhost", "mycomment")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "API error: DHCP reservation already exists") {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

func TestRemoveReservation(t *testing.T) {
	t.Run("validation error - empty MAC", func(t *testing.T) {
		c := NewClient("http://localhost:5380", "token")
		err := c.RemoveReservation("")
		if err == nil || err.Error() != "MAC address cannot be empty" {
			t.Errorf("expected MAC address cannot be empty error, got: %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Path != "/api/dhcp/scopes/removeReservedLease" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if r.Header.Get("Authorization") != "Bearer token123" {
				w.Write([]byte(`{"status": "error", "errorMessage": "unauthorized"}`))
				return
			}
			q := r.URL.Query()
			if q.Get("hardwareAddress") != "00:11:22:33:44:55" {
				w.Write([]byte(`{"status": "error", "errorMessage": "invalid params"}`))
				return
			}
			w.Write([]byte(`{"status": "ok"}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		err := c.RemoveReservation("00:11:22:33:44:55")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("failure - API error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status": "error", "errorMessage": "DHCP reservation not found"}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		err := c.RemoveReservation("00:11:22:33:44:55")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "API error: DHCP reservation not found") {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

func TestSetAppConfig(t *testing.T) {
	t.Run("validation error - empty configJSON", func(t *testing.T) {
		c := NewClient("http://localhost:5380", "token")
		err := c.SetAppConfig("")
		if err == nil || err.Error() != "configJSON cannot be empty" {
			t.Errorf("expected configJSON cannot be empty error, got: %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Path != "/api/apps/config/set" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if r.Header.Get("Authorization") != "Bearer token123" {
				w.Write([]byte(`{"status": "error", "errorMessage": "unauthorized"}`))
				return
			}
			q := r.URL.Query()
			if q.Get("name") != "Advanced Blocking" {
				w.Write([]byte(`{"status": "error", "errorMessage": "invalid app name"}`))
				return
			}
			if err := r.ParseForm(); err != nil {
				w.Write([]byte(`{"status": "error", "errorMessage": "failed to parse form"}`))
				return
			}
			if r.FormValue("config") != `{"enable": true}` {
				w.Write([]byte(`{"status": "error", "errorMessage": "invalid config"}`))
				return
			}
			w.Write([]byte(`{"status": "ok"}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		err := c.SetAppConfig(`{"enable": true}`)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("failure - API error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status": "error", "errorMessage": "invalid json config"}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		err := c.SetAppConfig(`{invalid}`)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "API error: invalid json config") {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

func TestFetchCurrentScope(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Path != "/api/dhcp/scopes/list" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Write([]byte(`{"status": "ok", "response": [{"name": "Default"}]}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		res, err := c.FetchCurrentScope()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		var parsed APIResponse
		if err := json.Unmarshal(res, &parsed); err != nil {
			t.Fatalf("expected valid JSON, got unmarshal error: %v", err)
		}
		if parsed.Status != "ok" {
			t.Errorf("expected status 'ok', got: %s", parsed.Status)
		}
	})

	t.Run("failure - server error status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`internal server error`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		_, err := c.FetchCurrentScope()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "HTTP error: 500") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestGetAppConfig(t *testing.T) {
	t.Run("success - string response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Path != "/api/apps/config/get" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if r.URL.Query().Get("name") != "Advanced Blocking" {
				w.Write([]byte(`{"status": "error", "errorMessage": "invalid app name"}`))
				return
			}
			w.Write([]byte(`{"status": "ok", "response": {"config": "{\"enableBlocking\": true}"}}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		res, err := c.GetAppConfig()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		expected := `{"enableBlocking": true}`
		if res != expected {
			t.Errorf("expected config %s, got: %s", expected, res)
		}
	})

	t.Run("success - object response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status": "ok", "response": {"config": {"enableBlocking": true}}}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		res, err := c.GetAppConfig()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if !strings.Contains(res, `"enableBlocking":true`) && !strings.Contains(res, `"enableBlocking": true`) {
			t.Errorf("expected json object, got: %s", res)
		}
	})

	t.Run("success - null response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status": "ok", "response": {"config": null}}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "token123")
		res, err := c.GetAppConfig()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if res != "" {
			t.Errorf("expected empty config, got: %s", res)
		}
	})
}
