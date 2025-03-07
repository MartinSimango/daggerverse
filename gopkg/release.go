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
	gpgKey *dagger.Secret,
	pass *dagger.Secret,
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
		WithMountedCache("/root/.npm", dag.CacheVolume("node-23")).
		WithMountedCache("/src/node_modules", dag.CacheVolume("node_modules-cache")).
		WithSecretVariable("GITHUB_TOKEN", token).
		WithExec([]string{"npm", "install", "--save-dev", "@semantic-release/git"}).
		WithExec([]string{"npm", "install", "--save-dev", "@semantic-release/changelog"}).
		WithExec([]string{"npm", "install", "--save-dev", "conventional-changelog-conventionalcommits"})

	if gpgKey != nil {
		container = container.
			WithSecretVariable("GPG_KEY", gpgKey).
			WithEnvVariable("GIT_EMAIL", "shukomango@gmail.com").
			WithEnvVariable("GIT_USERNAME", "MartinSimango").
			WithEnvVariable("GIT_AUTHOR_NAME", "Martin Simango").
			WithEnvVariable("GIT_AUTHOR_EMAIL", "shukomango@gmail.com").
			WithEnvVariable("GIT_COMMITTER_NAME", "Martin Simango").
			WithEnvVariable("GIT_COMMITTER_EMAIL", "shukomango@gmail.com").
			// WithSecretVariable("PASS", pass).
			WithExec([]string{"apt-get", "install", "-y", "gnupg"}).
			// WithExec([]string{"bash", "-c", "echo \"$GPG_KEY\" > /tmp/gpg.key"}).
			// WithExec([]string{"bash", "-c", "echo \"$PASS\" | gpg --import --batch --yes --passphrase-fd 0 /tmp/gpg.key"}).
			WithExec([]string{"bash", "-c", "echo \"$GPG_KEY\" | gpg --import"}).
			WithExec([]string{"bash", "-c", "git config commit.gpgsign true"}).
			WithExec([]string{"bash", "-c", "git config tag.gpgsign true"}).
			// WithExec([]string{"bash", "-c", "git config --global user.name \"Martin Simango\""}).
			// WithExec([]string{"bash", "-c", "git config --global user.email \"shukomango@gmail.com\""}).
			WithExec([]string{"bash", "-c", "printf \"trust\n5\ny\nquit\n\" | gpg --batch --command-fd 0 --edit-key 60BEEE74E301083F"}).
			WithExec([]string{"bash", "-c", fmt.Sprintf("git config --global user.signingkey %s", "60BEEE74E301083F")}).
			WithExec([]string{"bash", "-c", "git config --global gpg.program gpg"})
		// WithExec([]string{"bash", "-c", "git commit -a -m \"Test\""})
	}

	return container.
		WithExec(
			[]string{"bash", "-c", fmt.Sprintf("npx semantic-release %s --debug", releaseCommand)},
		).
		// WithExec([]string{"bash", "-c", "git log --show-signature"}).
		Stdout(ctx)
}
