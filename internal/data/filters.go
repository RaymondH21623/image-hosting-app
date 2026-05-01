package data

import (
	"shareapp/internal/validator"
	"strings"
)

type Filters struct {
	Page             int
	PageSize         int
	Sort             string
	SortSafelist     []string
	Filetype         string
	FiletypeSafelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be a positive integer")
	v.Check(f.Page <= 1000000, "page", "must be a maximum of 1 million")
	v.Check(f.PageSize > 0, "pageSize", "must be a positive integer")
	v.Check(f.PageSize <= 100, "pageSize", "must be a maximum of 100")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
	v.Check(validator.PermittedValue(f.Filetype, f.FiletypeSafelist...), "filetype", "invalid filetype value")
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}
