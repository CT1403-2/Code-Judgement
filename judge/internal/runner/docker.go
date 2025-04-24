package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/CT1403-2/Code-Judgement/judge/config"
	"github.com/CT1403-2/Code-Judgement/proto"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	timeOutError      = "timeout: sending signal TERM to command './main'\n"
	didntCompileError = "timeout: failed to run command './main': No such file or directory\n"
)

type dockerRunner struct {
	config *config.Config
}

type SuiteConfig struct {
	Code  []byte `json:"code"`
	Input string `json:"input"`
}

func (d dockerRunner) Run(ctx context.Context, question *proto.Question, submission *proto.Submission) (*proto.SubmissionState, error) {
	logger := logrus.WithFields(logrus.Fields{"submission_id": *submission.Id})
	logger.Info("Starting submission evaluation")

	docker, err := createDockerClient(ctx)
	if err != nil {
		return nil, err
	}
	defer docker.Close()

	suite := SuiteConfig{
		Code:  submission.Code,
		Input: *question.Input,
	}
	jsonSuite, err := json.Marshal(suite)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %w", err)
	}

	containerConfig, hostConfig := d.prepareContainerConfig(question, jsonSuite)
	containerName := fmt.Sprintf("submission-%s-judge", *submission.Id)
	resp, err := docker.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	containerID := resp.ID
	logger.WithField("container_id", containerID).Info("Container created")

	defer func() {
		removeOptions := container.RemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}
		if err := docker.ContainerRemove(context.Background(), containerID, removeOptions); err != nil {
			logger.WithError(err).Warn("Failed to remove container")
		} else {
			logger.Info("Container removed")
		}
	}()

	if err := docker.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}
	logger.Info("Container started")

	var inspect container.InspectResponse
	statusCh, errCh := docker.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	var statusCode int64

	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("error waiting for container: %w", err)
		}
	case status := <-statusCh:
		statusCode = status.StatusCode
		inspect, err = docker.ContainerInspect(ctx, containerID)
		if err != nil {
			return nil, fmt.Errorf("couldn't inspect container: %w", err)
		}
	}

	logger.WithField("status_code", statusCode).Info("Container execution completed")

	stdoutStr, stderrStr, err := getContainerLogs(ctx, docker, containerID)
	if err != nil {
		return nil, err
	}

	result := d.evaluateResult(statusCode, stdoutStr, stderrStr, question, inspect.State.OOMKilled)
	logger.WithField("result", result.String()).Info("Submission evaluated")

	return result, nil
}

func createDockerClient(ctx context.Context) (*client.Client, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	docker.NegotiateAPIVersion(ctx)
	return docker, nil
}

func (d dockerRunner) prepareContainerConfig(question *proto.Question, jsonSuite []byte) (*container.Config, *container.HostConfig) {
	containerConfig := &container.Config{
		Image:           d.config.Runner.Image,
		Cmd:             []string{"/bin/sh", "-c", "echo '" + string(jsonSuite) + "' > /playground/app/suite.json && ./run.sh"},
		Tty:             false,
		NetworkDisabled: true,
		Env:             []string{fmt.Sprintf("TIMEOUT=%d", question.Limitations.Duration/1000)},
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory:           question.Limitations.Memory * 1024 * 1024,
			MemorySwap:       question.Limitations.Memory * 1024 * 1024,
			MemorySwappiness: &[]int64{0}[0],
			CPUPeriod:        100000,
			CPUQuota:         100000,
		},
	}

	return containerConfig, hostConfig
}

func getContainerLogs(ctx context.Context, docker *client.Client, containerID string) (string, string, error) {
	out, err := docker.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer out.Close()

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, out)
	if err != nil {
		return "", "", fmt.Errorf("failed to read container output: %w", err)
	}

	return stdout.String(), stderr.String(), nil
}

func (d dockerRunner) evaluateResult(statusCode int64, stdoutStr, stderrStr string, question *proto.Question, isOOMKilled bool) *proto.SubmissionState {
	if isOOMKilled {
		return statePtr(proto.SubmissionState_SUBMISSION_STATE_MEMORY_LIMIT_EXCEEDED)
	}

	if statusCode == 0 && stdoutStr == *question.Output {
		return statePtr(proto.SubmissionState_SUBMISSION_STATE_OK)
	}

	if statusCode == 123 {
		if stderrStr == timeOutError {
			return statePtr(proto.SubmissionState_SUBMISSION_STATE_TIME_LIMIT_EXCEEDED)
		}

		if strings.HasSuffix(stderrStr, didntCompileError) {
			return statePtr(proto.SubmissionState_SUBMISSION_STATE_COMPILE_ERROR)
		}
	}

	return statePtr(proto.SubmissionState_SUBMISSION_STATE_WRONG_ANSWER)
}

func pullImage(ctx context.Context, docker *client.Client, img string) error {
	_, err := docker.ImagePull(ctx, img, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull docker image: %w", err)
	}
	return nil
}

func statePtr(e proto.SubmissionState) *proto.SubmissionState {
	return &e
}
