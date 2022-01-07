package common

import (
	"fmt"
	"regexp"
)

func GetFileNameFromPath(path string) (fileName string, err error) {
	re := regexp.MustCompile(`[^\\/:*?"<>|\r\n]+$`)
	gettFileName := re.FindStringSubmatch(path)
	if len(gettFileName) <= 0 {
		return "", fmt.Errorf("can't get file name from path object `%s`", path)
	}
	fileName = gettFileName[0]
	return fileName, err
}
