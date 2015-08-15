package parser

import (
	"fmt"
	"io"
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
	Parse(*ThreadSafeBufferedReader) error
	GetSubRules() []Parser
	GetName() string
	Rename(string) Parser
}

type stringRule struct {
	str  string
	name string
}

func (rule stringRule) Parse(input *ThreadSafeBufferedReader) error {
	//TODO make this more efficient
	oneByte := make([]byte, 1)
	for _, chr := range []byte(rule.str) {
		_, err := input.Read(oneByte)
		if err != nil {
			if err == io.EOF {
				return ParseError{
					input.Offset(),
					fmt.Sprintf("'%s'", rule.str),
					"EOF",
				}
			} else {
				return err
			}
		}
		if chr != oneByte[0] {
			return ParseError{
				input.Offset() - 1,
				fmt.Sprintf("'%s'", rule.str),
				fmt.Sprintf("expected '%c' found '%c'", chr, oneByte[0]),
			}
		}
	}
	return nil
}
func (rule stringRule) GetSubRules() []Parser {
	return []Parser{}
}
func (rule stringRule) GetName() string {
	return rule.name
}
func (rule *stringRule) Rename(name string) Parser {
	rule.name = name
	return rule
}

type sequenceRule struct {
	subRules []Parser
	name     string
}

func (rule sequenceRule) Parse(input *ThreadSafeBufferedReader) error {
	for _, subRule := range rule.subRules {
		err := subRule.Parse(input)
		if err != nil {
			switch err := err.(type) {
			default:
				return err
			case ParseError:
				return ParseError{
					err.Offset,
					fmt.Sprintf("%s>%s", rule.name, err.Rule),
					err.Msg,
				}
			}
		}
	}
	return nil
}
func (rule sequenceRule) GetSubRules() []Parser {
	return rule.subRules
}
func (rule sequenceRule) GetName() string {
	return rule.name
}
func (rule *sequenceRule) Rename(name string) Parser {
	rule.name = name
	return rule
}

type oneOfRule struct {
	subRules []Parser
	name     string
}

func (rule oneOfRule) Parse(input *ThreadSafeBufferedReader) error {
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
	return ParseError{
		highestErrOffset,
		fmt.Sprintf("%s>%s", rule.name, errSubRule),
		errSubMsg,
	}
}
func (rule oneOfRule) GetSubRules() []Parser {
	return rule.subRules
}
func (rule oneOfRule) GetName() string {
	return rule.name
}
func (rule *oneOfRule) Rename(name string) Parser {
	rule.name = name
	return rule
}

type atLeastNumOfRule struct {
	subRule Parser
	num     int
	name    string
}

func (rule atLeastNumOfRule) Parse(input *ThreadSafeBufferedReader) error {
	var err error
	for i := 0; i < rule.num; i++ {
		err = rule.subRule.Parse(input)
		if err != nil {
			switch err := err.(type) {
			default:
				return err
			case ParseError:
				return ParseError{
					err.Offset,
					fmt.Sprintf("%s>%s", rule.name, err.Rule),
					err.Msg,
				}
			}
		}
	}
	return nil
}
func (rule atLeastNumOfRule) GetSubRules() []Parser {
	return []Parser{rule.subRule}
}
func (rule atLeastNumOfRule) GetName() string {
	return rule.name
}
func (rule *atLeastNumOfRule) Rename(name string) Parser {
	rule.name = name
	return rule
}

type asManyAsNumOfRule struct {
	subRule Parser
	num     int
	name    string
}

func (rule asManyAsNumOfRule) Parse(input *ThreadSafeBufferedReader) error {
	var err error
	subInput := input.Clone()
	for i := 1; ; i++ {
		err = rule.subRule.Parse(subInput)
		if err != nil {
			subInput.Done()
			return nil
		} else {
			input.Done()
			*input = *subInput
			if i < rule.num {
				subInput = input.Clone()
			} else {
				return nil
			}
		}
	}
}
func (rule asManyAsNumOfRule) GetSubRules() []Parser {
	return []Parser{rule.subRule}
}
func (rule asManyAsNumOfRule) GetName() string {
	return rule.name
}
func (rule *asManyAsNumOfRule) Rename(name string) Parser {
	rule.name = name
	return rule
}
