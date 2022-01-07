package common

import (
	"fmt"
	"regexp"
	"time"
)

const layoutDateMeta = "20060102"

func GenerateFileMeta(now time.Time, fileName string) (fileMeta string) {
	re := regexp.MustCompile(`\.(pdf|PDF)$`)
	filetype := re.FindStringSubmatch(fileName)
	if len(filetype) <= 0 {
		filetype = append(filetype, "")
	}
	fileMeta = fmt.Sprintf("%s-%s%s%s", now.Format(layoutDateMeta), RandStringRunes(7), RandStringRunes(7), filetype[0])
	return fileMeta
}
