package alert

import (
	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	sirenv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/siren/v1beta1/sirenv1beta1grpc"
)

type AlertService struct {
	sirenClient  sirenv1beta1grpc.SirenServiceClient
	shieldClient shieldv1beta1rpc.ShieldServiceClient
}

func NewAlertService(
	sirenClient sirenv1beta1grpc.SirenServiceClient,
	shieldClient shieldv1beta1rpc.ShieldServiceClient,
) *AlertService {
	return &AlertService{
		sirenClient:  sirenClient,
		shieldClient: shieldClient,
	}
}
