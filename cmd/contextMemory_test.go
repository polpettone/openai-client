package cmd

import (
	"testing"
)

func TestAdd(t *testing.T) {

	contextMemory := NewContextMemory(3)

	contextMemory.Add("A")
	contextMemory.Add("B")
	contextMemory.Add("C")
	contextMemory.Add("D")
	contextMemory.Add("E")

	result := contextMemory.All()
	expected := "D\nE\nC\n"
	if result != expected {
		t.Errorf("Wanted %s, got %s", expected, result)
	}
}

func TestClear(t *testing.T) {

	contextMemory := NewContextMemory(10)

	contextMemory.Reset()

	result := contextMemory.All()
	expected := ""
	if result != expected {
		t.Errorf("Wanted %s, got |%s|", expected, result)
	}

}
