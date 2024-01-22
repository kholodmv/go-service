package interceptors

import (
	"github.com/kholodmv/go-service/internal/configs"
	"net"

	"github.com/kholodmv/go-service/internal/middleware/logger"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RegisterUnaryInterceptorChain(cfg configs.ServerConfig) grpc.ServerOption {
	log := logger.Initialize()
	var ipNet *net.IPNet
	if cfg.TrustedSubnet != "" {
		_, trustedNet, err := net.ParseCIDR(cfg.TrustedSubnet)
		ipNet = trustedNet
		if err != nil {
			log.Fatal("unable to parse trusted subnet", zap.Error(err))
		}
	}
	return grpc.ChainUnaryInterceptor(
		InitCheckSubnet(ipNet),
		InitVerifySignature(cfg.Key),
		logging.UnaryServerInterceptor(InterceptorLogger(log.Desugar())),
	)
}
