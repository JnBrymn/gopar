package parser

import (
	"bytes"
	"strings"
	"testing"
)

func expectNoErr(t *testing.T, rule Parser, inText string) {
	expected := []byte("XYZ123")
	input := NewReader(strings.NewReader(inText + "XYZ123"))
	err := rule.Parse(input)
	found := make([]byte, 6)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(found); bytes.Compare(found, expected) != 0 {
		t.Error("input reader set to wrong index")
	}
}

func expectErr(t *testing.T, rule Parser, inText, errText string) {
	input := NewReader(strings.NewReader(inText))
	err := rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != errText {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}
}
