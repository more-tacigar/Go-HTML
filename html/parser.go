// ==================================================
// Copyright (c) 2016 tacigar
// https://github.com/tacigar/Go-HTML
// ==================================================

package html

import (
	"io"
)

type parser struct {
	tokenizer *Tokenizer
	nextToken *Token
}

// HTMLをパースして先頭ノードを返す
func Parse(reader io.Reader) *Node {
	tokenizer := NewTokenizer(reader)
	parser := &parser{
		tokenizer: tokenizer,
		nextToken: nil,
	}
	parser.nextToken = tokenizer.Next()
	// あとは適当に再帰するだけ
	var parseImpl func() *Node
	parseImpl = func() *Node {
		switch parser.nextToken.Type {
		case TextToken:
			newNode := &Node{
				Parent:     nil, // 親側で設定
				Children:   []*Node{},
				Type:       TextNode,
				Data:       parser.nextToken.Data,
				Attributes: map[string]string{},
			}
			parser.nextToken = tokenizer.Next()
			return newNode
		case SelfClosingTagToken:
			newNode := &Node{
				Parent:     nil,
				Children:   []*Node{},
				Type:       ElementNode,
				Data:       parser.nextToken.Data,
				Attributes: parser.nextToken.Attributes,
			}
			parser.nextToken = tokenizer.Next()
			return newNode
		case StartTagToken:
			newNode := &Node{
				Parent:     nil,
				Children:   []*Node{},
				Type:       ElementNode,
				Data:       parser.nextToken.Data,
				Attributes: parser.nextToken.Attributes,
			}
			parser.nextToken = tokenizer.Next()
			for {
				if parser.nextToken.Type == EndTagToken &&
					parser.nextToken.Data == newNode.Data {
					parser.nextToken = tokenizer.Next()
					return newNode
				} else {
					child := parseImpl()
					child.Parent = newNode
					newNode.Children = append(newNode.Children, child)
				}
			}
		}
		return nil
	}
	return parseImpl()
}
