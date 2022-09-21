package cmd

import (
	"strings"
	"testing"
)

func TestWikiLinks(t *testing.T) {
	links := getWikiLinks()
	seen := make(map[string]struct{})
	for _, link := range links {
		if _, ok := seen[link]; ok {
			t.Errorf("fail found duplicate wiki link: %s", link)
		} else {
			seen[link] = struct{}{}
		}
	}
}

func TestThousandWords(t *testing.T) {
	commonWords := readWordsFile("../words/common_words.txt")
	// for i, word := range commonWords {
	// 	println(i, word)
	// }
	thousandWords(t, commonWords[:])
}

func TestCommonWords(t *testing.T) {
	commonWordsStr := CommonWords("../words/common_words.txt")
	commonWords := strings.Split(commonWordsStr, " ")

	thousandWords(t, commonWords)
}

func thousandWords(t *testing.T, words []string) {
	for i, word := range words {
		if word == "" {
			t.Errorf("common word index %d is empty", i)
		}
		for _, char := range word {
			if char == ' ' {
				t.Errorf("common word %s contains a space character", word)
			}
		}
	}

	if len(words) != 1000 {
		t.Errorf("common words not 1000 in length (%d)", len(words))
	}

	seen := make(map[string]struct{})
	for _, word := range words {
		if _, ok := seen[word]; ok {
			t.Errorf("common words contains duplicate: %s", word)
		} else {
			seen[word] = struct{}{}
		}
	}
}

func TestShuffleWords(t *testing.T) {
	words := readWordsFile("../words/common_words.txt")
	shuffledWords := shuffleWords(words)
	if words == shuffledWords {
		t.Errorf("shuffled words equal original words")
	}

}
