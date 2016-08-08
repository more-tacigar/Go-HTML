package html

import (
	"unicode"
)

// 現実のHTMLパーサの実装がどのようになっているのかはわからないが、今回の
// 用途的には空白文字だけからなる要素は特にテキスト要素として保存しておく必要
// はないので、とりあえず無視することとする。それのチェックを行う関数。
func isSpaceString(text string) bool {
	for _, runeValue := range text {
		if !unicode.IsSpace(runeValue) {
			return false
		}
	}
	return true
}
