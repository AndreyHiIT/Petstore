package models

type Pet struct {
	ID        int      `json:"id" db:"id"`
	Category  Category `json:"category" db:"category"`
	Name      string   `json:"name" db:"name"`
	PhotoUrls []string `json:"photourls" db:"photourls"`
	Tags      []Tag    `json:"tags" db:"tags"`
	Status    string   `json:"status" db:"status"`
}

type Category struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Tag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
