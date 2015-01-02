package util

const comma = byte(',')

// ConcatJSON concatenates 2 json objects efficiently
func ConcatJSON(blobs ...[]byte) []byte {
	if len(blobs) == 0 {
		return nil
	}
	if len(blobs) == 1 {
		return blobs[0]
	}
	var acc []byte
	last := len(blobs) - 1
	var closingRune = '}'
	a := 0
	for i, b := range blobs {
		if len(b) < 3 {
			if i == last && a > 0 {
				acc = append(acc, byte(closingRune))
			}
			continue
		}
		idx := 0
		if a > 0 {
			acc = append(acc, comma)
			idx = 1
		} else {
			if b[0] == '[' {
				closingRune = ']'
			}
		}
		if i != last {
			acc = append(acc, b[idx:len(b)-1]...)
		} else {
			acc = append(acc, b[idx:]...)
		}
		a++
	}
	if len(acc) == 0 {
		acc = []byte("{}")
	}
	return acc
}
