package decode

import "github.com/ionous/iffy/dl/composer"

type SwapType interface {
	composer.Composer
	Choices() (nameToType map[string]interface{})
}
type StrType interface {
	composer.Composer
	Choices() (tokenToValue map[string]string)
}

type NumType interface {
	composer.Composer
	Choices() []float64
}

// translate a choice, typically a $TOKEN, to a value.
// note: go-code doesnt currently have a way to find a string's label.
func FindChoice(op StrType, choice string) (ret string, found bool) {
	spec, keys := op.Compose(), op.Choices()
	if str, ok := keys[choice]; ok {
		ret, found = str, ok
	} else if !ok && spec.OpenStrings {
		ret = choice
	}
	return
}
