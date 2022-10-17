package schema

type Book struct {
	BookID int64  `json:"bookID"`
	Name    string `json:"name"`
	Author   string  `json:"author"`
	Publisher string `json:"publisher"`
}

type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}