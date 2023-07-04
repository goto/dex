package tests

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Params map[string]string

// WithUrlParams returns a pointer to a request object with the given URL params
// added to a new chi.Context object. for single param assignment see WithUrlParam.
func SetChiParams(r *http.Request, params Params) *http.Request {
	chiCtx := chi.NewRouteContext()
	req := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	for key, value := range params {
		chiCtx.URLParams.Add(key, value)
	}
	return req
}
