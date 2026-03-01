package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/docker"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DockerStartSchema returns the JSON Schema for the docker_start tool.
func DockerStartSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"container_id": map[string]any{
				"type":        "string",
				"description": "Container ID or name to start",
			},
		},
		"required": []any{"container_id"},
	})
	return s
}

// DockerStart returns a tool handler that starts a stopped container.
func DockerStart() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "container_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		containerID := helpers.GetString(req.Arguments, "container_id")

		output, err := docker.Run(ctx, "start", containerID)
		if err != nil {
			return helpers.ErrorResult("docker_error", err.Error()), nil
		}
		return helpers.TextResult(fmt.Sprintf("Container started: %s", output)), nil
	}
}
