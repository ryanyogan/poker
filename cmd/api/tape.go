package main

import "io"

// tape can be thought of as a standard tape (cassette), it may be re-wound,
// we have to do this as the file size may be smaller on the next write, if we
// were to perform an IO `Seek` on a different file size, we would get a
// slice out of bounds error
type tape struct {
	file io.ReadWriteSeeker
}

// Notice that we're only implementing Write now, as it encapsulates
// the Seek part. This means our FileSystemStore can just have a
// reference to a Writer instead.
func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Seek(0, 0)
	return t.file.Write(p)
}
