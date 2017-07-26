package file

import "github.com/WedgeNix/awsapi/types"

// Any represents any package-level type of AWS file.
type Any interface {
	__()
}

// BananasMon maps SKUs to their respective just-in-time data.
type BananasMon map[string]types.BananasMonSKU

// allows Monitor to be bound to the Any interface
func (_ BananasMon) __() {}

// BananasCfg maps SKUs to their respective just-in-time data.
type BananasCfg types.BananasCfg

// allows Monitor to be bound to the Any interface
func (_ BananasCfg) __() {}
