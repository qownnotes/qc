package snippet

import (
	"bytes"
	"sort"
	"strings"

	"github.com/qownnotes/qc/config"
	"github.com/qownnotes/qc/entity"
	"github.com/qownnotes/qc/websocket"
)

type Snippets struct {
	Snippets []entity.SnippetInfo
}

// Load reads snippets.
func (snippets *Snippets) Load() error {
	snippets.Snippets = websocket.FetchSnippetsData()
	snippets.Order()

	return nil
}

//// Save saves the snippets to toml file.
//func (snippets *Snippets) Save() error {
//	snippetFile := config.Conf.General.SnippetFile
//	f, err := os.Create(snippetFile)
//	defer f.Close()
//	if err != nil {
//		return fmt.Errorf("Failed to save snippet file. err: %s", err)
//	}
//	return toml.NewEncoder(f).Encode(snippets)
//}

// ToString returns the contents of toml file.
func (snippets *Snippets) ToString() (string, error) {
	var buffer bytes.Buffer
	//err := toml.NewEncoder(&buffer).Encode(snippets)
	//if err != nil {
	//	return "", fmt.Errorf("Failed to convert struct to TOML string: %v", err)
	//}
	return buffer.String(), nil
}

// Order snippets regarding SortBy option defined in config toml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (snippets *Snippets) Order() {
	sortBy := config.Conf.General.SortBy

	switch sortBy {
	case "command", "+command":
		sort.Sort(ByCommand(snippets.Snippets))
	case "-command":
		sort.Sort(sort.Reverse(ByCommand(snippets.Snippets)))

	case "description", "+description":
		sort.Sort(ByDescription(snippets.Snippets))
	case "-description":
		sort.Sort(sort.Reverse(ByDescription(snippets.Snippets)))

	case "output", "+output":
		sort.Sort(ByOutput(snippets.Snippets))
	case "-output":
		sort.Sort(sort.Reverse(ByOutput(snippets.Snippets)))

	case "-recency":
		snippets.reverse()
	}
}

func (snippets *Snippets) reverse() {
	for i, j := 0, len(snippets.Snippets)-1; i < j; i, j = i+1, j-1 {
		snippets.Snippets[i], snippets.Snippets[j] = snippets.Snippets[j], snippets.Snippets[i]
	}
}

type ByCommand []entity.SnippetInfo

func (a ByCommand) Len() int      { return len(a) }
func (a ByCommand) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool {
	return strings.ToLower(a[i].Command) > strings.ToLower(a[j].Command)
}

type ByDescription []entity.SnippetInfo

func (a ByDescription) Len() int      { return len(a) }
func (a ByDescription) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool {
	return strings.ToLower(a[i].Description) > strings.ToLower(a[j].Description)
}

type ByOutput []entity.SnippetInfo

func (a ByOutput) Len() int           { return len(a) }
func (a ByOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOutput) Less(i, j int) bool { return a[i].Output > a[j].Output }
