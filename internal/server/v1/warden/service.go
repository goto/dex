package warden

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/goto/dex/internal/server/reqctx"
)

const hostName = "https://go-cloud.golabs.io/api/v1/users/akarsh.satija@gojek.com/teams"

type Service struct {
	hostName string
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewService(doer Doer) *Service {
	return &Service{

		hostName: hostName,
	}
}

func (c *Service) TeamList(ctx context.Context) (any, error) {

	reqCtx := reqctx.From(ctx)

	fmt.Println("Email:", reqCtx)

	resp, err := http.Get("https://go-cloud.golabs.io/api/v1/users/sudheer.pal@gojek.com/teams")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil

}

type TeamListResponse struct {
	Success bool   `json:"success"`
	Data    []Team `json:"data"`
}

type Team struct {
	Identifier uuid.UUID `json:"identifier"`
	Name       string    `json:"name"`
}
