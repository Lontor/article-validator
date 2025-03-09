package parser

import (
	"strings"
	"testing"

	"github.com/Lontor/article-validator/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultParser_Parse(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectAuthors []string
		expectTitle   []string
		expectError   bool
	}{
		{
			name:          "mixed_authors_format",
			input:         "Waals H.; DE VRIES, M.L. Quantum Mechanics Approaches. – Springer, 2021",
			expectAuthors: []string{"Waals", "VRIES"},
			expectTitle:   []string{"Quantum", "Mechanics", "Approaches"},
		},
		{
			name:          "title_with_special_chars",
			input:         "CONNOR P. AI-based $%^&*() Methods (Part 1) // Proc. Conf. – 2024",
			expectAuthors: []string{"CONNOR"},
			expectTitle:   []string{"AI-based", "Methods"},
		},
		{
			name:        "invalid_reference",
			input:       "Numerical Methods for PDEs. – 2023. P. 45-67",
			expectError: true,
		},
		{
			name:          "cyrillic",
			input:         "ИВАНОВ А. А., ПЕТРОВ В. В. Методы численного анализа. – М.: Наука, 2020",
			expectAuthors: []string{"ИВАНОВ", "ПЕТРОВ"},
			expectTitle:   []string{"Методы", "численного", "анализа"},
		},
	}

	var p core.Parser = New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.input)

			if tt.expectError {
				assert.ErrorIs(t, err, ErrInvalidReference)
				return
			}

			require.NoError(t, err)

			for _, expectedAuthor := range tt.expectAuthors {
				assert.Contains(t, got.Authors, expectedAuthor, "missing expected author")
			}

			titleLower := strings.ToLower(got.Title)
			for _, keyword := range tt.expectTitle {
				assert.Contains(t, titleLower, strings.ToLower(keyword), "missing keyword in title")
			}
		})
	}
}
