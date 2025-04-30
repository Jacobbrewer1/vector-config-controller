//go:build mage

package main

import (
	"os"
	"strconv"
	"sync"

	"github.com/magefile/mage/mg"
)

// isCIRunner checks debug mode is enabled
var isDebugMode = sync.OnceValue(func() bool {
	gotStr := os.Getenv("RUNNER_DEBUG")
	got, _ := strconv.ParseBool(gotStr)

	// Enable debug logging if the RUNNER_DEBUG environment variable is set to true or if not running in a CI environment.
	return got || !isCIRunner()
})

// Init initializes the mage environment
func Init() error {
	// This is a workaround to prevent overriding local binaries
	if !isCIRunner() {
		return nil
	}

	mg.Deps(
		mg.F(Dep.Install, "github.com/bazelbuild/bazelisk@latest"),
	)
	return nil
}
