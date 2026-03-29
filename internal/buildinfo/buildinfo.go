package buildinfo

import "fmt"

var version = "dev"

func Summary() string {
	return fmt.Sprintf("mock %s", version)
}
