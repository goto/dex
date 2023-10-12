package warden_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/v1/warden"
	"github.com/goto/dex/internal/server/v1/warden/mocks"
	shareMocks "github.com/goto/dex/mocks"
)

func TestHandler_teamList(t *testing.T) {
	t.Run("should return 200 OK", func(t *testing.T) {
		// mocked response from warden returning list of teams for a user email
		wardenResponse := `{
            "success": true,
            "errors": [],
            "data": {
            "teams": [
                {
                    "name": "data_fabric",
                    "created_at": "2023-08-25T03:52:13.548Z",
                    "updated_at": "2023-08-25T03:52:13.548Z",
                    "owner_id": 433,
                    "parent_team_identifier": "2079834a-05c4-420d-bfc8-44b934adea9f",
                    "identifier": "b5aea046-dab3-4dac-b1ea-e1eef423226b",
                    "product_group_name": "data_engineering",
                    "product_group_id": "2079834a-05c4-420d-bfc8-44b934adea9f",
                    "labels": null,
                    "short_code": "T394"
                }
            ]},
            "status": "ok"
        }`
		dexTeamListResonses := `{             
			 "teams": [
					{
					"name": "data_fabric",
					"created_at": "2023-08-25T03:52:13.548Z",
					"updated_at": "2023-08-25T03:52:13.548Z",
					"owner_id": 433,
					"parent_team_identifier": "2079834a-05c4-420d-bfc8-44b934adea9f",
					"identifier": "b5aea046-dab3-4dac-b1ea-e1eef423226b",
					"product_group_name": "data_engineering",
					"product_group_id": "2079834a-05c4-420d-bfc8-44b934adea9f",
					"labels": null,
					"short_code": "T394"
					}
			]}`
		doer := mocks.NewDoer(t)
		doer.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(wardenResponse)),
		}, nil)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/users/me/warden_teams", nil)
		require.NoError(t, err)
		req.Header.Add("X-Auth-Email", "test@email.com")
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		warden.Routes(nil, nil)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, dexTeamListResonses, resp.Body.String())
	})

	t.Run("should return 401 when user email is not present in warden", func(t *testing.T) {
		wardenResponse := `{
            "success": false,
            "errors": [],
            "data": {
            "teams": [
                {
                    "name": "data_fabric",
                    "created_at": "2023-08-25T03:52:13.548Z",
                    "updated_at": "2023-08-25T03:52:13.548Z",
                    "owner_id": 433,
                    "parent_team_identifier": "2079834a-05c4-420d-bfc8-44b934adea9f",
                    "identifier": "b5aea046-dab3-4dac-b1ea-e1eef423226b",
                    "product_group_name": "data_engineering",
                    "product_group_id": "2079834a-05c4-420d-bfc8-44b934adea9f",
                    "labels": null,
                    "short_code": "T394"
                }
            ]},
            "status": "ok"
        }`
		shieldClient := shareMocks.NewShieldServiceClient(t)
		doer := mocks.NewDoer(t)
		doer.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(wardenResponse)),
		}, nil)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/users/me/warden_teams", nil)
		require.NoError(t, err)
		req.Header.Add("X-Auth-Email", "test@email.com")
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		warden.Routes(shieldClient, nil)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("should return Unauthorized when X-Auth-Email is not present is header", func(t *testing.T) {
		shieldClient := shareMocks.NewShieldServiceClient(t)
		// doer := mocks.NewDoer(t)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/users/me/warden_teams", nil)
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		warden.Routes(shieldClient, nil)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestHandler_updateGroup(t *testing.T) {
	t.Run("should return 200 OK", func(t *testing.T) {
		// mocked response from warden returning a team for a wardenUUID
		wardenResponse := `{
         "success": true,
         "data": {
             "name": "data_fabric",
             "created_at": "2023-08-25T03:52:13.548Z",
             "updated_at": "2023-08-25T03:52:13.548Z",
             "owner_id": 433,
             "parent_team_identifier": "2079834a-05c4-420d-bfc8-44b934adea9f",
             "identifier": "b5aea046-dab3-4dac-b1ea-e1eef423226b",
             "product_group_name": "data_engineering",
             "product_group_id": "2079834a-05c4-420d-bfc8-44b934adea9f",
             "labels": null,
             "short_code": "T394"
         },
         "errors": []
     }`

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

		doer := mocks.NewDoer(t)
		doer.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(wardenResponse)),
		}, nil)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, "/groups/e38527ee-a8cd-40f9-98a7-1f0bbd20909f/metadata", bytes.NewBufferString(`{"warden_team_id": "123"}`))
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		warden.Routes(shieldClient, nil)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, dexGroupMetadataResponse, resp.Body.String())
	})

	t.Run("should return 500 when warden_team_id is invalid", func(t *testing.T) {
		wardenResponse := `{
         "success": false,
         "data": {
             "name": "data_fabric",
             "created_at": "2023-08-25T03:52:13.548Z",
             "updated_at": "2023-08-25T03:52:13.548Z",
             "owner_id": 433,
             "parent_team_identifier": "2079834a-05c4-420d-bfc8-44b934adea9f",
             "identifier": "b5aea046-dab3-4dac-b1ea-e1eef423226b",
             "product_group_name": "data_engineering",
             "product_group_id": "2079834a-05c4-420d-bfc8-44b934adea9f",
             "labels": null,
             "short_code": "T394"
         },
         "errors": []
     }`
		shieldClient := shareMocks.NewShieldServiceClient(t)

		doer := mocks.NewDoer(t)
		doer.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(wardenResponse)),
		}, nil)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, "/groups/e38527ee-a8cd-40f9-98a7-1f0bbd20909f/metadata", bytes.NewBufferString(`{"warden_team_id": "123"}`))
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		warden.Routes(shieldClient, nil)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("should return error when warden_team_id is empty string", func(t *testing.T) {
		shieldClient := shareMocks.NewShieldServiceClient(t)
		// doer := mocks.NewDoer(t)

		resp := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodPatch, "/groups/e38527ee-a8cd-40f9-98a7-1f0bbd20909f/metadata", bytes.NewBufferString(`{"warden_team_id": ""}`))
		require.NoError(t, err)
		router := chi.NewRouter()
		router.Use(reqctx.WithRequestCtx())
		warden.Routes(shieldClient, nil)(router)
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
