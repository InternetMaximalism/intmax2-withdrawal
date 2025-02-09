package withdrawal_server

import (
	"intmax2-withdrawal/configs"
	"intmax2-withdrawal/internal/logger"

	"github.com/dimiro1/health"
)

type WithdrawalServer struct {
	log              logger.Logger
	config           *configs.Config
	dbApp            SQLDriverApp
	commands         Commands
	cookieForAuthUse bool
	hc               *health.Handler
}

// New initializes a new Server struct.
func New(
	log logger.Logger,
	config *configs.Config,
	dbApp SQLDriverApp,
	commands Commands,
	cookieForAuthUse bool,
	hc *health.Handler,
) *WithdrawalServer {
	const (
		srv  = "server"
		name = "withdrawal"
	)

	return &WithdrawalServer{
		log:              log.WithFields(logger.Fields{srv: name}),
		config:           config,
		dbApp:            dbApp,
		commands:         commands,
		cookieForAuthUse: cookieForAuthUse,
		hc:               hc,
	}
}

func (s *WithdrawalServer) Commands() Commands {
	return s.commands
}
