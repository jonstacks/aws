package cmd

import (
	"fmt"
	"os"
)

// HandleError is a common error handler for our commands.
func HandleError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(2)
	}
}
