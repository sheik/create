package build

import "fmt"

func GoFlags(version, buildTime string) string {
	return fmt.Sprintf("-ldflags '-X main.Version=%s -X main.BuildTime=%s'", version, buildTime)
}

func RPM(project, version string) string {
	return fmt.Sprintf("%s-%s-1.x86_64.rpm", project, version)
}
