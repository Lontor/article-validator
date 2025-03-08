package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Lontor/article-validator/internal/core"
)

type author struct {
	AuthorID string `json:"authorId"`
	Name     string `json:"name"`
}

type paper struct {
	PaperID    string   `json:"paperId"`
	Title      string   `json:"title"`
	Authors    []author `json:"authors"`
	MatchScore float64  `json:"matchScore"`
}

type matchResponse struct {
	Data []paper `json:"data"`
}

type SemanticScholarClient struct {
	endpoint   string
	maxRetries int
}

func New(endpoint string, maxRetries int) *SemanticScholarClient {
	return &SemanticScholarClient{endpoint: endpoint, maxRetries: maxRetries}
}

func (c *SemanticScholarClient) Name() string {
	return "Semantic Scholar"
}

func (c *SemanticScholarClient) Validate(ref core.Reference) (bool, error) {
	params := url.Values{}
	params.Set("query", fmt.Sprintf(`"%s"`, ref.Title))
	params.Set("fields", "title,authors.name")

	u, err := url.Parse(c.endpoint)
	if err != nil {
		return false, err
	}

	u.RawQuery = params.Encode()

	var resp *http.Response

	for retries := 0; retries < c.maxRetries; retries++ {
		resp, err = http.Get(u.String())
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}

	var data matchResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return false, err
	}

	if len(data.Data) == 0 {
		return false, nil
	}

	return true, nil
}
