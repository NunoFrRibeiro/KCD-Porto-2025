package main

import (
	"dagger/kcd/internal/dagger"
)

var REPO = "https://github.com/NunoFrRibeiro/KCD-Porto-2025"

type Kcd struct {
	// Project Source Directory
	// +private
	Source *dagger.Directory
	// +private
	Build *dagger.Build
}

func New(
	// Project Source Directory
	// +defaultPath="/"
	// +optional
	// +ignore=[".github", "tmp"]
	source *dagger.Directory,

	// Checkout the repository (at the designated ref) and use it as the source directory instead of the local one.
	// +optional
	ref string,
) (*Kcd, error) {
	if source == nil && ref != "" {
		source = dag.Git(REPO).
			Ref(ref).
			Tree()
	}

	return &Kcd{
		Source: source,
		Build:  dag.Build(source),
	}, nil
}
