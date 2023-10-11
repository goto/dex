package warden

import (
	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	chiv5 "github.com/go-chi/chi/v5"
)

func Routes(shieldClient shieldv1beta1rpc.ShieldServiceClient, doer HTTPClient, wardenAddr string) func(r chiv5.Router) {
	service := NewService(shieldClient, doer, wardenAddr)
	handler := NewHandler(service)
	return func(r chiv5.Router) {
		r.Get("/users/me/warden_teams", handler.teamList)

		r.Patch("/groups/{group_id}/metadata", handler.updateGroupMetadata)
	}
}
