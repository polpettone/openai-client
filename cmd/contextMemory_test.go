package cmd

import (
	"testing"
)

func TestAdd(t *testing.T) {

	contextMemory := NewContextMemory("id", 3)

	contextMemory.Add(&Entry{Value: "A", Tokens: 1})
	contextMemory.Add(&Entry{Value: "B", Tokens: 1})
	contextMemory.Add(&Entry{Value: "C", Tokens: 1})

	result := contextMemory.All()
	expected := "A\nB\nC\n"
	if result != expected {
		t.Errorf("Wanted %s, got %s", expected, result)
	}
}

func TestMaxTokens(t *testing.T) {

	contextMemory := NewContextMemory("id", 3)

	contextMemory.Add(&Entry{Value: "A", Tokens: 1})
	contextMemory.Add(&Entry{Value: "B", Tokens: 1})
	contextMemory.Add(&Entry{Value: "C", Tokens: 1})
	contextMemory.Add(&Entry{Value: "D", Tokens: 3})

	result := contextMemory.All()
	expected := "D\n"
	if result != expected {
		t.Errorf("Wanted %s, got %s", expected, result)
	}
}

func TestClear(t *testing.T) {

	contextMemory := NewContextMemory("id", 10)
	contextMemory.Reset()
	result := contextMemory.All()
	expected := ""

	if result != expected {
		t.Errorf("Wanted %s, got |%s|", expected, result)
	}
}

func TestTokenCount(t *testing.T) {
	contextMemory := NewContextMemory("id", 10)

	contextMemory.Add(&Entry{Value: "A", Tokens: 1})
	contextMemory.Add(&Entry{Value: "B", Tokens: 2})
	contextMemory.Add(&Entry{Value: "C", Tokens: 3})

	result := contextMemory.TokenCount()

	expected := 6
	if result != expected {
		t.Errorf("Wanted %d, got %d", expected, result)
	}
}
