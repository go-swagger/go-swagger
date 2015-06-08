package swag

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// GOPATHKey represents the env key for gopath
	GOPATHKey = "GOPATH"
)

// FindInSearchPath finds a package in a provided lists of paths
func FindInSearchPath(searchPath, pkg string) string {
	pathsList := filepath.SplitList(searchPath)
	for _, path := range pathsList {
		if evaluatedPath, err := filepath.EvalSymlinks(filepath.Join(path, "src", pkg)); err == nil {
			if _, err := os.Stat(evaluatedPath); err == nil {
				return evaluatedPath
			}
		}
	}
	return ""
}

// FindInGoSearchPath finds a package in the $GOPATH:$GOROOT
func FindInGoSearchPath(pkg string) string {
	return FindInSearchPath(FullGoSearchPath(), pkg)
}

// FullGoSearchPath gets the search paths for finding packages
func FullGoSearchPath() string {
	allPaths := os.Getenv(GOPATHKey)
	if allPaths != "" {
		allPaths = strings.Join([]string{allPaths, runtime.GOROOT()}, ":")
	} else {
		allPaths = runtime.GOROOT()
	}
	return allPaths
}
