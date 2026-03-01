package tools

import (
	"context"
	"fmt"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/docker"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DockerComposeUpSchema returns the JSON Schema for the docker_compose_up tool.
func DockerComposeUpSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"directory": map[string]any{
				"type":        "string",
				"description": "Directory containing docker-compose.yml",
			},
			"services": map[string]any{
				"type":        "string",
				"description": "Space-separated list of services to start (default: all)",
			},
			"detach": map[string]any{
				"type":        "boolean",
				"description": "Run containers in the background (default true)",
			},
			"build": map[string]any{
				"type":        "boolean",
				"description": "Build images before starting containers",
			},
		},
		"required": []any{"directory"},
	})
	return s
}

// DockerComposeUp returns a tool handler that runs docker compose up.
func DockerComposeUp() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "directory"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		directory := helpers.GetString(req.Arguments, "directory")
		services := helpers.GetString(req.Arguments, "services")
		detach := helpers.GetBool(req.Arguments, "detach")
		build := helpers.GetBool(req.Arguments, "build")

		args := []string{"up"}

		// Default to detached mode; only skip if explicitly set to false
		if detach || req.Arguments == nil || req.Arguments.Fields["detach"] == nil {
			args = append(args, "--detach")
		}
		if build {
			args = append(args, "--build")
		}
		if services != "" {
			for _, svc := range strings.Fields(services) {
				args = append(args, svc)
			}
		}

		output, err := docker.Compose(ctx, directory, args...)
		if err != nil {
			return helpers.ErrorResult("docker_error", err.Error()), nil
		}
		return helpers.TextResult(fmt.Sprintf("Compose up completed:\n%s", output)), nil
	}
}
