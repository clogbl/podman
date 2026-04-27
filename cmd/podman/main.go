package main

import (
	"os"

	_ "github.com/containers/podman/v5/cmd/podman/completion"
	_ "github.com/containers/podman/v5/cmd/podman/containers"
	_ "github.com/containers/podman/v5/cmd/podman/generate"
	_ "github.com/containers/podman/v5/cmd/podman/healthcheck"
	_ "github.com/containers/podman/v5/cmd/podman/images"
	_ "github.com/containers/podman/v5/cmd/podman/manifest"
	_ "github.com/containers/podman/v5/cmd/podman/networks"
	_ "github.com/containers/podman/v5/cmd/podman/play"
	_ "github.com/containers/podman/v5/cmd/podman/pods"
	_ "github.com/containers/podman/v5/cmd/podman/secrets"
	_ "github.com/containers/podman/v5/cmd/podman/system"
	_ "github.com/containers/podman/v5/cmd/podman/volumes"
	"github.com/containers/podman/v5/cmd/podman/registry"
	"github.com/containers/podman/v5/pkg/rootless"
	"github.com/sirupsen/logrus"
)

func main() {
	// rootless reexec must happen before anything else
	// to ensure the process is running in the correct user namespace
	if reexec := rootless.TryReexecRootless(); reexec {
		os.Exit(0)
	}

	if err := registry.PodmanConfig().ContainersConf.CheckCgroupsV1(); err != nil {
		logrus.Warnf("Failed to check cgroups v1: %v", err)
	}

	if err := registry.Execute(); err != nil {
		if registry.GetExitCode() == 0 {
			registry.SetExitCode(registry.ExecErrorCodeGeneric)
		}
		// Log at Error level so the message is visible by default without
		// needing --log-level=debug. Using %v (not %w) intentionally since
		// we are only logging, not wrapping.
		logrus.Errorf("%v", err)
	}

	os.Exit(registry.GetExitCode())
}
