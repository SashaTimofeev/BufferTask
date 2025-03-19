package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testTask/pkg/model"
	"time"
)

// Client for working with API
type Client struct {
	SaveFactURL string
	GetFactsURL string
	Token       string
	httpClient  *http.Client
}

// Create new client
func NewClient(saveFactURL, getFactsURL, token string) *Client {
	return &Client{
		SaveFactURL: saveFactURL,
		GetFactsURL: getFactsURL,
		Token:       token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Get facts by period
func (c *Client) GetFacts(ctx context.Context, periodStart, periodEnd, periodKey string, indicatorToMoID int) (string, error) {
	// Forming request data
	form := url.Values{}
	form.Set("period_start", periodStart)
	form.Set("period_end", periodEnd)
	form.Set("period_key", periodKey)
	form.Set("indicator_to_mo_id", strconv.Itoa(indicatorToMoID))

	// Create new POST request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.GetFactsURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	// Send request
	log.Printf("Sending request: %s\n", req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return "", fmt.Errorf("Error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response from server
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		return "", fmt.Errorf("Error reading response: %w", err)
	}

	log.Printf("Response status code: %d", resp.StatusCode)

	return string(bodyBytes), nil
}

func (c *Client) SaveFact(ctx context.Context, fact model.Fact) error {
	// Forming request data
	form := url.Values{}
	form.Set("period_start", fact.PeriodStart)
	form.Set("period_end", fact.PeriodEnd)
	form.Set("period_key", fact.PeriodKey)
	form.Set("indicator_to_mo_id", strconv.Itoa(fact.IndicatorToMoID))
	form.Set("indicator_to_mo_fact_id", strconv.Itoa(fact.IndicatorToMoFactID))
	form.Set("value", strconv.Itoa(fact.Value))
	form.Set("fact_time", fact.FactTime)
	form.Set("is_plan", strconv.Itoa(fact.IsPlan))
	form.Set("auth_user_id", strconv.Itoa(fact.AuthUserID))
	form.Set("comment", fact.Comment)

	// Create new POST request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.SaveFactURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return fmt.Errorf("Error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	// Do request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to save fact, status: %s, body: %s", resp.Status, string(bodyBytes))
		return fmt.Errorf("Failed to save fact, status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	log.Printf("Fact saved successfully: %s", fact.Comment)
	return nil
}
