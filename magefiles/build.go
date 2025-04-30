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
		"--platforms", fmt.Sprintf("@io_bazel_rules_go//go/toolchain:linux_%s", hostArch()),
		"//...",
	}

	if err := sh.Run("bazel", args...); err != nil {
		return fmt.Errorf("error building all code: %w", err)
	}

	fmt.Printf("[INFO] Build completed in %s\n", time.Since(start))
	return nil
}
