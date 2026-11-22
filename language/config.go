package language

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

// Config 配置当前项目的本地词典；由启动代码调用 Load，不在包初始化或请求中读取文件。
type Config struct {
	// Dir 是词典目录，未配置或为空时使用工作目录下的 etc/language。
	// 每个语言一个 JSON 文件，例如 zh-CN.json，内容为扁平的字符串键值对。
	Dir string `json:",optional"`
}

type dictionaries map[string]map[string]string

// 已发布的词典只读，完整加载成功后才替换；并发请求不共享可变 map。
var current atomic.Pointer[dictionaries]

// Load 在服务接收请求前加载项目词典。目录为空时读取 etc/language；
// 目录不存在、没有 JSON 词典或词典无效时返回错误，调用方应停止启动。
// 忽略子目录及非 JSON 文件，不递归加载。失败保留原词典，成功整体替换；不监听文件变更。
func Load(config Config) error {
	dir := config.Dir
	if dir == "" {
		dir = filepath.Join("etc", "language")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("load language directory %q: %w", dir, err)
	}
	loaded := make(dictionaries)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		locale := normalizeLocale(strings.TrimSuffix(entry.Name(), ".json"))
		if locale == "" {
			return fmt.Errorf("load language: empty locale in %q", entry.Name())
		}
		if _, exists := loaded[locale]; exists {
			return fmt.Errorf("load language: duplicate locale %q", locale)
		}
		file := filepath.Join(dir, entry.Name())
		raw, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read language file %q: %w", file, err)
		}
		messages, err := parseDictionary(raw)
		if err != nil {
			return fmt.Errorf("parse language file %q: %w", file, err)
		}
		loaded[locale] = messages
	}
	if len(loaded) == 0 {
		return fmt.Errorf("load language directory %q: no JSON dictionaries", dir)
	}
	current.Store(&loaded)
	return nil
}

// parseDictionary 严格检查字符串词典，拒绝重复键、null、嵌套值和尾随 JSON，避免静默覆盖或误配。
func parseDictionary(raw []byte) (map[string]string, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	if token != json.Delim('{') {
		return nil, fmt.Errorf("dictionary must be a JSON object")
	}
	messages := make(map[string]string)
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		name := key.(string)
		if _, exists := messages[name]; exists {
			return nil, fmt.Errorf("duplicate dictionary key %q", name)
		}
		value, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		text, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("dictionary value for %q must be a string", name)
		}
		messages[name] = text
	}
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return nil, fmt.Errorf("dictionary contains trailing data")
	}
	return messages, nil
}

// normalizeLocale 统一大小写和下划线写法，zh_CN 与 zh-CN 查找同一词典。
func normalizeLocale(locale string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(locale)), "_", "-")
}

func dictionaryFor(locale string) map[string]string {
	snapshot := current.Load()
	if snapshot == nil {
		return nil
	}
	locale = normalizeLocale(locale)
	dictionary := (*snapshot)[locale]
	if dictionary == nil && locale == "en-us" {
		dictionary = (*snapshot)["en"]
	}
	return dictionary
}

func translate(dictionary map[string]string, original string) string {
	if text, exists := dictionary[original]; exists {
		return text
	}
	return original
}
