package main

import (
	"context"
	"dagger/kcd/internal/dagger"
	"fmt"
)

// Runs GolangCILint for all sources
func (m *Kcd) Lint(
	ctx context.Context,
) (string, error) {
	adderResult, error := m.Build.Lint(ctx, m.Source.Directory("AdderBackend"))
	if error != nil {
		return "", error
	}

	counterResult, error := m.Build.Lint(ctx, m.Source.Directory("CounterBackend"))
	if error != nil {
		return "", error
	}

	result := adderResult + "\n" + counterResult
	return result, nil
}

// Runs all tests
func (m *Kcd) Test(
	ctx context.Context,
) (string, error) {
	adderResult, err := m.Build.UnitTests(ctx, m.Source.Directory("AdderBackend"))
	if err != nil {
		return "", err
	}

	counterResult, err := m.Build.UnitTests(ctx, m.Source.Directory("CounterBackend"))
	if err != nil {
		return "", err
	}

	result := adderResult + "\n" + counterResult
	return result, nil
}

// Run Trivy container scan
func (m *Kcd) RunTrivy(
	ctx context.Context,
) (string, error) {
	adderContainer := m.Build.Container(m.Source.Directory("AdderBackend"), 8080, "AdderBackend")
	adderResult, err := m.Build.RunTrivy(ctx, adderContainer)
	if err != nil {
		return "", nil
	}

	counterContainer := m.Build.Container(m.Source.Directory("CounterBackend"), 8081, "CounterBackend")
	counterResult, err := m.Build.RunTrivy(ctx, counterContainer)
	if err != nil {
		return "", nil
	}

	result := adderResult + "\n" + counterResult
	return result, nil
}

// Run ci-check
func (m *Kcd) Check(
	ctx context.Context,
	// Token with permissions to comment on PR
	githubToken *dagger.Secret,
	// GitHub git commit
	commit string,
	// LLM model used to debug tests
	// *optional
	// +default="gemini-2.0-flash"
	model string,
) (string, error) {
	lintResult, err := m.Lint(ctx)
	if err != nil {
		if githubToken != nil {
			debugPr := m.DebugPR(ctx, githubToken, commit, model)
			return "", fmt.Errorf("failed to lint.\nrunning debugger for %v %v", err, debugPr)
		}
		return "", err
	}

	testResult, err := m.Test(ctx)
	if err != nil {
		if githubToken != nil {
			debugPr := m.DebugPR(ctx, githubToken, commit, model)
			return "", fmt.Errorf("failed to test.\nrunning debugger for %v %v", err, debugPr)
		}
		return "", err
	}

	trivyResult, err := m.RunTrivy(ctx)
	if err != nil {
		if githubToken != nil {
			debugPr := m.DebugPR(ctx, githubToken, commit, model)
			return "", fmt.Errorf("failed to run trivy scan.\nrunning debugger for %v %v", err, debugPr)
		}
		return "", err
	}

	return fmt.Sprintf("lint result:\n%s\ntest result:\n%s\ntrivy scan result:\n%s\n", lintResult, testResult, trivyResult), nil
}
