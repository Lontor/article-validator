package core

import (
	"errors"
	"testing"
)

type mockParser struct {
	ref Reference
	err error
}

func (m *mockParser) Parse(raw string) (Reference, error) {
	return m.ref, m.err
}

type mockAPIClient struct {
	name  string
	valid bool
	err   error
}

func (m *mockAPIClient) Name() string { return m.name }
func (m *mockAPIClient) Validate(Reference) (bool, error) {
	return m.valid, m.err
}

func TestCore_Validate(t *testing.T) {
	tests := []struct {
		name      string
		parser    *mockParser
		clients   []APIClient
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "parser error",
			parser:    &mockParser{err: errors.New("parse error")},
			wantValid: false,
			wantErr:   true,
		},
		{
			name: "all clients invalid",
			parser: &mockParser{
				ref: Reference{Title: "Test", Authors: []string{"Author"}},
			},
			clients: []APIClient{
				&mockAPIClient{name: "API1", valid: false},
				&mockAPIClient{name: "API2", valid: false},
			},
			wantValid: false,
		},
		{
			name: "one client valid",
			parser: &mockParser{
				ref: Reference{Title: "Test", Authors: []string{"Author"}},
			},
			clients: []APIClient{
				&mockAPIClient{name: "API1", valid: false},
				&mockAPIClient{name: "API2", valid: true},
			},
			wantValid: true,
		},
		{
			name: "client error",
			parser: &mockParser{
				ref: Reference{Title: "Test", Authors: []string{"Author"}},
			},
			clients: []APIClient{
				&mockAPIClient{name: "API1", err: errors.New("API error")},
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.parser, tt.clients)
			resp, err := c.Validate("any string")

			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil && !tt.wantErr {
				t.Fatal("nil response")
			}
			if resp != nil && resp.IsValid != tt.wantValid {
				t.Errorf("isValid: got %v, want %v", resp.IsValid, tt.wantValid)
			}
		})
	}
}
