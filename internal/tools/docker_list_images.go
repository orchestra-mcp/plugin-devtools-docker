package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/docker"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DockerListImagesSchema returns the JSON Schema for the docker_list_images tool.
func DockerListImagesSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"all": map[string]any{
				"type":        "boolean",
				"description": "Show all images (including intermediate). Default false",
			},
		},
	})
	return s
}

// DockerListImages returns a tool handler that lists Docker images.
func DockerListImages() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		all := helpers.GetBool(req.Arguments, "all")

		args := []string{"images"}
		if all {
			args = append(args, "--all")
		}

		output, err := docker.Run(ctx, args...)
		if err != nil {
			return helpers.ErrorResult("docker_error", err.Error()), nil
		}
		if output == "" {
			output = "No images found"
		}
		return helpers.TextResult(output), nil
	}
}
