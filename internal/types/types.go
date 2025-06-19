package types

type Source struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Link string `json:"link"`
}

type Sources struct {
	Sources []Source `json:"sources"`
}
