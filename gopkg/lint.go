package main

import (
	"context"

	"github.com/MartinSimango/daggerverse/gopkg/internal/dagger"
)

func (m *Gopkg) lint(ctx context.Context, source *dagger.Directory, config string) (string, error) {
	return dag.Container().From("golang:1.24").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}).
		WithExec([]string{"golangci-lint", "run", "--config", config}).
		Stdout(ctx)
}
