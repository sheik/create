package git

import (
	"fmt"
	"strconv"
	"strings"
)

func IncrementMinorVersion(version string) string {
	parts := strings.Split(strings.Split(version, "_")[0], ".")
	if len(parts) != 3 {
		return version
	}
	if minorVersion, err := strconv.Atoi(parts[2]); err == nil {
		minorVersion += 1
		return fmt.Sprintf("%s.%s.%d", parts[0], parts[1], minorVersion)
	} else {
		return version
	}
}
