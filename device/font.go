package device

import (
	"fmt"
	"strconv"
)

type Font [][]byte

func (Font) FromInt(value int) ([][]byte, error) {
	if -99 > value || value > 99 {
		return nil, fmt.Errorf("Value exceeded (-100;100) range. Given value: %d", value)
	}

	return [][]byte{
		Rotate(DEFAULT_FONT[value/10]), // first digit
		Rotate(DEFAULT_FONT[value%10]), // second digit
	}, nil
}

func addTrailingZeroes(str string, targetLen int) string {
	for len(str) < targetLen {
		str = "0" + str
	}

	return str
}

func transpose(tab [][]rune) [][]rune {
	res := make([][]rune, 0)

	x := len(tab[0])
	y := len(tab)

	tmp := make([]rune, 0)
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			tmp = append(tmp, tab[j][x-i-1])
		}
		res = append(res, tmp)
		tmp = make([]rune, 0)
	}

	return res
}

func Rotate(tab []byte) []byte {
	strTab := make([][]rune, 0)
	for _, b := range tab {
		strTab = append(strTab, []rune(addTrailingZeroes(strconv.FormatInt(int64(b), 2), 8)))
	}

	strTab = transpose(strTab)
	res := make([]byte, 0)
	for _, binaryString := range strTab {
		i, _ := strconv.ParseInt(string(binaryString), 2, 64)
		res = append(res, byte(i))
	}

	return res
}

var DEFAULT_FONT = Font{
	{0x3E, 0x7F, 0x71, 0x59, 0x4D, 0x7F, 0x3E, 0x00}, // '0'
	{0x40, 0x42, 0x7F, 0x7F, 0x40, 0x40, 0x00, 0x00}, // '1'
	{0x62, 0x73, 0x59, 0x49, 0x6F, 0x66, 0x00, 0x00}, // '2'
	{0x22, 0x63, 0x49, 0x49, 0x7F, 0x36, 0x00, 0x00}, // '3'
	{0x18, 0x1C, 0x16, 0x53, 0x7F, 0x7F, 0x50, 0x00}, // '4'
	{0x27, 0x67, 0x45, 0x45, 0x7D, 0x39, 0x00, 0x00}, // '5'
	{0x3C, 0x7E, 0x4B, 0x49, 0x79, 0x30, 0x00, 0x00}, // '6'
	{0x03, 0x03, 0x71, 0x79, 0x0F, 0x07, 0x00, 0x00}, // '7'
	{0x36, 0x7F, 0x49, 0x49, 0x7F, 0x36, 0x00, 0x00}, // '8'
	{0x06, 0x4F, 0x49, 0x69, 0x3F, 0x1E, 0x00, 0x00}, // '9'
}
