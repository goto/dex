package alert_test

import (
	"context"
	"fmt"
	"testing"

	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	sirenv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/siren/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

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

	t.Run("should return error on invalid shield metadata", func(t *testing.T) {
		tests := []struct {
			name        string
			criticality alert.ChannelCriticality
			metadata    *structpb.Struct
		}{
			{
				name:        "empty metadata",
				criticality: alert.ChannelCriticalityWarning,
				metadata:    nil,
			},
			{
				name:        "empty metadata.alerting",
				criticality: alert.ChannelCriticalityWarning,
				metadata: newStruct(t, map[string]interface{}{
					"alerting": map[string]interface{}{},
				}),
			},
			{
				name:        "could not found metadata.alerting[criticality].slack",
				criticality: alert.ChannelCriticalityWarning,
				metadata: newStruct(t, map[string]interface{}{
					"alerting": map[string]interface{}{
						string(alert.ChannelCriticalityInfo): map[string]interface{}{},
					},
				}),
			},
			{
				name:        "empty metadata.alerting[criticality].slack",
				criticality: alert.ChannelCriticalityWarning,
				metadata: newStruct(t, map[string]interface{}{
					"alerting": map[string]interface{}{
						string(alert.ChannelCriticalityWarning): map[string]interface{}{
							"slack": map[string]interface{}{},
						},
					},
				}),
			},
			{
				name:        "empty metadata.alerting[criticality].slack.channel",
				criticality: alert.ChannelCriticalityWarning,
				metadata: newStruct(t, map[string]interface{}{
					"alerting": map[string]interface{}{
						string(alert.ChannelCriticalityWarning): map[string]interface{}{
							"slack": map[string]interface{}{
								"channel": "",
							},
						},
					},
				}),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				form := alert.SubscriptionForm{
					ProjectID: projectID,
					GroupID:   groupID,
				}
				shieldGroup := &shieldv1beta1.Group{
					Slug:     "test-group",
					Metadata: test.metadata,
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
				client := new(mocks.SirenServiceClient)
				defer client.AssertExpectations(t)

				service := alert.NewSubscriptionService(client, shield)
				_, err := service.CreateSubscription(ctx, form)
				assert.ErrorIs(t, err, alert.ErrNoShieldSlackChannel)
			})
		}
	})

	t.Run("should create subscription on success", func(t *testing.T) {
		channelName := "test-slack-channel"
		sirenNamespace := 5
		form := alert.SubscriptionForm{
			UserID:             "john.doe@example.com",
			AlertSeverity:      alert.AlertSeverityCritical,
			ChannelCriticality: alert.ChannelCriticalityInfo,
			GroupID:            groupID,
			ProjectID:          projectID,
			ResourceType:       "firehose",
			ResourceID:         "test-job",
		}
		shieldGroup := &shieldv1beta1.Group{
			Slug: "test-group",
			Metadata: newStruct(t, map[string]interface{}{
				"alerting": map[string]interface{}{
					string(form.ChannelCriticality): map[string]interface{}{
						"slack": map[string]interface{}{
							"channel": channelName,
						},
					},
				},
			}),
		}
		shieldProject := &shieldv1beta1.Project{
			Slug: "my-project-1",
			Metadata: newStruct(t, map[string]interface{}{
				"siren_namespace": sirenNamespace,
			}),
		}
		expectedSirenPayload := &sirenv1beta1.CreateSubscriptionRequest{
			Urn: fmt.Sprintf(
				"%s:%s:%s:%s",
				form.GroupID, form.AlertSeverity, form.ResourceType, form.ResourceID,
			),
			Namespace: uint64(sirenNamespace),
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{
					Id: 1,
					Configuration: newStruct(t, map[string]interface{}{
						"channel_name": channelName,
					}),
				},
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
				"channel_criticality": string(form.ChannelCriticality),
				"project_id":          form.ProjectID,
				"project_slug":        shieldProject.Slug,
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
		client := new(mocks.SirenServiceClient)
		client.
			On("CreateSubscription", ctx, expectedSirenPayload).
			Return(&sirenv1beta1.CreateSubscriptionResponse{Id: 5}, nil)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
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

	t.Run("should return error on invalid shield metadata", func(t *testing.T) {
		tests := []struct {
			name        string
			criticality alert.ChannelCriticality
			metadata    *structpb.Struct
		}{
			{
				name:        "empty metadata",
				criticality: alert.ChannelCriticalityWarning,
				metadata:    nil,
			},
			{
				name:        "empty metadata.alerting",
				criticality: alert.ChannelCriticalityWarning,
				metadata: newStruct(t, map[string]interface{}{
					"alerting": map[string]interface{}{},
				}),
			},
			{
				name:        "could not found metadata.alerting[criticality].slack",
				criticality: alert.ChannelCriticalityWarning,
				metadata: newStruct(t, map[string]interface{}{
					"alerting": map[string]interface{}{
						string(alert.ChannelCriticalityInfo): map[string]interface{}{},
					},
				}),
			},
			{
				name:        "empty metadata.alerting[criticality].slack",
				criticality: alert.ChannelCriticalityWarning,
				metadata: newStruct(t, map[string]interface{}{
					"alerting": map[string]interface{}{
						string(alert.ChannelCriticalityWarning): map[string]interface{}{
							"slack": map[string]interface{}{},
						},
					},
				}),
			},
			{
				name:        "empty metadata.alerting[criticality].slack.channel",
				criticality: alert.ChannelCriticalityWarning,
				metadata: newStruct(t, map[string]interface{}{
					"alerting": map[string]interface{}{
						string(alert.ChannelCriticalityWarning): map[string]interface{}{
							"slack": map[string]interface{}{
								"channel": "",
							},
						},
					},
				}),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				form := alert.SubscriptionForm{
					ProjectID: projectID,
					GroupID:   groupID,
				}
				shieldGroup := &shieldv1beta1.Group{
					Slug:     "test-group",
					Metadata: test.metadata,
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
				client := new(mocks.SirenServiceClient)
				defer client.AssertExpectations(t)

				service := alert.NewSubscriptionService(client, shield)
				err := service.UpdateSubscription(ctx, subscriptionID, form)
				assert.ErrorIs(t, err, alert.ErrNoShieldSlackChannel)
			})
		}
	})

	t.Run("should update subscription on success", func(t *testing.T) {
		channelName := "test-slack-channel"
		sirenNamespace := 5
		form := alert.SubscriptionForm{
			UserID:             "john.doe@example.com",
			AlertSeverity:      alert.AlertSeverityCritical,
			ChannelCriticality: alert.ChannelCriticalityInfo,
			GroupID:            groupID,
			ProjectID:          projectID,
			ResourceType:       "firehose",
			ResourceID:         "test-job",
		}
		shieldGroup := &shieldv1beta1.Group{
			Slug: "test-group",
			Metadata: newStruct(t, map[string]interface{}{
				"alerting": map[string]interface{}{
					string(form.ChannelCriticality): map[string]interface{}{
						"slack": map[string]interface{}{
							"channel": channelName,
						},
					},
				},
			}),
		}
		shieldProject := &shieldv1beta1.Project{
			Slug: "my-project-1",
			Metadata: newStruct(t, map[string]interface{}{
				"siren_namespace": sirenNamespace,
			}),
		}
		expectedSirenPayload := &sirenv1beta1.UpdateSubscriptionRequest{
			Id: uint64(subscriptionID),
			Urn: fmt.Sprintf(
				"%s:%s:%s:%s",
				form.GroupID, form.AlertSeverity, form.ResourceType, form.ResourceID,
			),
			Namespace: uint64(sirenNamespace),
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{
					Id: 1,
					Configuration: newStruct(t, map[string]interface{}{
						"channel_name": channelName,
					}),
				},
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
				"channel_criticality": string(form.ChannelCriticality),
				"project_id":          form.ProjectID,
				"project_slug":        shieldProject.Slug,
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
		client := new(mocks.SirenServiceClient)
		client.
			On("UpdateSubscription", ctx, expectedSirenPayload).
			Return(&sirenv1beta1.UpdateSubscriptionResponse{}, nil)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client, shield)
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

func newStruct(t *testing.T, d map[string]interface{}) *structpb.Struct {
	t.Helper()

	strct, err := structpb.NewStruct(d)
	require.NoError(t, err)
	return strct
}
