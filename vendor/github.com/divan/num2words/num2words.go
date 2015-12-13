/*
num2words implements numbers to words converter.
*/
package num2words

import "math"

// how many digit's groups to process
const groups_number int = 4

var _smallNumbers = []string{
	"zero", "one", "two", "three", "four",
	"five", "six", "seven", "eight", "nine",
	"ten", "eleven", "twelve", "thirteen", "fourteen",
	"fifteen", "sixteen", "seventeen", "eighteen", "nineteen",
}
var _tens = []string{
	"", "", "twenty", "thirty", "forty", "fifty",
	"sixty", "seventy", "eighty", "ninety",
}
var _scaleNumbers = []string{
	"", "thousand", "million", "billion",
}

type digitGroup int

// Convert converts number into the words representation.
func Convert(number int) string {
	// Zero rule
	if number == 0 {
		return _smallNumbers[0]
	}

	// Divide into three-digits group
	var groups [groups_number]digitGroup
	positive := math.Abs(float64(number))

	// Form three-digit groups
	for i := 0; i < groups_number; i++ {
		groups[i] = digitGroup(math.Mod(positive, 1000))
		positive /= 1000
	}

	var textGroup [groups_number]string
	for i := 0; i < groups_number; i++ {
		textGroup[i] = digitGroup2Text(groups[i])
	}
	combined := textGroup[0]

	for i := 1; i < groups_number; i++ {
		if groups[i] != 0 {
			prefix := textGroup[i] + " " + _scaleNumbers[i]

			if len(combined) != 0 {
				prefix += " "
			}

			combined = prefix + combined
		}
	}

	if number < 0 {
		combined = "minus " + combined
	}

	return combined
}

func intMod(x, y int) int {
	return int(math.Mod(float64(x), float64(y)))
}

func digitGroup2Text(group digitGroup) (ret string) {
	hundreds := group / 100
	tensUnits := intMod(int(group), 100)

	if hundreds != 0 {
		ret += _smallNumbers[hundreds] + " hundred"

		if tensUnits != 0 {
			ret += " "
		}
	}

	tens := tensUnits / 10
	units := intMod(tensUnits, 10)

	if tens >= 2 {
		ret += _tens[tens]

		if units != 0 {
			ret += " " + _smallNumbers[units]
		}
	} else if tensUnits != 0 {
		ret += _smallNumbers[tensUnits]
	}

	return
}
