package luhn

import (
	"regexp"
	"strconv"
	"strings"
)

// IntToString converts integer to string
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// Int32ToFloat32 converts integer32 to float32
func Int32ToFloat32(i int32) float32 {
	return float32(i)
}

// Int64ToFloat64 converts integer64 to float64
func Int64ToFloat64(i int64) float64 {
	return float64(i)
}

// StringToBool converts string to boolean
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// StringToInt converts string to integer
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// RuneToString converts rune to string
func RuneToString(r rune) string {
	return string(r)
}

// ByteToInteger converts byte to integer
func ByteToInteger(b byte) (int, error) {
	return strconv.Atoi(string(b))
}

func Validation(s string) (validStatus bool) {
	if len(s) > 1 {
		//check string on occurence of symbols except digits and spaces
		re := regexp.MustCompile(`^[0-9 ]+$`)
		if re.MatchString(s) {
			//delete spaces
			s = strings.Replace(s, " ", "", -1)
			if checkValue(s)%10 == 0 {
				validStatus = true
			}
		}
	}
	return
}

func checkValue(s string) int {
	var index, sum int
	numberLenght := len(s)
	//if string length is odd then sum takes value of the first byte and index starts from 1
	if numberLenght%2 != 0 {
		sum, _ = ByteToInteger(s[index])
		index = 1
	}
	// loop by all bytes in the string with step 2
	for i := index; i <= numberLenght-1; i += 2 {
		sum += sumDouble(s[i], s[i+1])
	}
	return sum
}

func sumDouble(b1 byte, b2 byte) int {
	val1, _ := ByteToInteger(b1)
	val2, _ := ByteToInteger(b2)
	val1 *= 2
	if val1 > 9 {
		val1 -= 9
	}
	return val1 + val2
}
