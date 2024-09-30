// A generated module for Buf functions
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

	// _, err := source.File(config).ID(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("no buf.yaml file found")
	// }

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

func (b *Buf) Container() *dagger.Container {
	ctr := dag.
		Container().
		From("bufbuild/buf:latest").
		WithWorkdir(WorkDir).
		WithMountedDirectory(WorkDir, b.Source)

	for _, p := range b.Packages {
		ctr.WithExec([]string{"go", "install", p})
	}

	return ctr
}

func (b *Buf) Lint() *dagger.Container {
	return b.
		Container().
		WithExec([]string{"buf", "lint", "--config", b.Config})
}
