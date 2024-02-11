package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for fExt, fSize := range files {
		for j := 1; j <= fSize; j++ {
			builtFname := fmt.Sprintf("file%d%s", j, fExt)
			fpath := filepath.Join(tempDir, builtFname)
			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{name: "DeleteExtensionNoMatch", cfg: config{
			ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     0,
			nNoDelete:   10,
			expected:    "",
		},
		{name: "DeleteExtensionMatch", cfg: config{
			ext: ".log", del: true},
			extNoDelete: "",
			nDelete:     10,
			nNoDelete:   0,
			expected:    "",
		},
		{name: "DeleteExtensionMixed", cfg: config{
			ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     5,
			nNoDelete:   5,
			expected:    ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)

			tc.cfg.wLog = &logBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})
			defer cleanup()

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()
			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead.\n", tc.expected, res)
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}
			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("Expected %d got %d instead.\n", tc.nNoDelete, len(filesLeft))
			}

			expLogLines := tc.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("Expected %d got %d instead.\n", expLogLines, len(lines))
			}
		})
	}
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{
			name:     "NoFilter",
			root:     "testdata",
			cfg:      config{ext: "", size: 0, list: true},
			expected: "testdata/dir.log\ntestdata/dir2/text.txt\n"},
		{
			name:     "FilterExtensionMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 0, list: true},
			expected: "testdata/dir.log\n"},
		{
			name:     "FilterExtensionSizeMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 10, list: true},
			expected: "testdata/dir.log\n"},
		{
			name:     "FilterExtensionSizeNoMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 20, list: true},
			expected: ""},
		{
			name:     "FilterExtensionNoMatch",
			root:     "testdata",
			cfg:      config{ext: ".gz", size: 0, list: true},
			expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}
			res := buffer.String()
			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead.", tc.expected, res)
			}
		})
	}
}
