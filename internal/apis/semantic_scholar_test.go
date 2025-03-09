package apis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Lontor/article-validator/internal/core"
)

func TestSemanticScholarClient_Validate(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		handler  func(w http.ResponseWriter, r *http.Request)
		want     bool
		wantErr  bool
	}{
		{
			name: "successful validation",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": []interface{}{
						map[string]string{"paperId": "123"},
					},
				})
			},
			want: true,
		},
		{
			name: "empty result",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
			},
			want: false,
		},
		{
			name: "http error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "invalid json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{"))
			},
			wantErr: true,
		},
		{
			name:     "invalid endpoint",
			handler:  func(w http.ResponseWriter, r *http.Request) {},
			endpoint: "h ttp://invalid",
			wantErr:  true,
		},
		{
			name: "not found status",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "body read error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "100")
				w.WriteHeader(http.StatusOK)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			if tt.handler != nil {
				server = httptest.NewServer(http.HandlerFunc(tt.handler))
				defer server.Close()
			}

			endpoint := server.URL
			if tt.endpoint != "" {
				endpoint = tt.endpoint
			}

			client := NewSemanticScholarClient(
				endpoint,
				3,
				&http.Client{Timeout: 100 * time.Millisecond},
			)

			got, err := client.Validate(core.Reference{Title: "Test"})
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemanticScholarClient_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewSemanticScholarClient(
		server.URL,
		0,
		&http.Client{
			Timeout: 100 * time.Millisecond,
		},
	)

	_, err := client.Validate(core.Reference{Title: "Test"})
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestSemanticScholarClient_RetryOn429(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []interface{}{map[string]string{"title": "Test"}},
		})
	}))
	defer server.Close()

	client := NewSemanticScholarClient(
		server.URL,
		3,
		&http.Client{
			Timeout: 1 * time.Second,
		},
	)

	_, err := client.Validate(core.Reference{Title: "Test"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}
