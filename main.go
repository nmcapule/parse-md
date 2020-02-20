package main

import (
	"fmt"
)

func startsWith(data, prefix []byte) bool {
	if len(prefix) > len(data) {
		return false
	}
	for i, ch := range prefix {
		if ch != data[i] {
			return false
		}
	}
	return true
}

type tokenizer struct {
	tokens [][]byte
	data   []byte
	cursor int
}

func newTokenizer(data []byte) *tokenizer {
	return &tokenizer{data: data}
}

func (t *tokenizer) skipSpace() {
	for {
		if t.eof() {
			break
		}
		if t.data[t.cursor] == ' ' {
			t.cursor += 1
		}
		if t.data[t.cursor] == '\n' {
			t.cursor += 1
		}
		if t.data[t.cursor] == '\t' {
			t.cursor += 1
		}
		break
	}
}

func (t *tokenizer) push(kind string, token []byte) {
	// t.tokens = append(t.tokens, []byte(kind+": "+string(token)))
	t.tokens = append(t.tokens, token)
}

func (t *tokenizer) tryConsume(token []byte) bool {
	if t.nextStartsWith(token) {
		t.push("keyword", token)
		t.cursor += len(token)
		return true
	}
	return false
}

func (t *tokenizer) nextStartsWith(token []byte) bool {
	return startsWith(t.data[t.cursor:], token)
}

func (t *tokenizer) readExpression() bool {
	start := t.cursor
	for {
		if t.eof() {
			token := t.data[start:t.cursor]
			t.push("expr", token)
			break
		}
		if ok := t.nextStartsWith([]byte("\n\n")); ok {
			token := t.data[start:t.cursor]
			t.push("expr", token)
			t.tryConsume([]byte("\n\n"))
			break
		}
		t.cursor += 1
	}
	t.skipSpace()

	return true
}

func (t *tokenizer) readHighlight() bool {
	t.skipSpace()
	// t.readExpression

	return true
}

func (t *tokenizer) readHeader() bool {
	t.tryConsume([]byte("####"))
	t.tryConsume([]byte("###"))
	t.tryConsume([]byte("##"))
	t.tryConsume([]byte("#"))

	t.skipSpace()

	t.readExpression()

	return true
}

func (t *tokenizer) readCodeBlock() bool {
	t.tryConsume([]byte("```"))

	start := t.cursor
	for {
		if t.eof() {
			token := t.data[start:t.cursor]
			t.push("code", token)
			break
		}
		if ok := t.nextStartsWith([]byte("```")); ok {
			token := t.data[start:t.cursor]
			t.push("code", token)
			t.tryConsume([]byte("```"))
			break
		}
		t.cursor += 1
	}

	t.skipSpace()

	return true
}

func (t *tokenizer) read() bool {
	t.skipSpace()

	if startsWith(t.data[t.cursor:], []byte("#")) {
		return t.readHeader()
	}
	if startsWith(t.data[t.cursor:], []byte("```")) {
		return t.readCodeBlock()
	}

	return t.readExpression()
}

func (t *tokenizer) eof() bool {
	return t.cursor >= len(t.data)
}

func main() {
	const test = `
	#**Sym**bol

	This is a **bold text**.
	as a multiline expr

	And this is an *italic text*.

	And

	` + "```" + `
	code block
	i am

	a fucking

	another

	code block





	shit
	` + "```" + `

	## Ohshit

	I am a [link](http://markdown.com)

	----

	## Hell no

	Yep. Hell no but
	multiline.
	`

	tok := newTokenizer([]byte(test))
	for !tok.eof() {
		tok.read()
	}

	for _, token := range tok.tokens {
		fmt.Println("[" + string(token) + "]")
	}
}
