package cmd

import (
	"strings"
	"testing"
)

func TestConfirm_Yes(t *testing.T) {
	origStdin := stdin
	defer func() { stdin = origStdin }()

	stdin = strings.NewReader("y\n")
	if !confirm("Continue?") {
		t.Error("expected true for 'y'")
	}
}

func TestConfirm_YesFull(t *testing.T) {
	origStdin := stdin
	defer func() { stdin = origStdin }()

	stdin = strings.NewReader("yes\n")
	if !confirm("Continue?") {
		t.Error("expected true for 'yes'")
	}
}

func TestConfirm_YesUppercase(t *testing.T) {
	origStdin := stdin
	defer func() { stdin = origStdin }()

	stdin = strings.NewReader("Y\n")
	if !confirm("Continue?") {
		t.Error("expected true for 'Y'")
	}
}

func TestConfirm_No(t *testing.T) {
	origStdin := stdin
	defer func() { stdin = origStdin }()

	stdin = strings.NewReader("n\n")
	if confirm("Continue?") {
		t.Error("expected false for 'n'")
	}
}

func TestConfirm_Empty(t *testing.T) {
	origStdin := stdin
	defer func() { stdin = origStdin }()

	stdin = strings.NewReader("\n")
	if confirm("Continue?") {
		t.Error("expected false for empty input (default N)")
	}
}

func TestConfirm_Garbage(t *testing.T) {
	origStdin := stdin
	defer func() { stdin = origStdin }()

	stdin = strings.NewReader("maybe\n")
	if confirm("Continue?") {
		t.Error("expected false for non-y input")
	}
}

func TestConfirm_EOF(t *testing.T) {
	origStdin := stdin
	defer func() { stdin = origStdin }()

	stdin = strings.NewReader("")
	if confirm("Continue?") {
		t.Error("expected false for EOF")
	}
}
