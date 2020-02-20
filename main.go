package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var (
	wordsRe = regexp.MustCompile(`^[\w \t\.\-]+`)
)

func scopedEqual(left, right []byte) bool {
	min := len(left)
	if len(right) < min {
		min = len(right)
	}

	return string(left[:min]) == string(right[:min])
}

func readWhile(data []byte, while byte) int {
	for i, ch := range data {
		if ch != while {
			return i
		}
	}
	return len(data)
}

func readUntilBytes(data []byte, until []byte) int {
	for i := range data {
		if scopedEqual(data[i:], until) {
			return i + len(until)
		}
	}
	return len(data)
}

func markdownSplitFunc(data []byte, eof bool) (int, []byte, error) {
	if eof && len(data) == 0 {
		return 0, nil, nil
	}

	if scopedEqual(data, []byte("\n\n")) {
		return 2, []byte("-newline-"), nil
	}
	if scopedEqual(data, []byte("'''")) {
		advance := readUntilBytes(data[3:], []byte("'''")) + 3
		return advance, data[:advance], nil
	}
	switch data[0] {
	case '*', '#', '-', '\n', '\t':
		fallthrough
	case '[', ']', '(', ')':
		advance := readWhile(data, data[0])
		return advance, data[:advance], nil
	}

	match := wordsRe.Find(data)

	return len(match), match, nil
}

func main() {
	const test = `
	#**Sym**bol


	This is a **bold text**.

	And this is an *italic text*.

	And

	'''
	code block
	i am
	'''

	## Ohshit

	I am a [link](http://markdown.com)

	----

	## Hell no

	Yep. Hell no but
	multiline.
	`

	scanner := bufio.NewScanner(strings.NewReader(test))
	scanner.Split(markdownSplitFunc)
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" {
			continue
		}
		fmt.Println("[" + strings.TrimSpace(scanner.Text()) + "]")
	}
}
