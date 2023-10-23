package warden

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
	Client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *Client) ListUserTeams(ctx context.Context, req TeamListRequest) ([]Team, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/teams", c.BaseURL, req.Email)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrUserEmailNotFound
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}
		bodyString := string(bodyBytes)
		return nil, errors.New(fmt.Sprintf("got non-200 http status code=(%d) body=%s", resp.StatusCode, bodyString))
	}

	var data teamListResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error decoding warden teamListResponse: %w", err)
	}
	return data.Data.Teams, nil
}

func (c *Client) TeamByUUID(ctx context.Context, req TeamByUUIDRequest) (*Team, error) {
	url := fmt.Sprintf("%s/api/v2/teams/%s", c.BaseURL, req.TeamUUID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrTeamUUIDNotFound
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}
		bodyString := string(bodyBytes)
		return nil, errors.New(fmt.Sprintf("got non-200 http status code=(%d) body=%s", resp.StatusCode, bodyString))
	}

	var data teamResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error decoding teamResponse: %w", err)
	}

	return &data.Data, nil
}
