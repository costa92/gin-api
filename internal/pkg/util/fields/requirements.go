package fields

import "github.com/costa92/go-web/internal/pkg/util/selection"

// Requirements are AND of all requirements.
type Requirements []Requirement

type Requirement struct {
	Operator selection.Operator
	Field    string
	Value    string
}
