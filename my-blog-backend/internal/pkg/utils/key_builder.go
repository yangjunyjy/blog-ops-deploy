package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"reflect"
	"strconv"
)

// KeyBuilder 缓存键构建器
type KeyBuilder struct {
	prefix string
	parts  []string
}

func NewKeyBuilder(prefix string) *KeyBuilder {
	return &KeyBuilder{
		prefix: prefix,
		parts:  []string{},
	}
}

// Add 添加键部分
func (k *KeyBuilder) Add(part interface{}) *KeyBuilder {
	var str string

	switch v := part.(type) {
	case string:
		str = v
	case int:
		str = strconv.Itoa(v)
	case int64:
		str = strconv.FormatInt(v, 10)
	case uint:
		str = strconv.FormatUint(uint64(v), 10)
	case uint64:
		str = strconv.FormatUint(v, 10)
	case bool:
		str = strconv.FormatBool(v)
	default:
		// 使用反射处理其他类型
		str = fmt.Sprintf("%v", v)
	}

	k.parts = append(k.parts, str)
	return k
}

// Build 构建完整的缓存键
func (k *KeyBuilder) Build() string {
	if k.prefix == "" && len(k.parts) == 0 {
		return ""
	}

	var key string
	if k.prefix != "" {
		key = k.prefix
	}

	for _, part := range k.parts {
		if key == "" {
			key = part
		} else {
			key = fmt.Sprintf("%s:%s", key, part)
		}
	}

	return key
}

// BuildWithHash 构建带哈希的缓存键
func (k *KeyBuilder) BuildWithHash() string {
	key := k.Build()

	// 如果键太长，使用哈希
	if len(key) > 200 {
		h := md5.New()
		h.Write([]byte(key))
		hash := hex.EncodeToString(h.Sum(nil))

		if k.prefix != "" {
			return fmt.Sprintf("%s:%s", k.prefix, hash)
		}
		return hash
	}

	return key
}

// Hash 计算哈希值
func Hash(value interface{}) string {
	h := fnv.New64a()

	switch v := value.(type) {
	case string:
		h.Write([]byte(v))
	case []byte:
		h.Write(v)
	case int:
		h.Write([]byte(strconv.Itoa(v)))
	default:
		// 使用反射
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		h.Write([]byte(fmt.Sprintf("%v", rv.Interface())))
	}

	return strconv.FormatUint(h.Sum64(), 10)
}

// GenerateCacheKey 生成缓存键
func GenerateCacheKey(prefix string, parts ...interface{}) string {
	builder := NewKeyBuilder(prefix)
	for _, part := range parts {
		builder.Add(part)
	}
	return builder.Build()
}
