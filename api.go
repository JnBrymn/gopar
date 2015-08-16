package gopar

import (
	"strings"
)

func S(str string) Parser {
	return &stringRule{str, "_String"}
}

func Seq(parsers ...Parser) Parser {
	if len(parsers) == 0 {
		panic("empty sequenceRule not allowed")
	}
	return &sequenceRule{parsers, "_Sequence"}
}

func OneOf(parsers ...Parser) Parser {
	if len(parsers) == 0 {
		panic("empty oneOfRule not allowed")
	}
	return &oneOfRule{parsers, "_OneOf"}
}

func OneOfChars(chars string) Parser {
	if len(chars) == 0 {
		panic("empty oneOfRule not allowed")
	}
	charParsers := make([]Parser,len(chars))
	for i,c := range chars {
		charParsers[i] = &stringRule{
			string(c),
			"",
		}
	}
	ruleName := "{" + strings.Join(strings.Split(chars,""),"|") + "}"
	return &oneOfRule{charParsers, ruleName}
}

func AtLeastNumOf(parser Parser, num int) Parser {
	return &atLeastNumOfRule{parser, num, "_AtLeastNumOf"}
}

func AsManyAsNumOf(parser Parser, num int) Parser {
	return &asManyAsNumOfRule{parser, num, "_AsManyAsNumOf"}
}

func ZeroOrMoreOf(parser Parser) Parser {
	return &asManyAsNumOfRule{parser, MaxInt, "_ZeroOrMoreOf"}
}

func OneOrMoreOf(parser Parser) Parser {
	return &sequenceRule{[]Parser{
		&atLeastNumOfRule{parser, 1, ""},
		&asManyAsNumOfRule{parser, MaxInt, ""},
	}, "OneOrMoreOf"}
}

func ZeroOrOneOf(parser Parser) Parser {
	return &asManyAsNumOfRule{parser, 1, "_ZeroOrOneOf"}
}

func P(parserName string) Parser {
	return &placeholderRule{patchRuleName:parserName}
}


