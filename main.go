package main

import (
	"fmt"

	"github.com/gojp/kana"
	"golang.org/x/text/unicode/runenames"
)

func main() {
	text := "日本語のテキストです"
	for _, r := range text {
		if kana.IsKanji(string(r)) {
			fmt.Printf("Kanji: %c\n", r)
		} else if kana.IsKatakana(string(r)) {
			fmt.Printf("Katakana: %c\n", r)
		} else if kana.IsHiragana(string(r)) {
			fmt.Printf("Hiragana: %c\n", r)
		} else {
			name := runenames.Name(r)
			fmt.Printf("Outro: %c (%s)\n", r, name)
		}
	}
}
