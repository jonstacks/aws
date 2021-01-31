package views

import "io"

// View is an interface which supports the Render method.
type View interface {
	Render(io.Writer)
}
