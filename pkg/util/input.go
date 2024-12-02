package util

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
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
	return DoEachRowBuff(r, f)
}

func DoEachRowReader(r io.Reader, f func(row []byte, nr int) error) (err error) {
	return DoEachRowBuff(bufio.NewReader(r), f)
}

func DoEachRowFile(filename string, rowFunc func(row []byte, nr int) error) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			err = errors.Join(err, err2)
		}
	}()
	return DoEachRowBuff(bufio.NewReader(f), rowFunc)
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
