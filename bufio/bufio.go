package bufio

import (
	"bufio"
	"io"

	"github.com/nexuer/utils/unsafe"
)

// ReadLineFunc read the io.Reader line by line and call f(c) to process each line of string
func ReadLineFunc(reader io.Reader, f func(num int, line string) error) error {
	scanner := bufio.NewScanner(reader)
	num := 0
	for scanner.Scan() {
		num++
		if err := f(num, unsafe.BytesToString(scanner.Bytes())); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// ReadLineBytesFunc read the io.Reader line by line and call f(c) to process each line of bytes
func ReadLineBytesFunc(reader io.Reader, f func(num int, line []byte) error) error {
	scanner := bufio.NewScanner(reader)
	num := 0
	for scanner.Scan() {
		num++
		if err := f(num, scanner.Bytes()); err != nil {
			return err
		}
	}
	return scanner.Err()
}
