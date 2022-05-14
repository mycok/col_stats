package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	op := flag.String("op", "", "Arithmetic operation to execute on the csv file data")
	column := flag.Int("col", 1, "column on which to execute the operation")
	flag.Parse()

	err := run(flag.Args(), *op, *column, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

}

func run(filenames []string, op string, column int, w io.Writer) error {
	var opFunc statsFunc

	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	mergedColValues := make([]float64, 0)

	for _, fname := range filenames {
		// Open the file for reading.
		f, err := os.Open(fname)
		if err != nil {
			return fmt.Errorf("Cannot open file: %w", err)
		}

		// Parse the csv column data into a slice of float64 numbers.
		data, err := csvToFloat(f, column)
		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}

		mergedColValues = append(mergedColValues, data...)
	}

	_, err := fmt.Fprintln(w, opFunc(mergedColValues))

	return err
}
