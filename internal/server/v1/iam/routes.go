package iam

import (
	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	chiv5 "github.com/go-chi/chi/v5"
)

func Routes(shieldClient shieldv1beta1rpc.ShieldServiceClient, wardenClient WardenClient) func(r chiv5.Router) {
	service := NewService(shieldClient, wardenClient)
	handler := NewHandler(service)
	return func(r chiv5.Router) {
		r.Get("/users/me/warden_teams", handler.listUserWardenTeams)

		r.Put("/groups/{group_id}/metadata/warden", handler.linkGroupToWarden)
	}
}
