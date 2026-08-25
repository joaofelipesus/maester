package logs

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func RemoveUpCalls(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, "\"path\":\"/up\"") {
			lineWithBreak := fmt.Sprintf("%s\n", line)
			writer.Write([]byte(lineWithBreak))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
