package config

import "errors"

// AddOption adds a new option and value to the configuration.
//
// If the section is nil then uses the section by default; if it does not exist,
// it is created in advance.
//
// It returns true if the option and value were inserted, and false if the value
// was overwritten.
func (cfg *Config) AddOption(section string, option string, value string) bool {
	cfg.AddSection(section) // Make sure section exists

	if section == "" {
		section = DefaultSection
	}

	_, ok := cfg.data[section][option]

	cfg.data[section][option] = &tValue{cfg.lastIDOption[section], value}
	cfg.lastIDOption[section]++

	return !ok
}

// RemoveOption removes a option and value from the configuration.
// It returns true if the option and value were removed, and false otherwise,
// including if the section did not exist.
func (cfg *Config) RemoveOption(section string, option string) bool {
	if _, ok := cfg.data[section]; !ok {
		return false
	}

	_, ok := cfg.data[section][option]
	delete(cfg.data[section], option)

	return ok
}

// HasOption checks if the configuration has the given option in the section.
// It returns false if either the option or section do not exist.
func (cfg *Config) HasOption(section string, option string) bool {
	if _, ok := cfg.data[section]; !ok {
		return false
	}

	_, okd := cfg.data[DefaultSection][option]
	_, oknd := cfg.data[section][option]

	return okd || oknd
}

// Options returns the list of options available in the given section.
// It returns an error if the section does not exist and an empty list if the
// section is empty. Options within the default section are also included.
func (cfg *Config) Options(section string) (options []string, err error) {
	if _, ok := cfg.data[section]; !ok {
		return nil, errors.New(sectionError(section).Error())
	}

	options = make([]string, len(cfg.data[DefaultSection])+len(cfg.data[section]))
	i := 0
	for s := range cfg.data[DefaultSection] {
		options[i] = s
		i++
	}
	for s := range cfg.data[section] {
		options[i] = s
		i++
	}

	return options, nil
}
