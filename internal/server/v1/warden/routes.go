package warden

import (
	chiv5 "github.com/go-chi/chi/v5"
)

func Routes() func(r chiv5.Router) {
	service := NewService(nil)
	handler := NewHandler(service)
	return func(r chiv5.Router) {
		r.Get("/users/me/warden_teams", handler.teamList)
	}
}
