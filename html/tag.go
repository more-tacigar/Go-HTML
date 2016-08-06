// ==================================================
// Copyright (c) 2016 tacigar
// https://github.com/tacigar/Go-HTML
// ==================================================

package html

import (
	"bufio"
	"strings"
)

type tagParser struct {
	reader   *bufio.Reader
	nextRune rune
	buffer   []rune
}

// Self-Closing Tagを判定する方法が全然思いつかないので、普通にタグ名で判定することにします
var selfClosingTags = [...]string{
	"area",
	"base",
	"br",
	"col",
	"command",
	"embed",
	"hr",
	"img",
	"input",
	"keygen",
	"link",
	"meta",
	"param",
	"source",
	"track",
	"wbr",
}

func newTagParser(tagStr string) *tagParser {
	return &tagParser{
		reader:   bufio.NewReader(strings.NewReader(tagStr)),
		nextRune: -1,
		buffer:   []rune{},
	}
}

// 空白文字(半角スペース、'\n'、'\r'、'\f'、'\t')を読み飛ばす
func (tagParser *tagParser) skipSpace() {
	for isSpace(tagParser.nextRune) {
		runeValue, _, _ := tagParser.reader.ReadRune()
		tagParser.nextRune = runeValue
	}
}

// tagParserで１文字先読みする
func (tagParser *tagParser) readNext() rune {
	runeValue, _, _ := tagParser.reader.ReadRune()
	tagParser.nextRune = runeValue
	return runeValue
}

func (tagParser *tagParser) readNextIdentifier() string {
	for {
		tagParser.buffer = append(tagParser.buffer, tagParser.nextRune)
		tagParser.readNext()
		if isSpace(tagParser.nextRune) || tagParser.nextRune == rune('>') ||
			tagParser.nextRune == rune('=') || tagParser.nextRune == rune('/') {
			return string(tagParser.buffer)
		}
	}
}

func parseTagToken(tagStr string) *Token {
	parser := newTagParser(tagStr)
	parser.readNext()                 // '<'を先読み
	parser.readNext()                 // '<'の次の文字を先読みしておく
	if parser.nextRune == rune('!') { // この場合はコメントかDOCTYPEなので無視
		return nil
	}
	parser.skipSpace()
	parser.buffer = []rune{}
	if parser.nextRune == rune('/') { // 終了タグの場合はシンプル
		tagName := parser.readNextIdentifier()
		return &Token{
			Data:       tagName[1:], // '/'は含めないこととする
			Type:       EndTagToken,
			Attributes: map[string]string{},
		}
	}
	tagName := parser.readNextIdentifier()
	parser.skipSpace()
	parser.buffer = []rune{}
	attributes := map[string]string{}
	for {
		if parser.nextRune == rune('/') || parser.nextRune == rune('>') {
			for _, selfClosingTag := range selfClosingTags {
				if selfClosingTag == tagName {
					return &Token{
						Data:       tagName,
						Type:       SelfClosingTagToken,
						Attributes: attributes,
					}
				}
			}
			return &Token{
				Data:       tagName,
				Type:       StartTagToken,
				Attributes: attributes,
			}
		}
		key := parser.readNextIdentifier()
		parser.skipSpace()
		if parser.nextRune != rune('=') {
			attributes[key] = ""
			continue
		}
		parser.readNext()
		parser.skipSpace()
		parser.buffer = []rune{}
		delimiter := parser.nextRune
		parser.readNext()
		for parser.nextRune != delimiter {
			parser.buffer = append(parser.buffer, parser.nextRune)
			parser.readNext()
		} // nextRuneは'\"'
		value := string(parser.buffer)
		parser.readNext()
		parser.skipSpace()
		attributes[key] = value
		parser.buffer = []rune{}
	}
}
