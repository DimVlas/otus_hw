package hw03frequencyanalysis

import (
	"regexp"
	"slices"
	"sort"
	"strings"
)

// Разбивает текст на слова, возвращает срез слов.
func splitWords(text string, pattern *regexp.Regexp) []string {
	if len(text) < 1 {
		return []string{}
	}

	if pattern == nil {
		return strings.Fields(text)
	} else {
		//res := pattern.Split(text, -1)
		return slices.DeleteFunc(pattern.Split(text, -1), func(s string) bool {
			return len(s) == 0
		})
	}
}

// Считает частоту элементов в строковом срезе.
// Возвращает slice структур wordWidth - слова с их весами.
// Slice отсортирован в зависимости от веса
// isNotCaseSens - true, не учитывать регистр
func wordsWidthsSort(s []string, isNotCaseSens bool) []wordWidth {
	if len(s) < 1 {
		return []wordWidth{}
	}

	m := make(map[string]wordWidth)

	for _, w := range s {
		word := w
		if isNotCaseSens {
			word = strings.ToUpper(word)
		}

		if width, ok := m[word]; ok {
			width.Width++
			m[word] = width
		} else {
			m[word] = wordWidth{Word: w, Width: 1}
		}
	}

	ww := make([]wordWidth, 0, len(m))

	for _, v := range m {
		ww = append(ww, v)
	}

	sort.Slice(ww, func(i, j int) bool {
		return ww[i].Width > ww[j].Width ||
			(ww[i].Width == ww[j].Width && ww[i].Word < ww[j].Word)
	})
	return ww
}

// Слово с его весом
type wordWidth struct {
	Word  string
	Width int
}

// func (w wordWidth) GetKey() string {
// 	return fmt.Sprintf("%04d%s", w.Width, w.Word)
// }

var pattern *regexp.Regexp = regexp.MustCompile(`[\s,.]+`) //"(?U)\\W+".

func Top10(text string) []string {
	if len(text) < 1 {
		return []string{}
	}

	words := splitWords(text, pattern)
	if len(words) < 1 {
		return []string{}
	}

	widths := wordsWidthsSort(words, false)

	res := make([]string, 0, 10)

	ln := len(widths)
	if ln > 10 {
		ln = 10
	}

	for _, ww := range widths[:ln] {
		res = append(res, ww.Word)
	}
	return res
}
