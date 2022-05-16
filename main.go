package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
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

	filesCh := make(chan string)
	resultCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	// Loop through all files passing them to the files channel where each
	// will be processed by a worker goroutine when one is available.
	go func() {
		defer close(filesCh)

		for _, fname := range filenames {
			filesCh <- fname
		}
	}()

	// Loop through all filenames and process each file concurrently.
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for fname := range filesCh {
				// Open the file for reading.
				f, err := os.Open(fname)
				if err != nil {
					errCh <- fmt.Errorf("cannot open file: %w", err)
					// Write the error to the error channel and  stop execution from
					// continuing by returning from the goroutine.
					return
				}

				// Parse the csv column data into a slice of float64 numbers.
				data, err := csvToFloat(f, column)
				if err != nil {
					errCh <- err
				}

				if err := f.Close(); err != nil {
					errCh <- err
				}

				resultCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resultCh:
			mergedColValues = append(mergedColValues, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(w, opFunc(mergedColValues))
			return err
		}
	}
}
