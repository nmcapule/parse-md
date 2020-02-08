package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Type int

const (
	TypeText Type = iota
	TypeH1
	TypeH2
	TypeH3
	TypeH4
	TypeBold
	TypeItalic
	TypeNewline
	TypeDivider
)

type Node struct {
	Value    string
	Type     Type
	Children []*Node
}

func split(data []byte, eof bool) (int, []byte, error) {
	var start int

	// Skip whitespaces.
	for i := 0; i < len(data); i++ {
		start = i
		if data[i] == '\n' || !unicode.IsSpace(rune(data[i])) {
			break
		}
	}

	for i := start + 1; i < len(data); i++ {
		if data[i] == data[i-1] {
			continue
		}
		if unicode.IsSpace(rune(data[i])) {
			return i + 1, data[start:i], nil
		}
		if !unicode.IsLetter(rune(data[i])) || (unicode.IsLetter(rune(data[i])) && !unicode.IsLetter(rune(data[i-1]))) {
			return i, data[start:i], nil
		}
	}

	return start, nil, nil
}

func tokenize(reader io.Reader) []string {
	scanner := bufio.NewScanner(reader)
	scanner.Split(split)

	var tokens []string
	for scanner.Scan() {
		tokens = append(tokens, string(scanner.Bytes()))
	}

	return tokens
}

// func parse(input string) *Node {
// 	var buf bytes.Buffer
// 	// var stack []*Node
// 	for _, ch := range input {

// 	}

// 	return nil
// }

func main() {
	const test = `
	#Symbol


	This is a **bold text**.

	And this is an *italic text*.

	----

	## Hell no

	Yep. Hell no but
	multiline.
	`

	tokens := tokenize(strings.NewReader(test))
	for _, token := range tokens {
		fmt.Printf("[%s]\n", token)
	}
}
