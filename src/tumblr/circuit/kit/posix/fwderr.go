package posix

import (
	"bufio"
	"io"
	"io/ioutil"
)

func ForwardStderrBatch(stderr io.ReadCloser) {
	go func() {
		all, _ := ioutil.ReadAll(stderr)
		println(string(all))
	}()
}

func ForwardStderr(prefix string, stderr io.ReadCloser) {
	go func() {
		r := bufio.NewReader(stderr)
		for {
			line, err := r.ReadString('\n')
			if line != "" {
				println(prefix, line)
			}
			if err != nil {
				break
			}
		}
		stderr.Close()
	}()
}

// lineReader is intended to read lines delimited by any contiguous block of \r's and \n's
type lineReader struct {
	io.Reader
	fill []byte
}

func newLineReader(r io.Reader) *lineReader {
	return &lineReader{Reader: r}
}

func (lr *lineReader) ReadLine() (line string, isprefix bool, err error) {
	panic("buggy")
	for {
		// Read a new nibble
		nibble := make([]byte, 1e3)
		n, err := lr.Reader.Read(nibble)
		nibble = nibble[:n]

		if n > 0 {
			// Look for first occurence of separator
			ex := -1
			for i, b := range nibble {
				if isSep(b) {
					ex = i
					break
				}
			}
			// If no separator was reached
			if ex < 0 {
				lr.fill = append(lr.fill, nibble...)
				if len(lr.fill) > 1e4 {
					var r []byte
					r, lr.fill = lr.fill, nil
					return string(r), true, err
				}
				continue
			}

			var r []byte
			r, lr.fill = append(lr.fill, nibble[:ex]...), nil
			
			// Skip over separators after first
			var nsep int
			for _, b := range nibble[ex:] {
				if !isSep(b) {
					break
				}
				nsep++
			}
			lr.fill = nibble[ex+nsep:]

			return string(r), false, nil
		}

		// Return whatever is in fill
		var r []byte
		r, lr.fill = lr.fill, nil
		return string(r), true, err
	}
	panic("unreach")
}

func isSep(b byte) bool {
	return b == '\r' || b == '\n'
}
