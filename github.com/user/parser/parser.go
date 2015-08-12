package parser

import (
	"fmt"
	"io"

	"github.com/user/tsbr"
)

type ParseError struct {
	Offset int
	Rule   string
	Msg    string
}

func (p ParseError) Error() string {
	return fmt.Sprintf("error at offset %d in rule %s. %s", p.Offset, p.Rule, p.Msg)
}

// takes input *ThreadSafeBufferedReader and error if bad parse
type Parser interface {
	Parse(*tsbr.ThreadSafeBufferedReader) error
}

type AsManyAsNumOfRule struct {
	subRule Parser
	num	int
}

func (rule AsManyAsNumOfRule) Parse(input *tsbr.ThreadSafeBufferedReader) error {
	var err error
	subInput := input.Clone()
	for i:=1; ; i++ {
		err = rule.subRule.Parse(subInput)
		if err != nil {
			subInput.Done()
			return nil
		} else {
			input.Done()
			*input = *subInput
			if i<rule.num {
				subInput = input.Clone()
			} else {
				return nil
			}
		}
	}
}

type AtLeastNumOfRule struct {
	subRule Parser
	num	int
}

func (rule AtLeastNumOfRule) Parse(input *tsbr.ThreadSafeBufferedReader) error {
	var err error
	for i:=0; i<rule.num; i++ {
		err = rule.subRule.Parse(input)
		if err != nil {
			return err
		}
	}
	return nil
}

type OneOfRule struct {
	subRules []Parser
}

func (rule OneOfRule) Parse(input *tsbr.ThreadSafeBufferedReader) error {
	var highestErrOffset int = -1
	errSubRule := ""
	errSubMsg := ""
	for _, subRule := range rule.subRules {
		subInput := input.Clone()
		err := subRule.Parse(subInput)
		if err != nil {
			subInput.Done()
			switch err := err.(type) {
			default:
				return err
			case ParseError:
				if err.Offset > highestErrOffset {
					highestErrOffset = err.Offset
					errSubRule = err.Rule
					errSubMsg = err.Msg
				}
			}
		} else {
			input.Done()
			*input = *subInput
			return nil
		}
	}
	return ParseError{highestErrOffset, "OneOf>" + errSubRule, errSubMsg}
}

type SequenceRule struct {
	subRules []Parser
}

func (rule SequenceRule) Parse(input *tsbr.ThreadSafeBufferedReader) error {
	for _, subRule := range rule.subRules {
		err := subRule.Parse(input)
		if err != nil {
			switch err := err.(type) {
			default:
				return err
			case ParseError:
				return ParseError{
					err.Offset,
					fmt.Sprintf("Sequence>%s", err.Rule),
					err.Msg,
				}
			}
		}
	}
	return nil
}

type StringRule struct {
	str string
}

func (rule StringRule) Parse(input *tsbr.ThreadSafeBufferedReader) error {
	//TODO make this more efficient
	oneByte := make([]byte, 1)
	for _, chr := range []byte(rule.str) {
		_, err := input.Read(oneByte)
		if err != nil {
			if err == io.EOF {
				return ParseError{
					input.Offset(),
					"String>'" + rule.str + "'",
					"EOF",
				}
			} else {
				return err
			}
		}
		if chr != oneByte[0] {
			return ParseError{
				input.Offset() - 1,
				"String>'" + rule.str + "'",
				fmt.Sprintf("expected '%c' found '%c'", chr, oneByte[0]),
			}
		}
	}
	return nil
}
