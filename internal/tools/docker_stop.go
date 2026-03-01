package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/docker"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DockerStopSchema returns the JSON Schema for the docker_stop tool.
func DockerStopSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"container_id": map[string]any{
				"type":        "string",
				"description": "Container ID or name to stop",
			},
			"timeout": map[string]any{
				"type":        "number",
				"description": "Seconds to wait before killing the container (default 10)",
			},
		},
		"required": []any{"container_id"},
	})
	return s
}

// DockerStop returns a tool handler that stops a running container.
func DockerStop() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "container_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		containerID := helpers.GetString(req.Arguments, "container_id")
		timeout := helpers.GetInt(req.Arguments, "timeout")

		args := []string{"stop"}
		if timeout > 0 {
			args = append(args, "--time", fmt.Sprintf("%d", timeout))
		}
		args = append(args, containerID)

		output, err := docker.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("docker_error", err.Error()), nil
		}
		return helpers.TextResult(fmt.Sprintf("Container stopped: %s", output)), nil
	}
}
