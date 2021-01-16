// The rbxattr package implements the serialized format of Roblox's instance
// attributes.
package rbxattr

import (
	"fmt"
	"io"
)

// Model is a low-level model of Roblox's instance attribute format.
type Model struct {
	Value ValueDictionary
}

// ReadFrom decodes bytes from r, setting Value on success.
func (f *Model) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = f.Value.ReadFrom(r)
	if err != nil {
		err = fmt.Errorf("format: %w", err)
	}
	return n, err
}

// WriteTo encodes Value into bytes written to w.
func (f *Model) WriteTo(w io.Writer) (n int64, err error) {
	n, err = f.Value.WriteTo(w)
	if err != nil {
		err = fmt.Errorf("format: %w", err)
	}
	return n, err
}
