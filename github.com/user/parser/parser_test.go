package parser

import (
	"strings"
	"testing"

	"github.com/user/tsbr"
)

func TestStringRule(t *testing.T) {
	rule := StringRule{"hello"}
	oneByte := make([]byte, 1)

	input := tsbr.NewReader(strings.NewReader("helloX"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("hell"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != "error at offset 4 in rule String>'hello'. EOF" {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}

	input = tsbr.NewReader(strings.NewReader("helso"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != "error at offset 3 in rule String>'hello'. expected 'l' found 's'" {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}
}

func TestStringRuleWithUnicode(t *testing.T) {
	rule := StringRule{"abcあいうえおdef"}
		oneByte := make([]byte, 1)

	input := tsbr.NewReader(strings.NewReader("abcあいうえおdefX"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("abcあいえおdef"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		// this is techinically wrong because it doesn't "see" the multi-byte
		// string char
		err.Error() != "error at offset 11 in rule String>'abcあいうえおdef'. expected '' found ''" {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}
}

func TestSequenceRule(t *testing.T) {
	rule := SequenceRule{[]Parser{
		StringRule{"hello"},
		StringRule{"goodbye"},
	}}
	oneByte := make([]byte, 1)

	input := tsbr.NewReader(strings.NewReader("hellogoodbyeX"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("hellgoodbye"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != "error at offset 4 in rule Sequence>String>'hello'. expected 'o' found 'g'" {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}

	input = tsbr.NewReader(strings.NewReader("hellogodbye"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != "error at offset 7 in rule Sequence>String>'goodbye'. expected 'o' found 'd'" {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}
}

func TestOneOfRule(t *testing.T) {
	rule := OneOfRule{[]Parser{
		StringRule{"hello"},
		StringRule{"goodbye"},
	}}
	oneByte := make([]byte, 1)

	input := tsbr.NewReader(strings.NewReader("helloX"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("goodbyeX"))
	err = rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("hell"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != "error at offset 4 in rule OneOf>String>'hello'. EOF" {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}

}

func TestOneOfThenSequenceRule(t *testing.T) {
	// this is "tricky" because this will match input="abc" but it's not
	// the first or the thing it tries. This tests the ability to "back up"
	// and try again
	rule := SequenceRule{[]Parser{
		OneOfRule{[]Parser{
			StringRule{"abx"},
			StringRule{"a"},
		}},
		StringRule{"bc"},
	}}
	oneByte := make([]byte, 1)

	input := tsbr.NewReader(strings.NewReader("abcX"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("abxbcX"))
	err = rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("aby"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != "error at offset 2 in rule Sequence>String>'bc'. expected 'c' found 'y'" {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}
}

func TestAtLeastNumOfRule(t *testing.T) {
	rule := AtLeastNumOfRule{
		StringRule{"abc"},
		3,
	}
	oneByte := make([]byte, 1)

	input := tsbr.NewReader(strings.NewReader("abcabcabcX"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("abcabcX"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != "error at offset 6 in rule String>'abc'. expected 'a' found 'X'" {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("abcaXcabc"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != "error at offset 4 in rule String>'abc'. expected 'b' found 'X'" {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}
}

func TestAsManyAsNumOfRule(t *testing.T) {
	rule := AsManyAsNumOfRule{
		StringRule{"abc"},
		3,
	}
	oneByte := make([]byte, 1)

	input := tsbr.NewReader(strings.NewReader("X"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("abcX"))
	err = rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("abcabcX"))
	err = rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}

	input = tsbr.NewReader(strings.NewReader("abcabcabcX"))
	err = rule.Parse(input)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(oneByte); oneByte[0] != byte('X') {
		t.Error("input reader set to wrong index")
	}
}
