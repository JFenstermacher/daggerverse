package main

import (
	"context"
	"fmt"
	"github.com/JFenstermacher/daggerverse/buf/internal/dagger"
	"strings"
)

const (
	WorkDir = "/workspace"
)

// Buf project
type Buf struct {
	// Go packages to install
	// Each will be run with 'go install <package>'
	// The following packages will be installed by default:
	//   * github.com/bufbuild/buf/cmd/buf@latest
	//   * google.golang.org/protobuf/cmd/protoc-gen-go@latest
	//   * connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	Packages []string

	// Project source directory
	Source *dagger.Directory

	// Path to config file
	Config string
}

func New(
	// Context
	ctx context.Context,

	// Project source directory
	source *dagger.Directory,

	// Packages to install
	// +optional
	// +default=[]
	packages []string,

	// Path to buf.yaml
	// +optional
	// +default="buf.yaml"
	config string,
) (*Buf, error) {
	defPackages := map[string]string{
		"github.com/bufbuild/buf/cmd/buf":                  "latest",
		"google.golang.org/protobuf/cmd/protoc-gen-go":     "latest",
		"connectrpc.com/connect/cmd/protoc-gen-connect-go": "latest",
	}

	_, err := source.File(config).ID(ctx)
	if err != nil {
		return nil, fmt.Errorf("no buf.yaml file found")
	}

	for _, p := range packages {
		parts := strings.Split(p, "@")
		name := parts[0]

		_, ok := defPackages[name]

		if ok {
			delete(defPackages, name)
		}
	}

	for k, v := range defPackages {
		packages = append(packages, fmt.Sprintf("%s@%s", k, v))
	}

	return &Buf{
		Packages: packages,
		Source:   source,
		Config:   config,
	}, nil
}

// Buf container with packages installed
func (b *Buf) Container() *dagger.Container {
	ctr := dag.
		Container().
		From("golang:latest").
		WithWorkdir(WorkDir).
		WithMountedDirectory(WorkDir, b.Source)

	for _, p := range b.Packages {
		ctr = ctr.WithExec([]string{"go", "install", p})
	}

	return ctr
}

// Lint protobuf files
func (b *Buf) Lint() *dagger.Container {
	return b.
		Container().
		WithExec([]string{"buf", "lint", "--config", b.Config})
}

// Formats protobuf files
func (b *Buf) Format() *dagger.Directory {
	out := b.
		Container().
		WithExec([]string{"buf", "format", "--config", b.Config}).
		Directory(WorkDir)

	return b.Source.Diff(out)
}

// Generate services and clients based on buf.gen.yaml
func (b *Buf) Generate() *dagger.Directory {
	out := b.
		Container().
		WithExec([]string{"buf", "generate", "--config", b.Config}).
		Directory(WorkDir)

	return b.Source.Diff(out)
}
