# New Release Version

Returns a new release version based on previous git tags that can be used in a new release.

This is a simple binary that can be used in CD pipelines to read version files (e.g. Makefile, versions.gradle, etc) and return a 'patch' incremented version.

If you need to bump the major or minor version simply increment the version in your version file.

This helps in continuous delivery if you want an automatic release when a change is merged to master.  Traditional approaches mean the version is stored in a file that is checked and updated after each release.  If you want automatic releases this means you will get another release triggered from the version update resulting in a cyclic release sitiation.  

Using a git tag to work out the next release version is better than traditional approaches of storing it in a VERSION file or updating a pom.xml.  If a major or minor version increase is required then still update the file and `new-release-version` will use you new version.

## Prerequisits

- `git` to be available on your `$PATH`

## Examples

- If your project is new or has no existing git tags then running `new-release-version` will return a default version of `0.0.1`

- If your latest git tag is `1.2.3` and you Makefile or pom.xml is `1.2.0-SNAPSHOT` then `new-release-version` will return `1.2.4`

- If your latest git tag is `1.2.3` and your Makefile or pom.xml is `2.0.0` then `new-release-version` will return `2.0.0`

- If you need to support an old release for example 7.0.x and tags for new realese 7.1.x already exist, the `-same-release` flag  will help to obtain version from 7.0.x release. If the pom file version is 7.0.0-SNAPSHOT and both the 7.1.0 and 7.0.2 tags exist the command `new-release-version` will return 7.1.1 but if we run `new-release-version -same-release` it will return 7.0.3

- If you need to get a release version `1.1.0` for older release and your last tag is `1.2.3` please change your Makefile or pom.xml to `1.1.0-SNAPSHOT` and run `new-release-version -same-release`

## Example Makefile

```Makefile
VERSION := 2.0.0
```

## Example versions.gradle

```gradle
ext {
    project.version = '1.0.0'
}
```

Then in your release pipeline do something like this:

```sh
    ➜ RELEASE_VERSION=$(new-release-version)
    ➜ echo "New release version ${RELEASE_VERSION}
    ➜ mvn versions:set -DnewVersion=${RELEASE_VERSION}
    ➜ git commit -a -m 'release ${RELEASE_VERSION}'
    ➜ git tag -fa v${RELEASE_VERSION} -m 'Release version ${RELEASE_VERSION}'
    ➜ git push origin v${RELEASE_VERSION}
```
