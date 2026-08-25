package logs

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRemoveUpCalls(t *testing.T) {
	var (
		input  bytes.Buffer
		output bytes.Buffer
	)
	expected := "line 1\nline 2\n"

	fmt.Fprintln(&input, "line 1")
	fmt.Fprintln(&input, "{\"path\":\"/up\"}")
	fmt.Fprintln(&input, "line 2")

	if err := RemoveUpCalls(&input, &output); err != nil {
		t.Fatal(err)
	}

	if expected != output.String() {
		t.Errorf("expected %q, got %q", expected, output.String())
	}
}
