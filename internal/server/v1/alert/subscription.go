package alert

import (
	"context"
	"fmt"

	sirenv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/siren/v1beta1/sirenv1beta1grpc"
	sirenv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/siren/v1beta1"
	"github.com/goto/dex/generated/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type SubscriptionService struct {
	sirenClient sirenv1beta1grpc.SirenServiceClient
}

func NewSubscriptionService(sirenClient sirenv1beta1grpc.SirenServiceClient) *SubscriptionService {
	return &SubscriptionService{
		sirenClient: sirenClient,
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
		request.Metadata[resourceIdMetadataKey] = resourceID
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

func (svc *SubscriptionService) CreateSubscription(ctx context.Context, form models.SubscriptionForm, channelName, userID string) (subscriptionID int, err error) {
	configuration, err := structpb.NewStruct(map[string]interface{}{
		"channel_name": channelName,
	})
	if err != nil {
		return 0, fmt.Errorf("error building configuration: %w", err)
	}
	metadata, err := structpb.NewStruct(map[string]interface{}{
		"channel_name": channelName,
	})
	if err != nil {
		return 0, fmt.Errorf("error building metadata: %w", err)
	}

	request := &sirenv1beta1.CreateSubscriptionRequest{
		Urn:       svc.buildSubscriptionURN(form),
		Namespace: 1,
		Receivers: []*sirenv1beta1.ReceiverMetadata{
			{
				Id:            1,
				Configuration: configuration,
			},
		},
		Match: map[string]string{
			"severity":   alertCriticalSeverityKey,
			"identifier": *form.ResourceID,
		},
		Metadata:  metadata,
		CreatedBy: userID,
	}

	resp, err := svc.sirenClient.CreateSubscription(ctx, request)
	if err != nil {
		return 0, err
	}

	return int(resp.Id), nil
}

func (svc *SubscriptionService) UpdateSubscription(ctx context.Context, subscriptionID int, form models.SubscriptionForm, channelName, userID string) error {
	configuration, err := structpb.NewStruct(map[string]interface{}{
		"channel_name": channelName,
	})
	if err != nil {
		return fmt.Errorf("error building configuration: %w", err)
	}
	metadata, err := structpb.NewStruct(map[string]interface{}{
		"channel_name": channelName,
	})
	if err != nil {
		return fmt.Errorf("error building metadata: %w", err)
	}

	request := &sirenv1beta1.UpdateSubscriptionRequest{
		Id:        uint64(subscriptionID),
		Urn:       svc.buildSubscriptionURN(form),
		Namespace: 1,
		Receivers: []*sirenv1beta1.ReceiverMetadata{
			{
				Id:            1,
				Configuration: configuration,
			},
		},
		Match: map[string]string{
			"severity":   alertCriticalSeverityKey,
			"identifier": *form.ResourceID,
		},
		Metadata:  metadata,
		UpdatedBy: userID,
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
		return err
	}

	return nil
}

func (svc *SubscriptionService) buildSubscriptionURN(form models.SubscriptionForm) string {
	return fmt.Sprintf("%s:%s:%s:%s", *form.GroupID, *form.AlertSeverity, *form.ResourceType, *form.ResourceID)
}
