package sorter

import (
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type fileMetadata struct {
	NameExt string
	Name    string
	Ext     string

	StrTokens   []string
	NumTokens   []int
	TokenIsNums []bool
}

func split(text string) []string {
	i := 0
	r := 0
	l := len(text)
	if l <= 0 {
		return []string{}
	}
	currentIsDigit := unicode.IsDigit(rune(text[0]))
	var tokens []string

	for r < l {
		c := rune(text[r])
		if unicode.IsDigit(c) {
			if !currentIsDigit {
				tokens = append(tokens, text[i:r])
				currentIsDigit = true
				i = r
			}
		} else {
			if currentIsDigit {
				tokens = append(tokens, text[i:r])
				currentIsDigit = false
				i = r
			}
		}
		r++
	}
	if i < l {
		tokens = append(tokens, text[i:l])
	}

	return tokens
}

func toFileMetadata(file fs.FileInfo) *fileMetadata {
	_ext := filepath.Ext(file.Name())
	_name := strings.TrimSuffix(file.Name(), _ext)

	_ext = strings.ToLower(_ext)
	_name = strings.ToLower(_name)

	tokens := split(_name)
	strTokens := make([]string, len(tokens))
	numTokens := make([]int, len(tokens))
	tokenIsNums := make([]bool, len(tokens))
	for i, token := range tokens {
		if numeric, err := strconv.Atoi(token); err == nil {
			numTokens[i] = numeric
			tokenIsNums[i] = true
		} else {
			strTokens[i] = token
			tokenIsNums[i] = false
		}
	}

	return &fileMetadata{
		NameExt: file.Name(),
		Name:    _name,
		Ext:     _ext,

		StrTokens:   strTokens,
		NumTokens:   numTokens,
		TokenIsNums: tokenIsNums,
	}
}
