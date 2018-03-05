package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// WriteFile saves the configuration representation to a file.
// The desired file permissions must be passed as in os.Open. The header is a
// string that is saved as a comment in the first line of the file.
func (cfg *Config) WriteFile(fname string, perm os.FileMode, header string) error {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	buf := bufio.NewWriter(file)
	if err = cfg.write(buf, header); err != nil {
		return err
	}
	buf.Flush()

	return file.Close()
}

func (cfg *Config) write(buf *bufio.Writer, header string) (err error) {
	if header != "" {
		// Add comment character after of each new line.
		if i := strings.Index(header, "\n"); i != -1 {
			header = strings.Replace(header, "\n", "\n"+cfg.comment, -1)
		}

		if _, err = buf.WriteString(cfg.comment + header + "\n"); err != nil {
			return err
		}
	}

	for _, orderedSection := range cfg.Sections() {
		for section, sectionMap := range cfg.data {
			if section == orderedSection {

				// Skip default section if empty.
				if section == DefaultSection && len(sectionMap) == 0 {
					continue
				}

				if _, err = buf.WriteString("\n[" + section + "]\n"); err != nil {
					return err
				}

				// Follow the input order in options.
				for i := 0; i < cfg.lastIDOption[section]; i++ {
					for option, tValue := range sectionMap {

						if tValue.position == i {
							if _, err = buf.WriteString(fmt.Sprint(
								option, cfg.separator, tValue.v, "\n")); err != nil {
								return err
							}
							cfg.RemoveOption(section, option)
							break
						}
					}
				}
			}
		}
	}

	if _, err = buf.WriteString("\n"); err != nil {
		return err
	}

	return nil
}
