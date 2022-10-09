package common

// Fact struct for marshal/unmarshal json data and work with db.
type Fact struct {
	ID    int      `json:"id"`
	Title string   `json:"title"`
	Desc  string   `json:"description"`
	Links []string `json:"links,omitempty"`
}

// FactsArr is a set of facts needed by POST method.
type FactsArr struct {
	Facts []Fact `json:"facts"`
}
