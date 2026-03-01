package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/docker"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DockerListContainersSchema returns the JSON Schema for the docker_list_containers tool.
func DockerListContainersSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"all": map[string]any{
				"type":        "boolean",
				"description": "Show all containers (including stopped). Default false",
			},
			"format": map[string]any{
				"type":        "string",
				"description": "Go template format string for output",
			},
		},
	})
	return s
}

// DockerListContainers returns a tool handler that lists Docker containers.
func DockerListContainers() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		all := helpers.GetBool(req.Arguments, "all")
		format := helpers.GetString(req.Arguments, "format")

		args := []string{"ps"}
		if all {
			args = append(args, "--all")
		}
		if format != "" {
			args = append(args, "--format", format)
		}

		output, err := docker.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("docker_error", err.Error()), nil
		}
		if output == "" {
			output = "No containers found"
		}
		return helpers.TextResult(output), nil
	}
}
