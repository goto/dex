package warden

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goto/dex/pkg/errors"
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
	const (
		endpoint      = "/api/v1"
		userPath      = "/users/"
		teamsEndpoint = "/teams"
	)
	url := fmt.Sprintf("%s%s%s%s%s", c.BaseURL, endpoint, userPath, req.Email, teamsEndpoint)
	fmt.Println("URL: ", url)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.ErrNotFound.WithMsgf("user with email %s not found", req.Email)
		}

		return nil, errors.ErrInternal.WithMsgf("failed to fetch teams: %v", resp.Status)
	}

	var data teamListResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Data.Teams, nil
}

func (c *Client) TeamByUUID(ctx context.Context, req TeamByUUIDRequest) (*Team, error) {
	const (
		endpoint = "/api/v2"
		teamPath = "/teams/"
	)

	url := fmt.Sprintf("%s%s%s%s", c.BaseURL, endpoint, teamPath, req.TeamUUID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.ErrNotFound.WithMsgf("team with uuid %s not found", req.TeamUUID)
		}

		return nil, errors.ErrInternal.WithMsgf("failed to fetch teams: %v", resp.Status)
	}

	var data teamResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data.Data, nil
}
