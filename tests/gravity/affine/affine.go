package affine

import (
	"golang.org/x/image/math/f64"
)

// NewIdentity returns the Identity matrix, i.e. a transformation
// where no change is made.
func NewIdentity() f64.Aff3 {
	return f64.Aff3{
		1, 0, 1,
		0, 1, 1,
	}
}
