// A generated module for Gopkg functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"

	"dagger/gopkg/internal/dagger"
)

type Gopkg struct{}

// Test the project
func (m *Gopkg) Test(ctx context.Context,
	// Source directory for the project
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	return m.test(ctx, source)
}

// Release a new version of the project using semantic-release
func (m *Gopkg) Release(
	ctx context.Context,
	// Source directory for the project
	// +defaultPath="."
	source *dagger.Directory,
	// GitHub token for the release
	token *dagger.Secret,
	// GPG private key for signing the semantic-release bot
	// +optional
	gpgKey *dagger.Secret,
	// +optional
	// +default=true
	// Dry run the release
	dryRun bool,
) (string, error) {
	return m.release(ctx, source, token, gpgKey, dryRun)
}
