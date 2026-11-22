package language

import (
	"bytes"
	"encoding/json"
	"strings"
)

// TranslateMessage 用项目词典直接匹配原始文案或业务 key，不进行额外 key 映射。
// 未加载词典、语言或词条不存在时返回原文，空翻译值按词典配置返回。
func TranslateMessage(message, locale string) string {
	return translate(dictionaryFor(locale), message)
}

// SwitchLanguage 在响应编码时自动翻译匹配的字符串值，不修改字段名或原始对象。
// 仅保护名为 id/code（大小写不敏感）的直接字符串值；保留 JSON 标签、数字精度和自定义序列化。
// string 和 []byte 保留原类型以供 AES 使用；[]byte 不翻译，未知语言或未加载时原样返回。
func SwitchLanguage(data any, locale string) any {
	dictionary := dictionaryFor(locale)
	if len(dictionary) == 0 {
		return data
	}
	switch value := data.(type) {
	case nil, []byte:
		return data
	case string:
		return translate(dictionary, value)
	default:
		// 捕获只读词典快照，保证同一 data 编码期间不会混用不同次加载的词条。
		return localizedJSON{data: value, dictionary: dictionary}
	}
}

type localizedJSON struct {
	data       any
	dictionary map[string]string
}

func (data localizedJSON) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(data.data)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	rewrite := jsonTranslation{decoder: decoder, raw: raw, dictionary: data.dictionary}
	if err := rewrite.value(false); err != nil {
		return nil, err
	}
	if rewrite.output == nil {
		return raw, nil
	}
	return append(rewrite.output, raw[rewrite.copied:]...), nil
}

type jsonTranslation struct {
	decoder    *json.Decoder
	raw        []byte
	dictionary map[string]string
	output     []byte
	copied     int64
}

func (rewrite *jsonTranslation) value(protected bool) error {
	previous := rewrite.decoder.InputOffset()
	token, err := rewrite.decoder.Token()
	if err != nil {
		return err
	}
	switch token := token.(type) {
	case json.Delim:
		for rewrite.decoder.More() {
			protectedChild := false
			if token == '{' {
				key, err := rewrite.decoder.Token()
				if err != nil {
					return err
				}
				name := key.(string)
				protectedChild = strings.EqualFold(name, "id") || strings.EqualFold(name, "code")
			}
			if err := rewrite.value(protectedChild); err != nil {
				return err
			}
		}
		_, err = rewrite.decoder.Token()
		return err
	case string:
		if protected {
			return nil
		}
		text := translate(rewrite.dictionary, token)
		if text != token {
			end := rewrite.decoder.InputOffset()
			// Token 可能先消费逗号或冒号，只替换当前字符串的 JSON 字节。
			start := previous + int64(bytes.IndexByte(rewrite.raw[previous:end], '"'))
			encoded, err := json.Marshal(text)
			if err != nil {
				return err
			}
			rewrite.output = append(rewrite.output, rewrite.raw[rewrite.copied:start]...)
			rewrite.output = append(rewrite.output, encoded...)
			rewrite.copied = end
		}
	}
	return nil
}
