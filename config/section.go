package config

// AddSection adds a new section to the configuration.
//
// If the section is nil then uses the section by default which it's already
// created.
//
// It returns true if the new section was inserted, and false if the section
// already existed.
func (cfg *Config) AddSection(section string) bool {
	// DefaultSection
	if section == "" {
		return false
	}

	if _, ok := cfg.data[section]; ok {
		return false
	}

	cfg.data[section] = make(map[string]*tValue)

	// Section order
	cfg.idSection[section] = cfg.lastIDSection
	cfg.lastIDSection++

	return true
}

// RemoveSection removes a section from the configuration.
// It returns true if the section was removed, and false if section did not exist.
func (cfg *Config) RemoveSection(section string) bool {
	_, ok := cfg.data[section]

	// Default section cannot be removed.
	if !ok || section == DefaultSection {
		return false
	}

	for o := range cfg.data[section] {
		delete(cfg.data[section], o) // *value
	}
	delete(cfg.data, section)

	delete(cfg.lastIDOption, section)
	delete(cfg.idSection, section)

	return true
}

// HasSection checks if the configuration has the given section.
// (The default section always exists.)
func (cfg *Config) HasSection(section string) bool {
	_, ok := cfg.data[section]
	return ok
}

// Sections returns the list of sections in the configuration.
// (The default section always exists.)
func (cfg *Config) Sections() (sections []string) {
	sections = make([]string, len(cfg.idSection))
	pos := 0 // Position in sections

	for i := 0; i < cfg.lastIDSection; i++ {
		for section, id := range cfg.idSection {
			if id == i {
				sections[pos] = section
				pos++
			}
		}
	}

	return sections
}
