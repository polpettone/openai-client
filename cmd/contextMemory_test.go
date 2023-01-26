package cmd

import (
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {

	contextMemory := NewContextMemory(3)

	contextMemory.Add("A")
	contextMemory.Add("B")
	contextMemory.Add("C")
	contextMemory.Add("D")
	contextMemory.Add("E")

	result := strings.Join(contextMemory.All(), ",")
	expected := "D,E,C"
	if result != expected {
		t.Errorf("Wanted %s, got %s", expected, result)
	}

}
