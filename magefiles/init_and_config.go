package main

import "github.com/magefile/mage/mg"

func Init() error {
	// This is a workaround to prevent overriding local binaries
	if !isCIRunner() {
		return nil
	}

	mg.Deps(
		mg.F(Dep.Install, "github.com/jacobbrewer1/goschema@latest"),
		mg.F(Dep.Install, "golang.org/x/tools/cmd/goimports@latest"),
		mg.F(Dep.Install, "github.com/bazelbuild/bazelisk@latest"),
	)
	return nil
}
