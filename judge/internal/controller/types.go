package controller

import (
	"context"
	"github.com/CT1403-2/Code-Judgement/judge/internal/runner"
	"github.com/CT1403-2/Code-Judgement/proto"

	"github.com/CT1403-2/Code-Judgement/judge/config"
)

//go:generate mockery --name=Controller --filename=controller.go --outpkg=mocks
type Controller interface {
	Run(ctx context.Context)
}

func New(cfg *config.Config, clientBuilder func(*config.Config) (proto.ManagerClient, error), runner runner.Runner) Controller {
	return &controller{
		config:        cfg,
		clientBuilder: clientBuilder,
		runner:        runner,
	}
}
