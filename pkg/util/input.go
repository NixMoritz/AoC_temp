package util

import (
	"bufio"
	"bytes"
	"io"
)

func DoEachRowAll(input []byte, f func(row []byte, rows [][]byte, nr, total int) error) (err error) {
	b := bytes.Split(input, []byte{'\n'})
	return DoEachRowAllBytes(b, f)
}

func DoEachRowAllBytes(rows [][]byte, f func(row []byte, rows [][]byte, nr, total int) error) (err error) {
	var (
		total = len(rows)
	)
	for i := 0; i < total; i++ {
		err = f(rows[i], rows, i, total)
		if err != nil {
			return
		}
	}
	return
}

func DoEachRowBytes(input []byte, f func(row []byte, nr int) error) (err error) {
	r := bufio.NewReader(bytes.NewReader(input))
	return DoEachRowReader(r, f)
}

func DoEachRowReader(r io.Reader, f func(row []byte, nr int) error) (err error) {
	return DoEachRowReader(bufio.NewReader(r), f)
}

func DoEachRowBuff(r *bufio.Reader, f func(row []byte, nr int) error) (err error) {
	var (
		row  []byte
		done bool
	)
	for i := 0; !done; i++ {
		row, err = r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return
		}
		if err == io.EOF {
			done = true
		}
		err = f(row, i)
		if err != nil {
			return
		}
	}
	return
}
