package options

import (
	"github.com/Kumengda/pilot/filter"
)

type Filter struct {
	Filters    filter.Filter
	FilterType []filter.FilterType
}

func (f *Filter) GetValue() interface{} {
	return f
}

func SetFilter(filter filter.Filter, filterType ...filter.FilterType) *Filter {
	return &Filter{Filters: filter, FilterType: filterType}
}
