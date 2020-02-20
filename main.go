package main

import (
	"fmt"
)

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

func (t *tokenizer) nextStartsWith(prefix []byte) bool {
	if len(prefix) > len(t.data[t.cursor:]) {
		return false
	}
	for i, ch := range prefix {
		if ch != t.data[i+t.cursor] {
			return false
		}
	}
	return true
}

func (t *tokenizer) readLink() bool {
	t.tryConsume([]byte("["))

	t.readExpression([]byte("]"))

	if ok := t.tryConsume([]byte("(")); !ok {
		return true
	}

	start := t.cursor
	for !t.nextStartsWith([]byte(")")) {
		t.cursor += 1
	}
	t.push("link", t.data[start:t.cursor])

	t.tryConsume([]byte(")"))

	return true
}

func equal(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for i, ch := range left {
		if ch != right[i] {
			return false
		}
	}
	return true
}

func (t *tokenizer) readExpression(delim []byte) bool {
	start := t.cursor
	for {
		if t.eof() {
			t.push("expr", t.data[start:t.cursor])
			break
		}
		if ok := t.nextStartsWith(delim); ok {
			t.push("expr", t.data[start:t.cursor])
			break
		}
		if ok := t.nextStartsWith([]byte("[")); ok {
			t.push("expr", t.data[start:t.cursor])
			t.readLink()
			start = t.cursor
			continue
		}
		if ok := t.nextStartsWith([]byte("**")); ok {
			t.tryConsume([]byte("**"))
			t.readExpression([]byte("**"))
			start = t.cursor
			continue
		}
		t.cursor += 1
	}
	t.tryConsume(delim)
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

	t.readExpression([]byte("\n\n"))

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

	if t.nextStartsWith([]byte("#")) {
		return t.readHeader()
	}
	if t.nextStartsWith([]byte("```")) {
		return t.readCodeBlock()
	}

	return t.readExpression([]byte("\n\n"))
}

func (t *tokenizer) eof() bool {
	return t.cursor >= len(t.data)
}

func main() {
	const test = `
	#**Sym**bol

	This is a **bold text**.
	as a multiline expr

	And this is a ***scandalous text***.

	And [this] is an *italic text*.

	And

	` + "```" + `
	code block
	i am
	a [fucking]
	another
	code block
	shit
	` + "```" + `

	## Ohshit

	I am a [link](http://markdown.com) here
	and multiline

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
		fmt.Println("<" + string(token) + ">")
	}
}
