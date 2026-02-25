## Release process

 1. export `GITHUB_TOKEN`
 2. Commit everything
 3. Tag latest commit: `git tag -s v1.0.0 -m "release v1.0.0"`
 4. Push tag: `git push origin v1.0.0`
 5. Upload release: `goreleaser release`

### Requirements

- Mage. You can install it with `go install github.com/magefile/mage`
- Docker with multi-platform enabled or configured, depending what Docker you use
  https://docs.docker.com/build/building/multi-platform
- Docker client has access _img.dev.storj.io_ registry (e.g. `docker login ...`). This is only
  required for publishing the images

## Update build image

Build image contains all the tools required to create local builds.

To publish a new `storjup/build` image:

 1. `mage dockerBuildBuild` (only build)
 2. `mage dockerBuildPublish` (build and publish)
 3. New tag is saved to `build.last`. Please commit that file with your PR.

## Update base image

Base image contains the storj-up binaries and required OS packages.

To publish a new `storjup/base` image:

 1. `mage dockerBaseBuild` (only build)
 2. `mage dockerBasePublish` (build and publish)
 3. New tag is saved to `base.last`. Please commit that file with your PR.

## Update Storj image

Storj image contains all the binaries required to run the satellite.

To publish a new `storjup/storj` image:

1. `mage dockerStorj <version> false "" ""` (only build, uses tags from .last files)
2. `mage dockerStorj <version> true "" ""` (build and publish, uses tags from .last files)
3. Use the new tag in recipe files

To use specific build/base image tags:

```bash
mage dockerStorj <version> <publish> <buildTag> <baseTag>
```

- `version`: Required. The storj version without `v` prefix (e.g., `1.147.5`)
- `publish`: `true` to push to registry, `false` for local build only
- `buildTag`: Optional. Build image tag; pass `""` to use value from `build.last`
- `baseTag`: Optional. Base image tag; pass `""` to use value from `base.last`

## Update Edge image

Edge image contains all the binaries required to run the edge services.

To publish a new `storjup/edge` image:

1. `mage dockerEdge <version> false "" ""` (only build, uses tags from .last files)
2. `mage dockerEdge <version> true "" ""` (build and publish, uses tags from .last files)
3. Use the new tag in recipe files

To use specific build/base image tags:

```bash
mage dockerEdge <version> <publish> <buildTag> <baseTag>
```

- `version`: Required. The edge version without `v` prefix (e.g., `1.111.0`)
- `publish`: `true` to push to registry, `false` for local build only
- `buildTag`: Optional. Build image tag; pass `""` to use value from `build.last`
- `baseTag`: Optional. Base image tag; pass `""` to use value from `base.last`
