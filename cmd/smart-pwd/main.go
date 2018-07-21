// Command smart-pwd intelligently prints the current working directory.
//
// The inital intended use case is a shell prompt, so speed is important. The
// intended user experience is that just enough information is provided to
// indicate the current location, while maintaining a level of brevity
// appropriate for a prompt.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// This code was originally a bash function from:
// https://github.com/kujenga/init/blob/526cf35e1aecfdc614fcff45f39dac0f91f11b73/lib/profile.sh#L33-L53
//
// smart_pwd() {
//     # set -x
//     # git prefix with repo name or pwd.
//     GIT=$(git rev-parse --show-prefix 2> /dev/null)
//     if [ $? -eq 0 ]; then
//         DIR="$(basename "$(git rev-parse --show-toplevel)")/$GIT"
//     else
//         DIR="$(pwd)"
//     fi
//     # strip trailing slash, replace $HOME with ~, shorten all but last.
//     echo "$DIR" | \
//         sed 's?/$??g' | \
//         sed "s?$HOME?~?g" | \
//         perl -F/ -ane 'print join( "/", map { $i++ < @F - 1 ?  substr $_,0,1 : $_ } @F)'
//     # set +x
// }

var pathSeparator = string([]rune{os.PathSeparator})

func main() {
	fmt.Println(getSmart())
}

func getSmart() string {
	dir := getDir()
	return smartenUp(dir)
}

// getDir gets the string representing the current directory in it's fully
// qualified form with no abbreviations.
func getDir() string {
	wd, err := os.Getwd()
	check(err, "getting current working directory")
	// Walk up the directory tree looking for a .git directory, for which
	// we would customize the printout.
	dir := wd
	// We terminate the loop when we are at the root, indicated by a
	// following spash as documented:
	// https://golang.org/pkg/path/filepath/#Dir
	for !strings.HasSuffix(dir, pathSeparator) {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			// If we find a directory with a .git directory, we
			// return the relative path from it's containing
			// directory.
			gitParent := filepath.Dir(dir)
			rel, err := filepath.Rel(gitParent, wd)
			check(err, "getting relative path to working directory")
			return rel
		}
		dir = filepath.Dir(dir)
	}
	return wd
}

// smartenUp takes a path and shortens all the components to just their
// first character, except for the last component. This is intended to present
// enough information that viewing the string would give the reader an
// indication where they are in their directory tree with minimal length,
// aiming at use in shell prompts
func smartenUp(s string) string {
	// Strip trailing slashes from the path.
	s = strings.TrimSuffix(s, "/")
	// Replace $HOME with ~ for brevity.
	home := os.Getenv("HOME")
	if strings.HasPrefix(s, home) {
		s = strings.Replace(s, home, "~", 1)
	}
	// Shorten all path components but the last.
	components := strings.Split(s, pathSeparator)
	for i := range components {
		if i == len(components)-1 {
			break
		}
		if components[i] != "" {
			components[i] = string([]rune(components[i])[0])
		}
	}
	// Use strings.Join here to preserve preceeding slash
	return strings.Join(components, pathSeparator)
}

// check exits the program with the specified message if the error is non-nil.
func check(err error, msg string) {
	if err != nil {
		log.Fatalf("smart-pwd: %s: %v", msg, err)
	}
}
