package household

import "testing"

func TestGenerateJoinCode(t *testing.T) {
	code, err := GenerateJoinCode()
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != codeLength {
		t.Fatalf("expected length %d, got %d", codeLength, len(code))
	}
	for _, c := range code {
		if !containsRune(codeChars, c) {
			t.Fatalf("invalid character %q in code %q", c, code)
		}
	}
}

func containsRune(chars string, r rune) bool {
	for _, c := range chars {
		if c == r {
			return true
		}
	}
	return false
}
