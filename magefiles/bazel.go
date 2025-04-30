package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/magefile/mage/sh"
)

// isCIRunner determines if the current process is running in a Continuous Integration (CI) environment.
// It checks the values of the "CI" and "GITHUB_ACTIONS" environment variables.
//
// GitHub Documentation: https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/store-information-in-variables#default-environment-variables
var isCIRunner = sync.OnceValue(func() bool {
	// Check if running in CI
	ciRunner, _ := strconv.ParseBool(os.Getenv("CI"))
	githubRunner, _ := strconv.ParseBool(os.Getenv("GITHUB_ACTIONS"))
	return ciRunner || githubRunner
})

// binDirectory returns the absolute path to the root bazel-bin directory.
var binDirectory = sync.OnceValue(func() string {
	dir, err := sh.Output("bazel", "info", "bazel-bin")
	if err != nil {
		panic(err)
	}

	return dir
})

// outputDirectory returns the absolute path to the root output path.
var outputDirectory = sync.OnceValue(func() string {
	dir, err := sh.Output("bazel", "info", "output_path")
	if err != nil {
		panic(err)
	}

	return dir
})

// hostArch returns the host CPU architecture. Use this instead of runtime.GOARCH as the latter succumbs to cross-compilation, so can lie.
var hostArch = sync.OnceValue(func() string {
	out, err := sh.Output("uname", "-m")
	if err != nil {
		panic(err)
	}

	switch out {
	case "arm64":
		return "arm64"
	case "x86_64":
		return "amd64"
	default:
		panic(fmt.Sprintf("unsupported host arch: %s", out))
	}
})
