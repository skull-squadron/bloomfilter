package bloomfilter

import (
	"encoding"
	"encoding/gob"
	"io"
)

// compile-time conformance tests

func (f *Filter) conformsToBinaryMarshaler() encoding.BinaryMarshaler     { return f }
func (f *Filter) conformsToBinaryUnmarshaler() encoding.BinaryUnmarshaler { return f }
func (f *Filter) conformsToTextMarshaler() encoding.TextMarshaler         { return f }
func (f *Filter) conformsToTextUnmarshaler() encoding.TextUnmarshaler     { return f }
func (f *Filter) conformsToReaderFrom() io.ReaderFrom                     { return f }
func (f *Filter) conformsToWriterTo() io.WriterTo                         { return f }
func (f *Filter) conformsToGobDecoder() gob.GobDecoder                    { return f }
func (f *Filter) conformsToGobEncoder() gob.GobEncoder                    { return f }
