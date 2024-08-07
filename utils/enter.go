package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strings"
)

func InList(list []string, key string) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

func InListByRegex(list []string, key string) (ok bool) {
	for _, s := range list {
		regex, err := regexp.Compile(s)
		if err != nil {
			logx.Error(err)
			return
		}
		if regex.MatchString(key) {
			return true
		}
	}
	return false
}

func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func GetFilePrefix(fileName string) (prefix string) {
	nameList := strings.Split(fileName, ".")
	for i := 0; i < len(nameList)-1; i++ {
		if i == len(nameList)-2 {
			prefix += nameList[i]
			continue
		}
		prefix += nameList[i] + "."

	}
	return prefix
}

// DeduplicationList 去重
func DeduplicationList[T string | int | uint | uint32](req []T) (response []T) {
	Map := make(map[T]bool)
	for _, val := range req {
		if !Map[val] {
			Map[val] = true
		}
	}
	response = make([]T, 0)
	for key, _ := range Map {
		response = append(response, key)
	}
	return
}
