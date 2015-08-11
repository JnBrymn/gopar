/*
Threadsafe buffered io.Reader that also tracks global offset for the initial
reader from which all readers are spawned.

Readers will read as long as it is safe to do so
*/
package tsbr

import (
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
	globalOffsets map[int]int
	mutex sync.Mutex
	nextSubscriberId int
}

func newSharedBufferedReader(reader io.Reader) *sharedBufferedReader {
	return &sharedBufferedReader{
		wrappedReader:reader,
		buffer:[]byte{},
		bytesRead: 0,
		globalOffsets:map[int]int{},
		nextSubscriberId: 1,
	}
}

func (sbr *sharedBufferedReader) read(b []byte, tsbrId int) (int,error) {
	sbr.mutex.Lock()
	defer sbr.mutex.Unlock()
	
	// fulfill request
	var err error
	n1 := copy(b, sbr.buffer[sbr.globalOffsets[tsbrId] - sbr.bytesRead:])
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
	sbr.globalOffsets[tsbrId] += n 
	
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
1) When a new TSBR is created it subscribes itself with parentId=0 (indicating
it has no parent).
2) When a TSBR gets Cloned, it subscribes a child using it's id.

In both cases the ID for the new TSBR is returned.
*/
func (sbr *sharedBufferedReader) subscribe(parentId int) int {
	// note this depends upon the fact that alway sbr.globalOffsets[nil]==0
	parentOffset := sbr.globalOffsets[parentId]
	childId := sbr.nextSubscriberId
	sbr.nextSubscriberId += 1
	sbr.globalOffsets[childId] = parentOffset
	return childId
}

func (sbr *sharedBufferedReader) offset(tsbrId int) int {
	return sbr.globalOffsets[tsbrId]
}

func (sbr *sharedBufferedReader) done(tsbrId int) {
	delete(sbr.globalOffsets,tsbrId)
}

type ThreadSafeBufferedReader struct {
	sbr *sharedBufferedReader
	id int
}

func NewReader(reader io.Reader) *ThreadSafeBufferedReader {
	tsbr := &ThreadSafeBufferedReader{newSharedBufferedReader(reader),0}
	tsbr.id = tsbr.sbr.subscribe(0)
	return tsbr
}

func (tsbr *ThreadSafeBufferedReader) Clone() *ThreadSafeBufferedReader {
	childTsbr := &ThreadSafeBufferedReader{tsbr.sbr,0}
	childTsbr.id = tsbr.sbr.subscribe(tsbr.id)
	return childTsbr
}

func (tsbr *ThreadSafeBufferedReader) Read(b []byte) (int, error) {
	return tsbr.sbr.read(b,tsbr.id)	
}

func (tsbr *ThreadSafeBufferedReader) Offset() int {
	return tsbr.sbr.offset(tsbr.id)
}

func (tsbr *ThreadSafeBufferedReader) Done() {
	tsbr.sbr.done(tsbr.id)
}