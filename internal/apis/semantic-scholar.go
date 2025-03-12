package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Lontor/article-validator/internal/core"
)

type SemanticScholarClient struct {
	endpoint   string
	maxRetries int
	client     *http.Client
}

func NewSemanticScholarClient(endpoint string, maxRetries int, client *http.Client) *SemanticScholarClient {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &SemanticScholarClient{
		endpoint:   endpoint,
		maxRetries: maxRetries,
		client:     client,
	}
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

	for retries := 0; retries <= c.maxRetries; retries++ {
		resp, err = c.client.Get(u.String())
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

	var result struct {
		Data []struct {
			PaperID string `json:"paperId"`
			Title   string `json:"title"`
			Authors []struct {
				AuthorID string `json:"authorId"`
				Name     string `json:"name"`
			} `json:"authors"`
			MatchScore float64 `json:"matchScore"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if len(result.Data) == 0 {
		return false, nil
	}

	return true, nil

}
