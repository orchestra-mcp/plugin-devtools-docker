package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/docker"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DockerLogsSchema returns the JSON Schema for the docker_logs tool.
func DockerLogsSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"container_id": map[string]any{
				"type":        "string",
				"description": "Container ID or name to get logs from",
			},
			"tail": map[string]any{
				"type":        "number",
				"description": "Number of lines to show from the end of the logs",
			},
			"since": map[string]any{
				"type":        "string",
				"description": "Show logs since timestamp (e.g. 2024-01-01T00:00:00) or relative (e.g. 10m)",
			},
			"follow": map[string]any{
				"type":        "boolean",
				"description": "Follow log output (not recommended for non-interactive use)",
			},
		},
		"required": []any{"container_id"},
	})
	return s
}

// DockerLogs returns a tool handler that gets container logs.
func DockerLogs() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "container_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		containerID := helpers.GetString(req.Arguments, "container_id")
		tail := helpers.GetInt(req.Arguments, "tail")
		since := helpers.GetString(req.Arguments, "since")
		follow := helpers.GetBool(req.Arguments, "follow")

		args := []string{"logs"}
		if tail > 0 {
			args = append(args, "--tail", fmt.Sprintf("%d", tail))
		}
		if since != "" {
			args = append(args, "--since", since)
		}
		if follow {
			args = append(args, "--follow")
		}
		args = append(args, containerID)

		output, err := docker.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("docker_error", err.Error()), nil
		}
		if output == "" {
			output = "No logs available"
		}
		return helpers.TextResult(output), nil
	}
}
