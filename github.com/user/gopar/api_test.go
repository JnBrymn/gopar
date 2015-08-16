package gopar

import (
	"testing"
)

func TestS(t *testing.T) {
	rule := S("hello")
	expectNoErr(t, rule, "hello")
	expectErr(t, rule, "hell", "error at offset 4 in rule 'hello'. EOF")
	expectErr(t, rule, "helso", "error at offset 3 in rule 'hello'. expected 'l' found 's'")
}

func TestSeq(t *testing.T) {
	rule := Seq(
		S("hello"),
		S("goodbye"),
	)
	expectNoErr(t, rule, "hellogoodbye")
	expectErr(t, rule, "hellgoodbye", "error at offset 4 in rule _Sequence>'hello'. expected 'o' found 'g'")
	expectErr(t, rule, "hellogodbye", "error at offset 7 in rule _Sequence>'goodbye'. expected 'o' found 'd'")
}

func TestOneOf(t *testing.T) {
	rule := OneOf(
		S("hello"),
		S("goodbye"),
	)
	expectNoErr(t, rule, "hello")
	expectNoErr(t, rule, "goodbye")
	expectErr(t, rule, "hell", "error at offset 4 in rule _OneOf>'hello'. EOF")
}

func TestAtLeastNumOf(t *testing.T) {
	rule := AtLeastNumOf(
		S("abc"),
		3,
	)
	expectNoErr(t, rule, "abcabcabc")
	expectErr(t, rule, "abcabcX", "error at offset 6 in rule _AtLeastNumOf>'abc'. expected 'a' found 'X'")
	expectErr(t, rule, "abcaXcabc", "error at offset 4 in rule _AtLeastNumOf>'abc'. expected 'b' found 'X'")
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
	expectErr(t, rule, "abcabcabcabc!", "error at offset 9 in rule _Sequence>'!'. expected '!' found 'a'")
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
	expectErr(t, rule, "", "error at offset 0 in rule OneOrMoreOf>>'abc'. EOF")
}

func TestZeroOrOneOf(t *testing.T) {
	rule := ZeroOrOneOf(
		S("abc"),
	)
	expectNoErr(t, rule, "")
	expectNoErr(t, rule, "abc")
}

func TestRenaming(t *testing.T) {
	num := OneOrMoreOf(
		OneOfChars("0123456789"),
	).Rename("Number")
	prod := Seq(num, S("*"), num).Rename("Product")
	sum := Seq(prod, S("+"), prod).Rename("Sum")
	expectNoErr(t, sum, "33*44+1*3")
	expectErr(t, sum, "3*4+*35", "error at offset 4 in rule Sum>Product>Number>>{0|1|2|3|4|5|6|7|8|9}>'0'. expected '0' found '*'")
}

func TestPlaceholders(t *testing.T) {
	word := OneOrMoreOf(
		OneOfChars("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
	).Rename("Word")

	quote := Seq(
		S("'"),
		P("Sentence"),
		S("'"),
	).Rename("Quote")

	sentence := Seq(
		OneOf(word, quote),
		ZeroOrMoreOf(Seq(
			S(" "),
			OneOf(word, quote).Rename("Test"),
		)),
		OneOfChars("!?"),
	).Rename("Sentence")

	Patch(sentence)

	expectNoErr(t, sentence, "asd qwer sdfg erty!")
	expectNoErr(t, sentence, "asd qwer 'sdfg erty werq!' he said!")
	expectErr(t, sentence, "asd qwer 'sdfg erty werq he said!", "error at offset 8 in rule Sentence>{!|?}>'!'. expected '!' found ' '")
}
