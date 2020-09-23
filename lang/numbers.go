package lang

import (
	"strconv"

	"github.com/divan/num2words"
)

func NumToWords(n int) (ret string, okay bool) {
	if s := num2words.Convert(int(n)); len(s) > 0 {
		ret, okay = s, true
	}
	return
}

// Currently good up to twenty.
// maybe https://github.com/donna-legal/word2number?
func WordsToNum(s string) (ret int, okay bool) {
	if cnt, e := strconv.Atoi(s); e == nil {
		ret, okay = cnt, true
	} else {
		smallNumbers := []string{
			"one",
			"two",
			"three",
			"four",
			"five",
			"six",
			"seven",
			"eight",
			"nine",
			"ten",
			"eleven",
			"twelve",
			"thirteen",
			"fourteen",
			"fifteen",
			"sixteen",
			"seventeen",
			"eighteen",
			"nineteen",
			"twenty",
		}
		for i, n := range smallNumbers {
			if s == n {
				ret, okay = i+1, true
				break
			}
		}
	}
	return
}
