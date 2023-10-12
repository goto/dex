package warden

import (
	"context"

	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/warden"
)

//go:generate mockery --with-expecter --keeptree --case snake --name Doer

type Service struct {
	shieldClient shieldv1beta1rpc.ShieldServiceClient
	wardenClient *warden.Client
}

func NewService(shieldClient shieldv1beta1rpc.ShieldServiceClient, wardenClient *warden.Client) *Service {
	return &Service{
		shieldClient: shieldClient,
		wardenClient: wardenClient,
	}
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

	metaData["team-id"] = wardenTeam.Data.Identifier
	metaData["product-group-id"] = wardenTeam.Data.ProductGroupID

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

func (svc *Service) TeamList(ctx context.Context, userEmail string) ([]warden.Team, error) {
	teams, err := svc.wardenClient.ListUserTeams(ctx, warden.TeamListRequest{
		Email: userEmail,
	})
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (svc *Service) TeamByUUID(ctx context.Context, teamByUUID string) (*warden.TeamResponse, error) {
	team, err := svc.wardenClient.WardenTeamByUUID(ctx, warden.TeamByUUIDRequest{
		TeamUUID: teamByUUID,
	})
	if err != nil {
		return nil, err
	}

	return &warden.TeamResponse{
		Success: true,
		Data:    *team,
	}, nil
}
