package types

// Currency
type Currency struct {
	MetaData map[string]interface{} `json:"metadata" bson:"-"`
}

// Model
type Model struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
