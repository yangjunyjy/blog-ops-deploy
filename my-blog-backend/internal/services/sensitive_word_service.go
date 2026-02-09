package services

import (
	"regexp"
	"strings"
	"sync"

	"my-blog-backend/internal/config"
)

// SensitiveWordService 敏感词服务
type SensitiveWordService interface {
	// CheckSensitiveWords 检测敏感词，返回是否包含敏感词和匹配的敏感词列表
	CheckSensitiveWords(content string) (bool, []string)

	// FilterSensitiveWords 过滤敏感词（替换为***）
	FilterSensitiveWords(content string) string

	// AddSensitiveWord 添加敏感词
	AddSensitiveWord(word string)

	// RemoveSensitiveWord 移除敏感词
	RemoveSensitiveWord(word string)

	// ReloadSensitiveWords 重新加载敏感词
	ReloadSensitiveWords(words []string)
}

type sensitiveWordServiceImpl struct {
	config          *config.Config
	sensitiveWords  map[string]bool
	sensitiveRegex  *regexp.Regexp
	sensitiveRWLock sync.RWMutex
}

func NewSensitiveWordService(cfg *config.Config) SensitiveWordService {
	svc := &sensitiveWordServiceImpl{
		config:         cfg,
		sensitiveWords: make(map[string]bool),
	}

	svc.ReloadSensitiveWords(cfg.Comment.SensitiveWords)

	return svc
}

// ReloadSensitiveWords 重新加载敏感词
func (s *sensitiveWordServiceImpl) ReloadSensitiveWords(words []string) {
	s.sensitiveRWLock.Lock()
	defer s.sensitiveRWLock.Unlock()

	// 清空现有敏感词
	s.sensitiveWords = make(map[string]bool)

	// 构建敏感词字典
	for _, word := range words {
		if word != "" {
			s.sensitiveWords[word] = true
		}
	}

	// 构建正则表达式（对每个敏感词进行转义）
	if len(s.sensitiveWords) > 0 {
		escapedWords := make([]string, 0, len(words))
		for _, word := range words {
			escapedWords = append(escapedWords, regexp.QuoteMeta(word))
		}
		pattern := strings.Join(escapedWords, "|")
		s.sensitiveRegex = regexp.MustCompile(pattern)
	}
}

// CheckSensitiveWords 检测敏感词
func (s *sensitiveWordServiceImpl) CheckSensitiveWords(content string) (bool, []string) {
	s.sensitiveRWLock.RLock()
	defer s.sensitiveRWLock.RUnlock()

	if s.sensitiveRegex == nil {
		return false, nil
	}

	matches := s.sensitiveRegex.FindAllString(content, -1)
	if len(matches) > 0 {
		// 去重
		uniqueMatches := make(map[string]bool)
		for _, match := range matches {
			uniqueMatches[match] = true
		}

		result := make([]string, 0, len(uniqueMatches))
		for match := range uniqueMatches {
			result = append(result, match)
		}

		return true, result
	}

	return false, nil
}

// FilterSensitiveWords 过滤敏感词
func (s *sensitiveWordServiceImpl) FilterSensitiveWords(content string) string {
	s.sensitiveRWLock.RLock()
	defer s.sensitiveRWLock.RUnlock()

	if s.sensitiveRegex == nil {
		return content
	}

	return s.sensitiveRegex.ReplaceAllStringFunc(content, func(match string) string {
		return strings.Repeat("*", len([]rune(match)))
	})
}

// AddSensitiveWord 添加敏感词
func (s *sensitiveWordServiceImpl) AddSensitiveWord(word string) {
	s.sensitiveRWLock.Lock()
	defer s.sensitiveRWLock.Unlock()

	if word != "" {
		s.sensitiveWords[word] = true

		// 重建正则表达式
		words := make([]string, 0, len(s.sensitiveWords))
		for w := range s.sensitiveWords {
			words = append(words, w)
		}
		escapedWords := make([]string, 0, len(words))
		for _, w := range words {
			escapedWords = append(escapedWords, regexp.QuoteMeta(w))
		}
		pattern := strings.Join(escapedWords, "|")
		s.sensitiveRegex = regexp.MustCompile(pattern)
	}
}

// RemoveSensitiveWord 移除敏感词
func (s *sensitiveWordServiceImpl) RemoveSensitiveWord(word string) {
	s.sensitiveRWLock.Lock()
	defer s.sensitiveRWLock.Unlock()

	delete(s.sensitiveWords, word)

	// 重建正则表达式
	words := make([]string, 0, len(s.sensitiveWords))
	for w := range s.sensitiveWords {
		words = append(words, w)
	}

	if len(words) > 0 {
		escapedWords := make([]string, 0, len(words))
		for _, w := range words {
			escapedWords = append(escapedWords, regexp.QuoteMeta(w))
		}
		pattern := strings.Join(escapedWords, "|")
		s.sensitiveRegex = regexp.MustCompile(pattern)
	} else {
		s.sensitiveRegex = nil
	}
}
