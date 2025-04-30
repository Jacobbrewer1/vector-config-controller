package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

// All builds all applications
func (b Build) All() error {
	fmt.Println("[INFO] Building all code")

	start := time.Now()

	args := []string{
		"build",
		"--platforms", "@io_bazel_rules_go//go/toolchain:linux_" + hostArch(),
		"//...",
	}

	if err := sh.Run("bazel", args...); err != nil {
		return fmt.Errorf("error building all code: %w", err)
	}

	fmt.Printf("[INFO] Build completed in %s\n", time.Since(start))
	return nil
}

// One builds a single application
func (b Build) One(service string) error {
	fmt.Printf("[INFO] Building %s\n", service)

	start := time.Now()

	args := []string{
		"build",
		"--platforms", "@io_bazel_rules_go//go/toolchain:linux_" + hostArch(),
		"//cmd/" + service,
	}

	if err := sh.Run("bazel", args...); err != nil {
		return fmt.Errorf("error building %s: %w", service, err)
	}

	fmt.Printf("[INFO] Build completed in %s\n", time.Since(start))
	return nil
}
