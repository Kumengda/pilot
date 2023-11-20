package filter

type FilterType int

const (
	Host FilterType = iota
	URL
	BODY
	HEADER
)
