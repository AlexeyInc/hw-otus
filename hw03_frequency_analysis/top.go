package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const maxNumOfWords int = 10

type wordFrequency struct {
	Key   string
	Value int
}

func Top10(input string) (result []string) {
	splitedInput := strings.Fields(input)

	wordsFreq := calcWordsFreq(splitedInput)

	sortedWords := sortWordsByFreq(wordsFreq)

	for i := 0; i < len(sortedWords); i++ {
		if maxNumOfWords < i+1 {
			break
		}
		result = append(result, sortedWords[i].Key)
	}
	return result
}

func calcWordsFreq(words []string) (wordsFreq []wordFrequency) {
	wordsMap := make(map[string]int)
	for _, k := range words {
		wordsMap[k]++
	}
	for k, v := range wordsMap {
		wordsFreq = append(wordsFreq, wordFrequency{k, v})
	}
	return wordsFreq
}

func sortWordsByFreq(words []wordFrequency) []wordFrequency {
	sort.Slice(words, func(i, j int) bool {
		if words[i].Value == words[j].Value {
			frstWord := []rune(words[i].Key)
			secWord := []rune(words[j].Key)
			var shorterWord []rune
			if len(frstWord) < len(secWord) {
				shorterWord = frstWord
			} else {
				shorterWord = secWord
			}
			for i := 0; i < len(shorterWord); i++ {
				if frstWord[i] != secWord[i] {
					return frstWord[i] < secWord[i]
				}
			}
			return len(frstWord) < len(secWord)
		}
		return words[i].Value > words[j].Value
	})
	return words
}
