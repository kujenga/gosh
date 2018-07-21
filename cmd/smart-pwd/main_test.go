package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t testing.TB) (string, func()) {
	// mock home directory
	home, err := ioutil.TempDir("", "TestGetDir-home")
	require.NoError(t, err)
	home, err = filepath.EvalSymlinks(home)
	require.NoError(t, err)
	// mock home directory
	os.Setenv("HOME", home)

	// setup directories to test against.
	for _, dir := range []string{
		"aaa/.git",
		"aaa/bbb/ccc",
		"xxx/yyy",
	} {
		require.NoError(t, os.MkdirAll(filepath.Join(home, dir), 0777))
	}
	return home, func() { os.RemoveAll(home) }
}

var sink string

func BenchmarkGetSmartPlain(b *testing.B) {
	home, cleanup := setup(b)
	defer cleanup()

	var s string
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		os.Chdir(filepath.Join(home, "aaa/bbb/ccc"))
		s = getSmart()
	}
	b.StopTimer()
	sink = s
}

func BenchmarkGetSmartGitRepo(b *testing.B) {
	home, cleanup := setup(b)
	defer cleanup()

	var s string
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		os.Chdir(filepath.Join(home, "xxx/yyy"))
		s = getSmart()
	}
	b.StopTimer()
	sink = s
}

func TestGetDir(t *testing.T) {
	home, cleanup := setup(t)
	defer cleanup()

	tcs := []struct {
		wd    string
		raw   string
		smart string
	}{
		{
			wd:    filepath.Join(home, "aaa"),
			raw:   "aaa",
			smart: "aaa",
		},
		{
			wd:    filepath.Join(home, "aaa/bbb/ccc"),
			raw:   "aaa/bbb/ccc",
			smart: "a/b/ccc",
		},
		{
			wd:    filepath.Join(home, "xxx"),
			raw:   filepath.Join(home, "xxx"),
			smart: "~/xxx",
		},
		{
			wd:    filepath.Join(home, "xxx/yyy"),
			raw:   filepath.Join(home, "xxx/yyy"),
			smart: "~/x/yyy",
		},
	}
	for _, tc := range tcs {
		// mock current working directory
		os.Chdir(tc.wd)
		// test output
		out := getDir()
		assert.Equal(t, tc.raw, out, "raw output should match expected")
		smart := smartenUp(out)
		assert.Equal(t, tc.smart, smart, "smart output should match expected")
	}
}

func TestSmartenUp(t *testing.T) {
	// mock home directory
	os.Setenv("HOME", "/home")

	tcs := []struct {
		in     string
		expect string
	}{
		{
			in:     "aaa/bbb/ccc",
			expect: "a/b/ccc",
		},
		{
			in:     "/home/aaa/bbb/ccc",
			expect: "~/a/b/ccc",
		},
	}
	for _, tc := range tcs {
		out := smartenUp(tc.in)
		assert.Equal(t, tc.expect, out, "smart version should match expected")
	}
}
