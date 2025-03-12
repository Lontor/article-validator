package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Lontor/article-validator/internal/core"
)

type CrossrefClient struct {
	endpoint   string
	maxRetries int
	client     *http.Client
}

func NewCrossrefClient(endpoint string, maxRetries int, client *http.Client) *CrossrefClient {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &CrossrefClient{endpoint: endpoint, maxRetries: maxRetries, client: client}
}

func (c *CrossrefClient) Name() string {
	return "Crossref"
}

func (c *CrossrefClient) Validate(ref core.Reference) (bool, error) {
	params := url.Values{}
	params.Set("query", fmt.Sprintf(`"%s"`, ref.Title))
	params.Set("query.author", strings.Join(ref.Authors, " "))
	params.Set("select", "title,container-title,author")
	params.Set("sort", "relevance")
	params.Set("rows", "1")

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
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	var result struct {
		Message struct {
			Item []struct {
				Title          []string `json:"title"`
				ContainerTitle []string `json:"container-title"`
				Author         []struct {
					Family string `json:"family"`
				} `json:"author"`
			} `json:"items"`
		} `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if len(result.Message.Item) == 0 {
		return false, nil
	}

	item := result.Message.Item[0]
	if len(item.Title) == 0 {
		return false, nil
	}

	refAuthors := normalizeStringSlice(ref.Authors)

	refTitle := normalizeString(ref.Title) + " " + refAuthors
	apiTitle := normalizeString(item.Title[0]) + " " + refAuthors
	apiContainerTitle := normalizeString(item.ContainerTitle[0]) + " " + refAuthors

	match := max(calculateSimilarity(refTitle, apiTitle), calculateSimilarity(refTitle, apiContainerTitle))

	return match >= 0.8, nil
}

func normalizeString(str string) string {
	str = strings.ToLower(str)

	reg := regexp.MustCompile(`[^\p{L}\s]`)
	str = reg.ReplaceAllString(str, "")

	str = strings.Join(strings.Fields(str), " ")

	return str
}

func normalizeStringSlice(slc []string) string {
	normalized := make([]string, len(slc))
	for i, el := range slc {
		normalized[i] = normalizeString(el)
	}

	return strings.Join(normalized, " ")
}

func calculateSimilarity(a, b string) float64 {
	setA := make(map[string]struct{})
	for _, word := range strings.Fields(a) {
		setA[word] = struct{}{}
	}

	setB := make(map[string]struct{})
	for _, word := range strings.Fields(b) {
		setB[word] = struct{}{}
	}

	intersection := 0
	for word := range setA {
		if _, exists := setB[word]; exists {
			intersection++
		}
	}

	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}
