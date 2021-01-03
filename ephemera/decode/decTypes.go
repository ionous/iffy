package decode

import "github.com/ionous/iffy/dl/composer"

type SwapType interface {
	composer.Composer
	Choices() (nameToType map[string]interface{})
}

type StrType interface {
	composer.Composer
	Choices() []string
}

type NumType interface {
	composer.Composer
	Choices() []float64
}

// translate a $TOKEN to a label, return its index if found.
// return -1 when not found and if the set is open return the original choice,
// otherwise return an empty string.
func FindChoice(op StrType, choice string) (ret string, found int) {
	spec := op.Compose()
	keys := op.Choices()
	found = -1
	for i, key := range keys {
		if key == choice {
			found = i
			ret = spec.Strings[i] // panic if out of range; the two arrays should match
			break
		}
	}
	if found < 0 && spec.OpenStrings {
		ret = choice
	}
	return
}
