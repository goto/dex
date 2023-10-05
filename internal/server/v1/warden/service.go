package warden

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/goto/dex/internal/server/reqctx"
)

const baseURL = "https://go-cloud.golabs.io"
const endpoint = "/api/v1"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (c *Service) TeamList(ctx context.Context) (*TeamData, error) {

	userPath := "/users/"
	teamsEndpoint := "/teams"
	reqCtx := reqctx.From(ctx)

	if reqCtx.UserEmail == "" {
		return nil, ErrUserNotFound
	}

	url := baseURL + endpoint + userPath + reqCtx.UserEmail + teamsEndpoint

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data TeamListResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.Success {
		return &data.Data, nil
	}

	return nil, ErrTeamNotFound
}

type TeamListResponse struct {
	Success bool     `json:"success"`
	Data    TeamData `json:"data"`
}

type TeamData struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	Name       string    `json:"name"`
	Identifier uuid.UUID `json:"identifier"`
}
