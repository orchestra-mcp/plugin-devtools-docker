package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/docker"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DockerInspectSchema returns the JSON Schema for the docker_inspect tool.
func DockerInspectSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"container_id": map[string]any{
				"type":        "string",
				"description": "Container or image ID/name to inspect",
			},
			"format": map[string]any{
				"type":        "string",
				"description": "Go template format string for output (e.g. '{{.State.Status}}')",
			},
		},
		"required": []any{"container_id"},
	})
	return s
}

// DockerInspect returns a tool handler that inspects a container or image.
func DockerInspect() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "container_id"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		containerID := helpers.GetString(req.Arguments, "container_id")
		format := helpers.GetString(req.Arguments, "format")

		args := []string{"inspect"}
		if format != "" {
			args = append(args, "--format", format)
		}
		args = append(args, containerID)

		output, err := docker.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("docker_error", err.Error()), nil
		}
		return helpers.TextResult(output), nil
	}
}
