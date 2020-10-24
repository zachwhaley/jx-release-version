package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/zachwhaley/new-release-version/adapters"
	"github.com/zachwhaley/new-release-version/domain"

	"github.com/coreos/go-semver/semver"
	"github.com/hashicorp/go-version"

	"bufio"
	"context"
	"encoding/json"
	"flag"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
)

// Version is the build version
var Version string

// GitTag is the git tag of the build
var GitTag string

// BuildDate is the date when the build was created
var BuildDate string

type Project struct {
	Version string `xml:"version"`
}

type NewRelVer struct {
	dryrun       bool
	debug        bool
	dir          string
	ghOwner      string
	ghRepository string
	samerelease  bool
	baseVersion  string
	minor        bool
}

func main() {

	debug := flag.Bool("debug", false, "prints debug into to console")
	dir := flag.String("folder", ".", "the folder to look for files that contain a pom.xml or Makefile with the project version to bump")
	owner := flag.String("gh-owner", "", "a github repository owner if not running from within a git project  e.g. fabric8io")
	repo := flag.String("gh-repository", "", "a git repository if not running from within a git project  e.g. fabric8")
	baseVersion := flag.String("base-version", "", "use this instead of Makefile, pom.xml, etc, e.g. -base-version=2.0.0-SNAPSHOT")
	samerelease := flag.Bool("same-release", false, "for support old releases: for example 7.0.x and tag for new release 7.1.x already exist, with `-same-release` argument next version from 7.0.x will be returned ")
	ver := flag.Bool("version", false, "prints the version")
	minor := flag.Bool("minor", false, "increase minor version instead of patch")
	flag.Parse()

	if *ver {
		printVersion()
		os.Exit(0)
	}

	r := NewRelVer{
		debug:        *debug,
		dir:          *dir,
		ghOwner:      *owner,
		ghRepository: *repo,
		samerelease:  *samerelease,
		baseVersion:  *baseVersion,
		minor:        *minor,
	}

	if r.debug {
		fmt.Println("available environment:")
		for _, e := range os.Environ() {
			fmt.Println(e)
		}
	}

	gitHubClient := adapters.NewGitHubClient(r.debug)
	v, err := r.getNewVersionFromTag(gitHubClient)
	if err != nil {
		fmt.Println("failed to get new version", err)
		os.Exit(-1)
	}
	fmt.Printf("%s", v)
}

func printVersion() {
	fmt.Printf(`Version: %s
Git Tag: %s
Build Date: %s
`, Version, GitTag, BuildDate)
}

func (r NewRelVer) findVersionFile(f string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Join(r.dir, f))
	if err != nil && r.debug {
		fmt.Printf("found %s\n", f)
	}
	return data, err
}

func (r NewRelVer) getVersion() (string, error) {
	if r.baseVersion != "" {
		return r.baseVersion, nil
	}

	if gradle, err := r.findVersionFile("versions.gradle"); err == nil {
		return r.getGradleVersion(gradle)
	}
	if pkgjson, err := r.findVersionFile("package.json"); err == nil {
		return r.getPackageJSONVersion(pkgjson)
	}
	if setupCfg, err := r.findVersionFile("setup.cfg"); err == nil {
		return r.getSetupCfgVersion(setupCfg)
	}
	if setupPy, err := r.findVersionFile("setup.py"); err == nil {
		return r.getSetupPyVersion(setupPy)
	}
	if makefile, err := r.findVersionFile("Makefile"); err == nil {
		return r.getMakefileVersion(makefile)
	}
	if cmake, err := r.findVersionFile("CMakeLists.txt"); err == nil {
		return r.getCMakeVersion(cmake)
	}

	return "0.0.0", errors.New("No recognised file to obtain current version from")
}

func (r NewRelVer) getGradleVersion(gradle []byte) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(gradle)))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "project.version") {
			re := regexp.MustCompile("^\\s*project\\.version\\s*=\\s*['\"]([.\\d]+(-\\w+)?)['\"]")
			matched := re.FindStringSubmatch(scanner.Text())
			if len(matched) > 0 {
				version := strings.TrimSpace(matched[1])
				return version, nil
			}
		}
	}
	return "0.0.0", errors.New("No version found")
}

func (r NewRelVer) getPackageJSONVersion(pkgjson []byte) (string, error) {
	var project Project
	_ = json.Unmarshal(pkgjson, &project)
	if project.Version != "" {
		return project.Version, nil
	}
	return "0.0.0", errors.New("No version found")
}

func (r NewRelVer) getSetupCfgVersion(setupCfg []byte) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(setupCfg)))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "version") {
			re := regexp.MustCompile("^version\\s*=\\s*([.\\d]+(-\\w+)?)")
			matched := re.FindStringSubmatch(scanner.Text())
			if len(matched) > 0 {
				version := strings.TrimSpace(matched[1])
				return version, nil
			}
		}
	}

	return "0.0.0", errors.New("No version found")
}

func (r NewRelVer) getSetupPyVersion(setup []byte) (string, error) {
	// Regex to find the call to `setup(..., version='1.2.3', ...)`
	re := regexp.MustCompile("setup\\((.|\\n)*version\\s*=\\s*'(\\d|\\.)*'([^\\)]|\\n)*\\)")
	setupCallBytes := re.Find(setup)

	if len(setupCallBytes) > 0 {

		// Regex to find the argument `version='1.2.3'`
		versionRe := regexp.MustCompile("version\\s*=\\s*'(\\d*|\\.)*'")

		version := string(versionRe.Find(setupCallBytes))

		parts := strings.Split(strings.Replace(version, " ", "", -1), "=")
		v := strings.TrimPrefix(strings.TrimSuffix(parts[1], "'"), "'")

		if v != "" {
			return v, nil
		}
	}
	return "0.0.0", errors.New("No version found")
}

func (r NewRelVer) getMakefileVersion(makefile []byte) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(makefile)))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "VERSION") || strings.HasPrefix(scanner.Text(), "VERSION ") || strings.HasPrefix(scanner.Text(), "VERSION:") || strings.HasPrefix(scanner.Text(), "VERSION=") {
			parts := strings.Split(scanner.Text(), "=")

			v := strings.TrimSpace(parts[1])
			if v != "" {
				return v, nil
			}
		}
	}
	return "0.0.0", errors.New("No version found")
}

func (r NewRelVer) getCMakeVersion(cmake []byte) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(cmake)))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), " VERSION ") {
			re := regexp.MustCompile("project\\s*(([^\\s]+)\\s+VERSION\\s+([.\\d]+(-\\w+)?).*)")
			matched := re.FindStringSubmatch(scanner.Text())
			v := strings.TrimSpace(matched[3])
			if v != "" {
				return v, nil
			}
		}
	}
	return "0.0.0", errors.New("No version found")
}

func (r NewRelVer) getLatestTag(gitClient domain.GitClient) (string, error) {
	// Get base version from file, will fallback to 0.0.0 if not found.
	baseVersion, _ := r.getVersion()

	if r.debug {
		fmt.Printf("base version: %s\n", baseVersion)
	}

	// if repo isn't provided by flags fall back to using current repo if run from a git project
	var versionsRaw []string
	if r.ghOwner != "" && r.ghRepository != "" {
		ctx := context.Background()

		tags, err := gitClient.ListTags(ctx, r.ghOwner, r.ghRepository)

		if err != nil {
			return "", err
		}
		if len(tags) == 0 {
			// if no current flags exist then lets start at base version
			return baseVersion, errors.New("No existing tags found")
		}

		// build an array of all the tags
		versionsRaw = make([]string, len(tags))
		for i, tag := range tags {
			if r.debug {
				fmt.Printf("found remote tag %s\n", tag.Name)
			}
			versionsRaw[i] = tag.Name
		}
	} else {
		_, err := exec.LookPath("git")
		if err != nil {
			return "", fmt.Errorf("error running git: %v", err)
		}
		cmd := exec.Command("git", "fetch", "--tags", "-v")
		cmd.Env = append(cmd.Env, os.Environ()...)
		cmd.Dir = r.dir
		err = cmd.Run()
		if err != nil {
			return "", fmt.Errorf("error fetching tags: %v", err)
		}

		cmd = exec.Command("git", "tag")
		cmd.Dir = r.dir
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		str := strings.TrimSuffix(string(out), "\n")
		tags := strings.Split(str, "\n")

		if len(tags) == 0 {
			// if no current flags exist then lets start at base version
			return baseVersion, errors.New("No existing tags found")
		}

		// build an array of all the tags
		versionsRaw = make([]string, len(tags))
		for i, tag := range tags {
			if r.debug {
				fmt.Printf("found tag %s\n", tag)
			}
			tag = strings.TrimPrefix(tag, "v")
			if tag != "" {
				versionsRaw[i] = tag
			}
		}
	}

	// turn the array into a new collection of versions that we can sort
	var versions []*version.Version
	for _, raw := range versionsRaw {
		// if same-release argument is set work only with versions which Major and Minor versions are the same
		if r.samerelease {
			same, _ := isMajorMinorTheSame(baseVersion, raw)
			if same {
				v, _ := version.NewVersion(raw)
				if v != nil {
					versions = append(versions, v)
				}
			}
		} else {
			v, _ := version.NewVersion(raw)
			if v != nil {
				versions = append(versions, v)
			}
		}
	}

	if len(versions) == 0 {
		// if no current flags exist then lets start at base version
		return baseVersion, errors.New("No existing tags found")
	}

	// return the latest tag
	col := version.Collection(versions)
	if r.debug {
		fmt.Printf("version collection %v\n", col)
	}

	sort.Sort(col)
	latest := len(versions)
	if versions[latest-1] == nil {
		return baseVersion, errors.New("No existing tags found")
	}
	return versions[latest-1].String(), nil
}

func (r NewRelVer) getNewVersionFromTag(gitClient domain.GitClient) (string, error) {

	// get the latest github tag
	tag, err := r.getLatestTag(gitClient)
	if err != nil && tag == "" {
		return "", err
	}
	sv, err := semver.NewVersion(tag)
	if err != nil {
		return "", err
	}

	if r.minor {
		sv.BumpMinor()
	} else {
		sv.BumpPatch()
	}

	majorVersion := sv.Major
	minorVersion := sv.Minor
	patchVersion := sv.Patch

	// check if major or minor version has been changed
	baseVersion, err := r.getVersion()
	if err != nil {
		return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion), nil
	}

	// first use go-version to turn into a proper version, this handles 1.0-SNAPSHOT which semver doesn't
	tmpVersion, err := version.NewVersion(baseVersion)
	if err != nil {
		return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion), nil
	}
	bsv, err := semver.NewVersion(tmpVersion.String())
	if err != nil {
		return "", err
	}
	baseMajorVersion := bsv.Major
	baseMinorVersion := bsv.Minor
	basePatchVersion := bsv.Patch

	if baseMajorVersion > majorVersion ||
		(baseMajorVersion == majorVersion &&
			(baseMinorVersion > minorVersion) || (baseMinorVersion == minorVersion && basePatchVersion > patchVersion)) {
		majorVersion = baseMajorVersion
		minorVersion = baseMinorVersion
		patchVersion = basePatchVersion
	}
	return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion), nil
}

func isMajorMinorTheSame(v1 string, v2 string) (bool, error) {
	sv1, err1 := semver.NewVersion(v1)
	if err1 != nil {
		return false, err1
	}
	sv2, err2 := semver.NewVersion(v2)
	if err2 != nil {
		return false, err2
	}
	if sv1.Major != sv2.Major {
		return false, nil
	}
	if sv1.Minor != sv2.Minor {
		return false, nil
	}
	return true, nil
}
