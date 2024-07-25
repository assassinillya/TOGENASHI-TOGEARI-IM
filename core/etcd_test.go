package core

import (
	"context"
	"log"
	"testing"
)

func TestEtcd(t *testing.T) {
	client := InitEtcd("127.0.0.1:2379")
	res, err := client.Put(context.Background(), "auth_api", "127.0.0.1:20021")
	log.Println(res, err)
	getResponse, err := client.Get(context.Background(), "auth_api")
	log.Println(getResponse, err)
	if err == nil && len(getResponse.Kvs) > 0 {
		value := getResponse.Kvs[0].Value
		log.Println(string(value))
	}
}
