package file

import "github.com/WedgeNix/awsapi/types"

// Any represents any package-level type of AWS file.
type Any interface {
	__()
}

// Monitor maps SKUs to their respective just-in-time data.
type Monitor map[string]types.MonitorSKU

// allows Monitor to be bound to the Any interface
func (_ Monitor) __() {}
