package main

import (
	// "bufio"
	// "flag"
	"fmt"
	// "gana"
	"io"
	"net/http"
	// "net/url"
	"os"
	"path/filepath"
	"strings"
)

func ganas() string {
	var re []rune
	for hira := 'あ'; hira != 'ん'; hira++ {
		re = append(re, hira)
	}
	for kata := 'ア'; kata != 'ン'; kata++ {
		re = append(re, kata)
	}
	return string(re)
}

func is_gana(word string) bool {
	allganas := ganas()
	for _, c := range word {
		fmt.Printf("%c\n", c)
		if !strings.ContainsAny(allganas, string(c)) {
			fmt.Printf("no found : %c\n", c)
			return false

		}
	}
	return true
}

var Cookies = []*http.Cookie{
	{Name: "__utma", Value: "192312886.309048475.1382586676.1382586676.1382591041.2"},
	{Name: "__utmb", Value: "192312886.3.10.1382591041"},
	{Name: "__utmc", Value: "192312886"},
	{Name: "__utmz", Value: "192312886.1382586676.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none)"},
}

func download(url string) bool {
	fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36")
	req.Header.Set("Host", "www.csse.monash.edu.au")
	req.Header.Set("Referer", "http://www.csse.monash.edu.au/~jwb/cgi-bin/wwwjdic.cgi?1F")

	mp3 := filepath.Join("test", "test.swf")
	os.MkdirAll(filepath.Dir(mp3), os.ModeDir)
	out, err := os.Create(mp3)
	defer out.Close()
	if err != nil {
		return false
	}
	io.Copy(out, resp.Body)
	return true
}

// var addCookieTests = []struct {
// 	Cookies []*http.Cookie
// 	Raw     string
// }{
// 	{
// 		[]*http.Cookie{},
// 		"",
// 	},
// 	{
// 		[]*http.Cookie{{Name: "cookie-1", Value: "v$1"}},
// 		"cookie-1=v$1",
// 	},
// 	{
// 		[]*http.Cookie{
// 			{Name: "cookie-1", Value: "v$1"},
// 			{Name: "cookie-2", Value: "v$2"},
// 			{Name: "cookie-3", Value: "v$3"},
// 		},
// 		"cookie-1=v$1; cookie-2=v$2; cookie-3=v$3",
// 	},
// }

// func TestAddCookie() {
// 	for i, tt := range addCookieTests {
// 		req, _ := http.NewRequest("GET", "http://example.com/", nil)
// 		for _, c := range tt.Cookies {
// 			req.AddCookie(c)
// 		}
// 		g := req.Header.Get("Cookie")
// 		fmt.Printf("Readed %s \n", g)
// 		if g != tt.Raw {

// 			fmt.Printf("Test %d:\nwant: %s\n got: %s\n", i, tt.Raw, g)
// 			continue
// 		}
// 	}
// }
func get_speech() {
	url := "http://www.csse.monash.edu.au/~jwb/audiock.swf?u=kana=%25E3%2581%259F%25E3%2581%25B9%25E3%2582%2582%25E3%2581%25AE%25E3%2582%2584%26kanji=%25E9%25A3%259F%25E3%2581%25B9%25E7%2589%25A9%25E5%25B1%258B"
	download(url)
}
// func main() {
// 	get_speech()
// }
func main() {					// for test some thing
	fmt.Printf("Hello, 世界\n")
}

