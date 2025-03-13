package main

import (
	"context"
	"fmt"

	"dagger/gopkg/internal/dagger"
)

func (m *Gopkg) release(
	ctx context.Context,
	source *dagger.Directory,
	token *dagger.Secret,
	dryRun bool,
) (string, error) {
	releaseCommand := "--dry-run"
	if !dryRun {
		releaseCommand = "--no-ci"
	}

	container := dag.Container().
		From("node:23.7.0-alpine").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		// Cache npm global packages
		WithMountedCache("/root/.npm", dag.CacheVolume("node-23")).
		// Cache project dependencies
		WithMountedCache("/src/node_modules", dag.CacheVolume("node_modules-cache")).
		// Cache npx binary files
		WithMountedCache("/root/.cache", dag.CacheVolume("npx-cache")).
		// Cache Alpine package manager (optional, but can help speed up package installs)
		WithMountedCache("/var/cache/apk", dag.CacheVolume("apk-cache")).
		// Inject secret for GitHub authentication
		WithSecretVariable("GITHUB_TOKEN", token).
		WithExec([]string{"npm", "install", "--save-dev", "@semantic-release/git"}).
		WithExec([]string{"npm", "install", "--save-dev", "@semantic-release/changelog"}).
		WithExec([]string{"npm", "install", "--save-dev", "conventional-changelog-conventionalcommits"}).
		WithExec([]string{"apk", "add", "git", "bash"})

	if m.GpgConfig != nil {
		container = m.setUpGpg(container.WithExec([]string{"apk", "add", "gnupg"}))
	}

	return container.WithExec(
		[]string{"bash", "-c", fmt.Sprintf("npx semantic-release %s --debug", releaseCommand)},
	).Stdout(ctx)
}

func (m *Gopkg) setUpGpg(container *dagger.Container) *dagger.Container {
	container = container.
		WithSecretVariable("GPG_KEY", m.GpgConfig.GpgKey).
		WithSecretVariable("GPG_KEY_ID", m.GpgConfig.GpgKeyId).
		WithEnvVariable("GIT_AUTHOR_NAME", m.GpgConfig.GitAuthorName).
		WithEnvVariable("GIT_AUTHOR_EMAIL", m.GpgConfig.GitAuthorEmail).
		WithEnvVariable("GIT_COMMITTER_NAME", m.GpgConfig.GitAuthorName).
		WithEnvVariable("GIT_COMMITTER_EMAIL", m.GpgConfig.GitAuthorEmail).
		WithExec([]string{"apk", "add", "gnupg", "bash", "git"}).
		WithExec([]string{"bash", "-c", "gpg-agent --daemon"}). // ensure gpg-agent is running
		WithExec([]string{"bash", "-c", "git config commit.gpgsign true"}).
		WithExec([]string{"bash", "-c", "git config --global user.signingkey \"$GPG_KEY_ID\""})

	return m.loadGpgKey(container)
}

func (m *Gopkg) loadGpgKey(container *dagger.Container) *dagger.Container {
	if m.GpgConfig.GpgPassword != nil {
		container = container.WithSecretVariable("GPG_PASSPHRASE", m.GpgConfig.GpgPassword).
			WithExec([]string{"bash", "-c", "echo \"$GPG_KEY\" > /tmp/gpg.key"}).
			WithExec([]string{"bash", "-c", "echo \"$GPG_PASSPHRASE\" | gpg --import --batch --yes --passphrase-fd 0 /tmp/gpg.key"}).
			WithExec([]string{"bash", "-c", "echo 'gpg --passphrase ${GPG_PASSPHRASE} --batch --yes --pinentry-mode=loopback --no-tty \"$@\"' > /tmp/gpg-with-passphrase && chmod +x /tmp/gpg-with-passphrase"}).
			WithExec([]string{"bash", "-c", "git config gpg.program \"/tmp/gpg-with-passphrase\""})
	} else {
		container = container.WithExec([]string{"bash", "-c", "echo \"$GPG_KEY\" | gpg --import"})
	}
	return container.WithExec(
		[]string{
			"bash",
			"-c",
			"printf \"trust\n5\ny\nquit\n\" | gpg --batch --command-fd 0 --edit-key \"$GPG_KEY_ID\"",
		},
	)
}
