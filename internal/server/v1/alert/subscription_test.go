package alert_test

import (
	"context"
	"fmt"
	"testing"

	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	sirenv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/siren/v1beta1"
	sirenReceiverPkg "github.com/goto/siren/core/receiver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/v1/alert"
	"github.com/goto/dex/mocks"
)

func TestSubscriptionServiceFindSubscription(t *testing.T) {
	ctx := context.TODO()
	subscriptionID := 105

	t.Run("should return subscription on success", func(t *testing.T) {
		subscription := &sirenv1beta1.Subscription{
			Id:        uint64(subscriptionID),
			Urn:       "sample-urn",
			Namespace: 1,
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{Id: 30},
			},
		}

		shield := new(mocks.ShieldServiceClient)
		client := new(mocks.SirenServiceClient)
		client.On("GetSubscription", ctx, &sirenv1beta1.GetSubscriptionRequest{Id: subscription.Id}).
			Return(&sirenv1beta1.GetSubscriptionResponse{
				Subscription: subscription,
			}, nil)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
		result, err := service.FindSubscription(ctx, subscriptionID)
		assert.NoError(t, err)
		assert.Equal(t, subscription, result)
	})

	t.Run("should return not found error if optimus return NotFound code", func(t *testing.T) {
		grpcError := status.Error(codes.NotFound, "Not Found")

		shield := new(mocks.ShieldServiceClient)
		client := new(mocks.SirenServiceClient)
		client.On("GetSubscription", ctx, &sirenv1beta1.GetSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, grpcError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
		_, err := service.FindSubscription(ctx, subscriptionID)
		assert.ErrorIs(t, err, alert.ErrSubscriptionNotFound)
	})

	t.Run("should return if client return error", func(t *testing.T) {
		expectedError := status.Error(codes.Internal, "Internal")

		shield := new(mocks.ShieldServiceClient)
		client := new(mocks.SirenServiceClient)
		client.On("GetSubscription", ctx, &sirenv1beta1.GetSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
		_, err := service.FindSubscription(ctx, subscriptionID)
		assert.ErrorIs(t, err, expectedError)
	})
}

func TestSubscriptionServiceGetSubscriptions(t *testing.T) {
	ctx := context.TODO()

	t.Run("should return subscription on success", func(t *testing.T) {
		groupID := "19293012i31"
		resourceID := "sample-resource-id-or-urn"
		resourceType := "firehose"

		subscriptions := []*sirenv1beta1.Subscription{
			{
				Id:        1,
				Urn:       "sample-urn-1",
				Namespace: 1,
				Receivers: []*sirenv1beta1.ReceiverMetadata{
					{Id: 30},
				},
			},
			{
				Id:        2,
				Urn:       "sample-urn-2",
				Namespace: 2,
				Receivers: []*sirenv1beta1.ReceiverMetadata{
					{Id: 33},
				},
			},
		}

		shield := new(mocks.ShieldServiceClient)
		client := new(mocks.SirenServiceClient)
		client.On("ListSubscriptions", ctx, &sirenv1beta1.ListSubscriptionsRequest{Metadata: map[string]string{
			"group_id":      groupID,
			"resource_id":   resourceID,
			"resource_type": resourceType,
		}}).
			Return(&sirenv1beta1.ListSubscriptionsResponse{
				Subscriptions: subscriptions,
			}, nil)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
		result, err := service.GetSubscriptions(ctx, groupID, resourceID, resourceType)
		assert.NoError(t, err)
		assert.Equal(t, subscriptions, result)
	})

	t.Run("should return if client return error", func(t *testing.T) {
		expectedError := status.Error(codes.Internal, "Internal")

		shield := new(mocks.ShieldServiceClient)
		client := new(mocks.SirenServiceClient)
		client.On("ListSubscriptions", ctx, &sirenv1beta1.ListSubscriptionsRequest{Metadata: map[string]string{}}).
			Return(nil, expectedError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
		_, err := service.GetSubscriptions(ctx, "", "", "")
		assert.ErrorIs(t, err, expectedError)
	})
}

func TestSubscriptionServiceCreateSubscription(t *testing.T) {
	var (
		ctx       = context.TODO()
		groupID   = "8a7219cd-53c9-47f1-9387-5cac7abe4dcb"
		projectID = "5dab4194-9516-421a-aafe-72fd3d96ec56"
	)

	t.Run("should return error if siren namespace cannot be retrieved from project", func(t *testing.T) {
		tests := []struct {
			name     string
			metadata *structpb.Struct
		}{
			{
				name:     "empty metadata",
				metadata: nil,
			},
			{
				name: "empty metadata.siren_namespace",
				metadata: newStruct(t, map[string]interface{}{
					"siren_namespace": nil,
				}),
			},
			{
				name: "invalid format for metadata.siren_namespace",
				metadata: newStruct(t, map[string]interface{}{
					"siren_namespace": "wrong-format",
				}),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				form := alert.SubscriptionForm{
					ProjectID: projectID,
					GroupID:   groupID,
				}
				shieldProject := &shieldv1beta1.Project{
					Slug:     "test-project",
					Metadata: test.metadata,
				}

				shield := new(mocks.ShieldServiceClient)
				shield.On("GetProject", ctx, &shieldv1beta1.GetProjectRequest{Id: projectID}).
					Return(&shieldv1beta1.GetProjectResponse{Project: shieldProject}, nil)
				defer shield.AssertExpectations(t)
				client := new(mocks.SirenServiceClient)
				defer client.AssertExpectations(t)

				service := alert.NewSubscriptionService(client, shield)
				_, err := service.CreateSubscription(ctx, form)
				assert.ErrorIs(t, err, alert.ErrNoShieldSirenNamespace)
			})
		}
	})

	t.Run("should return error on failing to get siren's receiver", func(t *testing.T) {
		tests := []struct {
			name          string
			receivers     []*sirenv1beta1.Receiver
			expectedError error
		}{
			{
				name:          "nil receivers",
				receivers:     nil,
				expectedError: alert.ErrNoSirenReceiver,
			},
			{
				name:          "empty receivers",
				receivers:     []*sirenv1beta1.Receiver{},
				expectedError: alert.ErrNoSirenReceiver,
			},
			{
				name: "more than one receivers",
				receivers: []*sirenv1beta1.Receiver{
					{Id: 1},
					{Id: 2},
				},
			},
			{
				name: "receiver is not slack_channel type",
				receivers: []*sirenv1beta1.Receiver{
					{Id: 1, Type: "invalid-type"},
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				form := alert.SubscriptionForm{
					ProjectID:          projectID,
					GroupID:            groupID,
					ChannelCriticality: alert.ChannelCriticalityInfo,
				}
				shieldGroup := &shieldv1beta1.Group{
					Slug: "test-group",
				}
				shieldProject := &shieldv1beta1.Project{
					Slug: "test-project",
					Metadata: newStruct(t, map[string]interface{}{
						"siren_namespace": 5,
					}),
				}

				shield := new(mocks.ShieldServiceClient)
				shield.On("GetProject", ctx, &shieldv1beta1.GetProjectRequest{Id: projectID}).
					Return(&shieldv1beta1.GetProjectResponse{Project: shieldProject}, nil)
				shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: form.GroupID}).
					Return(&shieldv1beta1.GetGroupResponse{Group: shieldGroup}, nil)
				defer shield.AssertExpectations(t)
				siren := new(mocks.SirenServiceClient)
				siren.On("ListReceivers", ctx, &sirenv1beta1.ListReceiversRequest{
					Labels: map[string]string{
						"team":     shieldGroup.Slug,
						"severity": string(form.ChannelCriticality),
					},
				}).Return(&sirenv1beta1.ListReceiversResponse{
					Receivers: test.receivers,
				}, nil)
				defer siren.AssertExpectations(t)

				service := alert.NewSubscriptionService(siren, shield)
				_, err := service.CreateSubscription(ctx, form)
				if test.expectedError != nil {
					assert.ErrorIs(t, err, test.expectedError)
				} else {
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("should create subscription on success", func(t *testing.T) {
		receiverID := uint64(15)
		sirenNamespace := 5
		channelName := "test-alert-channel"

		// inputs
		form := alert.SubscriptionForm{
			UserID:             "john.doe@example.com",
			AlertSeverity:      alert.AlertSeverityCritical,
			ChannelCriticality: alert.ChannelCriticalityInfo,
			GroupID:            groupID,
			ProjectID:          projectID,
			ResourceType:       "firehose",
			ResourceID:         "test-job",
		}

		// conditions
		shieldGroup := &shieldv1beta1.Group{
			Slug: "test-group",
		}
		shieldProject := &shieldv1beta1.Project{
			Slug: "my-project-1",
			Metadata: newStruct(t, map[string]interface{}{
				"siren_namespace": sirenNamespace,
			}),
		}
		sirenReceivers := []*sirenv1beta1.Receiver{
			{Id: receiverID, Type: sirenReceiverPkg.TypeSlackChannel, Configurations: newStruct(t, map[string]interface{}{
				"channel_name": channelName,
			})},
		}

		// expectations
		expectedSirenPayload := &sirenv1beta1.CreateSubscriptionRequest{
			Urn: fmt.Sprintf(
				"%s:%s:%s:%s",
				shieldGroup.GetSlug(), form.AlertSeverity, form.ResourceType, form.ResourceID,
			),
			Namespace: uint64(sirenNamespace),
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{Id: receiverID},
			},
			Match: map[string]string{
				"severity":   string(alert.AlertSeverityCritical),
				"identifier": "test-job",
			},
			Metadata: newStruct(t, map[string]interface{}{
				"group_id":            form.GroupID,
				"group_slug":          shieldGroup.Slug,
				"resource_type":       form.ResourceType,
				"resource_id":         form.ResourceID,
				"project_id":          form.ProjectID,
				"project_slug":        shieldProject.Slug,
				"channel_criticality": string(form.ChannelCriticality),
				"channel_name":        channelName,
			}),
			CreatedBy: form.UserID,
		}

		shield := new(mocks.ShieldServiceClient)
		shield.On("GetProject", ctx, &shieldv1beta1.GetProjectRequest{Id: projectID}).
			Return(&shieldv1beta1.GetProjectResponse{Project: shieldProject}, nil)
		shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: form.GroupID}).
			Return(&shieldv1beta1.GetGroupResponse{
				Group: shieldGroup,
			}, nil)
		defer shield.AssertExpectations(t)
		siren := new(mocks.SirenServiceClient)
		siren.On("ListReceivers", ctx, &sirenv1beta1.ListReceiversRequest{
			Labels: map[string]string{
				"team":     shieldGroup.Slug,
				"severity": string(form.ChannelCriticality),
			},
		}).Return(&sirenv1beta1.ListReceiversResponse{Receivers: sirenReceivers}, nil)
		siren.
			On("CreateSubscription", ctx, expectedSirenPayload).
			Return(&sirenv1beta1.CreateSubscriptionResponse{Id: 5}, nil)
		defer siren.AssertExpectations(t)

		service := alert.NewSubscriptionService(siren, shield)
		subsID, err := service.CreateSubscription(ctx, form)
		assert.NoError(t, err)
		assert.Equal(t, 5, subsID)
	})
}

func TestSubscriptionServiceUpdateSubscription(t *testing.T) {
	var (
		ctx            = context.TODO()
		subscriptionID = 205
		groupID        = "8a7219cd-53c9-47f1-9387-5cac7abe4dcb"
		projectID      = "5dab4194-9516-421a-aafe-72fd3d96ec56"
	)

	t.Run("should return error if siren namespace cannot be retrieved from project", func(t *testing.T) {
		tests := []struct {
			name     string
			metadata *structpb.Struct
		}{
			{
				name:     "empty metadata",
				metadata: nil,
			},
			{
				name: "empty metadata.siren_namespace",
				metadata: newStruct(t, map[string]interface{}{
					"siren_namespace": nil,
				}),
			},
			{
				name: "invalid format for metadata.siren_namespace",
				metadata: newStruct(t, map[string]interface{}{
					"siren_namespace": "wrong-format",
				}),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				form := alert.SubscriptionForm{
					ProjectID: projectID,
					GroupID:   groupID,
				}
				shieldProject := &shieldv1beta1.Project{
					Slug:     "test-project",
					Metadata: test.metadata,
				}

				shield := new(mocks.ShieldServiceClient)
				shield.On("GetProject", ctx, &shieldv1beta1.GetProjectRequest{Id: projectID}).
					Return(&shieldv1beta1.GetProjectResponse{Project: shieldProject}, nil)
				defer shield.AssertExpectations(t)
				client := new(mocks.SirenServiceClient)
				defer client.AssertExpectations(t)

				service := alert.NewSubscriptionService(client, shield)
				err := service.UpdateSubscription(ctx, subscriptionID, form)
				assert.ErrorIs(t, err, alert.ErrNoShieldSirenNamespace)
			})
		}
	})

	t.Run("should return error on failing to get siren's receiver", func(t *testing.T) {
		tests := []struct {
			name          string
			receivers     []*sirenv1beta1.Receiver
			expectedError error
		}{
			{
				name:          "nil receivers",
				receivers:     nil,
				expectedError: alert.ErrNoSirenReceiver,
			},
			{
				name:          "empty receivers",
				receivers:     []*sirenv1beta1.Receiver{},
				expectedError: alert.ErrNoSirenReceiver,
			},
			{
				name: "more than one receivers",
				receivers: []*sirenv1beta1.Receiver{
					{Id: 1},
					{Id: 2},
				},
			},
			{
				name: "receiver is not slack_channel type",
				receivers: []*sirenv1beta1.Receiver{
					{Id: 1, Type: "invalid-type"},
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				form := alert.SubscriptionForm{
					ProjectID:          projectID,
					GroupID:            groupID,
					ChannelCriticality: alert.ChannelCriticalityInfo,
				}
				shieldGroup := &shieldv1beta1.Group{
					Slug: "test-group",
				}
				shieldProject := &shieldv1beta1.Project{
					Slug: "test-project",
					Metadata: newStruct(t, map[string]interface{}{
						"siren_namespace": 5,
					}),
				}

				shield := new(mocks.ShieldServiceClient)
				shield.On("GetProject", ctx, &shieldv1beta1.GetProjectRequest{Id: projectID}).
					Return(&shieldv1beta1.GetProjectResponse{Project: shieldProject}, nil)
				shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: form.GroupID}).
					Return(&shieldv1beta1.GetGroupResponse{Group: shieldGroup}, nil)
				defer shield.AssertExpectations(t)
				siren := new(mocks.SirenServiceClient)
				siren.On("ListReceivers", ctx, &sirenv1beta1.ListReceiversRequest{
					Labels: map[string]string{
						"team":     shieldGroup.Slug,
						"severity": string(form.ChannelCriticality),
					},
				}).Return(&sirenv1beta1.ListReceiversResponse{
					Receivers: test.receivers,
				}, nil)
				defer siren.AssertExpectations(t)

				service := alert.NewSubscriptionService(siren, shield)
				err := service.UpdateSubscription(ctx, subscriptionID, form)
				if test.expectedError != nil {
					assert.ErrorIs(t, err, test.expectedError)
				} else {
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("should update subscription on success", func(t *testing.T) {
		receiverID := uint64(17)
		sirenNamespace := 5
		channelName := "test-channel-update"

		// inputs
		form := alert.SubscriptionForm{
			UserID:             "john.doe@example.com",
			AlertSeverity:      alert.AlertSeverityCritical,
			ChannelCriticality: alert.ChannelCriticalityInfo,
			GroupID:            groupID,
			ProjectID:          projectID,
			ResourceType:       "firehose",
			ResourceID:         "test-job",
		}

		// conditions
		shieldGroup := &shieldv1beta1.Group{
			Slug: "test-group",
		}
		shieldProject := &shieldv1beta1.Project{
			Slug: "my-project-1",
			Metadata: newStruct(t, map[string]interface{}{
				"siren_namespace": sirenNamespace,
			}),
		}
		sirenReceivers := []*sirenv1beta1.Receiver{
			{Id: receiverID, Type: sirenReceiverPkg.TypeSlackChannel, Configurations: newStruct(t, map[string]interface{}{
				"channel_name": channelName,
			})},
		}

		// expecations
		expectedSirenPayload := &sirenv1beta1.UpdateSubscriptionRequest{
			Id: uint64(subscriptionID),
			Urn: fmt.Sprintf(
				"%s:%s:%s:%s",
				shieldGroup.GetSlug(), form.AlertSeverity, form.ResourceType, form.ResourceID,
			),
			Namespace: uint64(sirenNamespace),
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{Id: receiverID},
			},
			Match: map[string]string{
				"severity":   string(alert.AlertSeverityCritical),
				"identifier": "test-job",
			},
			Metadata: newStruct(t, map[string]interface{}{
				"group_id":            form.GroupID,
				"group_slug":          shieldGroup.Slug,
				"resource_type":       form.ResourceType,
				"resource_id":         form.ResourceID,
				"project_id":          form.ProjectID,
				"project_slug":        shieldProject.Slug,
				"channel_criticality": string(form.ChannelCriticality),
				"channel_name":        channelName,
			}),
			UpdatedBy: form.UserID,
		}

		shield := new(mocks.ShieldServiceClient)
		shield.On("GetProject", ctx, &shieldv1beta1.GetProjectRequest{Id: projectID}).
			Return(&shieldv1beta1.GetProjectResponse{Project: shieldProject}, nil)
		shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: form.GroupID}).
			Return(&shieldv1beta1.GetGroupResponse{
				Group: shieldGroup,
			}, nil)
		defer shield.AssertExpectations(t)
		siren := new(mocks.SirenServiceClient)
		siren.On("ListReceivers", ctx, &sirenv1beta1.ListReceiversRequest{
			Labels: map[string]string{
				"team":     shieldGroup.Slug,
				"severity": string(form.ChannelCriticality),
			},
		}).Return(&sirenv1beta1.ListReceiversResponse{Receivers: sirenReceivers}, nil)
		siren.
			On("UpdateSubscription", ctx, expectedSirenPayload).
			Return(&sirenv1beta1.UpdateSubscriptionResponse{}, nil)
		defer siren.AssertExpectations(t)

		service := alert.NewSubscriptionService(siren, shield)
		err := service.UpdateSubscription(ctx, subscriptionID, form)
		assert.NoError(t, err)
	})
}

func TestSubscriptionServiceDeleteSubscription(t *testing.T) {
	ctx := context.TODO()
	subscriptionID := 203

	t.Run("should not return error success", func(t *testing.T) {
		shield := new(mocks.ShieldServiceClient)
		client := new(mocks.SirenServiceClient)
		client.On("DeleteSubscription", ctx, &sirenv1beta1.DeleteSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, nil)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
		err := service.DeleteSubscription(ctx, subscriptionID)
		assert.NoError(t, err)
	})

	t.Run("should return not found error if optimus return NotFound code", func(t *testing.T) {
		expectedError := status.Error(codes.NotFound, "Not Found")

		shield := new(mocks.ShieldServiceClient)
		client := new(mocks.SirenServiceClient)
		client.On("DeleteSubscription", ctx, &sirenv1beta1.DeleteSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
		err := service.DeleteSubscription(ctx, subscriptionID)
		assert.ErrorIs(t, err, alert.ErrSubscriptionNotFound)
	})

	t.Run("should return if client return error", func(t *testing.T) {
		expectedError := status.Error(codes.Internal, "Internal")

		shield := new(mocks.ShieldServiceClient)
		client := new(mocks.SirenServiceClient)
		client.On("DeleteSubscription", ctx, &sirenv1beta1.DeleteSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
		err := service.DeleteSubscription(ctx, subscriptionID)
		assert.ErrorIs(t, err, expectedError)
	})
}

func TestSubscriptionServiceGetAlertChannels(t *testing.T) {
	ctx := context.TODO()
	groupID := "deafcced-845c-4089-89f0-06621486cb0a"

	t.Run("should not return error if group is not found", func(t *testing.T) {
		notFoundError := status.Error(codes.NotFound, "Not Found")

		shield := new(mocks.ShieldServiceClient)
		shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: groupID}).
			Return(nil, notFoundError)
		defer shield.AssertExpectations(t)
		siren := new(mocks.SirenServiceClient)

		service := alert.NewSubscriptionService(siren, shield)
		_, err := service.GetAlertChannels(ctx, groupID)
		assert.ErrorIs(t, err, alert.ErrNoShieldGroup)
	})

	t.Run("should return alert channels", func(t *testing.T) {
		groupSlug := "test-group-30"
		shieldGroup := &shieldv1beta1.Group{
			Slug: groupSlug,
		}
		sirenReceivers := []*sirenv1beta1.Receiver{
			{
				Id:   54,
				Name: "test-receiver-info-1",
				Labels: map[string]string{
					"severity": string(alert.AlertSeverityInfo),
				},
				Configurations: newStruct(t, map[string]interface{}{
					"channel_name": "test-channel-info-1",
				}),
			},
			{
				Id:   55,
				Name: "test-receiver-critical-1",
				Labels: map[string]string{
					"severity": string(alert.AlertSeverityCritical),
				},
				Configurations: newStruct(t, map[string]interface{}{
					"channel_name": "test-channel-critical-1",
				}),
			},
			{
				Id:   56,
				Name: "test-receiver-warning-1",
				Labels: map[string]string{
					"severity": string(alert.AlertSeverityWarning),
				},
				Configurations: newStruct(t, map[string]interface{}{
					"channel_name": "test-channel-warning-1",
				}),
			},
		}

		expected := []models.AlertChannel{
			{
				ReceiverID:         fmt.Sprint(sirenReceivers[0].Id),
				ReceiverName:       sirenReceivers[0].Name,
				ChannelCriticality: models.NewChannelCriticality(models.ChannelCriticalityINFO),
				ChannelName:        "test-channel-info-1",
			},
			{
				ReceiverID:         fmt.Sprint(sirenReceivers[1].Id),
				ReceiverName:       sirenReceivers[1].Name,
				ChannelCriticality: models.NewChannelCriticality(models.ChannelCriticalityCRITICAL),
				ChannelName:        "test-channel-critical-1",
			},
			{
				ReceiverID:         fmt.Sprint(sirenReceivers[2].Id),
				ReceiverName:       sirenReceivers[2].Name,
				ChannelCriticality: models.NewChannelCriticality(models.ChannelCriticalityWARNING),
				ChannelName:        "test-channel-warning-1",
			},
		}

		shield := new(mocks.ShieldServiceClient)
		shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: groupID}).
			Return(&shieldv1beta1.GetGroupResponse{
				Group: shieldGroup,
			}, nil)
		defer shield.AssertExpectations(t)
		siren := new(mocks.SirenServiceClient)
		siren.On("ListReceivers", ctx, &sirenv1beta1.ListReceiversRequest{
			Labels: map[string]string{
				"team": groupSlug,
			},
		}).
			Return(&sirenv1beta1.ListReceiversResponse{
				Receivers: sirenReceivers,
			}, nil)
		defer siren.AssertExpectations(t)

		service := alert.NewSubscriptionService(siren, shield)
		actual, err := service.GetAlertChannels(ctx, groupID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestSubscriptionServiceSetAlertChannels(t *testing.T) {
	var (
		ctx         = context.TODO()
		groupID     = "8a7219cd-53c9-47f1-9387-5cac7abe4dcb"
		orgID       = "ea597d5c-1280-473b-ad28-7551c1336fe0"
		shieldGroup = &shieldv1beta1.Group{
			Slug:  "test-group-slug-12",
			OrgId: orgID,
		}
		shieldOrg = &shieldv1beta1.Organization{
			Slug: "test-org-slug-21",
		}
	)

	t.Run("should return error if group cannot be found", func(t *testing.T) {
		expectedErr := status.Error(codes.NotFound, "Not Found")

		shield := new(mocks.ShieldServiceClient)
		shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: groupID}).
			Return(nil, expectedErr)
		defer shield.AssertExpectations(t)
		siren := new(mocks.SirenServiceClient)

		service := alert.NewSubscriptionService(siren, shield)
		_, err := service.SetAlertChannels(ctx, groupID, []alert.AlertChannelForm{})
		assert.ErrorIs(t, err, alert.ErrNoShieldGroup)
	})

	t.Run("should return error if org cannot be found", func(t *testing.T) {
		expectedErr := status.Error(codes.NotFound, "Not Found")

		shield := new(mocks.ShieldServiceClient)
		shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: groupID}).
			Return(&shieldv1beta1.GetGroupResponse{Group: shieldGroup}, nil)
		shield.On("GetOrganization", ctx, &shieldv1beta1.GetOrganizationRequest{Id: shieldGroup.OrgId}).
			Return(nil, expectedErr)
		defer shield.AssertExpectations(t)
		siren := new(mocks.SirenServiceClient)

		service := alert.NewSubscriptionService(siren, shield)
		_, err := service.SetAlertChannels(ctx, groupID, []alert.AlertChannelForm{})
		assert.Error(t, err)
	})

	t.Run("should return error if parent slack receiver cannot be found or invalid", func(t *testing.T) {
		tests := []struct {
			description string
			receivers   []*sirenv1beta1.Receiver
		}{
			{
				description: "empty receivers",
				receivers:   []*sirenv1beta1.Receiver{},
			},
			{
				description: "wrong type and entity",
				receivers: []*sirenv1beta1.Receiver{
					{
						Type: sirenReceiverPkg.TypeSlackChannel,
						Labels: map[string]string{
							"entity":          "test-org,sample-org",
							"is_parent_slack": "true",
						},
					},
				},
			},
			{
				description: "wrong entity",
				receivers: []*sirenv1beta1.Receiver{
					{
						Type: sirenReceiverPkg.TypeSlack,
						Labels: map[string]string{
							"entity":          "test-org,sample-org",
							"is_parent_slack": "true",
						},
					},
				},
			},
			{
				description: "wrong type",
				receivers: []*sirenv1beta1.Receiver{
					{
						Type: sirenReceiverPkg.TypeSlackChannel,
						Labels: map[string]string{
							"entity":          "test-org,test-org-slug-21,sample-org",
							"is_parent_slack": "true",
						},
					},
				},
			},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				shield := new(mocks.ShieldServiceClient)
				shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: groupID}).
					Return(&shieldv1beta1.GetGroupResponse{Group: shieldGroup}, nil)
				shield.On("GetOrganization", ctx, &shieldv1beta1.GetOrganizationRequest{Id: shieldGroup.OrgId}).
					Return(&shieldv1beta1.GetOrganizationResponse{Organization: shieldOrg}, nil)
				defer shield.AssertExpectations(t)
				siren := new(mocks.SirenServiceClient)
				siren.On("ListReceivers", ctx, &sirenv1beta1.ListReceiversRequest{
					Labels: map[string]string{
						"is_parent_slack": "true",
					},
				}).Return(&sirenv1beta1.ListReceiversResponse{
					Receivers: test.receivers,
				}, nil)
				defer siren.AssertExpectations(t)

				service := alert.NewSubscriptionService(siren, shield)
				_, err := service.SetAlertChannels(ctx, groupID, []alert.AlertChannelForm{})
				assert.ErrorIs(t, err, alert.ErrNoSirenParentSlackReceiver)
			})
		}
	})

	t.Run("should create/update the correct receiver(s) and return the correct results", func(t *testing.T) {
		parentReceiverID := uint64(20)

		tests := []struct {
			description       string
			existingReceivers []*sirenv1beta1.Receiver
			forms             []alert.AlertChannelForm
			setupSiren        func(*mocks.SirenServiceClient)
			expected          []models.AlertChannel
		}{
			{
				description: "single new slack channel",
				forms: []alert.AlertChannelForm{
					{
						ChannelCriticality: alert.ChannelCriticalityInfo,
						ChannelName:        "test-channel-932",
						ChannelType:        "slack_channel",
					},
				},
				setupSiren: func(siren *mocks.SirenServiceClient) {
					siren.On("CreateReceiver", ctx, &sirenv1beta1.CreateReceiverRequest{
						Name:     fmt.Sprintf("%s-%s-slack_channel-info", shieldOrg.Slug, shieldGroup.Slug),
						Type:     "slack_channel",
						ParentId: parentReceiverID,
						Labels: map[string]string{
							"team":     shieldGroup.Slug,
							"org":      shieldOrg.Slug,
							"severity": "INFO",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"channel_name": "test-channel-932",
						}),
					}).Return(&sirenv1beta1.CreateReceiverResponse{Id: 30}, nil).Once()
				},
				expected: []models.AlertChannel{
					{
						ChannelCriticality: models.NewChannelCriticality(models.ChannelCriticality("INFO")),
						ChannelName:        "test-channel-932",
						ChannelType:        models.NewAlertChannelType(models.AlertChannelType("slack_channel")),
						ReceiverID:         "30",
						ReceiverName:       fmt.Sprintf("%s-%s-slack_channel-info", shieldOrg.Slug, shieldGroup.Slug),
					},
				},
			},
			{
				description: "single update slack channel",
				existingReceivers: []*sirenv1beta1.Receiver{
					{
						Id:       15,
						Name:     "old-name-1293",
						Type:     "slack_channel",
						ParentId: parentReceiverID,
						Labels: map[string]string{
							"org":      "test",
							"team":     "sample-team",
							"severity": "CRITICAL",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"channel_name": "old-slack-channel-30",
						}),
					},
				},
				forms: []alert.AlertChannelForm{
					{
						ChannelCriticality: alert.ChannelCriticalityCritical,
						ChannelName:        "new-channel-932",
						ChannelType:        "slack_channel",
					},
				},
				setupSiren: func(siren *mocks.SirenServiceClient) {
					siren.On("UpdateReceiver", ctx, &sirenv1beta1.UpdateReceiverRequest{
						Id:       15,
						Name:     "old-name-1293",
						ParentId: parentReceiverID,
						Labels: map[string]string{
							"org":      "test",
							"team":     "sample-team",
							"severity": "CRITICAL",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"channel_name": "new-channel-932",
						}),
					}).Return(&sirenv1beta1.UpdateReceiverResponse{Id: 15}, nil).Once()
				},
				expected: []models.AlertChannel{
					{
						ChannelCriticality: models.NewChannelCriticality(models.ChannelCriticality("CRITICAL")),
						ChannelName:        "new-channel-932",
						ChannelType:        models.NewAlertChannelType(models.AlertChannelType("slack_channel")),
						ReceiverID:         "15",
						ReceiverName:       "old-name-1293",
					},
				},
			},
			{
				description: "single new pagerduty channel",
				forms: []alert.AlertChannelForm{
					{
						ChannelCriticality:  alert.ChannelCriticalityCritical,
						PagerdutyServiceKey: "sample-service-key-192903",
						ChannelType:         "pagerduty",
					},
				},
				setupSiren: func(siren *mocks.SirenServiceClient) {
					siren.On("CreateReceiver", ctx, &sirenv1beta1.CreateReceiverRequest{
						Name: fmt.Sprintf("%s-%s-pagerduty-critical", shieldOrg.Slug, shieldGroup.Slug),
						Type: "pagerduty",
						Labels: map[string]string{
							"team":     shieldGroup.Slug,
							"org":      shieldOrg.Slug,
							"severity": "CRITICAL",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"service_key": "sample-service-key-192903",
						}),
					}).Return(&sirenv1beta1.CreateReceiverResponse{Id: 82}, nil).Once()
				},
				expected: []models.AlertChannel{
					{
						ChannelCriticality:  models.NewChannelCriticality(models.ChannelCriticality("CRITICAL")),
						PagerdutyServiceKey: "sample-service-key-192903",
						ChannelType:         models.NewAlertChannelType(models.AlertChannelType("pagerduty")),
						ReceiverID:          "82",
						ReceiverName:        fmt.Sprintf("%s-%s-pagerduty-critical", shieldOrg.Slug, shieldGroup.Slug),
					},
				},
			},
			{
				description: "single update pagerduty channel",
				existingReceivers: []*sirenv1beta1.Receiver{
					{
						Id:   75,
						Name: "old-name-9953",
						Type: "pagerduty",
						Labels: map[string]string{
							"org":      "test",
							"team":     "sample-team",
							"severity": "INFO",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"service_key": "old-service-key-5i31",
						}),
					},
				},
				forms: []alert.AlertChannelForm{
					{
						ChannelCriticality:  alert.ChannelCriticalityInfo,
						ChannelType:         "pagerduty",
						PagerdutyServiceKey: "new-service-key-98293",
					},
				},
				setupSiren: func(siren *mocks.SirenServiceClient) {
					siren.On("UpdateReceiver", ctx, &sirenv1beta1.UpdateReceiverRequest{
						Id:   75,
						Name: "old-name-9953",
						Labels: map[string]string{
							"org":      "test",
							"team":     "sample-team",
							"severity": "INFO",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"service_key": "new-service-key-98293",
						}),
					}).Return(&sirenv1beta1.UpdateReceiverResponse{Id: 75}, nil).Once()
				},
				expected: []models.AlertChannel{
					{
						ChannelCriticality:  models.NewChannelCriticality(models.ChannelCriticality("INFO")),
						PagerdutyServiceKey: "new-service-key-98293",
						ChannelType:         models.NewAlertChannelType(models.AlertChannelType("pagerduty")),
						ReceiverID:          "75",
						ReceiverName:        "old-name-9953",
					},
				},
			},
			{
				description: "multiple new pagerduty channel",
				forms: []alert.AlertChannelForm{
					{
						ChannelCriticality:  alert.ChannelCriticalityInfo,
						PagerdutyServiceKey: "sample-service-key-123",
						ChannelType:         "pagerduty",
					},
					{
						ChannelCriticality:  alert.ChannelCriticalityWarning,
						PagerdutyServiceKey: "sample-service-key-321",
						ChannelType:         "pagerduty",
					},
				},
				setupSiren: func(siren *mocks.SirenServiceClient) {
					siren.On("CreateReceiver", ctx, &sirenv1beta1.CreateReceiverRequest{
						Name: fmt.Sprintf("%s-%s-pagerduty-info", shieldOrg.Slug, shieldGroup.Slug),
						Type: "pagerduty",
						Labels: map[string]string{
							"team":     shieldGroup.Slug,
							"org":      shieldOrg.Slug,
							"severity": "INFO",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"service_key": "sample-service-key-123",
						}),
					}).Return(&sirenv1beta1.CreateReceiverResponse{Id: 11}, nil).Once()
					siren.On("CreateReceiver", ctx, &sirenv1beta1.CreateReceiverRequest{
						Name: fmt.Sprintf("%s-%s-pagerduty-warning", shieldOrg.Slug, shieldGroup.Slug),
						Type: "pagerduty",
						Labels: map[string]string{
							"team":     shieldGroup.Slug,
							"org":      shieldOrg.Slug,
							"severity": "WARNING",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"service_key": "sample-service-key-321",
						}),
					}).Return(&sirenv1beta1.CreateReceiverResponse{Id: 98}, nil).Once()
				},
				expected: []models.AlertChannel{
					{
						ChannelCriticality:  models.NewChannelCriticality(models.ChannelCriticality("INFO")),
						PagerdutyServiceKey: "sample-service-key-123",
						ChannelType:         models.NewAlertChannelType(models.AlertChannelType("pagerduty")),
						ReceiverID:          "11",
						ReceiverName:        fmt.Sprintf("%s-%s-pagerduty-info", shieldOrg.Slug, shieldGroup.Slug),
					},
					{
						ChannelCriticality:  models.NewChannelCriticality(models.ChannelCriticality("WARNING")),
						PagerdutyServiceKey: "sample-service-key-321",
						ChannelType:         models.NewAlertChannelType(models.AlertChannelType("pagerduty")),
						ReceiverID:          "98",
						ReceiverName:        fmt.Sprintf("%s-%s-pagerduty-warning", shieldOrg.Slug, shieldGroup.Slug),
					},
				},
			},
			{
				description: "multiple new slack channel",
				forms: []alert.AlertChannelForm{
					{
						ChannelCriticality: alert.ChannelCriticalityCritical,
						ChannelName:        "test-channel-123",
						ChannelType:        "slack_channel",
					},
					{
						ChannelCriticality: alert.ChannelCriticalityWarning,
						ChannelName:        "test-channel-321",
						ChannelType:        "slack_channel",
					},
				},
				setupSiren: func(siren *mocks.SirenServiceClient) {
					siren.On("CreateReceiver", ctx, &sirenv1beta1.CreateReceiverRequest{
						Name:     fmt.Sprintf("%s-%s-slack_channel-critical", shieldOrg.Slug, shieldGroup.Slug),
						Type:     "slack_channel",
						ParentId: parentReceiverID,
						Labels: map[string]string{
							"team":     shieldGroup.Slug,
							"org":      shieldOrg.Slug,
							"severity": "CRITICAL",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"channel_name": "test-channel-123",
						}),
					}).Return(&sirenv1beta1.CreateReceiverResponse{Id: 76}, nil).Once()
					siren.On("CreateReceiver", ctx, &sirenv1beta1.CreateReceiverRequest{
						Name:     fmt.Sprintf("%s-%s-slack_channel-warning", shieldOrg.Slug, shieldGroup.Slug),
						Type:     "slack_channel",
						ParentId: parentReceiverID,
						Labels: map[string]string{
							"team":     shieldGroup.Slug,
							"org":      shieldOrg.Slug,
							"severity": "WARNING",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"channel_name": "test-channel-321",
						}),
					}).Return(&sirenv1beta1.CreateReceiverResponse{Id: 10}, nil).Once()
				},
				expected: []models.AlertChannel{
					{
						ChannelCriticality: models.NewChannelCriticality(models.ChannelCriticality("CRITICAL")),
						ChannelName:        "test-channel-123",
						ChannelType:        models.NewAlertChannelType(models.AlertChannelType("slack_channel")),
						ReceiverID:         "76",
						ReceiverName:       fmt.Sprintf("%s-%s-slack_channel-critical", shieldOrg.Slug, shieldGroup.Slug),
					},
					{
						ChannelCriticality: models.NewChannelCriticality(models.ChannelCriticality("WARNING")),
						ChannelName:        "test-channel-321",
						ChannelType:        models.NewAlertChannelType(models.AlertChannelType("slack_channel")),
						ReceiverID:         "10",
						ReceiverName:       fmt.Sprintf("%s-%s-slack_channel-warning", shieldOrg.Slug, shieldGroup.Slug),
					},
				},
			},
			{
				description: "multiple channel type and mixed create/update",
				existingReceivers: []*sirenv1beta1.Receiver{
					{
						Id:       20,
						Name:     "old-name-1231",
						Type:     "slack_channel",
						ParentId: parentReceiverID,
						Labels: map[string]string{
							"org":      "test-12",
							"team":     "sample-team-120",
							"severity": "INFO",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"channel_name": "old-slack-channel-023",
						}),
					},
					{
						Id:   123,
						Name: "old-name-10203",
						Type: "pagerduty",
						Labels: map[string]string{
							"org":      "test-12",
							"team":     "sample-team-120",
							"severity": "INFO",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"service_key": "old-service-key-111",
						}),
					},
					{
						Id:   421,
						Name: "old-name-9201",
						Type: "pagerduty",
						Labels: map[string]string{
							"org":      "test-83",
							"team":     "sample-team-2481",
							"severity": "WARNING",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"service_key": "old-service-key-3921",
						}),
					},
				},
				forms: []alert.AlertChannelForm{
					{
						ChannelCriticality: alert.ChannelCriticalityWarning,
						ChannelName:        "test-channel-3942",
						ChannelType:        "slack_channel",
					},
					{
						ChannelCriticality:  alert.ChannelCriticalityWarning,
						PagerdutyServiceKey: "test-service-key-83891",
						ChannelType:         "pagerduty",
					},
					{
						ChannelCriticality: alert.ChannelCriticalityInfo,
						ChannelName:        "test-channel-1929",
						ChannelType:        "slack_channel",
					},
					{
						ChannelCriticality:  alert.ChannelCriticalityCritical,
						PagerdutyServiceKey: "test-service-key-582",
						ChannelType:         "pagerduty",
					},
				},
				setupSiren: func(siren *mocks.SirenServiceClient) {
					siren.On("CreateReceiver", ctx, &sirenv1beta1.CreateReceiverRequest{
						Name:     fmt.Sprintf("%s-%s-slack_channel-warning", shieldOrg.Slug, shieldGroup.Slug),
						Type:     "slack_channel",
						ParentId: parentReceiverID,
						Labels: map[string]string{
							"team":     shieldGroup.Slug,
							"org":      shieldOrg.Slug,
							"severity": "WARNING",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"channel_name": "test-channel-3942",
						}),
					}).Return(&sirenv1beta1.CreateReceiverResponse{Id: 84}, nil).Once()

					siren.On("UpdateReceiver", ctx, &sirenv1beta1.UpdateReceiverRequest{
						Id:   421,
						Name: "old-name-9201",
						Labels: map[string]string{
							"org":      "test-83",
							"team":     "sample-team-2481",
							"severity": "WARNING",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"service_key": "test-service-key-83891",
						}),
					}).Return(&sirenv1beta1.UpdateReceiverResponse{Id: 421}, nil).Once()

					siren.On("UpdateReceiver", ctx, &sirenv1beta1.UpdateReceiverRequest{
						Id:       20,
						Name:     "old-name-1231",
						ParentId: parentReceiverID,
						Labels: map[string]string{
							"org":      "test-12",
							"team":     "sample-team-120",
							"severity": "INFO",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"channel_name": "test-channel-1929",
						}),
					}).Return(&sirenv1beta1.UpdateReceiverResponse{Id: 20}, nil).Once()

					siren.On("CreateReceiver", ctx, &sirenv1beta1.CreateReceiverRequest{
						Name: fmt.Sprintf("%s-%s-pagerduty-critical", shieldOrg.Slug, shieldGroup.Slug),
						Type: "pagerduty",
						Labels: map[string]string{
							"team":     shieldGroup.Slug,
							"org":      shieldOrg.Slug,
							"severity": "CRITICAL",
						},
						Configurations: newStruct(t, map[string]interface{}{
							"service_key": "test-service-key-582",
						}),
					}).Return(&sirenv1beta1.CreateReceiverResponse{Id: 129}, nil).Once()
				},
				expected: []models.AlertChannel{
					{
						ChannelCriticality: models.NewChannelCriticality(models.ChannelCriticality("WARNING")),
						ChannelName:        "test-channel-3942",
						ChannelType:        models.NewAlertChannelType(models.AlertChannelType("slack_channel")),
						ReceiverID:         "84",
						ReceiverName:       fmt.Sprintf("%s-%s-slack_channel-warning", shieldOrg.Slug, shieldGroup.Slug),
					},
					{
						ChannelCriticality:  models.NewChannelCriticality(models.ChannelCriticality("WARNING")),
						PagerdutyServiceKey: "test-service-key-83891",
						ChannelType:         models.NewAlertChannelType(models.AlertChannelType("pagerduty")),
						ReceiverID:          "421",
						ReceiverName:        "old-name-9201",
					},
					{
						ChannelCriticality: models.NewChannelCriticality(models.ChannelCriticality("INFO")),
						ChannelName:        "test-channel-1929",
						ChannelType:        models.NewAlertChannelType(models.AlertChannelType("slack_channel")),
						ReceiverID:         "20",
						ReceiverName:       "old-name-1231",
					},
					{
						ChannelCriticality:  models.NewChannelCriticality(models.ChannelCriticality("CRITICAL")),
						PagerdutyServiceKey: "test-service-key-582",
						ChannelType:         models.NewAlertChannelType(models.AlertChannelType("pagerduty")),
						ReceiverID:          "129",
						ReceiverName:        fmt.Sprintf("%s-%s-pagerduty-critical", shieldOrg.Slug, shieldGroup.Slug),
					},
				},
			},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				parentReceiver := &sirenv1beta1.Receiver{
					Id:   parentReceiverID,
					Type: sirenReceiverPkg.TypeSlack,
					Labels: map[string]string{
						"entity":          "test-org,test-org-slug-21,sample-org",
						"is_parent_slack": "true",
					},
				}

				shield := new(mocks.ShieldServiceClient)
				shield.On("GetGroup", ctx, &shieldv1beta1.GetGroupRequest{Id: groupID}).
					Return(&shieldv1beta1.GetGroupResponse{Group: shieldGroup}, nil)
				shield.On("GetOrganization", ctx, &shieldv1beta1.GetOrganizationRequest{Id: shieldGroup.OrgId}).
					Return(&shieldv1beta1.GetOrganizationResponse{Organization: shieldOrg}, nil)
				defer shield.AssertExpectations(t)
				siren := new(mocks.SirenServiceClient)
				siren.On("ListReceivers", ctx, &sirenv1beta1.ListReceiversRequest{
					Labels: map[string]string{
						"is_parent_slack": "true",
					},
				}).Return(&sirenv1beta1.ListReceiversResponse{
					Receivers: []*sirenv1beta1.Receiver{parentReceiver},
				}, nil).Once()
				siren.On("ListReceivers", ctx, &sirenv1beta1.ListReceiversRequest{
					Labels: map[string]string{
						"team": shieldGroup.Slug,
					},
				}).
					Return(&sirenv1beta1.ListReceiversResponse{Receivers: test.existingReceivers}, nil).
					Once()

				test.setupSiren(siren)
				defer siren.AssertExpectations(t)

				service := alert.NewSubscriptionService(siren, shield)
				results, err := service.SetAlertChannels(ctx, groupID, test.forms)
				require.NoError(t, err)
				assert.Equal(t, test.expected, results)
			})
		}
	})
}

func newStruct(t *testing.T, d map[string]interface{}) *structpb.Struct {
	t.Helper()

	strct, err := structpb.NewStruct(d)
	require.NoError(t, err)
	return strct
}
