package html

import (
	"bufio"
	"io"
	"unicode"
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

func newTagParser(reader io.Reader) *tagParser {
	return &tagParser{
		reader:   bufio.NewReader(reader),
		nextRune: -1,
		buffer:   []rune{},
	}
}

// 空白文字(半角スペース、'\n'、'\r'、'\f'、'\t')を読み飛ばす
func (tagParser *tagParser) skipSpace() {
	for unicode.IsSpace(tagParser.nextRune) {
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
		if unicode.IsSpace(tagParser.nextRune) || tagParser.nextRune == '>' ||
			tagParser.nextRune == '=' || tagParser.nextRune == '/' {
			return string(tagParser.buffer)
		}
	}
}

func parseTagToken(reader io.Reader) *Token {
	parser := newTagParser(reader)
	parser.readNext()           // '<'を先読み
	parser.readNext()           // '<'の次の文字を先読みしておく
	if parser.nextRune == '!' { // この場合はコメントかDOCTYPEなので無視
		return nil
	}
	parser.skipSpace()
	parser.buffer = []rune{}
	if parser.nextRune == '/' { // 終了タグの場合はシンプル
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
		if parser.nextRune == '/' || parser.nextRune == '>' {
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
		if parser.nextRune != '=' { // 次が属性キーである場合
			attributes[key] = ""
			continue
		} // else 次は属性値
		parser.readNext()
		parser.skipSpace()
		parser.buffer = []rune{}
		delimiter := parser.nextRune // デリミタは'\''or'"'なので抜き出しとく
		parser.readNext()
		for parser.nextRune != delimiter { // 同じデリミタが来るまで
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
