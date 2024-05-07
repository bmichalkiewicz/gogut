package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"dagger.io/dagger/dag"
)

var ParentDir string

func main() {
	ParentDir, err := parentDir()
	if err != nil {
		log.Fatalf("error: %s", err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	workspace := dag.Host().Directory(ParentDir)

	defer dag.Close()

	if err := qualityControl(ctx, workspace); err != nil {
		log.Fatalf("error: %s", err.Error())
		os.Exit(1)
	}

	if err := build(ctx, workspace); err != nil {
		log.Fatalf("error: %s", err.Error())
		os.Exit(1)
	}
}

func build(ctx context.Context, workspace *dagger.Directory) error {
	// Initialize the build container
	builder := dag.Container().
		From("golang:1.21.7").
		WithMountedDirectory("/src", workspace).WithWorkdir("/src").
		Pipeline("Build")

	path := "binary/"

	// Commands
	builder = builder.WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", path + "gogut"})

	// Get reference to build output directory in container
	outputDir := builder.Directory(path)

	// Write contents of container build/ directory to the host
	_, err := outputDir.Export(ctx, filepath.Join(ParentDir, "binary"))
	if err != nil {
		return fmt.Errorf("error exporting output directory: %v", err)
	}

	return nil
}

func qualityControl(ctx context.Context, workspace *dagger.Directory) error {
	// Initialize the linter container
	linter := dag.Container().
		From("golang:1.21.7").
		WithMountedDirectory("/src", workspace).WithWorkdir("/src").
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}).
		Pipeline("Lint")

	out, err := linter.
		WithExec([]string{"golangci-lint", "run", "./..."}).
		WithExec([]string{"go", "mod", "verify"}).
		WithExec([]string{"go", "vet", "./..."}).
		WithExec([]string{"go", "run", "honnef.co/go/tools/cmd/staticcheck@latest", "-checks=all,-ST1000,-U1000", "./..."}).
		WithExec([]string{"go", "test", "-race", "-buildvcs", "-vet=off", "./..."}).
		Stdout(ctx)

	if err != nil {
		fmt.Println(out)
		return fmt.Errorf("problem with checking code: %v", err)
	}

	return nil
}

func parentDir() (string, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return rootDir, nil
}
