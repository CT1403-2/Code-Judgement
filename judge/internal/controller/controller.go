package controller

import (
	"context"
	"fmt"
	"github.com/CT1403-2/Code-Judgement/judge/internal/runner"
	"github.com/CT1403-2/Code-Judgement/proto"
	"github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/CT1403-2/Code-Judgement/judge/config"
)

type controller struct {
	config        *config.Config
	clientBuilder func(*config.Config) (proto.ManagerClient, error)
	runner        runner.Runner
	client        proto.ManagerClient
}

func (c *controller) Run(ctx context.Context) {
	client, err := c.clientBuilder(c.config)
	if err != nil {
		logrus.WithError(err).Fatal("failed to run controller")
		os.Exit(1)
	}
	c.client = client
	c.judgePendingSubmissions(ctx)
}

func (c *controller) judgePendingSubmissions(ctx context.Context) {
	for { // if possible, receive a stream of submissions here
		select {
		case <-ctx.Done():
			return

		default:
			time.Sleep(10 * time.Millisecond)
			ctxWithTimeout, cancel := context.WithTimeout(ctx, c.config.Manager.Timeout)
			pendingState := proto.SubmissionState_SUBMISSION_STATE_PENDING.String()
			submissions, err := c.client.GetSubmissions(ctxWithTimeout, &proto.GetSubmissionsRequest{Filters: []*proto.Filter{{Field: "state", Value: pendingState}}})
			if err != nil {
				logrus.WithError(err).Error("couldn't get submissions")
			}

			for _, submission := range submissions.Submissions {
				err = c.judgeSubmission(ctxWithTimeout, submission)
				if err != nil {
					logrus.WithError(err).Error("failed to judge submission")
				}
			}

			cancel()
		}
	}
}

func (c *controller) judgeSubmission(ctx context.Context, submission *proto.Submission) error {

	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.config.Manager.Timeout)
	defer cancel()

	response, err := c.client.GetQuestion(ctxWithTimeout, &proto.ID{Value: submission.QuestionId})
	if err != nil {
		return fmt.Errorf("failed to judge submission:\n %w", err)
	}

	state := proto.SubmissionState_SUBMISSION_STATE_JUDGING
	submission.State = &state
	_, err = c.client.UpdateSubmission(ctxWithTimeout, submission)
	if err != nil {
		return fmt.Errorf("failed to judge submission:\n %w", err)
	}
	// todo: check if is updated

	updatedState, err := c.runner.Run(ctx, response.Question, submission)
	if err != nil {
		return fmt.Errorf("failed to judge submission:\n %w", err)
	}

	submission.State = updatedState
	_, err = c.client.UpdateSubmission(ctxWithTimeout, submission)
	if err != nil {
		return fmt.Errorf("failed to judge submission:\n %w", err)
	}

	return nil
}
