# Releasing `gopkg`

This document is the step-by-step guide for cutting a release of
`github.com/docodex/gopkg`. Follow it top to bottom; do not skip the
pre-flight checks.

## Table of Contents

- [Versioning Policy](#versioning-policy)
- [Pre-flight Checks](#pre-flight-checks)
- [Cutting a Release](#cutting-a-release)
- [Post-release Verification](#post-release-verification)
- [Publishing a GitHub Release (optional)](#publishing-a-github-release-optional)
- [Common Pitfalls](#common-pitfalls)
- [Rolling Back a Broken Tag](#rolling-back-a-broken-tag)

## Versioning Policy

`gopkg` follows [Semantic Versioning 2.0.0](https://semver.org/) as interpreted
by the Go module system:

- **MAJOR** (`vX.0.0`) - incompatible API changes. Because the module path
  `github.com/docodex/gopkg` has no `/vN` suffix, only `v0` and `v1` may be
  published without a path change. A `v2.0.0` would require a `/v2` module
  path (see <https://go.dev/ref/mod#major-version-suffixes>).
- **MINOR** (`vX.Y.0`) - backwards-compatible additions (new exported
  identifiers, new packages).
- **PATCH** (`vX.Y.Z`) - backwards-compatible bug fixes only.

`v0.y.z` tags carry **no API stability guarantee** - use them while the
public surface is still being shaped. `v1.0.0` is a commitment that
exported identifiers will not change in a breaking way until `v2`.

Tag names must match `vMAJOR.MINOR.PATCH` exactly (lowercase `v`, no
leading zeros, no build metadata). Anything else is rejected by the Go
module proxy.

## Pre-flight Checks

Run every one of these on a clean checkout of `master` (or the release
branch) before creating the tag:

```sh
# 1. Working tree is clean and up to date with origin.
git checkout master
git fetch origin
git status # must print: nothing to commit, working tree clean
git log --oneline origin/master..HEAD # must be empty: no unpushed commits

# 2. Module builds, vets, and all tests pass under the race detector.
go build ./...
go vet ./...
go test ./...
go test -race ./... # go test -race -count=1 ./...

# 3. go.mod declares the minimum Go version you intend to support.
#    After release, bumping this is a breaking change.
grep '^go ' go.mod
```

If you rely on external modules (see `go.mod`), also confirm they are
using tagged releases, not pseudo-versions - a release that depends on
`@master` of another module is fragile.

## Cutting a Release

Use a **annotated** tag (`-a`), never a lightweight one. Annotated tags
carry author, date, and message, which `pkg.go.dev`, `GoProxy`, and
GitHub all consume.

```sh
# Replace v1.0.0 with the version you are cutting.
VERSION=v1.0.0

git tag -a "$VERSION" -m "$VERSION: first stable release"
git push origin "$VERSION"
```

After `git push`, the tag is immutable from the Go module proxy's
perspective within ~1 minute. Do not retag the same version.

## Post-release Verification

```sh
# Force the Go module proxy to fetch the new tag (one-time; populates cache).
GOPROXY=https://proxy.golang.org \
  go list -m "github.com/docodex/gopkg@$VERSION"

# In a throwaway directory, confirm a fresh consumer can depend on it.
mkdir /tmp/verify-gopkg && cd /tmp/verify-gopkg
go mod init verify
go get "github.com/docodex/gopkg@$VERSION"
cat go.mod
```

Also visit <https://pkg.go.dev/github.com/docodex/gopkg> - it may take a
few minutes for the new version to be indexed.

## Publishing a GitHub Release (optional)

A GitHub Release attaches changelog notes to the tag and makes the
version discoverable from the repository sidebar. Requires
[`gh`](https://cli.github.com).

```sh
gh release create "$VERSION" \
  --title "$VERSION" \
  --generate-notes
```

`--generate-notes` auto-populates the description from merged PRs since
the previous tag. Review and edit the generated notes before clicking
publish if you want a curated changelog.

## Common Pitfalls

- **Lightweight tag (`git tag v1.0.0`) instead of annotated (`-a`).**
  Go module tooling still accepts it, but `pkg.go.dev` shows a blank
  author/date. Always use `-a -m "..."`.
- **Wrong tag format** - `V1.0.0`, `1.0.0`, `v1.0`, `v1.0.0-rc1+build`.
  Pre-release tags like `v1.0.0-rc1` are fine; build-metadata suffixes
  (`+build`) are not comparable and Go treats them as equivalent to
  the base version.
- **Amending the tagged commit** after pushing the tag. Once
  `proxy.golang.org` has cached `(v1.0.0, <old sha>)`, any user who
  runs `go get` before your force-push gets the old commit; users who
  pull afterward get the new one. This de-synchronization is permanent.
  If you need a fix, cut `v1.0.1`.
- **Moving a tag.** `git tag -f` + `git push --force` will not refresh
  the module proxy cache. Same resolution: cut a new patch version.
- **Forgetting `/v2` for a major bump.** If you ever need `v2.0.0`, you
  must first change the module path in `go.mod` to
  `github.com/docodex/gopkg/v2` and update every internal import. See
  <https://go.dev/ref/mod#major-version-suffixes>.
- **Depending on pseudo-versions of other modules at release time.**
  Run `go list -m all | grep -- '-0\.'` - any line that prints means a
  dependency is pinned to a commit, not a tag. Prefer tagged versions
  for stable release.

## Rolling Back a Broken Tag

You cannot unpublish a tag once the Go proxy has cached it. The only
supported recovery is to release a new patch version that supersedes
the broken one:

```sh
# 1. Commit the fix on master.
git commit -am "fix: <summary>"
git push origin master

# 2. Cut the patch tag.
git tag -a v1.0.1 -m "v1.0.1: fix <summary>"
git push origin v1.0.1
```

For widely broken release, also add a
[`retract` directive](https://go.dev/ref/mod#go-mod-file-retract) to
`go.mod` in the following release so consumers of the broken tag are
warned:

```go
// go.mod
retract v1.0.0 // contains data-corrupting bug, use v1.0.1+
```
