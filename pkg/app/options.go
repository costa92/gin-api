package app

import (
	cliflag "github.com/costa92/component-base/pkg/cli/flag"
)

type CliOptions interface {
	Flags() (fss cliflag.NamedFlagSets)
	Validate() []error
}

type ConfigurableOptions interface {
	ApplyFlags() []error
}

type CompletableOptions interface {
	Complete() error
}
type PrintableOptions interface {
	String() string
}
