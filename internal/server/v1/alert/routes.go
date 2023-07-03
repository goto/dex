package alert

import (
	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	sirenv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/siren/v1beta1/sirenv1beta1grpc"

	"github.com/go-chi/chi/v5"
)

func SubscriptionRoutes(
	siren sirenv1beta1rpc.SirenServiceClient,
	shield shieldv1beta1rpc.ShieldServiceClient,
) func(chi.Router) {
	subSrv := NewSubscriptionService(siren)
	handler := NewHandler(subSrv, shield)

	return func(r chi.Router) {
		// CRUD operations
		r.Get("/", handler.GetSubscriptions)
		r.Get("/{subscription_id}", handler.FindSubscription)

		r.Post("/", handler.CreateSubscription)
		r.Put("/{subscription_id}", handler.UpdateSubscription)
		r.Delete("/{subscription_id}", handler.DeleteSubscription)
	}
}
