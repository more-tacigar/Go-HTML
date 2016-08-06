// ==================================================
// Copyright (c) 2016 tacigar
// https://github.com/tacigar/Go-HTML
// ==================================================

package html

// 空白文字かを判定する
func isSpace(runeValue rune) bool {
	switch runeValue {
	case rune('\n'), rune('\t'), rune('\f'), rune('\r'), rune(' '):
		return true
	default:
		return false
	}
}

// 現実のHTMLパーサの実装がどのようになっているのかはわからないが、今回の
// 用途的には空白文字だけからなる要素は特にテキスト要素として保存しておく必要
// はないので、とりあえず無視することとする。それのチェックを行う関数。
func isSpaceString(text string) bool {
	for _, runeValue := range text {
		if !isSpace(runeValue) {
			return false
		}
	}
	return true
}
