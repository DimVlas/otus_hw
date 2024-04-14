package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

// Разбивает текст на слова, возвращает срез слов.
func GetWords(text string) []string {
	if len(text) < 1 {
		return []string{}
	}

	return strings.Fields(text)
}

// Считает частоту элементов в строковом срезе.
// Возвращает map, где key элемент среза,
// а value сколько раз встречался элемент.
func CountOfElements(s []string) map[string]int {
	m := make(map[string]int)

	for _, w := range s {
		if cnt, ok := m[w]; ok {
			cnt++
			m[w] = cnt
		} else {
			m[w] = 1
		}
	}
	return m
}

type wordWidth struct {
	Word  string
	Width int
}

func (w wordWidth) GetKey() string {
	return fmt.Sprintf("%04d%s", w.Width, w.Word)
}

func Top10(text string) []string {
	if len(text) < 1 {
		return []string{}
	}

	words := GetWords(text)
	if len(words) < 1 {
		return []string{}
	}

	countMap := CountOfElements(words)
	countWords := make([]wordWidth, 0, len(countMap))

	for k, v := range countMap {
		countWords = append(countWords, wordWidth{Word: k, Width: v})
	}

	sort.Slice(countWords, func(i, j int) bool {
		return countWords[i].Width > countWords[j].Width ||
			(countWords[i].Width == countWords[j].Width && countWords[i].Word < countWords[j].Word)
	})

	res := make([]string, 0, 10)

	ln := len(countWords)
	if ln > 10 {
		ln = 10
	}

	for _, ww := range countWords[:ln] {
		res = append(res, ww.Word)
	}
	return res
}
