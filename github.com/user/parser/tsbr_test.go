package parser

import (
	"bytes"
	"io"
	"testing"
	//	"fmt"
	"sync"
)

func TestOneTsbr(t *testing.T) {
	original_input := bytes.NewReader([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19})
	input := NewReader(original_input)
	b := make([]byte, 6)
	n, err := input.Read(b)
	if err != nil {
		t.Error(err)
	}
	if n != 6 {
		t.Errorf("Unexpected n: %d", n)
	}
	if bytes.Compare(b, []byte{0, 1, 2, 3, 4, 5}) != 0 {
		t.Errorf("unexpected b: %v", b)
	}

	b = make([]byte, 3)
	n, err = input.Read(b)
	if err != nil {
		t.Error(err)
	}
	if n != 3 {
		t.Errorf("Unexpected n: %d", n)
	}
	if bytes.Compare(b, []byte{6, 7, 8}) != 0 {
		t.Errorf("unexpected b: %v", b)
	}
}

func TestManyTsbr(t *testing.T) {
	original_input := bytes.NewReader([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19})
	reader1 := NewReader(original_input)
	b := make([]byte, 6)
	n, err := reader1.Read(b)
	if err != nil {
		t.Error(err)
	}
	if n != 6 {
		t.Errorf("Unexpected n: %d", n)
	}
	if bytes.Compare(b, []byte{0, 1, 2, 3, 4, 5}) != 0 {
		t.Errorf("unexpected b: %v", b)
	}
	if reader1.Offset() != 6 {
		t.Errorf("unexpected offset: %d", reader1.Offset())
	}

	reader2 := reader1.Clone()
	b = make([]byte, 4)
	n, err = reader2.Read(b)
	if err != nil {
		t.Error(err)
	}
	if n != 4 {
		t.Errorf("Unexpected n: %d", n)
	}
	if bytes.Compare(b, []byte{6, 7, 8, 9}) != 0 {
		t.Errorf("unexpected b: %v", b)
	}
	if reader2.Offset() != 10 {
		t.Errorf("unexpected offset: %d", reader2.Offset())
	}

	b = make([]byte, 3)
	n, err = reader1.Read(b)
	if err != nil {
		t.Error(err)
	}
	if n != 3 {
		t.Errorf("Unexpected n: %d", n)
	}
	if bytes.Compare(b, []byte{6, 7, 8}) != 0 {
		t.Errorf("unexpected b: %v", b)
	}
	if reader1.Offset() != 9 {
		t.Errorf("unexpected offset: %d", reader1.Offset())
	}
}

func TestManyConcurrentTsbr(t *testing.T) {

	var wg sync.WaitGroup
	testSequentialInput := func(reader io.Reader, from, until byte) {
		defer wg.Done()
		expected := from
		for {
			bs := make([]byte, 10)
			reader.Read(bs)
			//			fmt.Printf("from=%d; bs=%v\n",from,bs)
			for _, b := range bs {
				if b != expected {
					t.Fatal("non-sequential number")
				}
				if b == until {
					return
				}
				expected += 1
			}
		}
	}

	n := 100
	bs := make([]byte, n)
	for i, _ := range bs {
		bs[i] = byte(i)
	}
	original_input := bytes.NewReader(bs)
	reader := NewReader(original_input)

	oneByte := make([]byte, 1)
	rs := []io.Reader{}
	for i := 1; i < len(bs); i++ {
		rs = append(rs, reader)
		reader.Read(oneByte)
		//		fmt.Print(oneByte)
		reader = reader.Clone()
	}

	for i, r := range rs {
		wg.Add(1)
		go testSequentialInput(r, byte(i+1), byte(n-1))
	}
}

func TestImmediateClone(t *testing.T) {
	original_input := bytes.NewReader([]byte{0, 1, 2})
	reader1 := NewReader(original_input)
	reader2 := reader1.Clone()

	oneByte := make([]byte, 1)
	n, err := reader2.Read(oneByte)
	if err != nil {
		t.Error(err)
	}
	if n != 1 {
		t.Error("expected one byte read")
	}

	n, err = reader1.Read(oneByte)
	if err != nil {
		t.Error(err)
	}
	if n != 1 {
		t.Error("expected one byte read")
	}
}
