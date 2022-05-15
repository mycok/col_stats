package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"testing/iotest"
)

func TestOperations(t *testing.T) {
	data := [][]float64{
		{10, 20, 15, 30, 45, 50, 100, 30},
		{5.5, 8, 2.2, 9.75, 8.45, 3, 2.5, 10.25, 4.75, 6.1, 7.67, 12.287, 5.47},
		{-10, -20},
		{102, 37, 44, 57, 67, 129},
	}

	testcases := []struct {
		name      string
		operation statsFunc
		expected  []float64
	}{
		{
			name:      "sum",
			operation: sum,
			expected:  []float64{300, 85.927, -30, 436},
		},
		{
			name:      "avg",
			operation: avg,
			expected:  []float64{37.5, 6.609769230769231, -15, 72.666666666666666},
		},
	}

	for _, tc := range testcases {
		for k, expValue := range tc.expected {
			name := fmt.Sprintf("%sData%d", tc.name, k)

			t.Run(name, func(t *testing.T) {
				result := tc.operation(data[k])

				if result != expValue {
					t.Errorf("Expected %g, Got: %g instead", expValue, result)
				}
			})
		}
	}

}

func TestCsvToFloat(t *testing.T) {
	csvData := `IP Address,Requests,Response Time
				192.168.0.199,2056,236
				192.168.0.88,899,220
				192.168.0.199,3054,226
				192.168.0.100,4133,218
				192.168.0.199,950,238`

	testcases := []struct {
		name     string
		column   int
		expected []float64
		expErr   error
		r        io.Reader
	}{
		{
			name:     "column2",
			column:   2,
			expected: []float64{2056, 899, 3054, 4133, 950},
			expErr:   nil,
			r:        bytes.NewBufferString(csvData),
		},
		{
			name:     "column3",
			column:   3,
			expected: []float64{236, 220, 226, 218, 238},
			expErr:   nil,
			r:        bytes.NewBufferString(csvData),
		},
		{
			name:     "FailsOnRead",
			column:   1,
			expected: nil,
			expErr:   iotest.ErrTimeout,
			r:        iotest.TimeoutReader(bytes.NewReader([]byte{0})),
		},
		{
			name:     "FailsOnNotANumber",
			column:   1,
			expected: nil,
			expErr:   ErrNotNumber,
			r:        bytes.NewBufferString(csvData),
		},
		{
			name:     "FailsOnInvalidColumn",
			column:   5,
			expected: nil,
			expErr:   ErrInvalidColumn,
			r:        bytes.NewBufferString(csvData),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			results, err := csvToFloat(tc.r, tc.column)

			// Handle cases where an error is expected.
			if tc.expErr != nil {
				// If we expect an error but get nil, we log an error msg and return
				// the function.
				if err == nil {
					t.Errorf("Expected error, Got nil instead")

					return
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %q, Got: %q instead", tc.expErr, err)
				}

				return
			}

			// Handle cases where an error is not expected.
			if err != nil {
				t.Errorf("Unexpected error: %q", err)

				return
			}

			for i, exp := range tc.expected {
				if exp != results[i] {
					t.Errorf("Expected: %g, Got: %g instead", exp, results[i])
				}
			}
		})
	}
}
