package stdlog

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type StdLog struct {
	srvName string
}

func New(serviceName string, serviceNameMaxLen int, c color.Attribute) *StdLog {
	srvName := serviceName + strings.Repeat(" ", serviceNameMaxLen-len(serviceName)) + "  |"

	sl := &StdLog{
		srvName: color.New(c).SprintFunc()(srvName),
	}

	return sl
}

func (sl *StdLog) Write(data []byte) (int, error) {
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}
		fmt.Printf("%s %s\n", sl.srvName, line)
	}

	return len(data), nil
}
