package entity

type SnippetInfo struct {
	Description string   `json:"description"`
	Command     string   `json:"command"`
	Tag         []string `json:"tags"`
	Output      string   `json:"output"`
}

type NoteFolderInfo struct {
	Name string `json:"text"`
	Id   int    `json:"value"`
}
