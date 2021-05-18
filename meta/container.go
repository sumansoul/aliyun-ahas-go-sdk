package meta

import (
	"bufio"
	"os"
	"strings"

	"github.com/alibaba/sentinel-golang/util"
)

const (
	DockerEnvFile      = "/.dockerenv"
	ProcSelfCgroupFile = "/proc/1/cgroup"
)

func IsInContainer() bool {
	// Check whether current process is in container.
	cgroupInfoFirstLine, err := readSingleLineOfFile(ProcSelfCgroupFile)
	if err == nil && len(cgroupInfoFirstLine) == 0 {
		return strings.Contains(cgroupInfoFirstLine, "docker") || strings.Contains(cgroupInfoFirstLine, "kube")
	}
	// Check docker env file
	exists, err := util.FileExists(DockerEnvFile)
	return exists
}

func readSingleLineOfFile(filename string) (string, error) {
	f, err := os.Open(ProcSelfCgroupFile)
	if err != nil {
		return "", err
	}
	br := bufio.NewReader(f)
	bytes, _, err := br.ReadLine()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
