package writer

import (
	"fmt"
	"io"
)

func WriteLines(wr io.Writer, lines []string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintln(wr, line); err != nil {
			return err
		}
	}
	return nil
}
