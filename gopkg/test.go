package main

import (
	"context"

	"dagger/gopkg/internal/dagger"
)

func (m *Gopkg) test(ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	return dag.Container().
		From("golang:1.24").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-124")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-124")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"go", "test", "./..."}).
		Stdout(ctx)
}
