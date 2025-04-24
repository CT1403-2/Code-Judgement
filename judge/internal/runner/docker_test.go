package runner

import (
	"context"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"

	"github.com/CT1403-2/Code-Judgement/judge/config"
	"github.com/CT1403-2/Code-Judgement/proto"
)

type DockerRunnerSuite struct {
	suite.Suite
	question *proto.Question
	config   *config.Config
}

func TestDockerRunnerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DockerRunnerSuite))
}

func (s *DockerRunnerSuite) SetupSuite() {
	input := "b102aaef-4ac6-4c15-a92c-919fdef0d5c0"
	expectedOutput := "b102aaef-4ac6-4c15-a92c-919fdef0d5c0"
	s.question = &proto.Question{
		Id:        stringPtr("q123"),
		Title:     "Echo",
		Statement: "Write a program that echos the input.",
		Input:     &input,
		Output:    &expectedOutput,
		Limitations: &proto.Limitations{
			Duration: 1000,
			Memory:   512,
		},
	}

	s.config = &config.Config{
		Runner: config.RunnerConfig{
			Image: "runner:v0.0.10",
		},
	}
}

func (s *DockerRunnerSuite) TestSuccessfulRun() {
	code, err := os.ReadFile("test_data/good_code")
	if err != nil {
		s.Failf("Failed to read test code file: %v", err.Error())
	}

	submission := &proto.Submission{
		Id:         stringPtr("good-submission"),
		QuestionId: "q123",
		Code:       code,
		State:      statePtr(proto.SubmissionState_SUBMISSION_STATE_JUDGING),
	}

	runner := New(s.config)

	ctx := context.Background()
	result, err := runner.Run(ctx, s.question, submission)

	if err != nil {
		s.Failf("Error running submission: %v", err.Error())
	}

	s.Equal(proto.SubmissionState_SUBMISSION_STATE_OK, *result)
}

func (s *DockerRunnerSuite) TestFalse() {
	code, err := os.ReadFile("test_data/false_code")
	if err != nil {
		s.Failf("Failed to read test code file: %v", err.Error())
	}

	submission := &proto.Submission{
		Id:         stringPtr("false-submission"),
		QuestionId: "q123",
		Code:       code,
		State:      statePtr(proto.SubmissionState_SUBMISSION_STATE_JUDGING),
	}

	runner := New(s.config)

	ctx := context.Background()
	result, err := runner.Run(ctx, s.question, submission)

	if err != nil {
		s.Failf("Error running submission: %v", err.Error())
	}

	s.Equal(proto.SubmissionState_SUBMISSION_STATE_WRONG_ANSWER, *result)
}

func (s *DockerRunnerSuite) TestNonCompilable() {
	code, err := os.ReadFile("test_data/non_compilable_code")
	if err != nil {
		s.Failf("Failed to read test code file: %v", err.Error())
	}

	submission := &proto.Submission{
		Id:         stringPtr("non-compilable-submission"),
		QuestionId: "q123",
		Code:       code,
		State:      statePtr(proto.SubmissionState_SUBMISSION_STATE_JUDGING),
	}

	runner := New(s.config)

	ctx := context.Background()
	result, err := runner.Run(ctx, s.question, submission)

	if err != nil {
		s.Failf("Error running submission: %v", err.Error())
	}

	s.Equal(proto.SubmissionState_SUBMISSION_STATE_COMPILE_ERROR, *result)
}

func (s *DockerRunnerSuite) TestTimeLimit() {
	code, err := os.ReadFile("test_data/time_limit_code")
	if err != nil {
		s.Failf("Failed to read test code file: %v", err.Error())
	}

	submission := &proto.Submission{
		Id:         stringPtr("time-limit-submission"),
		QuestionId: "q123",
		Code:       code,
		State:      statePtr(proto.SubmissionState_SUBMISSION_STATE_JUDGING),
	}

	runner := New(s.config)

	ctx := context.Background()
	result, err := runner.Run(ctx, s.question, submission)

	if err != nil {
		s.Failf("Error running submission: %v", err.Error())
	}

	s.Equal(proto.SubmissionState_SUBMISSION_STATE_TIME_LIMIT_EXCEEDED, *result)
}

func (s *DockerRunnerSuite) TestMemoryLimit() {
	code, err := os.ReadFile("test_data/memory_limit_code")
	if err != nil {
		s.Failf("Failed to read test code file: %v", err.Error())
	}

	input := "b102aaef-4ac6-4c15-a92c-919fdef0d5c0"
	expectedOutput := "b102aaef-4ac6-4c15-a92c-919fdef0d5c0"
	question := &proto.Question{
		Id:        stringPtr("q123"),
		Title:     "Echo",
		Statement: "Write a program that echos the input.",
		Input:     &input,
		Output:    &expectedOutput,
		Limitations: &proto.Limitations{
			Duration: 1000,
			Memory:   8,
		},
	}

	submission := &proto.Submission{
		Id:         stringPtr("memory-limit-submission"),
		QuestionId: "q123",
		Code:       code,
		State:      statePtr(proto.SubmissionState_SUBMISSION_STATE_JUDGING),
	}
	runner := New(s.config)

	ctx := context.Background()
	result, err := runner.Run(ctx, question, submission)

	if err != nil {
		s.Failf("Error running submission: %v", err.Error())
	}

	s.Equal(proto.SubmissionState_SUBMISSION_STATE_MEMORY_LIMIT_EXCEEDED, *result)
}

func stringPtr(s string) *string {
	return &s
}
