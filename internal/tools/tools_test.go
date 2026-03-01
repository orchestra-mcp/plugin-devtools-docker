package tools

// Tests for the docker plugin tool handlers.
//
// Since all tools delegate to the `docker` CLI, we focus on:
//   1. Validation errors (missing required args) — no Docker needed.
//   2. docker_error paths — any tool with valid args will fail if Docker is
//      not available or returns a non-zero exit code. We test with a
//      nonexistent container ID ("no-such-container-xyz") which always fails.
//
// Tests that assert on exact success output are skipped in CI where Docker
// is not running. Use the DOCKER_AVAILABLE env var or check by calling
// docker.Run(ctx, "info") as a precondition.

import (
	"context"
	"os/exec"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// ---------- helpers ----------

func callTool(t *testing.T, handler func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error), args map[string]any) *pluginv1.ToolResponse {
	t.Helper()
	var s *structpb.Struct
	if args != nil {
		var err error
		s, err = structpb.NewStruct(args)
		if err != nil {
			t.Fatalf("NewStruct: %v", err)
		}
	}
	resp, err := handler(context.Background(), &pluginv1.ToolRequest{Arguments: s})
	if err != nil {
		t.Fatalf("handler returned Go error: %v", err)
	}
	return resp
}

func isError(resp *pluginv1.ToolResponse) bool {
	return resp != nil && !resp.Success
}

func getText(resp *pluginv1.ToolResponse) string {
	if resp == nil {
		return ""
	}
	if r := resp.GetResult(); r != nil {
		if f := r.GetFields(); f != nil {
			if tf, ok := f["text"]; ok {
				return tf.GetStringValue()
			}
		}
	}
	return ""
}

// dockerAvailable returns true if the docker binary is on PATH.
func dockerAvailable() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

// ---------- docker_list_containers ----------

func TestDockerListContainers_NoArgs(t *testing.T) {
	// No required args — should proceed to docker ps (may fail if no Docker).
	resp := callTool(t, DockerListContainers(), map[string]any{})
	// Either succeeds or returns docker_error — both are valid.
	_ = resp
}

func TestDockerListContainers_AllFlag(t *testing.T) {
	resp := callTool(t, DockerListContainers(), map[string]any{"all": true})
	_ = resp
}

func TestDockerListContainers_WithFormat(t *testing.T) {
	resp := callTool(t, DockerListContainers(), map[string]any{
		"format": "{{.Names}}",
	})
	_ = resp
}

// ---------- docker_start ----------

func TestDockerStart_MissingContainerID(t *testing.T) {
	resp := callTool(t, DockerStart(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing container_id")
	}
}

func TestDockerStart_NonexistentContainer(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	resp := callTool(t, DockerStart(), map[string]any{
		"container_id": "no-such-container-orchestra-test-xyz",
	})
	if !isError(resp) {
		t.Error("expected docker_error for nonexistent container")
	}
}

// ---------- docker_stop ----------

func TestDockerStop_MissingContainerID(t *testing.T) {
	resp := callTool(t, DockerStop(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing container_id")
	}
}

func TestDockerStop_NonexistentContainer(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	resp := callTool(t, DockerStop(), map[string]any{
		"container_id": "no-such-container-orchestra-test-xyz",
	})
	if !isError(resp) {
		t.Error("expected docker_error for nonexistent container")
	}
}

// ---------- docker_restart ----------

func TestDockerRestart_MissingContainerID(t *testing.T) {
	resp := callTool(t, DockerRestart(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing container_id")
	}
}

func TestDockerRestart_NonexistentContainer(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	resp := callTool(t, DockerRestart(), map[string]any{
		"container_id": "no-such-container-orchestra-test-xyz",
	})
	if !isError(resp) {
		t.Error("expected docker_error for nonexistent container")
	}
}

// ---------- docker_logs ----------

func TestDockerLogs_MissingContainerID(t *testing.T) {
	resp := callTool(t, DockerLogs(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing container_id")
	}
}

func TestDockerLogs_NonexistentContainer(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	resp := callTool(t, DockerLogs(), map[string]any{
		"container_id": "no-such-container-orchestra-test-xyz",
	})
	if !isError(resp) {
		t.Error("expected docker_error for nonexistent container")
	}
}

func TestDockerLogs_WithTailAndSince(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	resp := callTool(t, DockerLogs(), map[string]any{
		"container_id": "no-such-container-orchestra-test-xyz",
		"tail":         float64(50),
		"since":        "1h",
	})
	// Will return error because container doesn't exist, but args are valid.
	if !isError(resp) {
		t.Error("expected docker_error")
	}
}

// ---------- docker_exec ----------

func TestDockerExec_MissingContainerID(t *testing.T) {
	resp := callTool(t, DockerExec(), map[string]any{"command": "ls"})
	if !isError(resp) {
		t.Error("expected validation_error for missing container_id")
	}
}

func TestDockerExec_MissingCommand(t *testing.T) {
	resp := callTool(t, DockerExec(), map[string]any{"container_id": "mycontainer"})
	if !isError(resp) {
		t.Error("expected validation_error for missing command")
	}
}

func TestDockerExec_NonexistentContainer(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	resp := callTool(t, DockerExec(), map[string]any{
		"container_id": "no-such-container-orchestra-test-xyz",
		"command":      "echo hello",
	})
	if !isError(resp) {
		t.Error("expected docker_error for nonexistent container")
	}
}

// ---------- docker_list_images ----------

func TestDockerListImages_NoArgs(t *testing.T) {
	// No required args — may succeed or fail depending on Docker availability.
	resp := callTool(t, DockerListImages(), map[string]any{})
	_ = resp
}

// ---------- docker_compose_up ----------

func TestDockerComposeUp_MissingDirectory(t *testing.T) {
	resp := callTool(t, DockerComposeUp(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing directory")
	}
}

func TestDockerComposeUp_NonexistentDirectory(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	resp := callTool(t, DockerComposeUp(), map[string]any{
		"directory": "/tmp/no-such-compose-dir-orchestra-xyz",
	})
	if !isError(resp) {
		t.Error("expected docker_error for nonexistent directory")
	}
}

// ---------- docker_compose_down ----------

func TestDockerComposeDown_MissingDirectory(t *testing.T) {
	resp := callTool(t, DockerComposeDown(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing directory")
	}
}

// ---------- docker_inspect ----------

func TestDockerInspect_MissingContainerID(t *testing.T) {
	resp := callTool(t, DockerInspect(), map[string]any{})
	if !isError(resp) {
		t.Error("expected validation_error for missing container_id")
	}
}

func TestDockerInspect_NonexistentContainer(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	resp := callTool(t, DockerInspect(), map[string]any{
		"container_id": "no-such-container-orchestra-test-xyz",
	})
	if !isError(resp) {
		t.Error("expected docker_error for nonexistent container")
	}
}

func TestDockerInspect_WithFormat(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}
	resp := callTool(t, DockerInspect(), map[string]any{
		"container_id": "no-such-container-orchestra-test-xyz",
		"format":       "{{.State.Status}}",
	})
	if !isError(resp) {
		t.Error("expected docker_error")
	}
}
