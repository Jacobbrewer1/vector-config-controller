package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var containerRegistry = sync.OnceValue(func() string {
	// Check if there's an override registry set in the environment
	overrideRegistry := os.Getenv("CONTAINER_REGISTRY")
	if overrideRegistry != "" {
		return overrideRegistry
	}

	// Default to the GitHub Container Registry
	repoPath := os.Getenv("GITHUB_REPOSITORY")
	if repoPath == "" {
		panic("GITHUB_REPOSITORY environment variable is not set")
	}

	repoPath = strings.ToLower(repoPath)

	// repo path current format is "owner/repo"

	// Append the GitHub container registry URL
	return fmt.Sprintf("ghcr.io/%s", repoPath)
})

var commitTag = sync.OnceValue(func() string {
	got, err := sh.Output("git", "describe", "--tags", "--abbrev=0")
	if err == nil && got != "" {
		return strings.TrimSpace(got)
	}

	got, err = sh.Output("git", "describe", "--tags")
	if err == nil && got != "" {
		return strings.TrimSpace(got)
	}

	// Fallback to the commit hash if no tags are found
	got, err = sh.Output("git", "rev-parse", "--short", "HEAD")
	if err == nil && got != "" {
		return strings.TrimSpace(got)
	}

	panic("could not determine git tag")
})

type Push mg.Namespace

func (Push) All() error {
	mg.Deps(Build.All)

	// Grab all the targets from the cmd directory
	targets, err := os.ReadDir("cmd")
	if err != nil {
		return fmt.Errorf("reading cmd directory: %w", err)
	}

	// Filter out the directories that are not services
	serviceTargets := make([]string, 0, len(targets))
	for _, target := range targets {
		if target.IsDir() {
			serviceTargets = append(serviceTargets, target.Name())
		}
	}

	if err := pushImages(serviceTargets, commitTag()); err != nil {
		return fmt.Errorf("pushing images: %w", err)
	}

	return nil
}

func pushImages(targets []string, tag string) error {
	errChan := make(chan error, 1)
	wg := new(sync.WaitGroup)
	wg.Add(len(targets))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, target := range targets {
		go func(target string) {
			defer wg.Done()

			fmt.Println("[INFO] Pushing image", target)
			if err := pushImage(ctx, "cmd/"+target, tag); err != nil {
				errChan <- fmt.Errorf("pushing image %s: %w", target, err)
			}
		}(target)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	fmt.Println("[INFO] Watching error channel for push jobs")
	errs := make([]error, 0, len(targets))
	for err := range errChan {
		errs = append(errs, err)
	}

	fmt.Println("[INFO] Waiting for all push jobs to finish")
	wg.Wait()

	return errors.Join(errs...)
}

func pushImage(ctx context.Context, target string, tag string) error {
	tarballPath := filepath.Join(binDirectory(), target, "oci.tar", "tarball.tar")
	if err := sh.Run(
		"docker",
		"load",
		"--input", tarballPath,
	); err != nil {
		return fmt.Errorf("loading tarball %s to image: %w", tarballPath, err)
	}

	target = strings.TrimPrefix(target, "cmd/")

	// Is the context cancelled?
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Prevent blocking
	}

	fmt.Println("[INFO] Tagging images with ", tag)

	// Tag the image
	if err := sh.Run(
		"docker",
		"tag",
		target,
		fmt.Sprintf("%s/%s:%s", containerRegistry(), target, tag),
	); err != nil {
		return fmt.Errorf("tagging image %s: %w", target, err)
	}

	// Tag the image with the Latest tag
	if err := sh.Run(
		"docker",
		"tag",
		target,
		fmt.Sprintf("%s/%s:latest", containerRegistry(), target),
	); err != nil {
		return fmt.Errorf("tagging image %s with latest tag: %w", target, err)
	}

	// Push the image
	if err := sh.Run(
		"docker",
		"push",
		fmt.Sprintf("%s/%s:%s", containerRegistry(), target, tag),
	); err != nil {
		return fmt.Errorf("pushing image %s: %w", target, err)
	}

	// Push the image with the Latest tag
	if err := sh.Run(
		"docker",
		"push",
		fmt.Sprintf("%s/%s:latest", containerRegistry(), target),
	); err != nil {
		return fmt.Errorf("pushing image %s with latest tag: %w", target, err)
	}

	return nil
}
