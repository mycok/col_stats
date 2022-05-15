package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	total := 0.0

	for _, v := range data {
		total += v
	}

	return total
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func csvToFloat(r io.Reader, column int) ([]float64, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	column--

	var data []float64

	for i := 0; ; i++ {
		row, err := cr.Read();
		if err != nil {
			// If there's no more rows to read from the file.
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("cannot read data from file: %w", err)
		}

		// Read and discard the first row of file.(this row normally contains titles).
		if i == 0 {
			continue
		}

		// Check the number of columns in the row against the provided column argument.
		// having less columns in the row than the provided column argument makes the column
		// argument invalid.
		if len(row) <= column {
			return nil, fmt.Errorf("%w: file has only %d columns", ErrInvalidColumn, len(row))
		}

		// Try to convert row value for the provided column from a string to a float.
		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		data = append(data, v)
	}

	return data, nil
}
