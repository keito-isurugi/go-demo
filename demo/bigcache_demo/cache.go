package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/allegro/bigcache/v3"
)

// Cache bigcacheの操作を抽象化したインターフェース
type Cache interface {
	Get(key string, mapper interface{}) error
	Set(key string, value interface{}) error
	Delete(key string) error
	Reset() error
}

type BigCache struct {
	Cache *bigcache.BigCache
}

// NewCache bigcacheのインスタンスを生成する関数
func NewCache(c *bigcache.BigCache) (Cache, error) {
	if c == nil {
		return nil, errors.New("cache instance is nil")
	}
	return &BigCache{
		Cache: c,
	}, nil
}

// 登録済みの型を保持するマップ
var registeredTypes = make(map[string]bool)

// registerType 型を一度だけ登録するための関数
func registerType(value interface{}) {
	if value == nil {
		return
	}
	typeName := fmt.Sprintf("%T", value)
	if !registeredTypes[typeName] {
		gob.Register(value)
		registeredTypes[typeName] = true
	}
}

// Get キャッシュからデータを取得する関数
func (c BigCache) Get(key string, mapper interface{}) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if mapper == nil {
		return errors.New("mapper cannot be nil")
	}

	registerType(mapper)

	entry, err := c.Cache.Get(key)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewReader(entry)).Decode(mapper)
}

func (c BigCache) Set(key string, value interface{}) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	vb, err := serialize(value)
	if err != nil {
		return err
	}
	return c.Cache.Set(key, vb)
}

// serialize キャッシュに保存するデータをシリアライズする関数
func serialize(value interface{}) ([]byte, error) {
	registerType(value)

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode: %w", err)
	}

	return buf.Bytes(), nil
}

// Delete キャッシュからデータを削除する関数
func (c BigCache) Delete(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	return c.Cache.Delete(key)
}

// Reset キャッシュをリセットする関数
func (c BigCache) Reset() error {
	return c.Cache.Reset()
}
