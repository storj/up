//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

//Test executes all unit and integration tests
func Test() error {
	err := sh.RunV("go", "test", "./...")
	return err
}

//Test executes all unit and integration tests
func Coverage() error {
	fmt.Println("Executing tests and generate coverate information")
	err := sh.RunV("go", "test", "-coverprofile=/tmp/coverage.out", "./...")
	if err != nil {
		return err
	}
	return sh.RunV("go", "tool", "cover", "-html=/tmp/coverage.out", "-o", "coverage.html")
}

//Lint executes all the linters with golangci-lint
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

//Format reformat code automatically
func Format() error {
	err := sh.RunV("gofmt", "-w", ".")
	if err != nil {
		return err
	}
	return sh.RunV("goimports", "-w", ".")

}

//GenBuild re-generates `./build` helper binary
func GenBuild() error {
	envs := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        "linux",
		"GOARCH":      "amd64",
	}
	return sh.RunWithV(envs, "mage", "-compile", "build")

}

func DockerBaseBuild() error {
	err := sh.RunV("docker", "build", "-t", "ghcr.io/elek/storj-base", "-f", "base.Dockerfile", ".")
	if err != nil {
		return err
	}
	return nil
}

func DockerBuildBuild() error {
	err := sh.RunV(
		"docker",
		"build",
		"-t", "ghcr.io/elek/storj-build",
		"-f", "build.Dockerfile", ".")
	if err != nil {
		return err
	}
	return nil
}

func DockerCoreBuild(version string) error {
	version = "1.39.6"
	mg.Deps(DockerBaseBuild)
	mg.Deps(DockerBuildBuild)
	err := sh.RunV("docker",
		"build",
		"-t", "ghcr.io/elek/storj:"+version,
		"--build-arg", "BRANCH=v"+version,
		"-f", "pkg/storj.Dockerfile", ".")
	if err != nil {
		return err
	}
	return nil
}

func DockerEdgeBuild(version string) error {
	version = "1.14.0"
	mg.Deps(DockerBaseBuild)
	mg.Deps(DockerBuildBuild)
	err := sh.RunV("docker",
		"build",
		"-t", "ghcr.io/elek/storj-edge:"+version,
		"--build-arg", "BRANCH=v"+version,
		"-f", "pkg/edge.Dockerfile", ".")
	if err != nil {
		return err
	}
	return nil
}

func Integration() error {
	return sh.RunV("test/test.sh")
}

func Publish() error {
	coreVersion := "1.39.6"
	edgeVersion := "1.14.0"
	err := DockerCoreBuild(coreVersion)
	if err != nil {
		return err
	}

	err = DockerEdgeBuild(edgeVersion)
	if err != nil {
		return err
	}

	err = Integration()
	if err != nil {
		return err
	}

	err = DockerBasePublish()
	if err != nil {
		return err
	}
	err = DockerBuildPublish()
	if err != nil {
		return err
	}
	err = DockerCorePublish(coreVersion)
	if err != nil {
		return err
	}
	err = DockerEdgePublish(edgeVersion)
	if err != nil {
		return err
	}
	return nil
}

func DockerCorePublish(version string) error {
	return sh.RunV("docker", "push", "ghcr.io/elek/storj:"+version)
}

func DockerEdgePublish(version string) error {
	return sh.RunV("docker", "push", "ghcr.io/elek/storj:"+version)
}

func DockerBuildPublish() error {
	return sh.RunV("docker", "push", "ghcr.io/elek/storj-build")
}

func DockerBasePublish() error {
	return sh.RunV("docker", "push", "ghcr.io/elek/storj-base")
}
