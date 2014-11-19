package main

import (
	"flag"
	"fmt"
	"gana"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"util"
)

func fetch_speech(word string, dir string) bool {
	url := "http://translate.google.com/translate_tts?ie=UTF-8&tl=ja&q=" + url.QueryEscape(word)
	fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.101 Safari/537.36")
	req.Header.Set("Host", "translate.google.com")

	audio_dir := filepath.Join(dir, "audios")
	mp3 := filepath.Join(audio_dir, word+".mp3")
	os.MkdirAll(audio_dir, os.ModeDir)
	out, err := os.Create(mp3)
	defer out.Close()
	if err != nil {
		return false
	}
	io.Copy(out, resp.Body)
	return true
}

func parse_keyword(line string) (key string, newline string) {
	words := strings.Split(line, "\t")
	if len(words) < 2 {
		fmt.Println("no effictive word !  ", line)
		return
	}
	key = words[1]
	for index, word := range words {
		if index == 0 {
			continue
		}
		if !gana.Isgana(word) {
			key = word
		}
	}
	newline = words[0] + "[sound:" + key + ".mp3]\t" + strings.Join(words[1:], ":")
	fmt.Println(key, newline)
	return
}

func main() {
	var list, outdir string
	var justfmt bool
	flag.StringVar(&list, "list", "", "file of the words")
	flag.StringVar(&outdir, "outdir", "", "the output of the mp3 place and new config file")
	flag.BoolVar(&justfmt, "justfmt", false, "just format the new configfile no audio file downloading")
	flag.Parse()
	list, _ = filepath.Abs(list)
	if len(outdir) <= 0 {
		outdir = filepath.Join(filepath.Dir(list), "words_audio")
	}
	os.RemoveAll(outdir)
	os.MkdirAll(outdir, os.ModeDir)
	lines, _ := util.ReadLines(list)
	out := filepath.Join(outdir, filepath.Base(list)+".new")
	fout, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	defer fout.Close()
	if err != nil {
		fmt.Println("create output file Error", err)
		return
	}

	for _, word := range lines {
		key, newline := parse_keyword(word)
		fout.WriteString(newline + "\r\n")
		if !justfmt {
			fetch_speech(key, outdir)
		}
	}
}
