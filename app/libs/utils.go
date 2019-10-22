/**
  create by yy on 2019-10-10
*/

package libs

import (
	"errors"
	"fmt"
	"gin_template/app/config"
	"os"
	"runtime"
	"strings"
)

func GetErrorString(err error) string {
	return fmt.Sprintf("error: %v", err)
}

func NewReportError(err error) error {
	if !config.Config.App.DEBUG {
		return err
	}
	_, fileName, line, _ := runtime.Caller(1)
	data := fmt.Sprintf("%v, report in: %v: in line %v", err, fileName, line)
	return errors.New(data)
}

func DebugPrint(format string, values ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	_, _ = fmt.Fprintf(os.Stderr, "[Guaidashu-debug] "+format, values...)
}
