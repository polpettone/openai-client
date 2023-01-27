package cmd

import (
	"testing"
)

func TestAdd(t *testing.T) {

	contextMemory := NewContextMemory(3)

	contextMemory.Add(&Entry{value: "A", tokens: 1})
	contextMemory.Add(&Entry{value: "B", tokens: 1})
	contextMemory.Add(&Entry{value: "C", tokens: 1})

	result := contextMemory.All()
	expected := "A\nB\nC\n"
	if result != expected {
		t.Errorf("Wanted %s, got %s", expected, result)
	}
}

func TestMaxTokens(t *testing.T) {

	contextMemory := NewContextMemory(3)

	contextMemory.Add(&Entry{value: "A", tokens: 1})
	contextMemory.Add(&Entry{value: "B", tokens: 1})
	contextMemory.Add(&Entry{value: "C", tokens: 1})
	contextMemory.Add(&Entry{value: "D", tokens: 3})

	result := contextMemory.All()
	expected := "D\n"
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

func TestTokenCount(t *testing.T) {
	contextMemory := NewContextMemory(10)

	contextMemory.Add(&Entry{value: "A", tokens: 1})
	contextMemory.Add(&Entry{value: "B", tokens: 2})
	contextMemory.Add(&Entry{value: "C", tokens: 3})

	result := contextMemory.TokenCount()

	expected := 6
	if result != expected {
		t.Errorf("Wanted %d, got %d", expected, result)
	}
}
