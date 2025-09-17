package main

import (
	"context"
	"fmt"

	"dagger/kcd/internal/dagger"
)

var (
	GH_REPO       = "https://github.com/NunoFrRibeiro/KCD-Porto-2025"
	COUNTER_IMAGE = "nunofilribeiro/counterbackend:v0.1.0"
	ADDER_IMAGE   = "nunofilribeiro/adderbackend:v0.1.0"
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

	return fmt.Sprintf("lint result: %s\ntest result: %s\n", lintResult, testResult), nil
}

// TODO: Possible to show or not
//
// // Deploys the docker images to a registry
// func (m *Kcd) Deploy(
// 	ctx context.Context,
// 	// Infisical Auth Client ID
// 	// +required
// 	infisicalId *dagger.Secret,
// 	// Infisical Auth Client Secret
// 	// +required
// 	infisicalSecret *dagger.Secret,
// 	// Infisical Project to fetch secrets
// 	// +required
// 	infisicalProject string,
// ) (string, error) {
// 	var result string
// 	if infisicalId != nil && infisicalProject != "" {
// 		registryUser, err := dag.Infisical(infisicalId, infisicalSecret).
// 			GetSecret("DH_USER", infisicalProject, "dev", dagger.InfisicalGetSecretOpts{
// 				SecretPath: "/",
// 			}).
// 			Plaintext(ctx)
// 		if err != nil {
// 			return "", err
// 		}
//
// 		registryPass := dag.Infisical(infisicalId, infisicalSecret).
// 			GetSecret("DH_PASS", infisicalProject, "dev", dagger.InfisicalGetSecretOpts{
// 				SecretPath: "/",
// 			})
//
// 		counterImage := m.Build.Container(
// 			m.Source.Directory("CounterBackend"),
// 			8081,
// 			"CounterBackend",
// 		)
// 		counterResult, err := dag.Container().
// 			WithRegistryAuth(DH_REPO, registryUser, registryPass).
// 			Publish(ctx, COUNTER_IMAGE, dagger.ContainerPublishOpts{
// 				PlatformVariants: []*dagger.Container{
// 					counterImage,
// 				},
// 			})
//
// 		adderImage := m.Build.Container(m.Source.Directory("AdderBackend"), 8080, "AdderBackend")
// 		adderResult, err := dag.Container().
// 			WithRegistryAuth(DH_REPO, registryUser, registryPass).
// 			Publish(ctx, ADDER_IMAGE, dagger.ContainerPublishOpts{
// 				PlatformVariants: []*dagger.Container{
// 					adderImage,
// 				},
// 			})
//
// 		result = counterResult + "\n" + adderResult
//
// 		return result, nil
// 	}
//
// 	return result, nil
// }
