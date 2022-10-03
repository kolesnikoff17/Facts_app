package common

type Fact struct {
	Title string   `json:"title"`
	Desc  string   `json:"description"`
	Links []string `json:"links,omitempty"`
}
