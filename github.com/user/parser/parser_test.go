package parser

import (
	"testing"
	"github.com/user/tsbr"
	"strings"
)

func TestStringRule(t *testing.T) {
	rule := StringRule{"hello"}
	
	input := tsbr.NewReader(strings.NewReader("hello"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error",err)
	}
	
	input = tsbr.NewReader(strings.NewReader("hell"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil && 
		err.Error() != "error at offset 4 in rule String>'hello'. EOF" {
		t.Errorf("unexpected error message: '%v'",err.Error())
	}
	
	input = tsbr.NewReader(strings.NewReader("helso"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil && 
		err.Error() != "error at offset 3 in rule String>'hello'. expected 'l' found 's'" {
		t.Errorf("unexpected error message: '%v'",err.Error())
	}
}


func TestSequenceRule(t *testing.T) {
	rule := SequenceRule{[]Parser{
		StringRule{"hello"},
		StringRule{"goodbye"},
	}}
	
	input := tsbr.NewReader(strings.NewReader("hellogoodbye"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error",err)
	}
	
	input = tsbr.NewReader(strings.NewReader("hellgoodbye"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil && 
		err.Error() != "error at offset 4 in rule Sequence>String>'hello'. expected 'o' found 'g'" {
		t.Errorf("unexpected error message: '%v'",err.Error())
	}
	
	input = tsbr.NewReader(strings.NewReader("hellogodbye"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil && 
		err.Error() != "error at offset 7 in rule Sequence>String>'goodbye'. expected 'o' found 'd'" {
		t.Errorf("unexpected error message: '%v'",err.Error())
	}
}

func TestOneOfRule(t *testing.T) {
	rule := OneOfRule{[]Parser{
		StringRule{"hello"},
		StringRule{"goodbye"},
	}}
	
	input := tsbr.NewReader(strings.NewReader("hello"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error",err)
	}
	
	input = tsbr.NewReader(strings.NewReader("goodbye"))
	err = rule.Parse(input)
	if err != nil {
		t.Error("unexpected error",err)
	}
	
	input = tsbr.NewReader(strings.NewReader("hell"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil && 
		err.Error() != "error at offset 4 in rule OneOf>String>'hello'. EOF" {
		t.Errorf("unexpected error message: '%v'",err.Error())
	}
	
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
	
	input := tsbr.NewReader(strings.NewReader("abc"))
	err := rule.Parse(input)
	if err != nil {
		t.Error("unexpected error",err)
	}
	
	input = tsbr.NewReader(strings.NewReader("abxbc"))
	err = rule.Parse(input)
	if err != nil {
		t.Error("unexpected error",err)
	}

	input = tsbr.NewReader(strings.NewReader("aby"))
	err = rule.Parse(input)
	if err == nil {
		t.Errorf("expected error, but was none")
	}
	if err != nil && 
		err.Error() != "error at offset 2 in rule Sequence>String>'bc'. expected 'c' found 'y'" {
		t.Errorf("unexpected error message: '%v'",err.Error())
	}
	
}

