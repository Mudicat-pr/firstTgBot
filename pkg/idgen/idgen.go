package idgen

import (
	"fmt"
	"strconv"
	"time"
)

type Number interface {
	int | int64
}

// num = userID
func IDgenerator[T Number](num T) int {
	n := int(num)
	t := time.Now()
	hms := t.Format("150405000")
	res, _ := strconv.Atoi(hms)
	uuidStr := fmt.Sprintf("%d", n%1000+res)

	if len(uuidStr) > 6 {
		uuidStr = uuidStr[len(uuidStr)-6:]
	}
	uuid, _ := strconv.Atoi(uuidStr)
	return uuid // user unique indetifier
}
