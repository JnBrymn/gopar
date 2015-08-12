package parser

import (
	"strings"
	"testing"
	"bytes"

	"github.com/user/tsbr"
)

func expectNoErr(t *testing.T, rule Parser, inText string) {
	expected := []byte("XYZ123")
	input := tsbr.NewReader(strings.NewReader(inText + "XYZ123"))
	err := rule.Parse(input)
	found := make([]byte,6)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if input.Read(found); bytes.Compare(found,expected) != 0 {
		t.Error("input reader set to wrong index")
	}
}

func expectErr(t *testing.T, rule Parser, inText, errText string) {
	input := tsbr.NewReader(strings.NewReader(inText))
	err := rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil &&
		err.Error() != errText {
		t.Errorf("unexpected error: message: '%v'", err.Error())
	}
}

func TestStringRule(t *testing.T) {
	rule := StringRule{"hello"}
	expectNoErr(t,rule,"hello")
	expectErr(t,rule,"hell","error at offset 4 in rule String>'hello'. EOF")
	expectErr(t,rule,"helso","error at offset 3 in rule String>'hello'. expected 'l' found 's'")
}

func TestStringRuleWithUnicode(t *testing.T) {
	rule := StringRule{"abcあいうえおdef"}
	expectNoErr(t,rule,"abcあいうえおdef")
	expectErr(t,rule,"abcあいえおdef","error at offset 11 in rule String>'abcあいうえおdef'. expected '' found ''")
}

func TestSequenceRule(t *testing.T) {
	rule := SequenceRule{[]Parser{
		StringRule{"hello"},
		StringRule{"goodbye"},
	}}
	expectNoErr(t,rule,"hellogoodbye")
	expectErr(t,rule,"hellgoodbye","error at offset 4 in rule Sequence>String>'hello'. expected 'o' found 'g'")
	expectErr(t,rule,"hellogodbye","error at offset 7 in rule Sequence>String>'goodbye'. expected 'o' found 'd'")
}

func TestOneOfRule(t *testing.T) {
	rule := OneOfRule{[]Parser{
		StringRule{"hello"},
		StringRule{"goodbye"},
	}}
	expectNoErr(t,rule,"hello")
	expectNoErr(t,rule,"goodbye")
	expectErr(t,rule,"hell","error at offset 4 in rule OneOf>String>'hello'. EOF")
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
	expectNoErr(t,rule,"abc")
	expectNoErr(t,rule,"abxbc")
	expectErr(t,rule,"aby","error at offset 2 in rule Sequence>String>'bc'. expected 'c' found 'y'")
}

func TestAtLeastNumOfRule(t *testing.T) {
	rule := AtLeastNumOfRule{
		StringRule{"abc"},
		3,
	}
	expectNoErr(t,rule,"abcabcabc")
	expectErr(t,rule,"abcabcX","error at offset 6 in rule String>'abc'. expected 'a' found 'X'")
	expectErr(t,rule,"abcaXcabc","error at offset 4 in rule String>'abc'. expected 'b' found 'X'")
}

func TestAsManyAsNumOfRule(t *testing.T) {
	rule := SequenceRule{[]Parser{
		AsManyAsNumOfRule{
			StringRule{"abc"},
			3,
		},
		StringRule{"!"},
	}}
	expectNoErr(t,rule,"!")
	expectNoErr(t,rule,"abc!")
	expectNoErr(t,rule,"abcabc!")
	expectNoErr(t,rule,"abcabcabc!")
	expectErr(t,rule,"abcabcabcabc!","error at offset 9 in rule Sequence>String>'!'. expected '!' found 'a'")
}