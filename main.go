package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func convertBases(words []string) []string {
	for i := 0; i < len(words); i++ {
		if words[i] == "(hex)" || words[i] == "(bin)" {
			if i > 0 {
				base := 16
				if words[i] == "(bin)" {
					base = 2
				}
				n, err := strconv.ParseInt(words[i-1], base, 64)
				if err == nil {
					words[i-1] = strconv.FormatInt(n, 10)
				}
			}
			words = append(words[:i], words[i+1:]...)
			i--
		}
	}
	return words
}

func applyCase(word, mode string) string {
	switch mode {
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		runes := []rune(strings.ToLower(word))
		if len(runes) == 0 {
			return word
		}
		return strings.ToUpper(string(runes[0])) + string(runes[1:])
	}
	return word
}

func changeCase(words []string) []string {
	for i := 0; i < len(words); i++ {
		mode := ""
		count := 1
		toRemove := 0

		switch words[i] {
		case "(up)", "(low)", "(cap)":
			mode = words[i][1 : len(words[i])-1]
			toRemove = 1
		case "(up,", "(low,", "(cap,":
			if i+1 < len(words) && strings.HasSuffix(words[i+1], ")") {
				n, err := strconv.Atoi(strings.TrimSuffix(words[i+1], ")"))
				if err == nil {
					mode = words[i][1 : len(words[i])-1]
					count = n
					toRemove = 2
				}
			}
		}

		if toRemove == 0 {
			continue
		}

		start := i - count
		if start < 0 {
			start = 0
		}
		for j := start; j < i; j++ {
			words[j] = applyCase(words[j], mode)
		}

		words = append(words[:i], words[i+toRemove:]...)
		i--
	}
	return words
}

func isPunct(r rune) bool {
	return strings.ContainsRune(".,!?:;", r)
}

func fixPunctuation(text string) string {
	runes := []rune(text)
	var result []rune

	for i, r := range runes {
		if r == ' ' && i+1 < len(runes) && isPunct(runes[i+1]) {
			continue
		}

		result = append(result, r)

		if isPunct(r) && i+1 < len(runes) && !isPunct(runes[i+1]) && runes[i+1] != ' ' {
			result = append(result, ' ')
		}
	}
	return string(result)
}

func fixQuotes(text string) string {
	runes := []rune(text)
	var result []rune
	open := false

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r != '\'' {
			result = append(result, r)
			continue
		}

		if !open {
			result = append(result, r)
			for i+1 < len(runes) && runes[i+1] == ' ' {
				i++
			}
		} else {
			for len(result) > 0 && result[len(result)-1] == ' ' {
				result = result[:len(result)-1]
			}
			result = append(result, r)
		}
		open = !open
	}
	return string(result)
}

func fixArticles(text string) string {
	words := strings.Fields(text)

	for i := 0; i < len(words)-1; i++ {
		if words[i] != "a" && words[i] != "A" {
			continue
		}

		next := []rune(strings.ToLower(strings.TrimLeft(words[i+1], "'")))
		if len(next) == 0 {
			continue
		}

		if strings.ContainsRune("aeiouh", next[0]) {
			words[i] += "n"
		}
	}
	return strings.Join(words, " ")
}

func process(text string) string {
	words := strings.Fields(text)
	words = convertBases(words)
	words = changeCase(words)
	text = strings.Join(words, " ")
	text = fixPunctuation(text)
	text = fixQuotes(text)
	text = fixArticles(text)
	return text
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: go run . input.txt output.txt")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur de lecture :", err)
		os.Exit(1)
	}

	result := process(string(data))

	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur d'écriture :", err)
		os.Exit(1)
	}
}
