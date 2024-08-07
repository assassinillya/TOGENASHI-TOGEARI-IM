package utils

import (
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	md := MD5([]byte("1234"))
	fmt.Println(md)
}

func TestDeduplicationList(t *testing.T) {
	fmt.Println(DeduplicationList([]string{"1", "1", "2", "4", "2", "3", "4", "5", "6", "7", "8", "9", "10"}))
}
