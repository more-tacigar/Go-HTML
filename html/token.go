// ==================================================
// Copyright (c) 2016 tacigar
// https://github.com/tacigar/Go-HTML
// ==================================================

package html

import (
	"bufio"
	"io"
)

// HTMLの構文解析においては、字句解析でやるような細かいトークンに分割する
// 必要性がないように感じたので、HTML要素をトークンとして扱うこととする。

type TokenType int

const (
	TextToken TokenType = iota
	StartTagToken
	EndTagToken
	SelfClosingTagToken
	MaybeSelfClosingTagToken
)

type Token struct {
	Data       string
	Type       TokenType
	Attributes map[string]string
}

type TokenizerState int

const (
	foundText TokenizerState = iota
	foundTag
)

type Tokenizer struct {
	reader   *bufio.Reader
	nextRune rune
	buffer   []rune
	state    TokenizerState
}

func NewTokenizer(reader io.Reader) *Tokenizer {
	tokenizer := &Tokenizer{
		reader:   bufio.NewReader(reader),
		nextRune: -1,
		state:    foundText,
		buffer:   []rune{},
	}
	tokenizer.readNext()
	return tokenizer
}

// Tokenizerで１文字先読みする
func (tokenizer *Tokenizer) readNext() (rune, error) {
	runeValue, _, err := tokenizer.reader.ReadRune()
	if err != nil {
		return -1, err
	}
	tokenizer.nextRune = runeValue
	return runeValue, nil
}

// Tokenizerから次のトークンを吐き出す
func (tokenizer *Tokenizer) Next() *Token {
	tokenizer.buffer = []rune{}
	switch tokenizer.state {
	case foundText:
		for {
			if tokenizer.nextRune == rune('<') {
				tokenizer.state = foundTag
				if isSpaceString(string(tokenizer.buffer)) {
					return tokenizer.Next()
				} else {
					return &Token{
						Data:       string(tokenizer.buffer),
						Type:       TextToken,
						Attributes: map[string]string{},
					}
				}
			} else {
				tokenizer.buffer = append(tokenizer.buffer, tokenizer.nextRune)
				_, err := tokenizer.readNext()
				if err == io.EOF {
					return nil
				}
			}
		}
	case foundTag:
		for {
			if tokenizer.nextRune == rune('>') {
				tokenizer.state = foundText
				tokenizer.buffer = append(tokenizer.buffer, tokenizer.nextRune)
				tokenizer.readNext()
				token := parseTagToken(string(tokenizer.buffer))
				if token != nil {
					return token
				} else { // コメントorDOCTYPEは無視する
					return tokenizer.Next()
				}
			} else {
				tokenizer.buffer = append(tokenizer.buffer, tokenizer.nextRune)
				tokenizer.readNext()
			}
		}
	}
	return nil
}
