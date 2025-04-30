//go:build mage

package main

import "github.com/magefile/mage/mg"

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
