package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	// Дополнительное задание: не учитывать регистр букв и знаки препинания по краям слова.
	isAsteriksTask = true
	patt           = `[a-zA-Zа-яА-Я0-9]+[\.\-\,]*[a-zA-Zа-яА-Я0-9]+|\-{2,}|[a-zA-Zа-яА-Я0-9]+`
	pattern        = regexp.MustCompile(patt)
)

// Разбивает текст на слова, возвращает срез слов.
func splitWords(text string, pattern *regexp.Regexp) []string {
	if len(text) < 1 {
		return []string{}
	}

	if pattern == nil {
		return strings.Fields(text)
	}

	return pattern.FindAllString(text, -1)
}

// Считает частоту элементов в строковом срезе.
// Возвращает slice структур wordWidth - слова с их весами.
// Slice отсортирован в зависимости от веса.
// isIgnoreCase - true, не учитывать регистр.
func wordsWidthsSort(s []string, isIgnoreCase bool) []wordWidth {
	if len(s) < 1 {
		return []wordWidth{}
	}

	m := make(map[string]wordWidth)

	for _, w := range s {
		word := w
		if isIgnoreCase {
			word = strings.ToLower(word)
		}

		if width, ok := m[word]; ok {
			width.Width++
			m[word] = width
		} else {
			m[word] = wordWidth{Word: word, Width: 1}
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

// Слово с его весом.
type wordWidth struct {
	Word  string
	Width int
}

func Top10(text string) []string {
	if len(text) < 1 {
		return []string{}
	}

	var words []string
	if isAsteriksTask {
		words = splitWords(text, pattern)
	} else {
		words = splitWords(text, nil)
	}
	if len(words) < 1 {
		return []string{}
	}

	widths := wordsWidthsSort(words, isAsteriksTask)

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
