package devtoolsdocker

import (
	"github.com/orchestra-mcp/plugin-devtools-docker/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all Docker tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	tp := &internal.ToolsPlugin{}
	tp.RegisterTools(builder)
}
