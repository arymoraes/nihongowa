package services

import (
	"strings"
)

var hiragana = map[string]string{
	"a": "あ", "i": "い", "u": "う", "e": "え", "o": "お",
	"ka": "か", "ki": "き", "ku": "く", "ke": "け", "ko": "こ",
	"ga": "が", "gi": "ぎ", "gu": "ぐ", "ge": "げ", "go": "ご",
	"sa": "さ", "shi": "し", "su": "す", "se": "せ", "so": "そ",
	"za": "ざ", "ji": "じ", "zu": "ず", "ze": "ぜ", "zo": "ぞ",
	"ta": "た", "chi": "ち", "tsu": "つ", "te": "て", "to": "と",
	"da": "だ", "de": "で", "do": "ど",
	"na": "な", "ni": "に", "nu": "ぬ", "ne": "ね", "no": "の",
	"ha": "は", "hi": "ひ", "fu": "ふ", "he": "へ", "ho": "ほ",
	"ba": "ば", "bi": "び", "bu": "ぶ", "be": "べ", "bo": "ぼ",
	"pa": "ぱ", "pi": "ぴ", "pu": "ぷ", "pe": "ぺ", "po": "ぽ",
	"ma": "ま", "mi": "み", "mu": "む", "me": "め", "mo": "も",
	"ya": "や", "yu": "ゆ", "yo": "よ",
	"ra": "ら", "ri": "り", "ru": "る", "re": "れ", "ro": "ろ",
	"wa": "わ", "wo": "を",
	"n":   "ん",
	"kya": "きゃ", "kyu": "きゅ", "kyo": "きょ",
	"gya": "ぎゃ", "gyu": "ぎゅ", "gyo": "ぎょ",
	"sha": "しゃ", "shu": "しゅ", "sho": "しょ",
	"ja": "じゃ", "ju": "じゅ", "jo": "じょ",
	"cha": "ちゃ", "chu": "ちゅ", "cho": "ちょ",
	"nya": "にゃ", "nyu": "にゅ", "nyo": "にょ",
	"hya": "ひゃ", "hyu": "ひゅ", "hyo": "ひょ",
	"bya": "びゃ", "byu": "びゅ", "byo": "びょ",
	"pya": "ぴゃ", "pyu": "ぴゅ", "pyo": "ぴょ",
	"mya": "みゃ", "myu": "みゅ", "myo": "みょ",
	"rya": "りゃ", "ryu": "りゅ", "ryo": "りょ",
	"vu":     "ゔ",
	"sakuon": "っ",
}

var katakana = map[string]string{
	"a": "ア", "i": "イ", "u": "ウ", "e": "エ", "o": "オ",
	"ka": "カ", "ki": "キ", "ku": "ク", "ke": "ケ", "ko": "コ",
	"ga": "ガ", "gi": "ギ", "gu": "グ", "ge": "ゲ", "go": "ゴ",
	"sa": "サ", "shi": "シ", "su": "ス", "se": "セ", "so": "ソ",
	"za": "ザ", "ji": "ジ", "zu": "ズ", "ze": "ゼ", "zo": "ゾ",
	"ta": "タ", "chi": "チ", "tsu": "ツ", "te": "テ", "to": "ト",
	"da": "ダ", "de": "デ", "do": "ド",
	"na": "ナ", "ni": "ニ", "nu": "ヌ", "ne": "ネ", "no": "ノ",
	"ha": "ハ", "hi": "ヒ", "fu": "フ", "he": "ヘ", "ho": "ホ",
	"ba": "バ", "bi": "ビ", "bu": "ブ", "be": "ベ", "bo": "ボ",
	"pa": "パ", "pi": "ピ", "pu": "プ", "pe": "ペ", "po": "ポ",
	"ma": "マ", "mi": "ミ", "mu": "ム", "me": "メ", "mo": "モ",
	"ya": "ヤ", "yu": "ユ", "yo": "ヨ",
	"ra": "ラ", "ri": "リ", "ru": "ル", "re": "レ", "ro": "ロ",
	"wa": "ワ", "wo": "ヲ",
	"n":   "ン",
	"kya": "キャ", "kyu": "キュ", "kyo": "キョ",
	"gya": "ギャ", "gyu": "ギュ", "gyo": "ギョ",
	"sha": "シャ", "shu": "シュ", "sho": "ショ",
	"ja": "ジャ", "ju": "ジュ", "jo": "ジョ",
	"cha": "チャ", "chu": "チュ", "cho": "チョ",
	"nya": "ニャ", "nyu": "ニュ", "nyo": "ニョ",
	"hya": "ヒャ", "hyu": "ヒュ", "hyo": "ヒョ",
	"bya": "ビャ", "byu": "ビュ", "byo": "ビョ",
	"pya": "ピャ", "pyu": "ピュ", "pyo": "ピョ",
	"mya": "ミャ", "myu": "ミュ", "myo": "ミョ",
	"rya": "リャ", "ryu": "リュ", "ryo": "リョ",
	"vu": "ヴ",
	"va": "ヴァ", "vi": "ヴィ", "ve": "ヴェ", "vo": "ヴォ",
	"wi": "ウィ", "we": "ウェ",
	"fa": "ファ", "fi": "フィ", "fe": "フェ", "fo": "フォ",
	"che": "チェ",
	"di":  "ディ", "du": "ドゥ",
	"ti": "ティ", "tu": "トゥ",
	"je":     "ジェ",
	"she":    "シェ",
	"sakuon": "ッ",
	"pause":  "ー",
}

func RomajiToJapanese(romaji string) string {
	romaji = strings.ToLower(romaji)
	currentAlphabet := hiragana
	hiraganaIsCurrent := true
	resultStr := ""
	i := 0
	for i < len(romaji) {
		if string(romaji[i]) == "*" { // switch alphabets
			if hiraganaIsCurrent {
				currentAlphabet = katakana
				hiraganaIsCurrent = false
			} else {
				currentAlphabet = hiragana
				hiraganaIsCurrent = true
			}
			i++
		} else if string(romaji[i]) == " " { // check wa rule
			if i+3 < len(romaji) && romaji[i:i+4] == " wa " { // ha/wa rule
				resultStr += " " + currentAlphabet["ha"] + " "
				i += 4
				continue
			}
			resultStr += " "
			i++
		} else if i+2 < len(romaji) && string(romaji[i]) == "n" && string(romaji[i+1]) == "n" {
			_, exists := currentAlphabet[romaji[i+1:i+3]]
			if !exists {
				resultStr += currentAlphabet["sakuon"]
				i++
			}
		} else {
			checkLen := min(3, len(romaji)-i)
			for checkLen > 0 {
				checkStr := romaji[i : i+checkLen]
				if val, exists := currentAlphabet[checkStr]; exists {
					resultStr += val
					i += checkLen
					if i < len(romaji) {
						if string(romaji[i]) == "o" && string(romaji[i-1]) == "o" && hiraganaIsCurrent { // oo = ou rule
							resultStr += currentAlphabet["u"]
							i++
						} else if string(romaji[i]) == "e" && string(romaji[i-1]) == "e" && hiraganaIsCurrent { // ee = ei rule
							resultStr += currentAlphabet["i"]
							i++
						} else if string(romaji[i]) == string(romaji[i-1]) && !hiraganaIsCurrent {
							if string(romaji[i]) == "n" {
								break
							} else if strings.ContainsAny(string(romaji[i]), "aeiou") {
								resultStr += currentAlphabet["pause"]
							} else {
								resultStr += currentAlphabet["sakuon"]
							}
							i++
						}
					}
					break
				} else if checkLen == 1 {
					if strings.ContainsAny(string(checkStr), "?!.") { // punctuation
						resultStr += "。"
					} else if !strings.ContainsAny(string(checkStr), "abcdefghijklmnopqrstuvwxyz") { // print any characters that aren't a letter
						resultStr += checkStr
					} else if i+1 < len(romaji) && checkStr == string(romaji[i+1]) { // little tsu rule
						resultStr += currentAlphabet["sakuon"]
					}
					i++
					break
				}
				checkLen--
			}
		}
	}
	return resultStr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
