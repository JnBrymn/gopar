package parser

import (
	"testing"
)

func TestS(t *testing.T) {
	rule := S("hello")
	expectNoErr(t, rule, "hello")
	expectErr(t, rule, "hell", "error at offset 4 in rule String>'hello'. EOF")
	expectErr(t, rule, "helso", "error at offset 3 in rule String>'hello'. expected 'l' found 's'")
}

func TestSeq(t *testing.T) {
	rule := Seq(
		S("hello"),
		S("goodbye"),
	)
	expectNoErr(t, rule, "hellogoodbye")
	expectErr(t, rule, "hellgoodbye", "error at offset 4 in rule Sequence>String>'hello'. expected 'o' found 'g'")
	expectErr(t, rule, "hellogodbye", "error at offset 7 in rule Sequence>String>'goodbye'. expected 'o' found 'd'")
}

func TestOneOf(t *testing.T) {
	rule := OneOf(
		S("hello"),
		S("goodbye"),
	)
	expectNoErr(t, rule, "hello")
	expectNoErr(t, rule, "goodbye")
	expectErr(t, rule, "hell", "error at offset 4 in rule OneOf>String>'hello'. EOF")
}

func TestAtLeastNumOf(t *testing.T) {
	rule := AtLeastNumOf(
		S("abc"),
		3,
	)
	expectNoErr(t, rule, "abcabcabc")
	expectErr(t, rule, "abcabcX", "error at offset 6 in rule AtLeastNumOf>String>'abc'. expected 'a' found 'X'")
	expectErr(t, rule, "abcaXcabc", "error at offset 4 in rule AtLeastNumOf>String>'abc'. expected 'b' found 'X'")
}

func TestAsManyAsNumOf(t *testing.T) {
	rule := Seq(
		AsManyAsNumOf(
			S("abc"),
			3,
		),
		S("!"),
	)
	expectNoErr(t, rule, "!")
	expectNoErr(t, rule, "abc!")
	expectNoErr(t, rule, "abcabc!")
	expectNoErr(t, rule, "abcabcabc!")
	expectErr(t, rule, "abcabcabcabc!", "error at offset 9 in rule Sequence>String>'!'. expected '!' found 'a'")
}

func TestZeroOrMoreOf(t *testing.T) {
	rule := ZeroOrMoreOf(
		S("abc"),
	)
	expectNoErr(t, rule, "")
	expectNoErr(t, rule, "abc")
	expectNoErr(t, rule, "abcabc")
	expectNoErr(t, rule, "abcabcabc")
}

func TestOneOrMoreOf(t *testing.T) {
	rule := OneOrMoreOf(
		S("abc"),
	)
	expectNoErr(t, rule, "abc")
	expectNoErr(t, rule, "abcabc")
	expectNoErr(t, rule, "abcabcabc")
	expectErr(t, rule, "","error at offset 0 in rule OneOrMoreOf>>String>'abc'. EOF")
}

func TestZeroOrOneOf(t *testing.T) {
	rule := ZeroOrOneOf(
		S("abc"),
	)
	expectNoErr(t, rule, "")
	expectNoErr(t, rule, "abc")
}


func TestJson(t *testing.T) {
	rule := ZeroOrOneOf(
		S("abc"),
	)
	expectNoErr(t, rule, "")
	expectNoErr(t, rule, "abc")
}