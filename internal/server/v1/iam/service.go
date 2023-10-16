package iam

import (
	"context"

	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/warden"
)

//go:generate mockery --with-expecter --keeptree --case snake --name WardenClient

type Service struct {
	shieldClient shieldv1beta1rpc.ShieldServiceClient
	wardenClient WardenClient
}

type WardenClient interface {
	ListUserTeams(ctx context.Context, req warden.TeamListRequest) ([]warden.Team, error)
	TeamByUUID(ctx context.Context, req warden.TeamByUUIDRequest) (*warden.Team, error)
}

func NewService(shieldClient shieldv1beta1rpc.ShieldServiceClient, wardenClient WardenClient) *Service {
	return &Service{
		shieldClient: shieldClient,
		wardenClient: wardenClient,
	}
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

func (svc *Service) UpdateGroupMetadata(ctx context.Context, groupID, wardenTeamID string) (map[string]any, error) {
	team, err := svc.wardenClient.TeamByUUID(ctx, warden.TeamByUUIDRequest{
		TeamUUID: wardenTeamID,
	})
	if err != nil {
		return nil, err
	}

	getGroupRes, err := svc.shieldClient.GetGroup(ctx, &shieldv1beta1.GetGroupRequest{
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

	metaData["team-id"] = team.Identifier
	metaData["team-name"] = team.Name
	metaData["product-group-id"] = team.ProductGroupID
	metaData["product-group-name"] = team.ProductGroupName

	updatedMetaData, err := structpb.NewStruct(metaData)
	if err != nil {
		return nil, err
	}

	UpdatedGroupRes, err := svc.shieldClient.UpdateGroup(ctx, &shieldv1beta1.UpdateGroupRequest{
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
