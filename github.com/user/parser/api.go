package parser

import (
	"strings"
)

func S(str string) Parser {
	return &stringRule{str, "String"}
}

func Seq(parsers ...Parser) Parser {
	if len(parsers) == 0 {
		panic("empty sequenceRule not allowed")
	}
	return &sequenceRule{parsers, "Sequence"}
}

func OneOf(parsers ...Parser) Parser {
	if len(parsers) == 0 {
		panic("empty oneOfRule not allowed")
	}
	return &oneOfRule{parsers, "OneOf"}
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
	return &oneOfRule{charParsers, strings.Join(strings.Split(chars,""),"|")}
}

func AtLeastNumOf(parser Parser, num int) Parser {
	return &atLeastNumOfRule{parser, num, "AtLeastNumOf"}
}

func AsManyAsNumOf(parser Parser, num int) Parser {
	return &asManyAsNumOfRule{parser, num, "AsManyAsNumOf"}
}

func ZeroOrMoreOf(parser Parser) Parser {
	return &asManyAsNumOfRule{parser, MaxInt, "ZeroOrMoreOf"}
}

func OneOrMoreOf(parser Parser) Parser {
	return &sequenceRule{[]Parser{
		&atLeastNumOfRule{parser, 1, ""},
		&asManyAsNumOfRule{parser, MaxInt, ""},
	}, "OneOrMoreOf"}
}

func ZeroOrOneOf(parser Parser) Parser {
	return &asManyAsNumOfRule{parser, 1, "ZeroOrOneOf"}
}


