package function

import (
	"crypto/md5"
	"fmt"
	"github.com/segmentio/ksuid"
	"math/rand"
)

//GenUUIDMd5 生成 store_user user_id
func GenUUIDMd5() string {
	return fmt.Sprintf("%x",md5.Sum([]byte(GenUUID())))
}

//GenUUID 生成唯一字符串
func GenUUID() string {
	id := ksuid.New()
	return id.String()
}

func CutString(str string, start, len int) string  {
	return (str[start:])[0:len]
}

//RandRange 生成范围内随机数
func RandRange(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}