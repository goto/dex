package iam_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/v1/iam"
	"github.com/goto/dex/internal/server/v1/iam/mocks"
	shareMocks "github.com/goto/dex/mocks"
	"github.com/goto/dex/pkg/errors"
	"github.com/goto/dex/warden"
)

func TestHandler_teamList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		dexTeamListResonses := `{             
			 "teams": [
					{
					"name": "data_fabric",
					"created_at": "2023-10-13T14:42:58+05:30",
					"updated_at": "2023-10-13T14:42:58+05:30",
					"owner_id": 433,
					"parent_team_identifier": "2079834a-05c4-420d-bfc8-44b934adea9f",
					"identifier": "b5aea046-dab3-4dac-b1ea-e1eef423226b",
					"product_group_name": "data_engineering",
					"product_group_id": "2079834a-05c4-420d-bfc8-44b934adea9f",
					"labels": null,
					"short_code": "T394"
					}
			]}`
		wardenClient := mocks.NewWardenClient(t)
		frozenTime := time.Unix(1697188378, 0)
		wardenClient.EXPECT().ListUserTeams(mock.Anything, warden.TeamListRequest{
			Email: "test@domain.com",
		}).Return([]warden.Team{
			{
				Name:                 "data_fabric",
				CreatedAt:            frozenTime,
				UpdatedAt:            frozenTime,
				OwnerID:              433,
				ParentTeamIdentifier: "2079834a-05c4-420d-bfc8-44b934adea9f",
				Identifier:           "b5aea046-dab3-4dac-b1ea-e1eef423226b",
				ProductGroupName:     "data_engineering",
				ProductGroupID:       "2079834a-05c4-420d-bfc8-44b934adea9f",
				Labels:               nil,
				ShortCode:            "T394",
			},
		}, nil)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/users/me/warden_teams", nil)
		require.NoError(t, err)
		req.Header.Add("X-Auth-Email", "test@domain.com")
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		iam.Routes(nil, wardenClient)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, dexTeamListResonses, resp.Body.String())
	})

	t.Run("EmailNotFound", func(t *testing.T) {
		wardenClient := mocks.NewWardenClient(t)
		wardenClient.EXPECT().ListUserTeams(mock.Anything, warden.TeamListRequest{
			Email: "test@domain.com",
		}).Return(nil, errors.ErrNotFound)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/users/me/warden_teams", nil)
		require.NoError(t, err)
		req.Header.Add("X-Auth-Email", "test@domain.com")
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		iam.Routes(nil, wardenClient)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"code":"not_found", "message":"Requested entity not found", "op":"", "status":404}`, resp.Body.String())
	})

	t.Run("WardenClientFailure", func(t *testing.T) {
		wardenClient := mocks.NewWardenClient(t)
		wardenClient.EXPECT().ListUserTeams(mock.Anything, warden.TeamListRequest{
			Email: "test@domain.com",
		}).Return(nil, errors.ErrInternal)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/users/me/warden_teams", nil)
		require.NoError(t, err)
		req.Header.Add("X-Auth-Email", "test@domain.com")
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		iam.Routes(nil, wardenClient)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"code":"internal_error", "message":"Some unexpected error occurred", "op":"", "status":500}`, resp.Body.String())
	})

	t.Run("MissingEmail", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/users/me/warden_teams", nil)
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		iam.Routes(nil, nil)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestHandler_updateGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		dexGroupMetadataResponse := `{
			"privacy": "public",
			"product-group-id": "2079834a-05c4-420d-bfc8-44b934adea9f",
			"team-id": "b5aea046-dab3-4dac-b1ea-e1eef423226b"
		}`
		groupID := "e38527ee-a8cd-40f9-98a7-1f0bbd20909f"
		metaData, _ := structpb.NewStruct(map[string]any{
			"privacy": "public",
		})

		updatedMetaData, _ := structpb.NewStruct(map[string]any{
			"privacy":          "public",
			"team-id":          "b5aea046-dab3-4dac-b1ea-e1eef423226b",
			"product-group-id": "2079834a-05c4-420d-bfc8-44b934adea9f",
		})

		wardenClient := mocks.NewWardenClient(t)
		frozenTime := time.Unix(1697188378, 0)
		wardenClient.EXPECT().TeamByUUID(mock.Anything, warden.TeamByUUIDRequest{
			TeamUUID: "123",
		}).Return(&warden.Team{
			Name:                 "data_fabric",
			CreatedAt:            frozenTime,
			UpdatedAt:            frozenTime,
			OwnerID:              433,
			ParentTeamIdentifier: "2079834a-05c4-420d-bfc8-44b934adea9f",
			Identifier:           "b5aea046-dab3-4dac-b1ea-e1eef423226b",
			ProductGroupName:     "data_engineering",
			ProductGroupID:       "2079834a-05c4-420d-bfc8-44b934adea9f",
			Labels:               nil,
			ShortCode:            "T394",
		}, nil)

		shieldClient := shareMocks.NewShieldServiceClient(t)
		shieldClient.EXPECT().GetGroup(mock.Anything, &shieldv1beta1.GetGroupRequest{
			Id: groupID,
		}).Return(&shieldv1beta1.GetGroupResponse{
			Group: &shieldv1beta1.Group{
				Id: groupID, Name: "test", Slug: "testSlug", OrgId: "123", Metadata: metaData,
			},
		}, nil)

		shieldClient.EXPECT().UpdateGroup(mock.Anything, &shieldv1beta1.UpdateGroupRequest{
			Id: groupID, Body: &shieldv1beta1.GroupRequestBody{
				Metadata: updatedMetaData, Name: "test", Slug: "testSlug", OrgId: "123",
			},
		}).Return(&shieldv1beta1.UpdateGroupResponse{Group: &shieldv1beta1.Group{
			Id: groupID, Name: "test", Slug: "testSlug", OrgId: "123", Metadata: updatedMetaData,
		}}, nil)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, "/groups/e38527ee-a8cd-40f9-98a7-1f0bbd20909f/metadata", bytes.NewBufferString(`{"warden_team_id": "123"}`))
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		iam.Routes(shieldClient, wardenClient)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, dexGroupMetadataResponse, resp.Body.String())
	})

	t.Run("MissingWardenTeamID", func(t *testing.T) {
		// response return by handler
		dexGroupMetadataResponse := `{"code":"", "message":"missing warden_team_id", "op":"", "status":400}`

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, "/groups/e38527ee-a8cd-40f9-98a7-1f0bbd20909f/metadata", bytes.NewBufferString(`{}`))
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		iam.Routes(nil, nil)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, dexGroupMetadataResponse, resp.Body.String())
	})

	t.Run("WardenIdNotFound", func(t *testing.T) {
		// response returned by warden client
		dexGroupMetadataResponse := `{"code":"not_found", "message":"Requested entity not found", "op":"", "status":404}`

		wardenClient := mocks.NewWardenClient(t)
		wardenClient.EXPECT().TeamByUUID(mock.Anything, warden.TeamByUUIDRequest{
			TeamUUID: "123",
		}).Return(nil, errors.ErrNotFound)

		shieldClient := shareMocks.NewShieldServiceClient(t)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, "/groups/e38527ee-a8cd-40f9-98a7-1f0bbd20909f/metadata", bytes.NewBufferString(`{"warden_team_id": "123"}`))
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		iam.Routes(shieldClient, wardenClient)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, dexGroupMetadataResponse, resp.Body.String())
	})

	t.Run("ShieldGetGroupFailure", func(t *testing.T) {
		// response returned by shield client
		dexGroupMetadataResponse := `{"code":"internal_error", "message":"Some unexpected error occurred", "op":"", "status":500}`
		groupID := "e38527ee-a8cd-40f9-98a7-1f0bbd20909f"

		wardenClient := mocks.NewWardenClient(t)
		frozenTime := time.Unix(1697188378, 0)
		wardenClient.EXPECT().TeamByUUID(mock.Anything, warden.TeamByUUIDRequest{
			TeamUUID: "123",
		}).Return(&warden.Team{
			Name:                 "data_fabric",
			CreatedAt:            frozenTime,
			UpdatedAt:            frozenTime,
			OwnerID:              433,
			ParentTeamIdentifier: "2079834a-05c4-420d-bfc8-44b934adea9f",
			Identifier:           "b5aea046-dab3-4dac-b1ea-e1eef423226b",
			ProductGroupName:     "data_engineering",
			ProductGroupID:       "2079834a-05c4-420d-bfc8-44b934adea9f",
			Labels:               nil,
			ShortCode:            "T394",
		}, nil)

		shieldClient := shareMocks.NewShieldServiceClient(t)
		shieldClient.EXPECT().GetGroup(mock.Anything, &shieldv1beta1.GetGroupRequest{
			Id: groupID,
		}).Return(nil, errors.ErrInternal)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, "/groups/e38527ee-a8cd-40f9-98a7-1f0bbd20909f/metadata", bytes.NewBufferString(`{"warden_team_id": "123"}`))
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		iam.Routes(shieldClient, wardenClient)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, dexGroupMetadataResponse, resp.Body.String())
	})

	t.Run("ShielUpdateGroupFailure", func(t *testing.T) {
		// response returned by shield client
		dexGroupMetadataResponse := `{"code":"internal_error", "message":"Some unexpected error occurred", "op":"", "status":500}`
		groupID := "e38527ee-a8cd-40f9-98a7-1f0bbd20909f"
		metaData, _ := structpb.NewStruct(map[string]any{
			"privacy": "public",
		})

		updatedMetaData, _ := structpb.NewStruct(map[string]any{
			"privacy":          "public",
			"team-id":          "b5aea046-dab3-4dac-b1ea-e1eef423226b",
			"product-group-id": "2079834a-05c4-420d-bfc8-44b934adea9f",
		})

		wardenClient := mocks.NewWardenClient(t)
		frozenTime := time.Unix(1697188378, 0)
		wardenClient.EXPECT().TeamByUUID(mock.Anything, warden.TeamByUUIDRequest{
			TeamUUID: "123",
		}).Return(&warden.Team{
			Name:                 "data_fabric",
			CreatedAt:            frozenTime,
			UpdatedAt:            frozenTime,
			OwnerID:              433,
			ParentTeamIdentifier: "2079834a-05c4-420d-bfc8-44b934adea9f",
			Identifier:           "b5aea046-dab3-4dac-b1ea-e1eef423226b",
			ProductGroupName:     "data_engineering",
			ProductGroupID:       "2079834a-05c4-420d-bfc8-44b934adea9f",
			Labels:               nil,
			ShortCode:            "T394",
		}, nil)

		shieldClient := shareMocks.NewShieldServiceClient(t)
		shieldClient.EXPECT().GetGroup(mock.Anything, &shieldv1beta1.GetGroupRequest{
			Id: groupID,
		}).Return(&shieldv1beta1.GetGroupResponse{
			Group: &shieldv1beta1.Group{
				Id: groupID, Name: "test", Slug: "testSlug", OrgId: "123", Metadata: metaData,
			},
		}, nil)

		shieldClient.EXPECT().UpdateGroup(mock.Anything, &shieldv1beta1.UpdateGroupRequest{
			Id: groupID, Body: &shieldv1beta1.GroupRequestBody{
				Metadata: updatedMetaData, Name: "test", Slug: "testSlug", OrgId: "123",
			},
		}).Return(nil, errors.ErrInternal)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, "/groups/e38527ee-a8cd-40f9-98a7-1f0bbd20909f/metadata", bytes.NewBufferString(`{"warden_team_id": "123"}`))
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		iam.Routes(shieldClient, wardenClient)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, dexGroupMetadataResponse, resp.Body.String())
	})
}
