// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build mage
// +build mage

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/magefile/mage/sh"
	"github.com/zeebo/errs"
)

// Coverage executes all unit test with coverage measurement.
//
//nolint:deadcode
func Coverage() error {
	fmt.Println("Executing tests and generate coverate information")
	err := sh.RunV("go", "test", "-coverprofile=/tmp/coverage.out", "./...")
	if err != nil {
		return err
	}
	return sh.RunV("go", "tool", "cover", "-html=/tmp/coverage.out", "-o", "coverage.html")
}

// Format reformats code automatically.
//
//nolint:deadcode
func Format() error {
	err := sh.RunV("gofmt", "-w", ".")
	if err != nil {
		return err
	}
	return sh.RunV("goimports", "-w", ".")
}

// GenBuild re-generates `./build` helper binary.
//
//nolint:deadcode
func GenBuild() error {
	envs := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        "linux",
		"GOARCH":      "amd64",
	}
	return sh.RunWithV(envs, "mage", "-compile", "build")
}

func withDockerTag(filename string, publish bool, action func(tag string) error) error {
	tag, err := getNextDockerTag(filename)
	if err != nil {
		return err
	}

	err = action(tag)
	if err != nil {
		return err
	}
	if publish {
		return writeDockerTag(filename, tag)
	}
	return nil
}

// DockerBaseBuild builds storj-base image.
//
//nolint:deadcode
func DockerBaseBuild() error {
	return dockerBase(false)
}

// DockerBasePublish pushes storj-base image.
//
//nolint:deadcode
func DockerBasePublish() error {
	return dockerBase(true)
}

func dockerBase(publish bool) error {
	return withDockerTag("storj-base.last", publish, func(tag string) error {
		return buildxRun(publish,
			"build",
			"--tag", "img.dev.storj.io/storjup/base:"+tag,
			"-f", "cmd/files/docker/base.Dockerfile", ".")
	})
}

// DockerBuildBuild builds the storj-build docker image.
//
//nolint:deadcode
func DockerBuildBuild() error {
	return dockerBuild(false)
}

// DockerBuildPublish pushes the storj-build docker image
//
//nolint:deadcode
func DockerBuildPublish() error {
	return errs.Combine(
		dockerBuild(true),
	)
}

func dockerBuild(publish bool) error {
	return withDockerTag("build.last", publish, func(tag string) error {
		return buildxRun(publish,
			"build",
			"--build-arg", "TYPE=github",
			"--build-arg", "BRANCH=main",
			"--build-arg", "REPO=https://github.com/storj/storj.git",
			"--tag", "img.dev.storj.io/storjup/build:"+tag,
			"-f", "cmd/files/docker/build.Dockerfile", ".")
	})
}

func dockerCore(version string, publish bool) error {
	err := buildxRun(
		publish,
		"build",
		"-t", "img.dev.storj.io/storjup/storj:"+version,
		"--build-arg", "BRANCH=v"+version,
		"--build-arg", "TYPE=github",
		"-f", "cmd/files/docker/storj.Dockerfile", ".")
	if err != nil {
		return err
	}
	return nil
}

func buildxRun(publish bool, args ...string) error {
	if publish {
		args = append(args, "--push")
	}

	hasPlatform := false
	for _, arg := range args {
		if strings.HasPrefix(arg, "--platform") {
			hasPlatform = true
		}
	}
	if !hasPlatform {
		args = append(args, "--platform=linux/amd64,linux/arm64")
	}

	args = append([]string{"docker", "buildx"}, args...)
	return sh.RunV(args[0], args[1:]...)
}

func dockerEdge(version string, publish bool) error {
	err := buildxRun(publish,
		"build",
		"-t", "img.dev.storj.io/storjup/edge:"+version,
		"--build-arg", "BRANCH=v"+version,
		"--build-arg", "TYPE=github",
		"-f", "cmd/files/docker/edge.Dockerfile", ".")
	if err != nil {
		return err
	}
	return nil
}

// Integration executes integration tests.
//
//nolint:deadcode
func Integration() error {
	return sh.RunV("bash", "test/test.sh")
}

// RebuildImages rebuilds all core and edge images.
//
//nolint:deadcode
func RebuildImages() error {
	versions, err := listContainerVersions("storj")
	if err != nil {
		return err
	}
	for _, v := range versions {
		err := dockerCore(v, true)
		if err != nil {
			return err
		}
	}

	versions, err = listContainerVersions("edge")
	if err != nil {
		return err
	}
	for _, v := range versions {
		err := dockerEdge(v, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// DockerEdge builds a Edge docker image for local use.
//
//nolint:deadcode
func DockerEdge(version string, publish bool) error {
	if version == "" {
		return errs.New("VERSION should be defined with environment variable")
	}
	return dockerEdge(version, publish)
}

// DockerStorj builds a Core docker image for local use.
//
//nolint:deadcode
func DockerStorj(version string, publish bool) error {
	if version == "" {
		return errs.New("VERSION should be defined with environment variable")
	}
	return dockerCore(version, publish)
}

// Images build missing images for existing git tags
//
//nolint:deadcode
func Images() error {
	err := doOnMissing("storj", "storj", func(container string, repo string, version string) error {
		err := dockerCore(version, true)
		if err != nil {
			return err
		}
		return dockerCorePublish(version)
	})
	if err != nil {
		return err
	}

	err = doOnMissing("edge", "gateway-mt", func(container string, repo string, version string) error {
		err := dockerEdge(version, true)
		if err != nil {
			return err
		}
		return dockerEdgePublish(version)
	})
	if err != nil {
		return err
	}

	return nil
}

// ListImages prints all the existing storj and storj-edge images in the repo.
//
//nolint:deadcode
func ListImages() error {
	versions, err := listContainerVersions("storj")
	if err != nil {
		return err
	}
	for _, v := range versions {
		fmt.Printf("storj:%s\n", v)
	}

	versions, err = listContainerVersions("edge")
	if err != nil {
		return err
	}
	for _, v := range versions {
		fmt.Printf("edge:%s\n", v)
	}
	return nil
}

func dockerPush(image string, tag string) error {
	err := sh.RunV("docker", "push", fmt.Sprintf("img.dev.storj.io/storjup/%s:%s", image, tag))
	if err != nil {
		return err
	}
	return err
}

func dockerCorePublish(version string) error {
	return dockerPush("storj", version)
}

func dockerEdgePublish(version string) error {
	return dockerPush("storj-edge", version)
}

// getNextDockerTag generates docker tag with the pattern yyyymmdd-n.
// last used tag is saved to the file and supposed to be committed.
func getNextDockerTag(tagFile string) (string, error) {
	datePattern := time.Now().Format("20060102")

	if _, err := os.Stat(tagFile); os.IsNotExist(err) {
		return datePattern + "-1", nil
	}

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

	}
	return datePattern + "-1", nil
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
	return ioutil.WriteFile(tagFile, []byte(tag), 0o644)
}

// ListVersions prints out the available container / release versions.
//
//nolint:deadcode
func ListVersions() error {
	fmt.Println("container: storj")
	coreContainers, err := listContainerVersions("storj")
	if err != nil {
		return err
	}
	for _, v := range coreContainers {
		fmt.Println("   " + v)
	}
	fmt.Println("container: edge")
	edgeContainers, err := listContainerVersions("edge")
	if err != nil {
		return err
	}
	for _, v := range edgeContainers {
		fmt.Println("   " + v)
	}
	fmt.Println("repo: storj/storj")
	versions, err := listReleaseVersions("storj")
	if err != nil {
		return err
	}
	for _, v := range versions {
		fmt.Println("   " + v + " container:" + findContainer(coreContainers, v))
	}
	fmt.Println("repo: storj/gateway-mt")
	versions, err = listReleaseVersions("gateway-mt")
	if err != nil {
		return err
	}
	for _, v := range versions {
		fmt.Println("   " + v + " container:" + findContainer(edgeContainers, v))
	}
	return nil
}

func findContainer(containers []string, v string) string {
	for _, c := range containers {
		if c == v {
			return c
		}
	}
	return "MISSING"
}

func listReleaseVersions(name string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/storj/%s/releases?per_page=10", name)
	rawVersions, err := callGithubAPIV3(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var releases []release
	err = json.Unmarshal(rawVersions, &releases)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, v := range releases {
		name := v.TagName
		if strings.Contains(name, "rc") {
			continue
		}
		if name[0] == 'v' {
			name = name[1:]
		}
		res = append(res, name)
	}
	sort.Strings(res)
	return res, nil
}

// listContainerVersions lists the available tags for one specific container.
func listContainerVersions(name string) ([]string, error) {
	ctx := context.Background()
	url := fmt.Sprintf("https://img.dev.storj.io/auth?service=img.dev.storj.io&scope=repository:storjup/%s:pull", name)
	tokenResponse, err := httpCall(ctx, "GET", url, nil)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	k := struct {
		Token string `json:"token"`
	}{}
	err = json.Unmarshal(tokenResponse, &k)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	url = fmt.Sprintf("https://img.dev.storj.io/v2/storjup/%s/tags/list", name)
	tagResponse, err := httpCall(ctx, "GET", url, nil, func(request *http.Request) {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", k.Token))
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}

	var versions version
	err = json.Unmarshal(tagResponse, &versions)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, version := range versions.Tags {
		if version == "latest" {
			continue
		}
		res = append(res, version)
	}
	return res, nil
}

// callGithubAPIV3 is a wrapper around the HTTP method call.
func callGithubAPIV3(ctx context.Context, method string, url string, body io.Reader) ([]byte, error) {
	token, err := getToken()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return httpCall(ctx, method, url, body, func(req *http.Request) {
		req.Header.Add("Authorization", "token "+token)
		req.Header.Add("Accept", "application/vnd.github.v3+json")
	})
}

type httpRequestOpt func(*http.Request)

func httpCall(ctx context.Context, method string, url string, body io.Reader, opts ...httpRequestOpt) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	for _, o := range opts {
		o(req)
	}
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

// release is a Github API response object.
type release struct {
	URL             string    `json:"url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	HTMLURL         string    `json:"html_url"`
	ID              int       `json:"id"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
	MentionsCount   int       `json:"mentions_count,omitempty"`
}

// version is a Docker v2 REST API response object.
type version struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
