package alert

import (
	"context"
	"fmt"

	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	sirenv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/siren/v1beta1/sirenv1beta1grpc"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	sirenv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/siren/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type SubscriptionService struct {
	sirenClient  sirenv1beta1grpc.SirenServiceClient
	shieldClient shieldv1beta1rpc.ShieldServiceClient
}

func NewSubscriptionService(
	sirenClient sirenv1beta1grpc.SirenServiceClient,
	shieldClient shieldv1beta1rpc.ShieldServiceClient,
) *SubscriptionService {
	return &SubscriptionService{
		sirenClient:  sirenClient,
		shieldClient: shieldClient,
	}
}

func (svc *SubscriptionService) FindSubscription(ctx context.Context, subscriptionID int) (*sirenv1beta1.Subscription, error) {
	request := &sirenv1beta1.GetSubscriptionRequest{
		Id: uint64(subscriptionID),
	}

	resp, err := svc.sirenClient.GetSubscription(ctx, request)
	if err != nil {
		stat := status.Convert(err)
		if stat.Code() == codes.NotFound {
			return nil, ErrSubscriptionNotFound
		}

		return nil, err
	}

	return resp.Subscription, nil
}

func (svc *SubscriptionService) GetSubscriptions(ctx context.Context, groupID, resourceID, resourceType string) ([]*sirenv1beta1.Subscription, error) {
	request := &sirenv1beta1.ListSubscriptionsRequest{
		Metadata: make(map[string]string),
	}
	if groupID != "" {
		request.Metadata[groupMetadataKey] = groupID
	}
	if resourceID != "" {
		request.Metadata[resourceIDMetadataKey] = resourceID
	}
	if resourceType != "" {
		request.Metadata[resourceTypeMetadataKey] = resourceType
	}

	resp, err := svc.sirenClient.ListSubscriptions(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.Subscriptions, nil
}

func (svc *SubscriptionService) CreateSubscription(ctx context.Context, form SubscriptionForm) (int, error) {
	project, group, namespaceID, err := svc.setupForm(ctx, form.ProjectID, form.GroupID)
	if err != nil {
		return 0, err
	}
	channelName, err := svc.getSlackChannelByCriticality(group, form.ChannelCriticality)
	if err != nil {
		return 0, fmt.Errorf("error getting slack channel: %w", err)
	}

	configuration, err := structpb.NewStruct(map[string]interface{}{
		"channel_name": channelName,
	})
	if err != nil {
		return 0, fmt.Errorf("error building configuration: %w", err)
	}
	metadata, err := buildSubscriptionMetadataMap(form, project.Slug, group.Slug)
	if err != nil {
		return 0, err
	}

	request := &sirenv1beta1.CreateSubscriptionRequest{
		Urn:       buildSubscriptionURN(form),
		Namespace: namespaceID,
		Receivers: []*sirenv1beta1.ReceiverMetadata{
			{
				Id:            1,
				Configuration: configuration,
			},
		},
		Match: map[string]string{
			"severity":   alertCriticalSeverityKey,
			"identifier": form.ResourceID,
		},
		Metadata:  metadata,
		CreatedBy: form.UserID,
	}

	resp, err := svc.sirenClient.CreateSubscription(ctx, request)
	if err != nil {
		return 0, err
	}

	return int(resp.Id), nil
}

func (svc *SubscriptionService) UpdateSubscription(ctx context.Context, subscriptionID int, form SubscriptionForm) error {
	project, group, namespaceID, err := svc.setupForm(ctx, form.ProjectID, form.GroupID)
	if err != nil {
		return err
	}
	channelName, err := svc.getSlackChannelByCriticality(group, form.ChannelCriticality)
	if err != nil {
		return fmt.Errorf("error getting slack channel: %w", err)
	}

	configuration, err := structpb.NewStruct(map[string]interface{}{
		"channel_name": channelName,
	})
	if err != nil {
		return fmt.Errorf("error building configuration: %w", err)
	}
	metadata, err := buildSubscriptionMetadataMap(form, project.Slug, group.Slug)
	if err != nil {
		return err
	}

	request := &sirenv1beta1.UpdateSubscriptionRequest{
		Id:        uint64(subscriptionID),
		Urn:       buildSubscriptionURN(form),
		Namespace: namespaceID,
		Receivers: []*sirenv1beta1.ReceiverMetadata{
			{
				Id:            1,
				Configuration: configuration,
			},
		},
		Match: map[string]string{
			"severity":   alertCriticalSeverityKey,
			"identifier": form.ResourceID,
		},
		Metadata:  metadata,
		UpdatedBy: form.UserID,
	}

	_, err = svc.sirenClient.UpdateSubscription(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (svc *SubscriptionService) DeleteSubscription(ctx context.Context, subscriptionID int) error {
	request := &sirenv1beta1.DeleteSubscriptionRequest{
		Id: uint64(subscriptionID),
	}
	_, err := svc.sirenClient.DeleteSubscription(ctx, request)
	if err != nil {
		stat := status.Convert(err)
		if stat.Code() == codes.NotFound {
			return ErrSubscriptionNotFound
		}

		return err
	}

	return nil
}

func (svc *SubscriptionService) setupForm(
	ctx context.Context,
	projectID, groupID string,
) (*shieldv1beta1.Project, *shieldv1beta1.Group, uint64, error) {
	project, err := svc.getProject(ctx, projectID)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("error getting shield's project: %w", err)
	}
	namespaceID, err := svc.getSirenNamespaceID(project)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("error getting siren namespace: %w", err)
	}

	group, err := svc.getGroup(ctx, groupID)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("error getting shield's group: %w", err)
	}

	return project, group, namespaceID, nil
}

func (*SubscriptionService) getSlackChannelByCriticality(group *shieldv1beta1.Group, criticality ChannelCriticality) (string, error) {
	groupMetadata := group.GetMetadata().AsMap()

	valueMap := groupMetadata
	for _, key := range []string{"alerting", string(criticality), "slack"} {
		valInterface, exists := valueMap[key]
		if !exists {
			return "", ErrNoShieldSlackChannel
		}

		var ok bool
		valueMap, ok = valInterface.(map[string]interface{})
		if !ok {
			return "", ErrNoShieldSlackChannel
		}
	}

	slackChannelAny, exists := valueMap["channel"]
	if !exists {
		return "", ErrNoShieldSlackChannel
	}

	channelName, ok := slackChannelAny.(string)
	if !ok {
		return "", ErrNoShieldSlackChannel
	}

	return channelName, nil
}

func (svc *SubscriptionService) getGroup(ctx context.Context, groupID string) (*shieldv1beta1.Group, error) {
	resp, err := svc.shieldClient.GetGroup(ctx, &shieldv1beta1.GetGroupRequest{
		Id: groupID,
	})
	if err != nil {
		return nil, err
	}

	return resp.Group, nil
}

func (svc *SubscriptionService) getProject(ctx context.Context, projectID string) (*shieldv1beta1.Project, error) {
	resp, err := svc.shieldClient.GetProject(ctx, &shieldv1beta1.GetProjectRequest{
		Id: projectID,
	})
	if err != nil {
		return nil, err
	}

	return resp.Project, nil
}

func (*SubscriptionService) getSirenNamespaceID(project *shieldv1beta1.Project) (uint64, error) {
	projectMetadata := project.GetMetadata().AsMap()

	namespaceIDAny, exists := projectMetadata["siren_namespace"]
	if !exists {
		return 0, ErrNoShieldSirenNamespace
	}
	namespaceID, ok := namespaceIDAny.(float64)
	if !ok {
		return 0, ErrNoShieldSirenNamespace
	}

	return uint64(namespaceID), nil
}

func buildSubscriptionMetadataMap(form SubscriptionForm, projectSlug, groupSlug string) (*structpb.Struct, error) {
	metadata, err := structpb.NewStruct(map[string]interface{}{
		"channel_criticality": string(form.ChannelCriticality),
		"group_id":            form.GroupID,
		"resource_type":       form.ResourceType,
		"resource_id":         form.ResourceID,
		"project_id":          form.ProjectID,
		"group_slug":          groupSlug,
		"project_slug":        projectSlug,
	})
	if err != nil {
		return nil, fmt.Errorf("error building metadata: %w", err)
	}

	return metadata, nil
}

func buildSubscriptionURN(form SubscriptionForm) string {
	return fmt.Sprintf("%s:%s:%s:%s", form.GroupID, form.AlertSeverity, form.ResourceType, form.ResourceID)
}
