package config

import (
	"regexp"
	"strings"
)

const (
	// DefaultSection Default section name.
	DefaultSection = "DEFAULT"
	// DepthValues Maximum allowed depth when recursively substituing variable names.
	DepthValues = 200

	// DefaultComment 默认评论
	DefaultComment = "# "
	// AlternativeComment 替代评论
	AlternativeComment = "; "
	// DefaultSeparator 默认间隔
	DefaultSeparator = ":"
	// AlternativeSeparator 替代间隔
	AlternativeSeparator = "="
)

var (
	// Strings accepted as boolean.
	boolString = map[string]bool{
		"t":     true,
		"true":  true,
		"y":     true,
		"yes":   true,
		"on":    true,
		"1":     true,
		"f":     false,
		"false": false,
		"n":     false,
		"no":    false,
		"off":   false,
		"0":     false,
	}

	varRegExp = regexp.MustCompile(`%\(([a-zA-Z0-9_.\-]+)\)s`) // %(variable)s
)

// Config is the representation of configuration settings.
type Config struct {
	comment   string
	separator string

	// === Sections order
	lastIDSection int            // Last section identifier
	idSection     map[string]int // Section : position

	// The last option identifier used for each section.
	lastIDOption map[string]int // Section : last identifier

	// Section -> option : value
	data map[string]map[string]*tValue
}

// Hold the input position for a value.
type tValue struct {
	position int    // Option order
	v        string // value
}

// New creates an empty configuration representation.
// This representation can be filled with AddSection and AddOption and then
// saved to a file using WriteFile.
//
// === Arguments
//
// comment: has to be `DefaultComment` or `AlternativeComment`
// separator: has to be `DefaultSeparator` or `AlternativeSeparator`
// preSpace: indicate if is inserted a space before of the separator
// postSpace: indicate if is added a space after of the separator
func New(comment, separator string, preSpace, postSpace bool) *Config {
	if comment != DefaultComment && comment != AlternativeComment {
		panic("comment character not valid")
	}

	if separator != DefaultSeparator && separator != AlternativeSeparator {
		panic("separator character not valid")
	}

	// === Get spaces around separator
	if preSpace {
		separator = " " + separator
	}

	if postSpace {
		separator += " "
	}
	// ===

	c := new(Config)

	c.comment = comment
	c.separator = separator
	c.idSection = make(map[string]int)
	c.lastIDOption = make(map[string]int)
	c.data = make(map[string]map[string]*tValue)

	c.AddSection(DefaultSection) // Default section always exists.

	return c
}

// NewDefault creates a configuration representation with values by default.
func NewDefault() *Config {
	return New(DefaultComment, DefaultSeparator, false, true)
}

func stripComments(l string) string {
	// Comments are preceded by space or TAB
	for _, c := range []string{" ;", "\t;", " #", "\t#"} {
		if i := strings.Index(l, c); i != -1 {
			l = l[0:i]
		}
	}
	return l
}
