package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/docker"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// DockerComposeDownSchema returns the JSON Schema for the docker_compose_down tool.
func DockerComposeDownSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"directory": map[string]any{
				"type":        "string",
				"description": "Directory containing docker-compose.yml",
			},
			"volumes": map[string]any{
				"type":        "boolean",
				"description": "Remove named volumes declared in the volumes section",
			},
			"remove_orphans": map[string]any{
				"type":        "boolean",
				"description": "Remove containers for services not defined in the compose file",
			},
		},
		"required": []any{"directory"},
	})
	return s
}

// DockerComposeDown returns a tool handler that runs docker compose down.
func DockerComposeDown() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "directory"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}
		directory := helpers.GetString(req.Arguments, "directory")
		volumes := helpers.GetBool(req.Arguments, "volumes")
		removeOrphans := helpers.GetBool(req.Arguments, "remove_orphans")

		args := []string{"down"}
		if volumes {
			args = append(args, "--volumes")
		}
		if removeOrphans {
			args = append(args, "--remove-orphans")
		}

		output, err := docker.Compose(ctx, directory, args...)
		if err != nil {
			return helpers.ErrorResult("docker_error", err.Error()), nil
		}
		return helpers.TextResult(fmt.Sprintf("Compose down completed:\n%s", output)), nil
	}
}
