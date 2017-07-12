package depbleed

import (
	"fmt"
	"go/token"
	"go/types"
)

// Leak represents a leaking type.
type Leak struct {
	Object   types.Object
	Position token.Position
	err      error
}

// Error constructs an error string.
func (l Leak) Error() string {
	return fmt.Sprintf("%s: %s", l.Object.Name(), l.err)
}

// Leaks represents a slice of Leak instances.
type Leaks []Leak

// Len gives the length of the leak slice/
func (slice Leaks) Len() int {
	return len(slice)
}

// Less returns if slice[i] should be before slice[j].
func (slice Leaks) Less(i, j int) bool {
	return slice[i].Position.Filename < slice[j].Position.Filename ||
		(slice[i].Position.Filename == slice[j].Position.Filename &&
			slice[i].Position.Line < slice[j].Position.Line) ||
		(slice[i].Position.Filename == slice[j].Position.Filename &&
			slice[i].Position.Line == slice[j].Position.Line &&
			slice[i].Position.Column < slice[j].Position.Column)
}

// Swap swaps two elements.
func (slice Leaks) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
