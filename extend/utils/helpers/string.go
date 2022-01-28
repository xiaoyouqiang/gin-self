package helpers

import (
	"crypto/md5"
	"fmt"
	"gin-self/extend/self_db"
	"gin-self/model/mysql/yema_store_users_model"
	"github.com/segmentio/ksuid"
	"math/rand"
	"strconv"
	"strings"
	"time"
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

//CreateInviteCodeForUser10020 生成 User10020 邀请码
func CreateInviteCodeForUser10020(id int64) string {
	if id == 0 {
		return ""
	}

	str := fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(int(id))))) +
		"_" +
		strconv.Itoa(int(time.Now().Unix())) +
		strconv.Itoa(int(RandRange(1000, 9999)))

	str = strings.NewReplacer("/","","+","").Replace(str)

	str = strings.ToUpper(CutString(str,0, 6))

	query := yema_store_users_model.NewQueryBuilder()
	user, _ := query.WhereOp("invite_code", self_db.Equal, str).First()
	if user.Id != 0 {
		//重复的 继续尝试 生成
		return CreateInviteCodeForUser10020(id)
	}

	return str
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