package file

import "github.com/WedgeNix/awsapi/types"

// Any represents any package-level type of AWS file.
type Any interface {
	__()
}

//D2sVendorDays maps Vendor to processed days for D2s.
type D2sVendorDays map[string]*types.D2sVendor

// D2sVendorDaysName hods file name
const D2sVendorDaysName = "drive-2-sku/vendor_days.json"

func (_ D2sVendorDays) __() {}

// BananasMon maps SKUs to their respective just-in-time data.
type BananasMon struct {
	AvgWait   float64
	OrdSKUCnt float64
	SKUs      types.SKUs
}

func (_ BananasMon) __() {}

// BananasCfg maps SKUs to their respective just-in-time data.
type BananasCfg types.BananasCfg

// BananasCfgName holds the file name.
const BananasCfgName = "hit-the-bananas/config.json"

func (_ BananasCfg) __() {}
