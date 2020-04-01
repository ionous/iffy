package qna

type NumberList struct {
	i    int
	list []float64
}

type TextList struct {
	i    int
	list []string
}

type ObjectList struct {
	i    int
	list []string
}

func NewNumberList(list []float64) *NumberList {
	return &NumberList{list: list}
}

func (it *NumberList) HasNext() bool {
	return it.i < len(it.list)
}

func (it *NumberList) GetNumber() (ret float64, err error) {
	err = Assign(&ret, it.list[it.i])
	it.i++
	return
}

func NewTextList(list []string) *TextList {
	return &TextList{list: list}
}

func (it *TextList) HasNext() bool {
	return it.i < len(it.list)
}

func (it *TextList) GetText() (ret string, err error) {
	err = Assign(&ret, it.list[it.i])
	it.i++
	return
}
