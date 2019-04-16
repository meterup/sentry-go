package sentry

import (
	"bytes"
	"io/ioutil"
	"sync"
)

type SourceReader struct {
	mu    sync.Mutex
	cache map[string][][]byte
}

func NewSourceReader() SourceReader {
	return SourceReader{
		cache: make(map[string][][]byte),
	}
}

func (sr *SourceReader) ReadContextLines(filename string, line, context int) ([][]byte, int) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	lines, ok := sr.cache[filename]

	if !ok {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			sr.cache[filename] = nil
			return nil, 0
		}
		lines = bytes.Split(data, []byte{'\n'})
		sr.cache[filename] = lines
	}

	return calculateContextLines(lines, line, context)
}

func calculateContextLines(lines [][]byte, line, context int) ([][]byte, int) {
	// Stacktrace lines are 1-indexed, slices are 0-indexed
	line--

	if lines == nil || line >= len(lines) || line < 0 {
		return nil, 0
	}

	if context < 0 {
		context = 0
	}

	start := line - context

	if start < 0 {
		start = 0
	}

	end := line + context + 1

	if end > len(lines) {
		end = len(lines)
	}

	return lines[start:end], (end - start)
}
