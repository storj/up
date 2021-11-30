// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"github.com/valyala/fastjson"
	"github.com/zeebo/errs"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/magefile/mage/sh"
)

// Test executes all unit and integration tests.
func Test() error {
	err := sh.RunV("go", "test", "./...")
	return err
}

// Coverage executes all unit test with coverage measurement.
func Coverage() error {
	fmt.Println("Executing tests and generate coverate information")
	err := sh.RunV("go", "test", "-coverprofile=/tmp/coverage.out", "./...")
	if err != nil {
		return err
	}
	return sh.RunV("go", "tool", "cover", "-html=/tmp/coverage.out", "-o", "coverage.html")
}

// Lint executes all the linters with golangci-lint
func Lint() error {
	return sh.RunV("./scripts/lint.sh")
}

// Format reformat code automatically
func Format() error {
	err := sh.RunV("gofmt", "-w", ".")
	if err != nil {
		return err
	}
	return sh.RunV("goimports", "-w", ".")

}

// GenBuild re-generates `./build` helper binary
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

func dockerCoreBuild(version string) error {
	err := sh.RunV("docker",
		"build",
		"-t", "ghcr.io/elek/storj:"+version,
		"--build-arg", "BRANCH=v"+version,
		"--build-arg", "TYPE=github",
		"-f", "cmd/files/docker/storj.Dockerfile", ".")
	if err != nil {
		return err
	}
	return nil
}

func dockerEdgeBuild(version string) error {
	err := sh.RunV("docker",
		"build",
		"-t", "ghcr.io/elek/storj-edge:"+version,
		"--build-arg", "BRANCH=v"+version,
		"--build-arg", "TYPE=github",
		"-f", "cmd/files/docker/edge.Dockerfile", ".")
	if err != nil {
		return err
	}
	return nil
}

func Integration() error {
	return sh.RunV("test/test.sh")
}

func Images() error {
	err := doOnMissing("storj", "storj", func(container string, repo string, version string) error {
		err := dockerCoreBuild(version)
		if err != nil {
			return err
		}
		return DockerCorePublish(version)
	})
	if err != nil {
		return err
	}

	err = doOnMissing("storj-edge", "gateway-mt", func(container string, repo string, version string) error {
		err := dockerEdgeBuild(version)
		if err != nil {
			return err
		}
		return DockerEdgePublish(version)
	})
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
// last used tag is saved to the file and supposed to be committed
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

func doOnMissing(containerName string, repoName string, action func(string, string, string) error) error {
	containerVersions := make(map[string]bool)
	versions, err := listContainerVersions(containerName)
	if err != nil {
		return err
	}
	for _, v := range versions {
		containerVersions[v] = true
	}

	releases, err := listReleaseVersions(repoName)
	if err != nil {
		return err
	}
	for _, v := range releases {
		if _, found := containerVersions[v]; !found {
			err = action(containerName, repoName, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// writeDockerTag persist the last used docker tag to a file.
func writeDockerTag(tagFile string, tag string) error {
	return ioutil.WriteFile(tagFile, []byte(tag), 0644)
}

// ListVersions prints out the available container / release versions
func ListVersions() error {
	fmt.Println("container: elek/storj")
	versions, err := listContainerVersions("storj")
	if err != nil {
		return err
	}
	for _, v := range versions {
		fmt.Println("   " + v)
	}
	fmt.Println("container: elek/storj-edge")
	versions, err = listContainerVersions("storj-edge")
	if err != nil {
		return err
	}
	for _, v := range versions {
		fmt.Println("   " + v)
	}
	fmt.Println("repo: storj/storj")
	versions, err = listReleaseVersions("storj")
	if err != nil {
		return err
	}
	for _, v := range versions {
		fmt.Println("   " + v)
	}
	fmt.Println("repo: storj/gateway-mt")
	versions, err = listReleaseVersions("gateway-mt")
	if err != nil {
		return err
	}
	for _, v := range versions {
		fmt.Println("   " + v)
	}
	return nil
}

func listReleaseVersions(name string) ([]string, error) {
	var parser fastjson.Parser
	url := fmt.Sprintf("https://api.github.com/repos/storj/%s/releases?per_page=10", name)
	rawVersions, err := callGithubAPIV3(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	parsed, err := parser.ParseBytes(rawVersions)
	if err != nil {
		return nil, err
	}
	versions, err := parsed.Array()
	if err != nil {
		return nil, err
	}
	var res []string
	for _, v := range versions {
		val := string(v.GetStringBytes("name"))
		if strings.Contains(val, "rc") {
			continue
		}
		if val[0] == 'v' {
			val = val[1:]
		}
		res = append(res, val)
	}
	return res, nil
}

func listContainerVersions(name string) ([]string, error) {
	var parser fastjson.Parser
	url := fmt.Sprintf("https://api.github.com/users/elek/packages/container/%s/versions", name)
	rawVersions, err := callGithubAPIV3(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	parsed, err := parser.ParseBytes(rawVersions)
	if err != nil {
		return nil, err
	}
	versions, err := parsed.Array()
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, v := range versions {
		for _, t := range v.GetArray("metadata", "container", "tags") {
			val, err := t.StringBytes()
			if err != nil {
				return nil, err
			}
			if string(val) == "latest" {
				continue
			}
			res = append(res, string(val))
		}
	}
	return res, nil
}

// callGithubAPIV3 is a wrapper around the HTTP method call.
func callGithubAPIV3(ctx context.Context, method string, url string, body io.Reader) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	token, err := getToken()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if resp.StatusCode > 299 {
		return nil, errs.Combine(errs.New("%s url is failed (%s): %s", method, resp.Status, url), resp.Body.Close())
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	return responseBody, errs.Combine(err, resp.Body.Close())
}

// getToken retrieves the GITHUB_TOKEN for API usage.
func getToken() (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token, nil
	}
	return "", fmt.Errorf("GITHUB_TOKEN environment variable must set")
}
