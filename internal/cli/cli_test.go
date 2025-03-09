package cli

import (
	"bytes"
	"errors"
	"testing"

	"github.com/Lontor/article-validator/internal/core"
)

type mockCore struct {
	resp *core.ValidationResponse
	err  error
}

func (m *mockCore) Validate(string) (*core.ValidationResponse, error) {
	return m.resp, m.err
}

func TestCLI_Run(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		core    *mockCore
		wantOut string
	}{
		{
			name:    "no args",
			args:    []string{},
			wantOut: "No references provided",
		},
		{
			name:    "parser error",
			args:    []string{"invalid-ref"},
			core:    &mockCore{err: errors.New("parse error")},
			wantOut: "[PARSER ERROR] invalid-ref: parse error\n",
		},
		{
			name: "valid reference",
			args: []string{"valid-ref"},
			core: &mockCore{
				resp: &core.ValidationResponse{
					IsValid: true,
					Results: []core.ValidationResult{
						{APIName: "API1", Valid: true},
						{APIName: "API2", Valid: false},
					},
				},
			},
			wantOut: `[VALID] valid-ref
  ✅ API1
  ❌ API2
`,
		},
		{
			name: "api error",
			args: []string{"ref-with-error"},
			core: &mockCore{
				resp: &core.ValidationResponse{
					IsValid: true,
					Results: []core.ValidationResult{
						{APIName: "API1", Valid: false, Error: errors.New("timeout")},
					},
				},
			},
			wantOut: `[VALID] ref-with-error
  ❌ API1 | Error: timeout
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.core)
			c.SetArgs(tt.args)

			out := &bytes.Buffer{}
			c.SetOutput(out)

			c.Run()

			if got := out.String(); got != tt.wantOut {
				t.Errorf("\ngot:\n%s\nwant:\n%s", got, tt.wantOut)
			}
		})
	}
}
