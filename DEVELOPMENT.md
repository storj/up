## Release process

 1. export `GITHUB_TOKEN`
 2. Commit everything
 3. Tag latest commit: `git tag -s v1.0.0 -m "release v1.0.0"`
 4. Upload release: `goreleaser --rm-dist`
 5. Don't forget to push (both `main` and the tag)
