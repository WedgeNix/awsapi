package dir

import (
	"github.com/wedgenix/awsapi/file"
)

// Any represents any package-level type of AWS directory.
type Any interface {
	__()
}

type Monitor map[string]*file.Monitor

func (_ Monitor) __() {}
