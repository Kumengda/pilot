package filter

import (
	"fmt"
	"regexp"
)

var DropAll = regexp.MustCompile("^$")

type Filter []*regexp.Regexp

func NewFilterWithContains(keyword ...string) Filter {
	var regexps []*regexp.Regexp
	for _, v := range keyword {
		regexps = append(regexps, regexp.MustCompile(fmt.Sprintf(".*%s.*", v)))
	}
	return regexps
}

func NewFilterWithEqual(keyword ...string) Filter {
	var regexps []*regexp.Regexp
	for _, v := range keyword {
		regexps = append(regexps, regexp.MustCompile(fmt.Sprintf("^%s$", v)))
	}
	return regexps
}

func NewFilterWithPrefix(keyword ...string) Filter {
	var regexps []*regexp.Regexp
	for _, v := range keyword {
		regexps = append(regexps, regexp.MustCompile(fmt.Sprintf("^%s.*", v)))
	}
	return regexps
}

func NewFilterWithSuffix(keyword ...string) Filter {
	var regexps []*regexp.Regexp
	for _, v := range keyword {
		regexps = append(regexps, regexp.MustCompile(fmt.Sprintf(".*%s$", v)))
	}
	return regexps
}

func NewFilterAllowAll() Filter {
	return []*regexp.Regexp{regexp.MustCompile("[\\s\\S]*")}
}

func NewFilterWithCustom(regexps []*regexp.Regexp) Filter {
	return regexps
}
