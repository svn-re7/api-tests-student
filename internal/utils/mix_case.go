package utils

import "unicode"

// MixCase чередует нижний и верхний регистр букв, включая кириллицу.
// Пробелы, цифры и знаки не меняют порядок чередования.
func MixCase(text string) string {
	runes := []rune(text)
	upper := false
	for i, char := range runes {
		if !unicode.IsLetter(char) {
			continue
		}
		if upper {
			runes[i] = unicode.ToUpper(char)
		} else {
			runes[i] = unicode.ToLower(char)
		}
		upper = !upper
	}
	return string(runes)
}
