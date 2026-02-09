package smtp

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/big"
	"my-blog-backend/internal/config"
	"my-blog-backend/internal/pkg/logger"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

var globalEmailServer *Sender

// Mailer 定义发送邮件服务接口，便于扩展
type Mailer interface {
	Send(to, name string) error
}

// Sender 邮件发送器
type Sender struct {
	config     config.SmtpConfig
	codeCfg    config.VerificationCodeConfig
	auth       smtp.Auth
	mu         sync.Mutex
	clientPool chan *smtp.Client // 简单的客户端池
	stopChan   chan struct{}
}

// loginAuth 实现 LOGIN 认证方式
type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unexpected server challenge: %s", fromServer)
		}
	}
	return nil, nil
}

// NewEmailSender 创建新的邮件发送器
func NewEmailSender(config config.SmtpConfig, codeConfig config.VerificationCodeConfig) (*Sender, error) {
	// 尝试不同的认证方式
	var auth smtp.Auth

	// 首先尝试 PLAIN 认证
	auth = smtp.PlainAuth("", config.FromEmail, config.Password, config.SMTPHost)

	sender := &Sender{
		config:     config,
		codeCfg:    codeConfig,
		auth:       auth,
		clientPool: make(chan *smtp.Client, config.PoolSize), // 创建连接池
		stopChan:   make(chan struct{}),
	}

	logger.Info("创建邮件发送器",
		logger.String("host", config.SMTPHost),
		logger.Int("port", config.SMTPPort),
		logger.String("from", config.FromEmail))

	return sender, nil
}

// createClient 创建SMTP客户端
func (es *Sender) createClient() (*smtp.Client, error) {
	addr := net.JoinHostPort(es.config.SMTPHost, fmt.Sprintf("%d", es.config.SMTPPort))

	var conn net.Conn
	var err error

	// 设置连接超时
	dialer := &net.Dialer{Timeout: 30 * time.Second}

	logger.Info("连接SMTP服务器", logger.String("address", addr))

	if es.config.SMTPPort == 465 {
		// SSL/TLS连接
		tlsConfig := &tls.Config{
			ServerName:         es.config.SMTPHost,
			InsecureSkipVerify: false, // 调试时可以设置为 true
		}
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	} else {
		// 普通连接
		conn, err = dialer.Dial("tcp", addr)
	}

	if err != nil {
		return nil, fmt.Errorf("连接SMTP服务器失败: %v", err)
	}

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, es.config.SMTPHost)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建SMTP客户端失败: %v", err)
	}

	// 发送EHLO命令并检查服务器支持的功能
	if err := client.Hello("localhost"); err != nil {
		client.Close()
		return nil, fmt.Errorf("EHLO失败: %v", err)
	}

	// 检查服务器是否支持STARTTLS
	if es.config.SMTPPort != 465 {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{
				ServerName: es.config.SMTPHost,
			}
			if err := client.StartTLS(tlsConfig); err != nil {
				client.Close()
				return nil, fmt.Errorf("STARTTLS失败: %v", err)
			}
			logger.Info("已启用STARTTLS加密")
		}
	}

	// 检查服务器支持的认证方式
	_, authExtensions := client.Extension("AUTH")
	logger.Info("服务器支持的认证方式", logger.String("auth", authExtensions))

	// 根据服务器支持的认证方式选择合适的认证方法
	if strings.Contains(authExtensions, "PLAIN") {
		es.auth = smtp.PlainAuth("", es.config.FromEmail, es.config.Password, es.config.SMTPHost)
		logger.Info("使用PLAIN认证")
	} else if strings.Contains(authExtensions, "LOGIN") {
		es.auth = LoginAuth(es.config.FromEmail, es.config.Password)
		logger.Info("使用LOGIN认证")
	} else if strings.Contains(authExtensions, "CRAM-MD5") {
		es.auth = smtp.CRAMMD5Auth(es.config.FromEmail, es.config.Password)
		logger.Info("使用CRAM-MD5认证")
	} else {
		client.Close()
		return nil, fmt.Errorf("服务器不支持任何已知的认证方式: %s", authExtensions)
	}

	// 认证
	if err := client.Auth(es.auth); err != nil {
		client.Close()
		return nil, fmt.Errorf("SMTP认证失败: %v", err)
	}

	logger.Info("SMTP认证成功")
	return client, nil
}

// getClient 从池中获取客户端或创建新客户端
func (es *Sender) getClient() (*smtp.Client, error) {
	select {
	case client := <-es.clientPool:
		// 检查连接是否仍然有效
		if err := client.Noop(); err != nil {
			client.Close()
			return es.createClient() // 创建新连接替代失效的连接
		}
		return client, nil
	default:
		// 池中没有可用连接，创建新连接
		return es.createClient()
	}
}

// returnClient 将客户端返回池中
func (es *Sender) returnClient(client *smtp.Client) {
	// 检查连接是否仍然有效
	if err := client.Noop(); err != nil {
		client.Close()
		return
	}

	select {
	case es.clientPool <- client:
		// 成功放回池中
	default:
		// 池已满，关闭连接
		client.Quit()
		client.Close()
	}
}

// generateVerificationCode 生成验证码
func (es *Sender) generateVerificationCode() (string, error) {
	charset := es.codeCfg.Charset
	if charset == "" {
		charset = "0123456789" // 默认数字验证码
	}

	// 检查字符集是否有效
	if len(charset) == 0 {
		return "", fmt.Errorf("字符集为空，无法生成验证码")
	}

	// 检查验证码长度配置
	if es.codeCfg.Length <= 0 {
		es.codeCfg.Length = 6
	}

	code := make([]byte, es.codeCfg.Length)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("生成随机数失败: %v", err)
		}
		code[i] = charset[num.Int64()]
	}

	generatedCode := string(code)
	logger.Info("生成验证码成功", zap.String("code", generatedCode))
	return generatedCode, nil
}

// SendVerificationCode 发送验证码到单个邮箱
func (es *Sender) SendVerificationCode(toEmail, toName string) (string, error) {
	code, err := es.generateVerificationCode()
	if err != nil {
		logger.Error("生成验证码失败!", logger.Err("error", err))
		return "", err
	}

	err = es.SendVerificationCodeWithCustomCode(toEmail, toName, code)
	if err != nil {
		logger.Error("发送验证码失败!", logger.Err("error", err))
		return "", fmt.Errorf("发送验证码失败: %v", err)
	}
	return code, nil
}

// SendVerificationCodeWithCustomCode 使用自定义验证码发送
func (es *Sender) SendVerificationCodeWithCustomCode(toEmail, toName, code string) error {
	// 构建邮件内容
	from := fmt.Sprintf("%s <%s>", encodeHeader(es.config.FromName), es.config.FromEmail)
	to := fmt.Sprintf("%s <%s>", encodeHeader(toName), toEmail)
	subject := encodeHeader(es.config.Subject)

	// 使用模板填充内容
	htmlContent := fmt.Sprintf(es.codeCfg.HTMLTemplate, toName, code, es.codeCfg.ExpiresIn.Minutes())
	textContent := fmt.Sprintf(es.codeCfg.TextTemplate, toName, code, es.codeCfg.ExpiresIn.Minutes())

	// 构建邮件消息
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/alternative; boundary=\"BOUNDARY\"\r\n"+
		"\r\n"+
		"--BOUNDARY\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"Content-Transfer-Encoding: base64\r\n"+
		"\r\n"+
		"%s\r\n"+
		"\r\n"+
		"--BOUNDARY\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"Content-Transfer-Encoding: base64\r\n"+
		"\r\n"+
		"%s\r\n"+
		"\r\n"+
		"--BOUNDARY--\r\n",
		from, to, subject,
		base64.StdEncoding.EncodeToString([]byte(textContent)),
		base64.StdEncoding.EncodeToString([]byte(htmlContent)))

	// 发送邮件
	return es.sendMail(toEmail, message)
}

// sendMail 发送邮件核心方法
func (es *Sender) sendMail(toEmail, message string) error {
	// 获取SMTP客户端
	client, err := es.getClient()
	if err != nil {
		return fmt.Errorf("获取SMTP客户端失败: %v", err)
	}
	defer es.returnClient(client)

	// 设置发件人
	if err := client.Mail(es.config.FromEmail); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	// 发送邮件数据
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("准备发送数据失败: %v", err)
	}

	// 写入邮件内容
	_, err = w.Write([]byte(message))
	if err != nil {
		w.Close()
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	// 关闭写入器
	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭数据写入失败: %v", err)
	}

	logger.Info("邮件发送成功", logger.String("to", toEmail))
	return nil
}

// SendBatchVerificationCodes 批量发送验证码
func (es *Sender) SendBatchVerificationCodes(recipients []struct {
	Email string
	Name  string
}) (map[string]string, map[string]error) {
	successCodes := make(map[string]string)
	failures := make(map[string]error)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, recipient := range recipients {
		wg.Add(1)
		go func(email, name string) {
			defer wg.Done()

			code, err := es.SendVerificationCode(email, name)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				failures[email] = err
			} else {
				successCodes[email] = code
			}
		}(recipient.Email, recipient.Name)
	}

	wg.Wait()
	return successCodes, failures
}

// Close 关闭连接池
func (es *Sender) Close() {
	close(es.stopChan)

	// 关闭所有连接
	for {
		select {
		case client := <-es.clientPool:
			client.Quit()
			client.Close()
		default:
			return
		}
	}
}

// GetDefaultHTMLTemplate 获取默认HTML模板
func GetDefaultHTMLTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>验证码邮件</title>
</head>
<body>
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
        <h2>尊敬的 %s，您好！</h2>
        <p>您的验证码是：<strong style="font-size: 24px; color: #ff6600;">%s</strong></p>
        <p>该验证码 <strong>%.0f 分钟</strong>内有效，请及时使用。</p>
        <p>如非本人操作，请忽略此邮件。</p>
        <hr>
        <p style="color: #999; font-size: 12px;">此为系统自动发送邮件，请勿直接回复</p>
    </div>
</body>
</html>`
}

// GetDefaultTextTemplate 获取默认纯文本模板
func GetDefaultTextTemplate() string {
	return `尊敬的 %s，您好！
您的验证码是：%s
该验证码 %.0f 分钟内有效，请及时使用。
如非本人操作，请忽略此邮件。

此为系统自动发送邮件，请勿直接回复`
}

// GetDefaultVerificationCodeConfig 获取默认验证码配置
func GetDefaultVerificationCodeConfig() config.VerificationCodeConfig {
	return config.VerificationCodeConfig{
		Length:       6,
		Charset:      "0123456789",
		ExpiresIn:    5 * time.Minute,
		HTMLTemplate: GetDefaultHTMLTemplate(),
		TextTemplate: GetDefaultTextTemplate(),
	}
}

// InitEmailServer 初始化邮件服务器
func InitEmailServer(config config.SmtpConfig, codeConfig config.VerificationCodeConfig) error {
	logger.Info("初始化邮件服务")

	// 检查并填充默认值
	if codeConfig.Charset == "" {
		codeConfig.Charset = "0123456789"
		logger.Info("使用默认字符集")
	}
	if codeConfig.Length == 0 {
		codeConfig.Length = 6
		logger.Info("使用默认长度")
	}
	if codeConfig.ExpiresIn == 0 {
		codeConfig.ExpiresIn = 5 * time.Minute
		logger.Info("使用默认过期时间")
	}
	if codeConfig.HTMLTemplate == "" {
		codeConfig.HTMLTemplate = GetDefaultHTMLTemplate()
		logger.Info("使用默认HTML模板")
	}
	if codeConfig.TextTemplate == "" {
		codeConfig.TextTemplate = GetDefaultTextTemplate()
		logger.Info("使用默认文本模板")
	}

	sender, err := NewEmailSender(config, codeConfig)
	if err != nil {
		logger.Error("初始化邮件服务失败", logger.Err("error", err))
		return fmt.Errorf("初始化邮件服务失败:%v", err)
	}
	globalEmailServer = sender
	return nil
}

// GetEmailServer 获取邮件服务器实例
func GetEmailServer() *Sender {
	return globalEmailServer
}

// checkCodeConfig 检查验证码配置
func checkCodeConfig(codeConfig config.VerificationCodeConfig) bool {
	// 只要设置了基本参数就认为配置有效，模板可以留空使用默认值
	if codeConfig.Charset == "" || codeConfig.Length == 0 {
		return false
	}
	return true
}

// encodeHeader 编码邮件头
func encodeHeader(value string) string {
	// 如果包含非ASCII字符，则进行编码
	for _, r := range value {
		if r > 127 {
			return fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(value)))
		}
	}
	return value
}
