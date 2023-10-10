package warden

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
)

//go:generate mockery --with-expecter --keeptree --case snake --name Doer

const (
	baseURL = "https://go-cloud.golabs.io"
)

type Service struct {
	shieldClient shieldv1beta1rpc.ShieldServiceClient
	doer         HttpClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewService(shieldClient shieldv1beta1rpc.ShieldServiceClient, doer HttpClient) *Service {
	return &Service{
		shieldClient: shieldClient,
		doer:         doer,
	}
}

func (c *Service) TeamList(ctx context.Context, userEmail string) (*TeamData, error) {
	endpoint := "/api/v1"
	userPath := "/users/"
	teamsEndpoint := "/teams"

	url := baseURL + endpoint + userPath + userEmail + teamsEndpoint

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doer.Do(req)
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

	return nil, ErrEmailNotOnWarden
}

func (c *Service) UpdateGroupMetadata(ctx context.Context, groupID, wardenTeamID string) (map[string]any, error) {
	wardenTeam, err := c.TeamByUUID(ctx, wardenTeamID)
	if err != nil {
		return nil, err
	}

	getGroupRes, err := c.shieldClient.GetGroup(ctx, &shieldv1beta1.GetGroupRequest{
		Id: groupID,
	})
	if err != nil {
		return nil, err
	}

	group := getGroupRes.Group

	metaData := group.Metadata.AsMap()
	if metaData == nil {
		metaData = make(map[string]any)
	}

	metaData["team-id"] = wardenTeam.Identifier
	metaData["product-group-id"] = wardenTeam.ProductGroupID

	updatedMetaData, err := structpb.NewStruct(metaData)
	if err != nil {
		return nil, err
	}

	UpdatedGroupRes, err := c.shieldClient.UpdateGroup(ctx, &shieldv1beta1.UpdateGroupRequest{
		Id: groupID,
		Body: &shieldv1beta1.GroupRequestBody{
			Metadata: updatedMetaData,
			Name:     group.Name,
			Slug:     group.Slug,
			OrgId:    group.OrgId,
		},
	})
	if err != nil {
		return nil, err
	}

	return UpdatedGroupRes.Group.Metadata.AsMap(), nil
}

func (c *Service) TeamByUUID(ctx context.Context, teamByUUID string) (*Team, error) {
	endpoint := "/api/v2"
	teamPath := "/teams/"

	url := baseURL + endpoint + teamPath + teamByUUID

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data TeamResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.Success {
		return &data.Data, nil
	}

	return nil, ErrEmailNotOnWarden
}

type TeamResponse struct {
	Success bool   `json:"success"`
	Data    Team   `json:"data"`
	Message string `json:"message"`
}

type TeamListResponse struct {
	Success bool     `json:"success"`
	Data    TeamData `json:"data"`
}

type TeamData struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	Name                 string    `json:"name"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	OwnerID              int       `json:"owner_id"`
	ParentTeamIdentifier string    `json:"parent_team_identifier"`
	Identifier           string    `json:"identifier"`
	ProductGroupName     string    `json:"product_group_name"`
	ProductGroupID       string    `json:"product_group_id"`
	Labels               any       `json:"labels"`
	ShortCode            string    `json:"short_code"`
}
