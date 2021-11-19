package config

// Flag is global flag variable
var Flag FlagConfig

// FlagConfig is a struct of flag
type FlagConfig struct {
	Debug     bool
	Port      int
	Query     string
	Command   bool
	FilterTag string
	Color     bool
}
