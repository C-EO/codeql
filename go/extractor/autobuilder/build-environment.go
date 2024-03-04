package autobuilder

import (
	"fmt"
	"log"
	"os"

	"github.com/github/codeql-go/extractor/diagnostics"
	"github.com/github/codeql-go/extractor/project"
	"github.com/github/codeql-go/extractor/toolchain"
	"golang.org/x/mod/semver"
)

const minGoVersion = "1.11"
const maxGoVersion = "1.21"

type versionInfo struct {
	goModVersion      string // The version of Go found in the go directive in the `go.mod` file.
	goModVersionFound bool   // Whether a `go` directive was found in the `go.mod` file.
	goEnvVersion      string // The version of Go found in the environment.
	goEnvVersionFound bool   // Whether an installation of Go was found in the environment.
}

func (v versionInfo) String() string {
	return fmt.Sprintf(
		"go.mod version: %s, go.mod directive found: %t, go env version: %s, go installation found: %t",
		v.goModVersion, v.goModVersionFound, v.goEnvVersion, v.goEnvVersionFound)
}

// Check if `version` is lower than `minGoVersion`. Note that for this comparison we ignore the
// patch part of the version, so 1.20.1 and 1.20 are considered equal.
func belowSupportedRange(version string) bool {
	return semver.Compare(semver.MajorMinor("v"+version), "v"+minGoVersion) < 0
}

// Check if `version` is higher than `maxGoVersion`. Note that for this comparison we ignore the
// patch part of the version, so 1.20.1 and 1.20 are considered equal.
func aboveSupportedRange(version string) bool {
	return semver.Compare(semver.MajorMinor("v"+version), "v"+maxGoVersion) > 0
}

// Check if `version` is lower than `minGoVersion` or higher than `maxGoVersion`. Note that for
// this comparison we ignore the patch part of the version, so 1.20.1 and 1.20 are considered
// equal.
func outsideSupportedRange(version string) bool {
	return belowSupportedRange(version) || aboveSupportedRange(version)
}

// Assuming `v.goModVersionFound` is false, emit a diagnostic and return the version to install,
// or the empty string if we should not attempt to install a version of Go.
func getVersionWhenGoModVersionNotFound(v versionInfo) (msg, version string) {
	if !v.goEnvVersionFound {
		// There is no Go version installed in the environment. We have no indication which version
		// was intended to be used to build this project. Go versions are generally backwards
		// compatible, so we install the maximum supported version.
		msg = "No version of Go installed and no `go.mod` file found. Requesting the maximum " +
			"supported version of Go (" + maxGoVersion + ")."
		version = maxGoVersion
		diagnostics.EmitNoGoModAndNoGoEnv(msg)
	} else if outsideSupportedRange(v.goEnvVersion) {
		// The Go version installed in the environment is not supported. We have no indication
		// which version was intended to be used to build this project. Go versions are generally
		// backwards compatible, so we install the maximum supported version.
		msg = "No `go.mod` file found. The version of Go installed in the environment (" +
			v.goEnvVersion + ") is outside of the supported range (" + minGoVersion + "-" +
			maxGoVersion + "). Requesting the maximum supported version of Go (" + maxGoVersion +
			")."
		version = maxGoVersion
		diagnostics.EmitNoGoModAndGoEnvUnsupported(msg)
	} else {
		// The version of Go that is installed is supported. We have no indication which version
		// was intended to be used to build this project. We assume that the installed version is
		// suitable and do not install a version of Go.
		msg = "No `go.mod` file found. Version " + v.goEnvVersion + " installed in the " +
			"environment is supported. Not requesting any version of Go."
		version = ""
		diagnostics.EmitNoGoModAndGoEnvSupported(msg)
	}

	return msg, version
}

// Assuming `v.goModVersion` is above the supported range, emit a diagnostic and return the
// version to install, or the empty string if we should not attempt to install a version of Go.
func getVersionWhenGoModVersionTooHigh(v versionInfo) (msg, version string) {
	if !v.goEnvVersionFound {
		// The version in the `go.mod` file is above the supported range. There is no Go version
		// installed. We install the maximum supported version as a best effort.
		msg = "The version of Go found in the `go.mod` file (" + v.goModVersion +
			") is above the supported range (" + minGoVersion + "-" + maxGoVersion +
			"). No version of Go installed. Requesting the maximum supported version of Go (" +
			maxGoVersion + ")."
		version = maxGoVersion
		diagnostics.EmitGoModVersionTooHighAndNoGoEnv(msg)
	} else if aboveSupportedRange(v.goEnvVersion) {
		// The version in the `go.mod` file is above the supported range. The version of Go that
		// is installed is above the supported range. We do not install a version of Go.
		msg = "The version of Go found in the `go.mod` file (" + v.goModVersion +
			") is above the supported range (" + minGoVersion + "-" + maxGoVersion +
			"). The version of Go installed in the environment (" + v.goEnvVersion +
			") is above the supported range (" + minGoVersion + "-" + maxGoVersion +
			"). Not requesting any version of Go."
		version = ""
		diagnostics.EmitGoModVersionTooHighAndEnvVersionTooHigh(msg)
	} else if belowSupportedRange(v.goEnvVersion) {
		// The version in the `go.mod` file is above the supported range. The version of Go that
		// is installed is below the supported range. We install the maximum supported version as
		// a best effort.
		msg = "The version of Go found in the `go.mod` file (" + v.goModVersion +
			") is above the supported range (" + minGoVersion + "-" + maxGoVersion +
			"). The version of Go installed in the environment (" + v.goEnvVersion +
			") is below the supported range (" + minGoVersion + "-" + maxGoVersion +
			"). Requesting the maximum supported version of Go (" + maxGoVersion + ")."
		version = maxGoVersion
		diagnostics.EmitGoModVersionTooHighAndEnvVersionTooLow(msg)
	} else if semver.Compare("v"+maxGoVersion, "v"+v.goEnvVersion) > 0 {
		// The version in the `go.mod` file is above the supported range. The version of Go that
		// is installed is supported and below the maximum supported version. We install the
		// maximum supported version as a best effort.
		msg = "The version of Go found in the `go.mod` file (" + v.goModVersion +
			") is above the supported range (" + minGoVersion + "-" + maxGoVersion +
			"). The version of Go installed in the environment (" + v.goEnvVersion +
			") is below the maximum supported version (" + maxGoVersion +
			"). Requesting the maximum supported version of Go (" + maxGoVersion + ")."
		version = maxGoVersion
		diagnostics.EmitGoModVersionTooHighAndEnvVersionBelowMax(msg)
	} else {
		// The version in the `go.mod` file is above the supported range. The version of Go that
		// is installed is the maximum supported version. We do not install a version of Go.
		msg = "The version of Go found in the `go.mod` file (" + v.goModVersion +
			") is above the supported range (" + minGoVersion + "-" + maxGoVersion +
			"). The version of Go installed in the environment (" + v.goEnvVersion +
			") is the maximum supported version (" + maxGoVersion +
			"). Not requesting any version of Go."
		version = ""
		diagnostics.EmitGoModVersionTooHighAndEnvVersionMax(msg)
	}

	return msg, version
}

// Assuming `v.goModVersion` is below the supported range, emit a diagnostic and return the
// version to install, or the empty string if we should not attempt to install a version of Go.
func getVersionWhenGoModVersionTooLow(v versionInfo) (msg, version string) {
	if !v.goEnvVersionFound {
		// There is no Go version installed. The version in the `go.mod` file is below the
		// supported range. Go versions are generally backwards compatible, so we install the
		// minimum supported version.
		msg = "The version of Go found in the `go.mod` file (" + v.goModVersion +
			") is below the supported range (" + minGoVersion + "-" + maxGoVersion +
			"). No version of Go installed. Requesting the minimum supported version of Go (" +
			minGoVersion + ")."
		version = minGoVersion
		diagnostics.EmitGoModVersionTooLowAndNoGoEnv(msg)
	} else if outsideSupportedRange(v.goEnvVersion) {
		// The version of Go that is installed is outside of the supported range. The version
		// in the `go.mod` file is below the supported range. Go versions are generally
		// backwards compatible, so we install the minimum supported version.
		msg = "The version of Go found in the `go.mod` file (" + v.goModVersion +
			") is below the supported range (" + minGoVersion + "-" + maxGoVersion +
			"). The version of Go installed in the environment (" + v.goEnvVersion +
			") is outside of the supported range (" + minGoVersion + "-" + maxGoVersion + "). " +
			"Requesting the minimum supported version of Go (" + minGoVersion + ")."
		version = minGoVersion
		diagnostics.EmitGoModVersionTooLowAndEnvVersionUnsupported(msg)
	} else {
		// The version of Go that is installed is supported. The version in the `go.mod` file is
		// below the supported range. We do not install a version of Go.
		msg = "The version of Go installed in the environment (" + v.goEnvVersion +
			") is supported and is high enough for the version found in the `go.mod` file (" +
			v.goModVersion + "). Not requesting any version of Go."
		version = ""
		diagnostics.EmitGoModVersionTooLowAndEnvVersionSupported(msg)
	}

	return msg, version
}

// Assuming `v.goModVersion` is in the supported range, emit a diagnostic and return the version
// to install, or the empty string if we should not attempt to install a version of Go.
func getVersionWhenGoModVersionSupported(v versionInfo) (msg, version string) {
	if !v.goEnvVersionFound {
		// There is no Go version installed. The version in the `go.mod` file is supported.
		// We install the version from the `go.mod` file.
		msg = "No version of Go installed. Requesting the version of Go found in the `go.mod` " +
			"file (" + v.goModVersion + ")."
		version = v.goModVersion
		diagnostics.EmitGoModVersionSupportedAndNoGoEnv(msg)
	} else if outsideSupportedRange(v.goEnvVersion) {
		// The version of Go that is installed is outside of the supported range. The version in
		// the `go.mod` file is supported. We install the version from the `go.mod` file.
		msg = "The version of Go installed in the environment (" + v.goEnvVersion +
			") is outside of the supported range (" + minGoVersion + "-" + maxGoVersion + "). " +
			"Requesting the version of Go from the `go.mod` file (" +
			v.goModVersion + ")."
		version = v.goModVersion
		diagnostics.EmitGoModVersionSupportedAndGoEnvUnsupported(msg)
	} else if semver.Compare("v"+v.goModVersion, "v"+v.goEnvVersion) > 0 {
		// The version of Go that is installed is supported. The version in the `go.mod` file is
		// supported and is higher than the version that is installed. We install the version from
		// the `go.mod` file.
		msg = "The version of Go installed in the environment (" + v.goEnvVersion +
			") is lower than the version found in the `go.mod` file (" + v.goModVersion +
			"). Requesting the version of Go from the `go.mod` file (" + v.goModVersion + ")."
		version = v.goModVersion
		diagnostics.EmitGoModVersionSupportedHigherGoEnv(msg)
	} else {
		// The version of Go that is installed is supported. The version in the `go.mod` file is
		// supported and is lower than or equal to the version that is installed. We do not install
		// a version of Go.
		msg = "The version of Go installed in the environment (" + v.goEnvVersion +
			") is supported and is high enough for the version found in the `go.mod` file (" +
			v.goModVersion + "). Not requesting any version of Go."
		version = ""
		diagnostics.EmitGoModVersionSupportedLowerEqualGoEnv(msg)
	}

	return msg, version
}

// Check the versions of Go found in the environment and in the `go.mod` file, and return a
// version to install. If the version is the empty string then no installation is required.
// We never return a version of Go that is outside of the supported range.
//
// +-----------------------+-----------------------+-----------------------+-----------------------------------------------------+------------------------------------------------+
// | Found in go.mod >     | *None*                | *Below min supported* | *In supported range*                                | *Above max supported                           |
// | Installed \/          |                       |                       |                                                     |                                                |
// |-----------------------|-----------------------|-----------------------|-----------------------------------------------------|------------------------------------------------|
// | *None*                | Install max supported | Install min supported | Install version from go.mod                         | Install max supported                          |
// | *Below min supported* | Install max supported | Install min supported | Install version from go.mod                         | Install max supported                          |
// | *In supported range*  | No action             | No action             | Install version from go.mod if newer than installed | Install max supported if newer than installed  |
// | *Above max supported* | Install max supported | Install min supported | Install version from go.mod                         | No action                                      |
// +-----------------------+-----------------------+-----------------------+-----------------------------------------------------+------------------------------------------------+
func getVersionToInstall(v versionInfo) (msg, version string) {
	if !v.goModVersionFound {
		return getVersionWhenGoModVersionNotFound(v)
	}

	if aboveSupportedRange(v.goModVersion) {
		return getVersionWhenGoModVersionTooHigh(v)
	}

	if belowSupportedRange(v.goModVersion) {
		return getVersionWhenGoModVersionTooLow(v)
	}

	return getVersionWhenGoModVersionSupported(v)
}

// Output some JSON to stdout specifying the version of Go to install, unless `version` is the
// empty string.
func outputEnvironmentJson(version string) {
	var content string
	if version == "" {
		content = `{ "go": {} }`
	} else {
		content = `{ "go": { "version": "` + version + `" } }`
	}
	_, err := fmt.Fprint(os.Stdout, content)

	if err != nil {
		log.Println("Failed to write environment json to stdout: ")
		log.Println(err)
	}
}

// Get the version of Go to install and output it to stdout as json.
func IdentifyEnvironment() {
	var v versionInfo
	workspaces := project.GetWorkspaceInfo(false)

	// Remove temporary extractor files (e.g. auto-generated go.mod files) when we are done
	defer project.RemoveTemporaryExtractorFiles()

	// Find the greatest Go version required by any of the workspaces.
	greatestGoVersion := project.RequiredGoVersion(&workspaces)
	v.goModVersion, v.goModVersionFound = greatestGoVersion.Version, greatestGoVersion.Found

	// Find which, if any, version of Go is installed on the system already.
	v.goEnvVersionFound = toolchain.IsInstalled()
	if v.goEnvVersionFound {
		v.goEnvVersion = toolchain.GetEnvGoVersion()[2:]
	}

	// Determine which version of Go we should recommend to install.
	msg, versionToInstall := getVersionToInstall(v)
	log.Println(msg)

	outputEnvironmentJson(versionToInstall)
}
