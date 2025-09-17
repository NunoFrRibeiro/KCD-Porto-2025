package main

import (
	"context"
	"fmt"

	"dagger/build/internal/dagger"
)

type Build struct {
	Source *dagger.Directory
}

func New(
	source *dagger.Directory,
) *Build {
	return &Build{
		Source: source,
	}
}

// Run the projects unit tests
func (m *Build) UnitTests(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	return dag.Golang().
		WithSource(source).
		Test(ctx)
}

// Runs GolangCILint against the source
func (m *Build) Lint(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	return dag.Golang().
		WithSource(source).
		GolangciLint(ctx)
}

// Formatter
func (m *Build) Format() *dagger.Directory {
	return dag.Golang().
		WithSource(m.Source).
		Fmt().
		GolangciLintFix()
}

// Checker
func (m *Build) Check(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	lint, err := m.Lint(ctx, source)
	if err != nil {
		return "", err
	}
	test, err := m.UnitTests(ctx, source)
	if err != nil {
		return "", fmt.Errorf("Error is: %v", err)
	}
	return "Lint result: " + lint + "\n" + "Test result: " + test, nil
}

// Builds the source binary
func (m *Build) Build(
	source *dagger.Directory,
) *dagger.Directory {
	return dag.Golang().Build([]string{}, dagger.GolangBuildOpts{
		Source: source,
	})
}

// Returns the source binary
func (m *Build) Binary(
	source *dagger.Directory,
	binaryName string,
) *dagger.File {
	binary := m.Build(source)
	return binary.File(binaryName)
}

// Returns a container with the built binary
func (m *Build) Container(
	source *dagger.Directory,
	// Port to open on container
	// +required
	port int,
	binaryName string,
) *dagger.Container {
	binary := m.Binary(source, binaryName)
	binaryStr := fmt.Sprintf("/bin/%s", "service")

	return dag.Container().
		From("ubuntu:24.10").
		WithFile(binaryStr, binary).
		WithEntrypoint([]string{binaryStr}).
		WithExposedPort(port)
}

// Stateless checker
func (m *Build) CheckDirectory(
	ctx context.Context,
	// Directory to run checks on
	source *dagger.Directory,
) (string, error) {
	// m.Source = source
	return m.Check(ctx, source)
}

// Stateless formatter
func (m *Build) FormatDirectory(
	// Directory to format
	source *dagger.Directory,
) *dagger.Directory {
	m.Source = source
	return m.Format()
}

// Stateless formatter
func (m *Build) FormatFile(
	// Directory with go module
	source *dagger.Directory,
	// File path to format
	path string,
) *dagger.Directory {
	return dag.
		Container().
		From("golang:1.24").
		WithExec([]string{"go", "install", "golang.org/x/tools/gopls@latest"}).
		WithWorkdir("/app").
		WithDirectory("/app", source).
		WithExec([]string{"gopls", "format", "-w", path}).
		WithExec([]string{"gopls", "imports", "-w", path}).
		Directory("/app")
}
