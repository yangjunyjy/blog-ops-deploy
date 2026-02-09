package captcha

import (
	"crypto/rand"
	"image/color"
	"math/big"
	"sync"

	"github.com/mojocn/base64Captcha"
)

var (
	store       base64Captcha.Store
	captchaOnce sync.Once
)

// InitCaptchaStore 初始化验证码存储
func InitCaptchaStore(expiration int) {
	captchaOnce.Do(func() {
		store = base64Captcha.NewMemoryStore(expiration, 10000)
	})
}

// GetStore 获取验证码存储实例
func GetStore() base64Captcha.Store {
	if store == nil {
		// 默认5分钟过期
		InitCaptchaStore(300)
	}
	return store
}

// GenerateCaptcha 生成图形验证码
func GenerateCaptcha(captchaType, length, width, height int) (id, b64s string, err error) {
	driver := createDriver(captchaType, length, width, height)
	c := base64Captcha.NewCaptcha(driver, GetStore())
	id, b64s, _, err = c.Generate()

	// 确保 base64 数据有正确的 data URL 前缀
	if b64s != "" && !contains(b64s, "data:image") {
		b64s = "data:image/png;base64," + b64s
	}

	return
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id, answer string) bool {
	return GetStore().Verify(id, answer, true)
}

// createDriver 创建验证码驱动
func createDriver(_, length, width, height int) base64Captcha.Driver {
	// 使用数字验证码驱动（最简单稳定）
	driver := base64Captcha.NewDriverDigit(
		height, // 高度
		width,  // 宽度
		length, // 验证码长度
		0.25,   // 最大倾斜角度
		50,     // 干扰线数量
	)

	// 如果你确实需要字符串验证码，可以这样配置：
	// driver := createStringDriver(height, width, length)

	return driver
}

// 创建字符串验证码驱动（如果需要）
func createStringDriver(height, width, length int) base64Captcha.Driver {
	// 创建一个简单的字体存储
	fonts := []string{
		"comic.ttf", // 这些字体需要存在于系统中，或者使用默认字体
	}

	// 使用 DriverString 结构体
	driver := &base64Captcha.DriverString{
		Height:          height,
		Width:           width,
		NoiseCount:      0, // 干扰线数量
		ShowLineOptions: 0, // 显示线条选项
		Length:          length,
		Source:          "0123456789", // 验证码来源，只使用数字
		BgColor: &color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
		Fonts: fonts,
	}

	// 如果没有字体文件，可以注释掉 Fonts 字段或使用空切片
	// driver.Fonts = []string{}

	return driver.ConvertFonts()
}

// GenerateRandomString 生成随机字符串（用于字母数字验证码）
func GenerateRandomString(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}
