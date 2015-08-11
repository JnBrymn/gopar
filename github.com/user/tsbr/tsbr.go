/*
Threadsafe buffered io.Reader that also tracks global offset for the initial
reader from which all readers are spawned.

Readers will read as long as it is safe to do so
*/
package tsbr

import (
	"fmt"
	"io"
	"sync"
)

const MaxInt = int(^uint(0) >> 1)

type Offsetter interface {
	Offset() int
}

type sharedBufferedReader struct {
	wrappedReader io.Reader
	buffer []byte
	bytesRead int
	globalOffsets map[*ThreadSafeBufferedReader]int
	mutex sync.Mutex
}

func newSharedBufferedReader(reader io.Reader) *sharedBufferedReader {
	return &sharedBufferedReader{
		wrappedReader:reader,
		buffer:[]byte{},
		bytesRead: 0,
		globalOffsets:map[*ThreadSafeBufferedReader]int{},
	}
}

func (sbr *sharedBufferedReader) read(b []byte, tsbr *ThreadSafeBufferedReader) (int,error) {
	sbr.mutex.Lock()
	defer sbr.mutex.Unlock()
	
	// fulfill request
	var err error
	n1 := copy(b, sbr.buffer[sbr.globalOffsets[tsbr] - sbr.bytesRead:])
	n2 := 0
	if n1 < len(b) {
		//still more to read
		bEnd := b[n1:]
		n2, err = sbr.wrappedReader.Read(bEnd)
		sbr.buffer = append(sbr.buffer,bEnd[:n2]...)
		if err != nil {
			return n1+n2, err
		}
	}
	n := n1+n2
	sbr.globalOffsets[tsbr] += n 
	
	// shrink buffer
	lowestGlobalOffset := MaxInt
	for _, globalOffset := range sbr.globalOffsets {
		if globalOffset < lowestGlobalOffset {
			lowestGlobalOffset = globalOffset
		}
	}
	
	sbr.buffer = sbr.buffer[lowestGlobalOffset-sbr.bytesRead:]
	sbr.bytesRead = lowestGlobalOffset
	
	return n, nil
}

/*
This is used in two cases:
1) When a new TSBR is created it subscribes itself with parent=nil (indicating
it has no parent).
2) When a TSBR gets Cloned, it subscribes it's child using this. Both parent 
and child are not nil.
*/
func (sbr *sharedBufferedReader) subscribe(child, parent *ThreadSafeBufferedReader) {
	// note this depends upon the fact that alway sbr.globalOffsets[nil]==0
	parentOffset := sbr.globalOffsets[parent]
	sbr.globalOffsets[child] = parentOffset
}

func (sbr *sharedBufferedReader) offset(tsbr *ThreadSafeBufferedReader) int {
	return sbr.globalOffsets[tsbr]
}

func (sbr *sharedBufferedReader) done(tsbr *ThreadSafeBufferedReader) {
	delete(sbr.globalOffsets,tsbr)
}

type ThreadSafeBufferedReader struct {
	sbr *sharedBufferedReader
}

func NewReader(reader io.Reader) *ThreadSafeBufferedReader {
	tsbr := &ThreadSafeBufferedReader{newSharedBufferedReader(reader)}
	tsbr.sbr.subscribe(tsbr,nil)
	return tsbr
}

func (tsbr *ThreadSafeBufferedReader) Clone() *ThreadSafeBufferedReader {
	childTsbr := &ThreadSafeBufferedReader{tsbr.sbr}
	tsbr.sbr.subscribe(childTsbr,tsbr)
	return childTsbr
}

func (tsbr *ThreadSafeBufferedReader) Read(b []byte) (int, error) {
	return tsbr.sbr.read(b,tsbr)	
}

func (tsbr *ThreadSafeBufferedReader) Offset() int {
	return tsbr.sbr.offset(tsbr)
}

func (tsbr *ThreadSafeBufferedReader) Done() {
	tsbr.sbr.done(tsbr)
}

func main() {
	fmt.Println("sdf")
}