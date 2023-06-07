## Release process

 1. export `GITHUB_TOKEN`
 2. Commit everything
 3. Tag latest commit: `git tag -s v1.0.0 -m "release v1.0.0"`
 4. Upload release: `goreleaser --rm-dist`
 5. Don't forget to push (both `main` and the tag)

## Update build image

Build image contains all the tools required to create local builds.

To publish a new `storjup/build` image:

 0. `go install github.com/magefile/mage@v1.11.0`
 1. `mage dockerBaseBuild` (only build)
 2. `mage dockerBasePublish` (build and publish)
 3. New tag is saved to `build.last`. Please commit that file with your PR.
 4. Use the new tag in `edge.Dockerfile` and `storj.Dockerfile`

Note: This process is assuming that you already authorized yourself with `img.dev.storj.io` with `docker login`.

Note: publishing base image is very similar