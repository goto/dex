package server

import (
	"context"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	sirenv1beta1 "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/siren/v1beta1/sirenv1beta1grpc"
	"github.com/MakeNowJust/heredoc"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/goto/dex/internal/server"
	"github.com/goto/dex/pkg/logger"
	"github.com/goto/dex/pkg/telemetry"
)

func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server <command>",
		Aliases: []string{"s"},
		Short:   "Server management",
		Long:    "Server management commands.",
		Example: heredoc.Doc(`
			$ dex server start
			$ dex server start -c ./config.yaml
		`),
		Annotations: map[string]string{
			"group": "core",
		},
	}

	cmd.PersistentFlags().StringP("config", "c", "dex_server.yml", "Path to configuration file")
	cmd.AddCommand(startCommand())
	return cmd
}

func startCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the server",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(cmd)
		if err != nil {
			return err
		}

		zapLog, err := logger.New(&cfg.Log)
		if err != nil {
			return err
		}

		telemetry.Init(cmd.Context(), cfg.Telemetry, zapLog)
		nrApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName(cfg.Telemetry.ServiceName),
			newrelic.ConfigLicense(cfg.Telemetry.NewRelicAPIKey),
		)
		if err != nil {
			return err
		}
		return runServer(cmd.Context(), nrApp, zapLog, cfg)
	}

	return cmd
}

func runServer(baseCtx context.Context, nrApp *newrelic.Application, zapLog *zap.Logger, cfg serverConfig) error {
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()

	shieldConn, err := grpc.Dial(cfg.Shield.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	entropyConn, err := grpc.Dial(cfg.Entropy.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	sirenConn, err := grpc.Dial(cfg.Siren.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	return server.Serve(ctx, cfg.Service.Addr(), nrApp, zapLog,
		shieldv1beta1.NewShieldServiceClient(shieldConn),
		entropyv1beta1.NewResourceServiceClient(entropyConn),
		sirenv1beta1.NewSirenServiceClient(sirenConn),
		cfg.Odin.Addr,
	)
}
