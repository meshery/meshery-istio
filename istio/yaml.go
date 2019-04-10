package istio

import (
	"bufio"
	"bytes"
	"io"
)

// YAMLDecoder reads chunks of objects and returns ErrShortBuffer if
// the data is not sufficient.
// borrowed from APIMachinery
type YAMLDecoder struct {
	r         io.ReadCloser
	scanner   *bufio.Scanner
	remaining []byte
}

// NewDocumentDecoder decodes YAML documents from the provided
// stream in chunks by converting each document (as defined by
// the YAML spec) into its own chunk. io.ErrShortBuffer will be
// returned if the entire buffer could not be read to assist
// the caller in framing the chunk.
func NewDocumentDecoder(r io.ReadCloser) io.ReadCloser {
	b := make([]byte, 4096)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(b, 256*1024) // overriding: the size of the buffer used was small when loading large sections from istio deployment yaml
	scanner.Split(splitYAMLDocument)
	return &YAMLDecoder{
		r:       r,
		scanner: scanner,
	}
}

// Read reads the previous slice into the buffer, or attempts to read
// the next chunk.
func (d *YAMLDecoder) Read(data []byte) (n int, err error) {
	left := len(d.remaining)
	if left == 0 {
		// return the next chunk from the stream
		if !d.scanner.Scan() {
			err := d.scanner.Err()
			if err == nil {
				err = io.EOF
			}
			return 0, err
		}
		out := d.scanner.Bytes()
		d.remaining = out
		left = len(out)
	}

	// fits within data
	if left <= len(data) {
		copy(data, d.remaining)
		d.remaining = nil
		return left, nil
	}

	// caller will need to reread
	copy(data, d.remaining[:len(data)])
	d.remaining = d.remaining[len(data):]
	return len(data), io.ErrShortBuffer
}

func (d *YAMLDecoder) Close() error {
	return d.r.Close()
}

const yamlSeparator = "\n---"
const separator = "---"

// splitYAMLDocument is a bufio.SplitFunc for splitting YAML streams into individual documents.
func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// We have a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// we can't read any more characters
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
