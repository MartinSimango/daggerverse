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

	"github.com/MartinSimango/daggerverse/gopkg/internal/dagger"
)

type GitGpgConfig struct {
	GpgKey         *dagger.Secret
	GpgKeyId       *dagger.Secret
	GpgPassword    *dagger.Secret
	GitAuthorName  string
	GitAuthorEmail string
}

type Gopkg struct {
	GpgConfig *GitGpgConfig
}

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
	// +optional
	// +default=true
	// Dry run the release
	dryRun bool,
) (string, error) {
	return m.release(ctx, source, token, dryRun)
}

/*
GopkgFlow runs a release flow for a Go project using semantic-release
This runs the following steps:
1. Test the project
2. Release a new version of the project using semantic-release
*/
func (m *Gopkg) GopkgFlow(
	ctx context.Context,
	// Source directory for the project
	// +defaultPath="."
	source *dagger.Directory,
	// GitHub token for the release
	token *dagger.Secret,
	// +optional
	// +default=true
	// Dry run the release
	dryRun bool,
) (string, error) {
	_, err := m.test(ctx, source)
	if err != nil {
		return "", err
	}
	return m.release(ctx, source, token, dryRun)
}

// WithGitGpgConfig sets the GPG configuration for the git repository. To be used with Release.
func (m *Gopkg) WithGitGpgConfig(
	// GPG private key for signing the semantic-release bot
	gpgKey *dagger.Secret,
	// GPG key ID
	gpgKeyId *dagger.Secret,
	// GPG password for the private key, if any
	// +optional
	gpgPassword *dagger.Secret,
	// Git author name for the release commits
	gitAuthorName string,
	// Git author email for the release commits
	gitAuthorEmail string,
) *Gopkg {
	m.GpgConfig = &GitGpgConfig{
		GpgKey:         gpgKey,
		GpgKeyId:       gpgKeyId,
		GpgPassword:    gpgPassword,
		GitAuthorName:  gitAuthorName,
		GitAuthorEmail: gitAuthorEmail,
	}

	return m
}
