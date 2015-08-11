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
		//TODO this fails because of the tsbr swap (*input = *subInput) in 
		//OneOfRule. This says "make the value of the input the same as the
		//value of the subInput" but the only value they hold is a pointer to
		//their sbr - which was already the same thing anyway.
		//TO FIX: rather than the sbr having a map[*tsbr]int it should be
		//map[int]int where the key is just some unique identifier of that tsbr
		//probably in the subscribe part is where they get their id
		t.Error("unexpected error",err)
	}
	
//	input = tsbr.NewReader(strings.NewReader("goodbye"))
//	err = rule.Parse(input)
//	if err != nil {
//		t.Errorf("unexpected error")
//	}
	
//	input = tsbr.NewReader(strings.NewReader("hell"))
//	err = rule.Parse(input)
//	if err == nil {
//		t.Errorf("expected error, but was none")
//	}
//	if err != nil && 
//		err.Error() != "error at offset 4 in rule OneOf>String>'hello'. EOF" {
//		t.Errorf("unexpected error message: '%v'",err.Error())
//	}
	
}

