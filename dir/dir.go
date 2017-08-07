package dir

import (
	"github.com/WedgeNix/awsapi/file"
)

// Any represents any package-level type of AWS directory.
type Any interface {
	__()
}

// BananasMon is a mapping of file names to monitor files.
type BananasMon map[string]file.BananasMon

type Path string

// BananasMonName holds the filenames in the directory.
const BananasMonName = Path("hit-the-bananas/mon/*.json")

func (_ BananasMon) __() {}
