package common

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func GetRandomRouteId() string {

	now := time.Now().UnixNano()
	r1 := RandInt(10001, 99999)
	strb1 := strconv.FormatInt(now, 10)
	strb2 := strconv.Itoa(r1)
	str := fmt.Sprintf("%s%s", strb1, strb2)

	log.Printf("random :s", str)

	return str

}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min+1) + min
}
