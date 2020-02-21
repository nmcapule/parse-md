package main

import (
	"fmt"
	"regexp"
)

var (
	whitespaceRe = regexp.MustCompile("^\\s+")
	blockDelimRe = regexp.MustCompile("^\n\n")
	escapeRe     = regexp.MustCompile("^\\\\.")
	highlightRe  = regexp.MustCompile("^\\*\\*")
	headerRe     = regexp.MustCompile("^#+")
	codeRe       = regexp.MustCompile("^`")
	codeBlockRe  = regexp.MustCompile("^```")
	rulerRe      = regexp.MustCompile("^\\-{3,}")
)

type Token struct {
	Kind     string
	Value    []byte
	Parent   *Token
	Children []*Token
}

func (t *Token) PrintTree(level int) {
	for i := 0; i <= level; i++ {
		fmt.Printf("----")
	}
	fmt.Println(t.String())
	for _, node := range t.Children {
		node.PrintTree(level + 1)
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s: %s", t.Kind, string(t.Value))
}

type Parser struct {
	tokens []*Token
	data   []byte
	idx    int
}

func NewParser(data []byte) *Parser {
	return &Parser{data: data}
}

func (p *Parser) Print() {
	for _, node := range p.tokens {
		node.PrintTree(0)
	}
}

func (p *Parser) Rest() []byte {
	return p.data[p.idx:]
}

func (p *Parser) EOF() bool {
	return p.idx >= len(p.data)
}

func (p *Parser) skipRegex(re *regexp.Regexp) {
	_, _, right := p.findIndex(re)
	p.idx += right
}

func (p *Parser) findIndex(re *regexp.Regexp) ([]byte, int, int) {
	idxs := re.FindIndex(p.Rest())
	if idxs == nil {
		return nil, 0, 0
	}
	return p.Rest()[idxs[0]:idxs[1]], idxs[0], idxs[1]
}

func (p *Parser) startsRegex(re *regexp.Regexp) bool {
	return re.Match(p.Rest())
}

func (p *Parser) Parse() {
	for !p.EOF() {
		p.skipRegex(whitespaceRe)
		if headerRe.Match(p.Rest()) {
			p.tokens = append(p.tokens, p.parseHeader())
		} else if rulerRe.Match(p.Rest()) {
			p.tokens = append(p.tokens, &Token{Kind: "ruler"})
			p.skipRegex(rulerRe)
		} else {
			p.tokens = append(p.tokens, &Token{
				Kind:     "block",
				Children: p.parseExpression(blockDelimRe),
			})
		}
	}
}

func (p *Parser) parseHeader() *Token {
	match, _, right := p.findIndex(headerRe)
	p.idx += right

	token := &Token{
		Kind:     "header",
		Value:    match,
		Children: p.parseExpression(blockDelimRe),
	}
	return token
}

func (p *Parser) parseExpression(delimRe *regexp.Regexp) []*Token {
	var tokens []*Token

	p.skipRegex(whitespaceRe)

	start := p.idx
	for !p.startsRegex(delimRe) && !p.EOF() {
		p.skipRegex(escapeRe)
		if p.startsRegex(codeBlockRe) {
			if p.idx > start {
				tokens = append(tokens, &Token{Kind: "plain", Value: p.data[start:p.idx]})
			}
			tokens = append(tokens, p.parseCode(codeBlockRe))
			start = p.idx
			continue
		}
		if p.startsRegex(codeRe) {
			if p.idx > start {
				tokens = append(tokens, &Token{Kind: "plain", Value: p.data[start:p.idx]})
			}
			tokens = append(tokens, p.parseCode(codeRe))
			start = p.idx
			continue
		}
		if p.startsRegex(highlightRe) {
			if p.idx > start {
				tokens = append(tokens, &Token{Kind: "plain", Value: p.data[start:p.idx]})
			}
			tokens = append(tokens, p.parseHighlight())
			start = p.idx
			continue
		}
		p.idx += 1
	}

	if p.idx > start {
		tokens = append(tokens, &Token{Kind: "plain", Value: p.data[start:p.idx]})
	}

	p.skipRegex(delimRe)
	p.skipRegex(whitespaceRe)

	return tokens
}

func (p *Parser) parseHighlight() *Token {
	p.skipRegex(highlightRe)
	return &Token{
		Kind:     "highlight",
		Children: p.parseExpression(highlightRe),
	}
}

func (p *Parser) parseCode(delim *regexp.Regexp) *Token {
	p.skipRegex(delim)

	start := p.idx
	for !p.startsRegex(delim) && !p.EOF() {
		p.skipRegex(escapeRe)
		p.idx += 1
	}

	value := p.data[start:p.idx]

	p.skipRegex(delim)

	return &Token{Kind: "code", Value: value}
}

func main() {
	const test = `
	#**Sym**bol

	This is a **bold text**.
	as a multiline expr

	And this is a ***scandalous text***.
	asdsad
	asdasds i ` + "`code`" + `.

	And \[this] is an *italic text*.

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

	I am a \[link](http://markdown.com) here
	and multiline

	----

	## Hell no

	Yep. Hell no but
	multiline.
	`

	parser := NewParser([]byte(test))
	parser.Parse()
	parser.Print()
}
