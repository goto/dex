package warden

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/internal/server/reqctx"
)

const (
	baseURL = "https://go-cloud.golabs.io"
)

type Service struct {
	shieldClient shieldv1beta1rpc.ShieldServiceClient
}

func NewService(shieldClient shieldv1beta1rpc.ShieldServiceClient) *Service {
	return &Service{
		shieldClient: shieldClient,
	}
}

func (*Service) TeamList(ctx context.Context) (*TeamData, error) {
	endpoint := "/api/v1"
	userPath := "/users/"
	teamsEndpoint := "/teams"

	reqCtx := reqctx.From(ctx)

	if reqCtx.UserEmail == "" {
		return nil, ErrUserNotFound
	}

	url := baseURL + endpoint + userPath + reqCtx.UserEmail + teamsEndpoint

	resp, err := http.Get(url) //nolint
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

func (c *Service) UpdateGroupMetadata(ctx context.Context, groupID, wardenTeamID string) (map[string]any, error) {
	shielTeam, err := c.TeamByUUID(ctx, wardenTeamID)
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

	metaData["team-id"] = shielTeam.Identifier
	metaData["product-group-id"] = shielTeam.ProductGroupID

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

func (*Service) TeamByUUID(_ context.Context, teamByUUID string) (*Team, error) {
	endpoint := "/api/v2"
	teamPath := "/teams/"

	url := baseURL + endpoint + teamPath + teamByUUID

	resp, err := http.Get(url) //nolint
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

	return nil, ErrTeamNotFound
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
