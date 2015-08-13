package parser

func S(str string) Parser {
	return StringRule{str}
}

func Seq(parsers ...Parser) Parser {
	if len(parsers) == 0 {
		panic("empty sequence rule not allowed")
	}
	return SequenceRule{parsers}
}

func OneOf(parsers ...Parser) Parser {
	if len(parsers) == 0 {
		panic("empty sequence rule not allowed")
	}
	return OneOfRule{parsers}
}

func AtLeastNumOf(parser Parser, num int) Parser {
	return AtLeastNumOfRule{parser, num}
}

func AsManyAsNumOf(parser Parser, num int) Parser {
	return AsManyAsNumOfRule{parser, num}
}

func ZeroOrMoreOf(parser Parser) Parser {
	return AsManyAsNumOfRule{parser, MaxInt}
}

func OneOrMoreOf(parser Parser) Parser {
	return SequenceRule{[]Parser{
		AtLeastNumOfRule{parser, 1},
		AsManyAsNumOfRule{parser, MaxInt},
	}}
}

func ZeroOrOneOf(parser Parser) Parser {
	return AsManyAsNumOfRule{parser, 1}
}