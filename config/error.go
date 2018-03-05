package config

type sectionError string

// Error 错误
func (err sectionError) Error() string {
	return "section not found: " + string(err)
}

type optionError string

// Error 错误
func (err optionError) Error() string {
	return "option not found: " + string(err)
}
