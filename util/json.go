package util

import "bytes"

const comma = byte(',')

var closers = map[byte]byte{
	'{': '}',
	'[': ']',
}

// ConcatJSON concatenates multiple json objects efficiently
func ConcatJSON(blobs ...[]byte) []byte {
	if len(blobs) == 0 {
		return nil
	}
	if len(blobs) == 1 {
		return blobs[0]
	}

	last := len(blobs) - 1
	var opening, closing byte
	a := 0
	setClosing := false
	idx := 0
	buf := bytes.NewBuffer(nil)

	for i, b := range blobs {
		if len(b) > 0 && !setClosing { // is this an array or an object?
			setClosing = true
			opening, closing = b[0], closers[b[0]]
		}

		if len(b) < 3 { // yep empty but also the last one, so closing this thing
			if i == last && a > 0 {
				buf.WriteByte(closing)
			}
			continue
		}

		idx = 0
		if a > 0 { // we need to join with a comma for everything beyond the first non-empty item
			buf.WriteByte(comma)
			idx = 1 // this is not the first or the last so we want to drop the leading bracket
		}

		if i != last { // not the last one, strip brackets
			buf.Write(b[idx : len(b)-1])
		} else { // last one, strip only the leading bracket
			buf.Write(b[idx:])
		}
		a++
	}
	// somehow it ended up being empty, so provide a default value
	if buf.Len() == 0 {
		buf.WriteByte(opening)
		buf.WriteByte(closing)
	}
	return buf.Bytes()
}
