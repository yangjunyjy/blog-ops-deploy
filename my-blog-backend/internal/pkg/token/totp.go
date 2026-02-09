package token

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"

	"my-blog-backend/internal/config"
)

// TOTP 基于时间的一次性密码生成和验证器
type TOTP struct {
	config     *config.AuthConfig
	secret     []byte
	timeStep   int
	windowSize int
}

// NewTOTP 创建TOTP实例
func NewTOTP(config *config.AuthConfig) (*TOTP, error) {
	if config == nil || !config.TOTP.Enabled {
		return nil, fmt.Errorf("TOTP未启用")
	}

	secret := config.TOTP.Secret
	if secret == "" {
		return nil, fmt.Errorf("TOTP密钥未配置")
	}

	// 标准化密钥，确保长度合适
	secret = strings.TrimSpace(secret)
	if len(secret) < 16 {
		// 如果密钥太短，可以重复填充，但更好的做法是生成新密钥
		return nil, fmt.Errorf("TOTP密钥太短，至少需要16个字符")
	}

	return &TOTP{
		config:     config,
		secret:     []byte(secret),
		timeStep:   config.TOTP.TimeStep,
		windowSize: config.TOTP.WindowSize,
	}, nil
}

// Generate 生成当前时间窗口的TOTP
func (t *TOTP) Generate() (string, error) {
	return t.GenerateAtTime(time.Now())
}

// GenerateAtTime 生成指定时间的TOTP
func (t *TOTP) GenerateAtTime(timestamp time.Time) (string, error) {
	if !t.config.TOTP.Enabled {
		return "", fmt.Errorf("TOTP未启用")
	}

	counter := t.getCounter(timestamp)
	token := t.hotp(counter)
	
	// 格式化为6位数字（标准TOTP格式）
	return fmt.Sprintf("%06d", token), nil
}

// Verify 验证TOTP是否有效
func (t *TOTP) Verify(token string) (bool, error) {
	return t.VerifyAtTime(token, time.Now())
}

// VerifyAtTime 在指定时间验证TOTP是否有效
func (t *TOTP) VerifyAtTime(token string, timestamp time.Time) (bool, error) {
	if !t.config.TOTP.Enabled {
		return false, fmt.Errorf("TOTP未启用")
	}

	token = strings.TrimSpace(token)
	if len(token) != 6 {
		return false, nil
	}

	// 检查当前时间窗口及前后窗口
	for i := -t.windowSize; i <= t.windowSize; i++ {
		adjustedTime := timestamp.Add(time.Duration(i*t.timeStep) * time.Second)
		counter := t.getCounter(adjustedTime)
		expectedToken := t.hotp(counter)
		expectedTokenStr := fmt.Sprintf("%06d", expectedToken)
		
		if hmac.Equal([]byte(token), []byte(expectedTokenStr)) {
			return true, nil
		}
	}

	return false, nil
}

// GenerateTokenForAPI 为API请求生成动态Token
// 这个Token会包含时间戳信息，用于API请求验证
func (t *TOTP) GenerateTokenForAPI() (string, error) {
	timestamp := time.Now()
	counter := t.getCounter(timestamp)
	token := t.hotp(counter)
	
	// 生成包含时间和计数器的Token格式: timestamp-counter-token
	// 实际使用时，客户端只需要发送token，服务器会根据时间验证
	return fmt.Sprintf("%d-%06d", timestamp.Unix(), token), nil
}

// VerifyTokenForAPI 验证API请求的动态Token
func (t *TOTP) VerifyTokenForAPI(token string) (bool, error) {
	parts := strings.Split(token, "-")
	if len(parts) != 2 {
		return false, nil
	}
	
	// 第一部分是时间戳，第二部分是token
	tokenValue := parts[1]
	return t.Verify(tokenValue)
}

// getCounter 获取时间计数器
func (t *TOTP) getCounter(timestamp time.Time) int64 {
	// Unix时间戳除以时间步长
	return timestamp.Unix() / int64(t.timeStep)
}

// hotp 生成HMAC-based One-Time Password
func (t *TOTP) hotp(counter int64) int32 {
	// 将计数器转换为8字节的大端序
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, uint64(counter))

	// 计算HMAC-SHA1
	hmacHash := hmac.New(sha1.New, t.secret)
	hmacHash.Write(counterBytes)
	hash := hmacHash.Sum(nil)

	// 动态截取
	offset := hash[len(hash)-1] & 0x0F
	binaryCode := (int32(hash[offset]) & 0x7F) << 24 |
		(int32(hash[offset+1]) & 0xFF) << 16 |
		(int32(hash[offset+2]) & 0xFF) << 8 |
		(int32(hash[offset+3]) & 0xFF)

	// 取模得到6位数字
	mod := int32(math.Pow10(6))
	return binaryCode % mod
}

// GenerateSecret 生成安全的TOTP密钥（用于初始化）
func GenerateSecret() (string, error) {
	// 生成32字节的随机密钥（符合TOTP标准）
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	// 转换为Base32格式（TOTP标准格式）
	encoded := base32.StdEncoding.EncodeToString(bytes)
	// 去除填充字符
	encoded = strings.TrimRight(encoded, "=")
	
	return encoded, nil
}