// Command smart-pwd intelligently prints the current working directory.
//
// The inital intended use case is a shell prompt, so speed is important. The
// intended user experience is that just enough information is provided to
// indicate the current location, while maintaining a level of brevity
// appropriate for a prompt.
//
// There are likely some future optimizations that can be done, in particular
// the two calls to git subshells can likely be optimized.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
	fmt.Println(smartenUp(getRaw()))
}

func getRaw() string {
	gitPrefix, err := gitRevParse("--show-prefix")
	if err == nil {
		gitTopLevel, err := gitRevParse("--show-toplevel")
		check(err, "checking top level git directory")
		// return the path with the git repo as the root.
		_, file := filepath.Split(gitTopLevel)
		return filepath.Join(file, gitPrefix)
	}
	// fallback to the current working directory.
	wd, err := os.Getwd()
	check(err, "getting current working directory")
	return wd
}

func smartenUp(s string) string {
	// Strip trailing slashes from the path.
	s = strings.TrimSuffix(s, "/")
	// Replace $HOME with ~ for brevity.
	home := os.Getenv("HOME")
	s = strings.Replace(s, home, "~", 1)
	// Shorten all path components but the last.
	components := strings.Split(s, pathSeparator)
	for i := range components {
		if i == len(components)-1 {
			break
		}
		components[i] = string([]rune(components[i])[0])
	}
	return filepath.Join(components...)
}

// gitRevParse is a helper function for returning the output of `git rev-parse`
// with various flags passed in.
func gitRevParse(flag string) (string, error) {
	out, err := exec.Command("git", "rev-parse", flag).CombinedOutput()
	if err != nil {
		return "", err
	}
	s := string(out)
	s = strings.TrimSpace(s)
	return s, nil
}

// check exits the program with the specified message if the error is non-nil.
func check(err error, msg string) {
	if err != nil {
		log.Fatalf("smart-pwd: %s: %v", msg, err)
	}
}
