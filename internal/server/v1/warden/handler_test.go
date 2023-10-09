package warden

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/v1/warden/mocks"
	shareMocks "github.com/goto/dex/mocks"
)

func TestHandler_teamList(t *testing.T) {
	t.Run("should return 200 OK", func(t *testing.T) {
		doer := mocks.NewDoer(t)
		doer.EXPECT().Do(mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewBufferString(`{
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
				]
				},
				"status": "ok"
				}`)),
		}, nil)

		r := chi.NewRouter()
		r.Use(reqctx.WithRequestCtx())
		r.Route("/dex/warden", Routes(nil, doer))

		req, err := http.NewRequest(http.MethodGet, "/dex/warden/users/me/warden_teams", nil)
		assert.NoError(t, err)
		req.Header.Add("X-Auth-Email", "test@email.com")

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{				"teams": [
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
			]}`, resp.Body.String())
	})
}

func TestHandler_updateGroup(t *testing.T) {
	t.Run("should return 200 OK", func(t *testing.T) {
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
			Body: io.NopCloser(bytes.NewBufferString(`{
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
			}`)),
		}, nil)

		r := chi.NewRouter()
		r.Use(reqctx.WithRequestCtx())
		r.Route("/dex/warden", Routes(shieldClient, doer))

		req, err := http.NewRequest(http.MethodPatch, "/dex/warden/groups/e38527ee-a8cd-40f9-98a7-1f0bbd20909f/metadata", bytes.NewBufferString(`{"warden_team_id": "123"}`))
		assert.NoError(t, err)

		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{
			"privacy": "public",
			"product-group-id": "2079834a-05c4-420d-bfc8-44b934adea9f",
			"team-id": "b5aea046-dab3-4dac-b1ea-e1eef423226b"
		}`, resp.Body.String())
	})
}
