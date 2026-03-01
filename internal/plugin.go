package internal

import (
	"github.com/orchestra-mcp/plugin-devtools-docker/internal/tools"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// ToolsPlugin registers all Docker and Docker Compose tools.
type ToolsPlugin struct{}

// RegisterTools registers all 10 tools with the plugin builder.
func (tp *ToolsPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	// Docker container operations (7 tools)
	builder.RegisterTool("docker_list_containers",
		"List Docker containers (use --all to include stopped)",
		tools.DockerListContainersSchema(), tools.DockerListContainers())

	builder.RegisterTool("docker_start",
		"Start a stopped Docker container",
		tools.DockerStartSchema(), tools.DockerStart())

	builder.RegisterTool("docker_stop",
		"Stop a running Docker container",
		tools.DockerStopSchema(), tools.DockerStop())

	builder.RegisterTool("docker_restart",
		"Restart a Docker container",
		tools.DockerRestartSchema(), tools.DockerRestart())

	builder.RegisterTool("docker_logs",
		"Get logs from a Docker container",
		tools.DockerLogsSchema(), tools.DockerLogs())

	builder.RegisterTool("docker_exec",
		"Execute a command inside a running Docker container",
		tools.DockerExecSchema(), tools.DockerExec())

	builder.RegisterTool("docker_inspect",
		"Inspect a Docker container or image (detailed JSON metadata)",
		tools.DockerInspectSchema(), tools.DockerInspect())

	// Docker image operations (1 tool)
	builder.RegisterTool("docker_list_images",
		"List Docker images",
		tools.DockerListImagesSchema(), tools.DockerListImages())

	// Docker Compose operations (2 tools)
	builder.RegisterTool("docker_compose_up",
		"Start services defined in docker-compose.yml",
		tools.DockerComposeUpSchema(), tools.DockerComposeUp())

	builder.RegisterTool("docker_compose_down",
		"Stop and remove containers defined in docker-compose.yml",
		tools.DockerComposeDownSchema(), tools.DockerComposeDown())
}
