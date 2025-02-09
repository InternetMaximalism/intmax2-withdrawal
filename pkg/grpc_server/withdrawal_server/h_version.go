package withdrawal_server

import (
	"context"
	"intmax2-withdrawal/configs/buildvars"
	"intmax2-withdrawal/internal/open_telemetry"
	node "intmax2-withdrawal/internal/pb/gen/withdrawal_service/node"
	"intmax2-withdrawal/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *WithdrawalServer) Version(ctx context.Context, _ *node.VersionRequest) (*node.VersionResponse, error) {
	const (
		hName     = "Handler Version"
		version   = "version"
		buildTime = "build_time"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(version, buildvars.Version),
			attribute.String(buildTime, buildvars.BuildTime),
		))
	defer span.End()

	info := s.Commands().GetVersion(buildvars.Version, buildvars.BuildTime).Do(spanCtx)
	return &node.VersionResponse{
		Version:   info.Version,
		Buildtime: info.BuildTime,
	}, utils.OK(spanCtx)
}
