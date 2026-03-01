package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/docker"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DockerExecSchema returns the JSON Schema for the docker_exec tool.
func DockerExecSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"container_id": map[string]any{
				"type":        "string",
				"description": "Container ID or name to execute command in",
			},
			"command": map[string]any{
				"type":        "string",
				"description": "Command to execute inside the container (passed to sh -c)",
			},
			"workdir": map[string]any{
				"type":        "string",
				"description": "Working directory inside the container",
			},
		},
		"required": []any{"container_id", "command"},
	})
	return s
}

// DockerExec returns a tool handler that executes a command in a running container.
func DockerExec() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "container_id", "command"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		containerID := helpers.GetString(req.Arguments, "container_id")
		command := helpers.GetString(req.Arguments, "command")
		workdir := helpers.GetString(req.Arguments, "workdir")

		args := []string{"exec"}
		if workdir != "" {
			args = append(args, "--workdir", workdir)
		}
		args = append(args, containerID, "sh", "-c", command)

		output, err := docker.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("docker_error", err.Error()), nil
		}
		if output == "" {
			output = "Command executed successfully (no output)"
		}
		return helpers.TextResult(output), nil
	}
}
