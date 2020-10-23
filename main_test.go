package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/trendmicro/release-version/domain"
	"testing"
)

func TestGradle(t *testing.T) {
	r := RelVer{
		dir: "test-resources/java",
	}
	v, err := r.getVersion()

	assert.NoError(t, err)

	assert.Equal(t, "1.2.3", v, "error with getVersion for a versions.gradle")
}

func TestPackageJSON(t *testing.T) {
	r := RelVer{
		dir: "test-resources/nodejs",
	}
	v, err := r.getVersion()

	assert.NoError(t, err)

	assert.Equal(t, "1.2.3", v, "error with getVersion for a package.json")
}

func TestSetupCfg(t *testing.T) {

	r := RelVer{
		dir: "test-resources/python/setup.cfg",
	}
	v, err := r.getVersion()

	assert.NoError(t, err)

	assert.Equal(t, "1.2.3", v, "error with getVersion for a setup.cfg")
}

func TestSetupPy(t *testing.T) {

	r := RelVer{
		dir: "test-resources/python/standard",
	}
	v, err := r.getVersion()

	assert.NoError(t, err)

	assert.Equal(t, "4.5.6", v, "error with getVersion for a setup.py")
}

func TestSetupPyNested(t *testing.T) {

	r := RelVer{
		dir: "test-resources/python/nested",
	}
	v, err := r.getVersion()

	assert.NoError(t, err)

	assert.Equal(t, "4.5.6", v, "error with getVersion for a setup.py")
}

func TestSetupPyOneLine(t *testing.T) {

	r := RelVer{
		dir: "test-resources/python/one_line",
	}
	v, err := r.getVersion()

	assert.NoError(t, err)

	assert.Equal(t, "4.5.6", v, "error with getVersion for a setup.py")
}

func TestMakefile(t *testing.T) {
	r := RelVer{
		dir: "test-resources/make",
	}

	v, err := r.getVersion()

	assert.NoError(t, err)

	assert.Equal(t, "1.2.0-SNAPSHOT", v, "error with getVersion for a Makefile")
}

func TestCMakefile(t *testing.T) {

	r := RelVer{
		dir: "test-resources/cmake",
	}

	v, err := r.getVersion()

	assert.NoError(t, err)

	assert.Equal(t, "1.2.0-SNAPSHOT", v, "error with getVersion for a CMakeLists.txt")
}

/* Disable git tag test until we have some tags
func TestGetNewVersionFromTagCurrentRepo(t *testing.T) {
	r := RelVer{
		dryrun: false,
		dir:    "test-resources/make",
	}

	tags := createTags()

	mockClient := &mocks.GitClient{}
	mockClient.On("ListTags", context.Background(), r.ghOwner, r.ghRepository).Return(tags, nil)
	v, err := r.getNewVersionFromTag(mockClient)

	assert.NoError(t, err)
	assert.Equal(t, "1.2.1", v, "error bumping a patch version")
}
*/

/* Disable GitHub test
func TestGetGitTag(t *testing.T) {
	r := RelVer{
		ghOwner:      "jenkins-x",
		ghRepository: "release-version",
	}

	gitHubClient := adapters.NewGitHubClient(r.debug)

	expectedVersion, err := getLatestTag(r, gitHubClient)
	assert.NoError(t, err)

	r = RelVer{}

	v, err := getLatestTag(r, gitHubClient)

	assert.NoError(t, err)

	assert.Equal(t, expectedVersion, v, "error with getLatestTag for a Makefile")
}

func TestGetNewMinorVersionFromGitHubTag(t *testing.T) {

	r := RelVer{
		ghOwner:      "rawlingsj",
		ghRepository: "semver-release-version",
		minor:        true,
	}

	tags := createTags()

	mockClient := &mocks.GitClient{}
	mockClient.On("ListTags", context.Background(), r.ghOwner, r.ghRepository).Return(tags, nil)

	v, err := getNewVersionFromTag(r, mockClient)

	assert.NoError(t, err)
	assert.Equal(t, "1.1.0", v, "error bumping a minor version")
}

func TestGetNewPatchVersionFromGitHubTag(t *testing.T) {

	r := RelVer{
		ghOwner:      "rawlingsj",
		ghRepository: "semver-release-version",
	}

	tags := createTags()

	mockClient := &mocks.GitClient{}
	mockClient.On("ListTags", context.Background(), r.ghOwner, r.ghRepository).Return(tags, nil)

	v, err := getNewVersionFromTag(r, mockClient)

	assert.NoError(t, err)
	assert.Equal(t, "1.0.18", v, "error bumping a patch version")
}
*/

func createTags() []domain.Tag {
	var tags []domain.Tag
	tags = append(tags, domain.Tag{Name: "v1.0.0"})
	tags = append(tags, domain.Tag{Name: "v1.0.1"})
	tags = append(tags, domain.Tag{Name: "v1.0.10"})
	tags = append(tags, domain.Tag{Name: "v1.0.11"})
	tags = append(tags, domain.Tag{Name: "v1.0.12"})
	tags = append(tags, domain.Tag{Name: "v1.0.13"})
	tags = append(tags, domain.Tag{Name: "v1.0.14"})
	tags = append(tags, domain.Tag{Name: "v1.0.15"})
	tags = append(tags, domain.Tag{Name: "v1.0.16"})
	tags = append(tags, domain.Tag{Name: "v1.0.17"})
	tags = append(tags, domain.Tag{Name: "v1.0.2"})
	tags = append(tags, domain.Tag{Name: "v1.0.3"})
	tags = append(tags, domain.Tag{Name: "v1.0.4"})
	tags = append(tags, domain.Tag{Name: "v1.0.5"})
	tags = append(tags, domain.Tag{Name: "v1.0.6"})
	tags = append(tags, domain.Tag{Name: "v1.0.7"})
	tags = append(tags, domain.Tag{Name: "v1.0.8"})
	tags = append(tags, domain.Tag{Name: "v1.0.9"})

	return tags
}
