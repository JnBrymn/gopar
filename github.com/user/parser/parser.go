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

type FromToOfRule struct {
	rule Parser
}

func (fromToRule *FromToOfRule) Parse(input *tsbr.ThreadSafeBufferedReader) error {
	return nil
}


type OneOfRule struct {
	subRules []Parser
}

func (oneOfRule OneOfRule) Parse(input *tsbr.ThreadSafeBufferedReader) error {
	var highestErrOffset int = -1
	errSubRule := ""
	errSubMsg := ""
	for _, rule := range oneOfRule.subRules {
		subInput := input.Clone()
		err := rule.Parse(subInput)
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
	SubRules []Parser
}

func (seqRule SequenceRule) Parse(input *tsbr.ThreadSafeBufferedReader) error {
	for _, rule := range seqRule.SubRules {
		err := rule.Parse(input)
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

func (strRule StringRule) Parse(input *tsbr.ThreadSafeBufferedReader) error {
	//TODO make this more efficient
	oneByte := make([]byte, 1)
	for _, chr := range []byte(strRule.str) {
		_, err := input.Read(oneByte)
		if err != nil {
			if err == io.EOF {
				return ParseError{
					input.Offset(),
					"String>'" + strRule.str + "'",
					"EOF",
				}
			} else {
				return err
			}
		}
		if chr != oneByte[0] {
			return ParseError{
				input.Offset() - 1,
				"String>'" + strRule.str + "'",
				fmt.Sprintf("expected '%c' found '%c'", chr, oneByte[0]),
			}
		}
	}
	return nil
}
