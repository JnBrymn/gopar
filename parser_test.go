package gopar

import "testing"

func TestStringRule(t *testing.T) {
	rule := &stringRule{"hello", "String"}
	expectNoErr(t, rule, "hello")
	expectErr(t, rule, "hell", "error at offset 4 in rule 'hello'. EOF")
	expectErr(t, rule, "helso", "error at offset 3 in rule 'hello'. expected 'l' found 's'")
}

func TestStringRuleWithUnicode(t *testing.T) {
	rule := &stringRule{"abcあいうえおdef", "String"}
	expectNoErr(t, rule, "abcあいうえおdef")
	expectErr(t, rule, "abcあいえおdef", "error at offset 11 in rule 'abcあいうえおdef'. expected '' found ''")
}

func TestSequenceRule(t *testing.T) {
	rule := &sequenceRule{[]Parser{
		&stringRule{"hello", "String"},
		&stringRule{"goodbye", "String"},
	}, "_Sequence"}
	expectNoErr(t, rule, "hellogoodbye")
	expectErr(t, rule, "hellgoodbye", "error at offset 4 in rule _Sequence>'hello'. expected 'o' found 'g'")
	expectErr(t, rule, "hellogodbye", "error at offset 7 in rule _Sequence>'goodbye'. expected 'o' found 'd'")
}

func TestOneOfRule(t *testing.T) {
	rule := &oneOfRule{[]Parser{
		&stringRule{"hello", "String"},
		&stringRule{"goodbye", "String"},
	}, "OneOf"}
	expectNoErr(t, rule, "hello")
	expectNoErr(t, rule, "goodbye")
	expectErr(t, rule, "hell", "error at offset 4 in rule OneOf>'hello'. EOF")
}

func TestOneOfThenSequenceRule(t *testing.T) {
	// this is "tricky" because this will match input="abc" but it's not
	// the first or the thing it tries. This tests the ability to "back up"
	// and try again
	rule := &sequenceRule{[]Parser{
		&oneOfRule{[]Parser{
			&stringRule{"abx","String"},
			&stringRule{"a","String"},
		},"OneOf"},
		&stringRule{"bc","String"},
	},"_Sequence"}
	expectNoErr(t, rule, "abc")
	expectNoErr(t, rule, "abxbc")
	expectErr(t, rule, "aby", "error at offset 2 in rule _Sequence>'bc'. expected 'c' found 'y'")
}

func TestAtLeastNumOfRule(t *testing.T) {
	rule := &atLeastNumOfRule{
		&stringRule{"abc","String"},
		3,
		"AtLeastNumOf",
	}
	expectNoErr(t, rule, "abcabcabc")
	expectErr(t, rule, "abcabcX", "error at offset 6 in rule AtLeastNumOf>'abc'. expected 'a' found 'X'")
	expectErr(t, rule, "abcaXcabc", "error at offset 4 in rule AtLeastNumOf>'abc'. expected 'b' found 'X'")
}

func TestAsManyAsNumOfRule(t *testing.T) {
	rule := &sequenceRule{[]Parser{
		&asManyAsNumOfRule{
			&stringRule{"abc","String"},
			3,
			"",
		},
		&stringRule{"!","String"},
	},"_Sequence"}
	expectNoErr(t, rule, "!")
	expectNoErr(t, rule, "abc!")
	expectNoErr(t, rule, "abcabc!")
	expectNoErr(t, rule, "abcabcabc!")
	expectErr(t, rule, "abcabcabcabc!", "error at offset 9 in rule _Sequence>'!'. expected '!' found 'a'")
}

