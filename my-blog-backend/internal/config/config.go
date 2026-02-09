package config

import (
	"time"
)

// Config 应用配置
type Config struct {
	App         AppConfig      `yaml:"app" env:"APP"`
	Server      ServerConfig   `yaml:"server" env:"SERVER"`
	Database    DatabaseConfig `yaml:"database" env:"DATABASE"`
	Redis       RedisConfig    `yaml:"redis" env:"REDIS"`
	Auth        AuthConfig     `yaml:"auth" env:"AUTH"`
	Session     SessionConfig  `yaml:"session" env:"SESSION"`
	OnceToken   OnceTokenConfig `yaml:"onceToken" env:"ONCE_TOKEN"`
	Log         LogConfig      `yaml:"log" env:"LOG"`
	EmailServer EmailConfig    `yaml:"emailServer" json:"emailServer"`
	Comment     CommentConfig  `yaml:"comment" env:"COMMENT"`
	Upload      UploadConfig   `yaml:"upload" env:"UPLOAD"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `yaml:"name" env:"NAME" env-default:"my-blog"`
	Version string `yaml:"version" env:"VERSION" env-default:"1.0.0"`
	Env     string `yaml:"env" env:"ENV" env-default:"development"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `yaml:"host" env:"HOST" env-default:"0.0.0.0"`
	Port         int           `yaml:"port" env:"PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"30s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"30s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"120s"`

	TLS struct {
		Enabled  bool   `yaml:"enabled" env:"TLS_ENABLED" env-default:"false"`
		CertFile string `yaml:"cert_file" env:"TLS_CERT_FILE"`
		KeyFile  string `yaml:"key_file" env:"TLS_KEY_FILE"`
	} `yaml:"tls"`

	CORS struct {
		Enabled          bool     `yaml:"enabled" env:"CORS_ENABLED" env-default:"true"`
		AllowedOrigins   []string `yaml:"allowed_origins" env:"CORS_ALLOWED_ORIGINS" env-separator:","`
		AllowedMethods   []string `yaml:"allowed_methods" env:"CORS_ALLOWED_METHODS" env-separator:","`
		AllowedHeaders   []string `yaml:"allowed_headers" env:"CORS_ALLOWED_HEADERS" env-separator:","`
		ExposeHeaders    []string `yaml:"expose_headers" env:"CORS_EXPOSE_HEADERS" env-separator:","`
		AllowCredentials bool     `yaml:"allow_credentials" env:"CORS_ALLOW_CREDENTIALS" env-default:"true"`
		MaxAge           int      `yaml:"max_age" env:"CORS_MAX_AGE" env-default:"86400"`
	} `yaml:"cors"`

	RateLimit struct {
		Enabled  bool          `yaml:"enabled" env:"RATE_LIMIT_ENABLED" env-default:"true"`
		Rate     int           `yaml:"rate" env:"RATE_LIMIT_RATE" env-default:"100"`
		Burst    int           `yaml:"burst" env:"RATE_LIMIT_BURST" env-default:"200"`
		Duration time.Duration `yaml:"duration" env:"RATE_LIMIT_DURATION" env-default:"1m"`
	} `yaml:"rate_limit"`

	HealthCheck struct {
		Enabled bool `yaml:"enabled" env:"HEALTH_CHECK_ENABLED" env-default:"true"`
		Port    int  `yaml:"port" env:"HEALTH_CHECK_PORT" env-default:"8081"`
	} `yaml:"health_check"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `yaml:"type" env:"TYPE" env-default:"mysql"`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"PORT" env-default:"3306"`
	Database string `yaml:"database" env:"DATABASE"`
	Username string `yaml:"username" env:"USERNAME"`
	Password string `yaml:"password" env:"PASSWORD"`
	Charset  string `yaml:"charset" env:"CHARSET" env-default:"utf8mb4"`
	DSN      string `yaml:"dsn" env:"DSN"` // 直接使用DSN，优先级高于其他配置

	Pool struct {
		MaxIdleConns    int           `yaml:"max_idle_conns" env:"MAX_IDLE_CONNS" env-default:"10"`
		MaxOpenConns    int           `yaml:"max_open_conns" env:"MAX_OPEN_CONNS" env-default:"100"`
		ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"CONN_MAX_LIFETIME" env-default:"1h"`
		ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env:"CONN_MAX_IDLE_TIME" env-default:"30m"`
	} `yaml:"pool"`

	ReadReplica struct {
		Enabled  bool   `yaml:"enabled" env:"READ_REPLICA_ENABLED" env-default:"false"`
		Host     string `yaml:"host" env:"READ_REPLICA_HOST"`
		Port     int    `yaml:"port" env:"READ_REPLICA_PORT"`
		Database string `yaml:"database" env:"READ_REPLICA_DATABASE"`
		Username string `yaml:"username" env:"READ_REPLICA_USERNAME"`
		Password string `yaml:"password" env:"READ_REPLICA_PASSWORD"`
		DSN      string `yaml:"dsn" env:"READ_REPLICA_DSN"`
	} `yaml:"read_replica"`

	LogLevel string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret string        `yaml:"jwtSecret" env:"JWT_SECRET"`
	JWTExpire time.Duration `yaml:"jwtExpire" env:"JWT_EXPIRE" env-default:"24h"`
	Issuer    string        `yaml:"issuer" env:"ISSUER" env-default:"my-blog"`

	Password struct {
		MinLength    int  `yaml:"min_length" env:"MIN_LENGTH" env-default:"6"`
		RequireUpper bool `yaml:"require_upper" env:"REQUIRE_UPPER" env-default:"true"`
		RequireLower bool `yaml:"require_lower" env:"REQUIRE_LOWER" env-default:"true"`
		RequireDigit bool `yaml:"require_digit" env:"REQUIRE_DIGIT" env-default:"true"`
	} `yaml:"password"`

	// TOTP配置（基于时间的动态Token）
	TOTP struct {
		Enabled      bool   `yaml:"enabled" env:"TOTP_ENABLED" env-default:"false"`
		Secret       string `yaml:"secret" env:"TOTP_SECRET"`           // 全局共享密钥
		TimeStep     int    `yaml:"timeStep" env:"TOTP_TIME_STEP" env-default:"30"` // 时间步长（秒）
		WindowSize   int    `yaml:"windowSize" env:"TOTP_WINDOW_SIZE" env-default:"1"` // 验证窗口大小
	} `yaml:"totp"`
}

// SessionConfig Session配置
type SessionConfig struct {
	CookieName string `yaml:"cookieName" env:"COOKIE_NAME" env-default:"session_id"`
	MaxAge     int    `yaml:"maxAge" env:"MAX_AGE" env-default:"86400"`
	Path       string `yaml:"path" env:"PATH" env-default:"/"`
	Domain     string `yaml:"domain" env:"DOMAIN" env-default:""`
	Secure     bool   `yaml:"secure" env:"SECURE" env-default:"false"`
	HttpOnly   bool   `yaml:"httpOnly" env:"HTTP_ONLY" env-default:"true"`
}

// OnceTokenConfig 一次性Token配置
type OnceTokenConfig struct {
	ExpireTime int `yaml:"expireTime" env:"EXPIRE_TIME" env-default:"300"`
}

// LogConfig 日志配置
type LogConfig struct {
	Devolpment    bool   `yaml:"devolpment" json:"devolpment"`       //运行模式
	Level         string `yaml:"level" json:"level"`                 //日志级别
	Encoding      string `yaml:"encoding" json:"encoding"`           //输出格式，console,json
	Directory     string `yaml:"directory" json:"directory"`         //日志目录
	FileName      string `yaml:"fileName" json:"fileName"`           //日志名
	MaxSize       uint   `yaml:"maxSize" json:"maxSize"`             //日志文件最大大小
	MaxBackups    uint   `yaml:"maxBackups" json:"maxBackups"`       //保留旧文件最大个数
	MaxAge        uint   `yaml:"maxAge" json:"maxAge"`               //保留旧文件最大天数
	Compress      bool   `yaml:"compress" json:"compress"`           //是否压缩旧文件
	EnableConsole bool   `yaml:"enableConsole" json:"enableConsole"` //是否启用控制台输出
	FileConsole   bool   `yaml:"fileConsole" json:"fileConsole"`     //是否启用文件输出
	EnableCaller  bool   `yaml:"enableCaller" json:"enableCaller"`   //是否添加调用者信息
}

// CacheConfig 缓存配置
type CacheConfig struct {
	// 基本配置
	Type         string        `yaml:"type" env:"TYPE" env-default:"redis"`
	Prefix       string        `yaml:"prefix" env:"PREFIX" env-default:"app"`
	DefaultTTL   time.Duration `yaml:"default_ttl" env:"DEFAULT_TTL" env-default:"1h"`
	MaxKeySize   int           `yaml:"max_key_size" env:"MAX_KEY_SIZE" env-default:"1024"`
	MaxValueSize int           `yaml:"max_value_size" env:"MAX_VALUE_SIZE" env-default:"1048576"` // 1MB
	Compression  bool          `yaml:"compression" env:"COMPRESSION" env-default:"false"`
	Serializer   string        `yaml:"serializer" env:"SERIALIZER" env-default:"json"`

	// Redis配置
	Redis RedisConfig `yaml:"redis" env:"REDIS"`

	// 内存缓存配置
	Memory struct {
		MaxMemory       int64         `yaml:"max_memory" env:"MAX_MEMORY" env-default:"104857600"` // 100MB
		CleanupInterval time.Duration `yaml:"cleanup_interval" env:"CLEANUP_INTERVAL" env-default:"1m"`
		MaxItems        int           `yaml:"max_items" env:"MAX_ITEMS" env-default:"10000"`
	} `yaml:"memory"`

	// 多级缓存配置
	Levels []CacheLevelConfig `yaml:"levels" env:"LEVELS"`

	// 监控配置
	Metrics struct {
		Enabled bool   `yaml:"enabled" env:"ENABLED" env-default:"true"`
		Prefix  string `yaml:"prefix" env:"PREFIX" env-default:"cache"`
	} `yaml:"metrics"`

	// 熔断器配置
	CircuitBreaker struct {
		Enabled          bool          `yaml:"enabled" env:"ENABLED" env-default:"true"`
		Name             string        `yaml:"name" env:"NAME" env-default:"cache"`
		MaxRequests      uint32        `yaml:"max_requests" env:"MAX_REQUESTS" env-default:"10"`
		Interval         time.Duration `yaml:"interval" env:"INTERVAL" env-default:"10s"`
		Timeout          time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
		TripThreshold    float64       `yaml:"trip_threshold" env:"TRIP_THRESHOLD" env-default:"0.5"`
		SuccessThreshold uint32        `yaml:"success_threshold" env:"SUCCESS_THRESHOLD" env-default:"5"`
	} `yaml:"circuit_breaker"`

	// 上下文超时
	ContextTimeout time.Duration `yaml:"context_timeout" env:"CONTEXT_TIMEOUT" env-default:"5s"`
}

// CacheLevelConfig 缓存级别配置
type CacheLevelConfig struct {
	Name            string        `yaml:"name" json:"name"`
	Type            string        `yaml:"type" json:"type"`
	Priority        int           `yaml:"priority" json:"priority"`
	Prefix          string        `yaml:"prefix" json:"prefix"`
	TTL             time.Duration `yaml:"ttl" json:"ttl"`
	MaxSize         int64         `yaml:"max_size" json:"max_size"`
	CleanupInterval time.Duration `yaml:"cleanup_interval" json:"cleanup_interval"`
	Redis           *RedisConfig  `yaml:"redis" json:"redis"`
}

func (config *LogConfig) SetDefault() {
	if config.Level == "" {
		config.Level = "info"
	}
	if config.Encoding == "" {
		config.Encoding = "json"
	}
	if config.FileName == "" {
		config.FileName = "app.log"
	}
	if config.MaxSize == 0 {
		config.MaxSize = 100 // 100MB
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 30 // 保留30个备份
	}
	if config.MaxAge == 0 {
		config.MaxAge = 90 // 保留90天
	}
}

// SmtpConfig 使用189.cn邮箱实现
type SmtpConfig struct {
	SMTPHost  string `yaml:"smtpHost" json:"smtpHost"`   // SMTP服务器地址
	SMTPPort  int    `yaml:"smtpPort" json:"smtpPort"`   // SMTP端口
	FromEmail string `yaml:"fromEmail" json:"fromEmail"` // 发件人邮箱
	FromName  string `yaml:"fromName" json:"fromName"`   // 发件人名称
	Password  string `yaml:"password" json:"password"`   // 授权码/密码
	Subject   string `yaml:"subject" json:"subject"`     // 邮件主题
	PoolSize  int    `yaml:"poolSize" json:"poolSize"`   // 连接池大小
}

type VerificationCodeConfig struct {
	Length       int           `yaml:"length" env:"LENGTH" env-default:"6"`            // 验证码长度
	Charset      string        `yaml:"charset" env:"CHARSET" env-default:"0123456789"` // 字符集
	ExpiresIn    time.Duration `yaml:"expiresIn" env:"EXPIRES_IN" env-default:"5m"`    // 过期时间
	HTMLTemplate string        `yaml:"htmlTemplate" env:"HTML_TEMPLATE"`               // HTML模板
	TextTemplate string        `yaml:"textTemplate" env:"TEXT_TEMPLATE"`               // 纯文本模板
}

type EmailConfig struct {
	SmtpServerConfig SmtpConfig             `yaml:"smtpServerConfig" json:"smtpServerConfig"`
	CodeConfig       VerificationCodeConfig `yaml:"codeConfig" json:"codeConfig"`
}

// CommentConfig 评论配置
type CommentConfig struct {
	// 审核配置
	Enabled     bool `yaml:"enabled" env:"ENABLED" env-default:"true"`           // 是否启用评论审核
	AutoApprove bool `yaml:"autoApprove" env:"AUTO_APPROVE" env-default:"false"` // 是否自动通过（无敏感词直接通过）

	// 敏感词配置
	SensitiveWordCheck bool     `yaml:"sensitiveWordCheck" env:"SENSITIVE_WORD_CHECK" env-default:"true"` // 是否启用敏感词检测
	SensitiveWords     []string `yaml:"sensitiveWords" env:"SENSITIVE_WORDS" env-separator:"|"`           // 敏感词列表

	// 多级评论配置
	MaxLevel     int  `yaml:"maxLevel" env:"MAX_LEVEL" env-default:"3"`            // 最大评论层级
	MaxDepth     int  `yaml:"maxDepth" env:"MAX_DEPTH" env-default:"3"`            // 最大回复深度
	ReplyEnabled bool `yaml:"replyEnabled" env:"REPLY_ENABLED" env-default:"true"` // 是否允许回复评论

	// 限制配置
	MaxLength     int `yaml:"maxLength" env:"MAX_LENGTH" env-default:"1000"`        // 单条评论最大长度
	MinLength     int `yaml:"minLength" env:"MIN_LENGTH" env-default:"1"`           // 单条评论最小长度
	MaxPerDay     int `yaml:"maxPerDay" env:"MAX_PER_DAY" env-default:"50"`         // 每人每天最大评论数
	MaxPerArticle int `yaml:"maxPerArticle" env:"MAX_PER_ARTICLE" env-default:"10"` // 每人每篇文章最大评论数
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	// 上传方式: local 或 remote
	Type string `yaml:"type" env:"TYPE" env-default:"local"`

	// 本地上传配置
	Local struct {
		UploadPath string `yaml:"uploadPath" env:"UPLOAD_PATH" env-default:"./uploads"`       // 上传目录
		AllowedExts []string `yaml:"allowedExts" env:"ALLOWED_EXTS" env-separator:","`        // 允许的文件扩展名
		MaxSize    int64  `yaml:"maxSize" env:"MAX_SIZE" env-default:"10485760"`              // 最大文件大小(10MB)
	} `yaml:"local"`

	// 远程上传配置 (如 OSS, COS, 七牛云等)
	Remote struct {
		Provider string `yaml:"provider" env:"PROVIDER" env-default:"oss"`  // 服务提供商: oss, cos, qiniu
		Endpoint string `yaml:"endpoint" env:"ENDPOINT"`                    // 服务端点
		AccessKey string `yaml:"accessKey" env:"ACCESS_KEY"`               // 访问密钥
		SecretKey string `yaml:"secretKey" env:"SECRET_KEY"`               // 密钥
		Bucket    string `yaml:"bucket" env:"BUCKET"`                      // 存储桶名称
		Region    string `yaml:"region" env:"REGION"`                      // 区域
		CDN       string `yaml:"cdn" env:"CDN"`                            // CDN域名
	} `yaml:"remote"`

	// URL 配置
	URLPrefix string `yaml:"urlPrefix" env:"URL_PREFIX"` // URL前缀，用于拼接完整URL
}
