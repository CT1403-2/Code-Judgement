package runner

import (
	"context"
	"github.com/CT1403-2/Code-Judgement/judge/config"
	"github.com/CT1403-2/Code-Judgement/proto"
)

//go:generate mockery --name=Runner --filename=runner.go --outpkg=mocks
type Runner interface {
	Run(ctx context.Context, question *proto.Question, submission *proto.Submission) (*proto.SubmissionState, error)
}

func New(cfg *config.Config) Runner {
	return &dockerRunner{
		config: cfg,
	}
}
