package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testcases := []struct {
		name      string
		column    int
		operation string
		expected  string
		files     []string
		expErr    error
	}{
		{
			name:      "RunAvg1File",
			column:    3,
			operation: "avg",
			expected:  "227.6\n",
			files:     []string{"testdata/example.csv"},
			expErr:    nil,
		},
		{
			name:      "RunAvgForMultipleFiles",
			column:    3,
			operation: "avg",
			expected:  "233.84\n",
			files:     []string{"testdata/example.csv", "testdata/example2.csv"},
			expErr:    nil,
		},
		{
			name:      "RunFailsOnRead",
			column:    2,
			operation: "avg",
			expected:  "",
			files:     []string{"testdata/example.csv", "testdata/fakefile.csv"},
			expErr:    os.ErrNotExist,
		},
		{
			name:      "RunFailsOnColumn",
			column:    0,
			operation: "avg",
			expected:  "",
			files:     []string{"testdata/example.csv"},
			expErr:    ErrInvalidColumn,
		},
		{
			name:      "RunFailsOnFilesPassed",
			column:    2,
			operation: "avg",
			expected:  "",
			files:     []string{},
			expErr:    ErrNoFiles,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			err := run(tc.files, tc.operation, tc.column, &buffer)
			// Handle cases when an error is expected.
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected an error but instead got nil")

					return
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %q, Got: %q instead", tc.expected, err)
				}

				return
			}

			// Handle cases where an error is not expected.
			if err != nil {
				t.Fatalf("Unexpected error: %q", err)

				return
			}

			if tc.expected != buffer.String() {
				t.Errorf("Expected %s, Got: %s instead", tc.expected, buffer.String())
			}
		})
	}
}

func BenchmarkRun(b *testing.B) {
	filenames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := run(filenames, "avg", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}
