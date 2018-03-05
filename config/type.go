package config

import (
	"errors"
	"strconv"
	"strings"
)

// Bool has the same behaviour as String but converts the response to bool.
// See "boolString" for string values converted to bool.
func (cfg *Config) Bool(section string, option string) (value bool, err error) {
	sv, err := cfg.String(section, option)
	if err != nil {
		return false, err
	}

	value, ok := boolString[strings.ToLower(sv)]
	if !ok {
		return false, errors.New("could not parse bool value: " + sv)
	}

	return value, nil
}

// Float has the same behaviour as String but converts the response to float.
func (cfg *Config) Float(section string, option string) (value float64, err error) {
	sv, err := cfg.String(section, option)
	if err == nil {
		value, err = strconv.ParseFloat(sv, 64)
	}

	return value, err
}

// Int has the same behaviour as String but converts the response to int.
func (cfg *Config) Int(section string, option string) (value int, err error) {
	sv, err := cfg.String(section, option)
	if err == nil {
		value, err = strconv.Atoi(sv)
	}

	return value, err
}

// RawString gets the (raw) string value for the given option in the section.
// The raw string value is not subjected to unfolding, which was illustrated in
// the beginning of this documentation.
//
// It returns an error if either the section or the option do not exist.
func (cfg *Config) RawString(section string, option string) (value string, err error) {
	if _, ok := cfg.data[section]; ok {
		if tValue, ok := cfg.data[section][option]; ok {
			return tValue.v, nil
		}
		return "", errors.New(optionError(option).Error())
	}
	return "", errors.New(sectionError(section).Error())
}

// String gets the string value for the given option in the section.
// If the value needs to be unfolded (see e.g. %(host)s example in the beginning
// of this documentation), then String does this unfolding automatically, up to
// DepthValues number of iterations.
//
// It returns an error if either the section or the option do not exist, or the
// unfolding cycled.
func (cfg *Config) String(section string, option string) (value string, err error) {
	value, err = cfg.RawString(section, option)
	if err != nil {
		return "", err
	}

	var i int

	for i = 0; i < DepthValues; i++ { // keep a sane depth
		vr := varRegExp.FindString(value)
		if len(vr) == 0 {
			break
		}

		// Take off leading '%(' and trailing ')s'
		noption := strings.TrimLeft(vr, "%(")
		noption = strings.TrimRight(noption, ")s")

		// Search variable in default section
		nvalue, _ := cfg.data[DefaultSection][noption]
		if _, ok := cfg.data[section][noption]; ok {
			nvalue = cfg.data[section][noption]
		}
		if nvalue.v == "" {
			return "", errors.New(optionError(noption).Error())
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		value = strings.Replace(value, vr, nvalue.v, -1)
	}

	if i == DepthValues {
		return "", errors.New("possible cycle while unfolding variables: " +
			"max depth of " + strconv.Itoa(DepthValues) + " reached")
	}

	return value, nil
}
