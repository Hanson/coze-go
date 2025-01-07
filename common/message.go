package common

type Message struct {
	Role        string            `json:"role"`
	Type        string            `json:"type"`
	Content     string            `json:"content"`
	ContentType string            `json:"content_type"`
	MetaData    map[string]string `json:"meta_data"`
}
