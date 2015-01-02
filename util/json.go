package util

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

	var acc []byte
	last := len(blobs) - 1
	var openingRune, closingRune byte
	a := 0
	setClosing := false
	idx := 0
	for i, b := range blobs {

		if len(b) > 0 && !setClosing { // is this an array or an object?
			setClosing = true
			openingRune, closingRune = b[0], closers[b[0]]
		}

		if len(b) < 3 { // yep empty but also the last one, so closing this thing
			if i == last && a > 0 {
				acc = append(acc, closingRune)
			}
			continue
		}

		idx = 0
		if a > 0 { // we need to join with a comma for everything beyond the first non-empty item
			acc = append(acc, comma)
			idx = 1
		}

		if i != last { // not the last one, strip ending bracket
			acc = append(acc, b[idx:len(b)-1]...)
		} else { // last one, strip the leading bracket
			acc = append(acc, b[idx:]...)
		}
		a++
	}
	// somehow it ended up being empty, so provide a default value
	if len(acc) == 0 {
		acc = []byte{openingRune, closingRune}
	}
	return acc
}
