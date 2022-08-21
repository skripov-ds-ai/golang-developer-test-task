package redclient

import (
	"context"
	"fmt"
	"golang-developer-test-task/structs"

	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/go-redis/redismock/v8"
	"github.com/mailru/easyjson"
)

func TestAddValue(t *testing.T) {
	info := structs.Info{
		GlobalID:       42,
		SystemObjectID: "777",
		ID:             1,
		IDEn:           9,
		Mode:           "abc",
		ModeEn:         "cba",
	}
	bs, _ := easyjson.Marshal(info)

	db, mock := redismock.NewClientMock()
	mock.ExpectSet(info.SystemObjectID, bs, 0).SetVal("OK")
	mock.ExpectSet(fmt.Sprintf("global_id:%d", info.GlobalID), info.SystemObjectID, 0).SetVal("OK")
	mock.ExpectSet(fmt.Sprintf("id:%d", info.ID), info.SystemObjectID, 0).SetVal("OK")
	mock.ExpectSet(fmt.Sprintf("id_en:%d", info.IDEn), info.SystemObjectID, 0).SetVal("OK")
	mock.ExpectRPush(fmt.Sprintf("mode:%s", info.Mode), info.SystemObjectID).SetVal(0)
	mock.ExpectRPush(fmt.Sprintf("mode_en:%s", info.ModeEn), info.SystemObjectID).SetVal(0)

	client := &RedisClient{*db}
	err := client.AddValue(context.Background(), info)

	if err != nil {
		t.Fatal(err)
	}
}

func TestFindValuesNotFoundSingle(t *testing.T) {
	db, mock := redismock.NewClientMock()
	key := "42"
	mock.ExpectGet(key).SetErr(redis.Nil)
	client := &RedisClient{*db}
	_, _, err := client.FindValues(context.Background(), key, false, 5, 0)

	if err != redis.Nil {
		t.Fatal(err)
	}
}

func TestFindValuesNotFoundMultiple(t *testing.T) {
	db, mock := redismock.NewClientMock()
	key := "42"
	mock.ExpectLLen(key).SetErr(redis.Nil)
	client := &RedisClient{*db}
	_, _, err := client.FindValues(context.Background(), key, true, 5, 0)

	if err != redis.Nil {
		t.Fatal(err)
	}
}
