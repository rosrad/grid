package gana

import (
	"strings"
)

func ganas() string {
	var re []rune
	for hira := 'あ'; hira != 'ん'+1; hira++ {
		re = append(re, hira)
	}
	for kata := 'ア'; kata != 'ン'+1; kata++ {
		re = append(re, kata)
	}
	return string(re)
}

func Isgana(word string) bool {
	allganas := ganas()
	for _, c := range word {
		if !strings.ContainsAny(allganas, string(c)) {
			return false
		}
	}
	return true
}
