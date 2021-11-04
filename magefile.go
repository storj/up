//go:build mage
// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

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
	tag, err := getNextDockerTag("storj-build.last")
	if err != nil {
		return err
	}
	err = sh.RunV("docker", "build", "-t", "ghcr.io/elek/storj-base:"+tag, "-f", "base.Dockerfile", ".")
	if err != nil {
		return err
	}
	return nil
}

func DockerBuildBuild() error {
	tag, err := getNextDockerTag("storj-build.last")
	if err != nil {
		return err
	}
	err = sh.RunV(
		"docker",
		"build",
		"--build-arg", "BRANCH=main",
		"--build-arg", "TYPE=github",
		"--build-arg", "REPO=https://github.com/storj/storj.git",
		"-t", "ghcr.io/elek/storj-build:"+tag,
		"-f", "cmd/files/docker/build.Dockerfile", ".")
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

func dockerPushWithNextTag(image string) error {
	tagFile := fmt.Sprintf("%s.last", image)
	tag, err := getNextDockerTag(tagFile)
	if err != nil {
		return err
	}
	err = sh.RunV("docker", "push", fmt.Sprintf("ghcr.io/elek/%s:%s", image, tag))
	if err != nil {
		return err
	}
	return writeDockerTag(tagFile, tag)
}

func dockerPush(image string, tag string) error {
	err := sh.RunV("docker", "push", fmt.Sprintf("ghcr.io/elek/%s:%s", image, tag))
	if err != nil {
		return err
	}
	return err
}

func DockerCorePublish(version string) error {
	return dockerPush("storj", version)
}

func DockerEdgePublish(version string) error {
	return dockerPush("storj-edge", version)
}

func DockerBuildPublish() error {
	return dockerPushWithNextTag("storj-build")
}

func DockerBasePublish() error {
	return dockerPushWithNextTag("storj-base")
}

// getNextDockerTag generates docker tag with the pattern yyyymmdd-n.
//last used tag is saved to the file and supposed to be committed
func getNextDockerTag(tagFile string) (string, error) {
	datePattern := time.Now().Format("20060102")

	if _, err := os.Stat(tagFile); os.IsNotExist(err) {
		return datePattern + "-1", nil
	} else {
		content, err := ioutil.ReadFile(tagFile)
		if err != nil {
			return "", err
		}
		parts := strings.Split(string(content), "-")
		if parts[0] == datePattern {
			i, err := strconv.Atoi(parts[1])
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s-%d", datePattern, i+1), err

		} else {
			return datePattern + "-1", nil
		}
	}
}

// writeDockerTag persist the last used docker tag to a file.
func writeDockerTag(tagFile string, tag string) error {
	return ioutil.WriteFile(tagFile, []byte(tag), 0644)
}
