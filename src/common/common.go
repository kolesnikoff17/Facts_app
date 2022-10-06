package common

type Fact struct {
	Id    int      `json:"id"`
	Title string   `json:"title"`
	Desc  string   `json:"description"`
	Links []string `json:"links,omitempty"`
}

type FactsArr struct {
	Facts []Fact `json:"facts"`
}
