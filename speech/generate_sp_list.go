package main
//package speech

import (
	"fmt"
	"strings"
	"bufio"
	"bytes"
	"os"
	"flag"
)
func ScanSpecil(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), nil, nil
	}

	if i := bytes.IndexByte(data, ':'); i >= 0 {
		items := strings.FieldsFunc(string(data[0:i]), func (r rune) bool {
			return strings.Contains("., ", string(r))
		})
		return i + 1, []byte(items[len(items)-2]), nil
	}
	// Request more data.
	return 0, nil, nil
}


func GenerateSlice(path string) ([]string,int){
	ifr, _ := os.Open(path)
	defer ifr.Close()
	ifbr := bufio.NewScanner(ifr)
	ifbr.Split(ScanSpecil)
	res := make([]string, 0)
	for ifbr.Scan() {
		fmt.Println("items: ", ifbr.Text())
		res = append(res, ifbr.Text())
	}
	return res,8
}

func Combine(sl []string, width int) []string {
	length:=len(sl)
	fmt.Println("Length of the Slice : ", length)
	fmt.Println("Width : ", width)
	for i:=0; i< length; i=i+width {
		for j:=1; j<width; j++ {
			pair := []string{sl[i], sl[i+j], "t"}
			item := strings.Join(pair," ")
			fmt.Println(item)
		}
		for k:=i+width; k<length; k=k+width {
			pair := []string{sl[i], sl[k], "f"}
			item := strings.Join(pair," ")
			fmt.Println(item)
		}
	}
	return make([]string, 10)
}


func main() {
	var infile string
	flag.StringVar(&infile, "i", "", "input file ")
	flag.Parse()
	fmt.Println("Test Combine ", infile)
	Combine(GenerateSlice(infile))
}










