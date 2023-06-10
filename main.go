package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/gojp/kana"
	"golang.org/x/net/html"
	"golang.org/x/text/unicode/runenames"
)

var (
	visited   = make(map[string]struct{})
	runeCount = make(map[rune]uint)
)

func printKana(text string) {
	for _, r := range text {
		s := string(r)
		isContable := kana.IsHiragana(s) || kana.IsKatakana(s) || kana.IsKanji(s)
		//isContable := kana.IsHiragana(s)
		if kana.IsKanji(s) {
			fmt.Printf("Kanji: %v\n", s)
		} else if kana.IsKatakana(s) {
			fmt.Printf("Katakana: %v\n", s)
		} else if kana.IsHiragana(s) {
			fmt.Printf("Hiragana: %v\n", s)
		} else {
			name := runenames.Name(r)
			fmt.Printf("Other: %v (%v)\n", r, name)
		}

		if isContable {
			_, ok := runeCount[r]
			if ok {
				x := runeCount[r]
				x++
				runeCount[r] = x
				continue
			}
			runeCount[r] = 1
		}
	}
}

func extractLinks(body string) []string {
	links := make([]string, 0)
	reader := strings.NewReader(body)
	tokenizer := html.NewTokenizer(reader)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			return links
		}

		token := tokenizer.Token()

		if tokenType == html.StartTagToken && token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
	}
}

func crawler(url string) {

	// prevent infinite loop
	_, ok := visited[url]
	if ok {
		return
	}

	fmt.Println("crawling:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error fetching url:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body:", err)
		return
	}

	bstr := string(body)

	printKana(bstr)

	links := extractLinks(bstr)
	for _, link := range links {
		if !strings.HasPrefix(link, "http") { // http or https
			link = url + link
		}
		//crawler(link)
	}
}

func main() {
	url := "https://www.jst.go.jp"
	crawler(url)

	for k, v := range runeCount {
		fmt.Printf("%v: %v\n", string(k), v)
	}

	// convert map to slice

	type kv struct {
		Kana  string
		Value uint
	}

	kanaSlice := make([]kv, 0, len(runeCount))
	i := 0
	for k := range runeCount {
		ka := kv{
			Kana:  string(k),
			Value: runeCount[k],
		}
		kanaSlice = append(kanaSlice, ka)
		i++
	}

	// sort by value
	sort.Slice(kanaSlice, func(i, j int) bool {
		return kanaSlice[i].Value > kanaSlice[j].Value
	})

	fmt.Println("Top 100:")

	// print top
	for i := 0; i < len(kanaSlice) && i < 100; i++ {
		fmt.Printf("%v: %v\n", kanaSlice[i].Kana, kanaSlice[i].Value)
	}

}
