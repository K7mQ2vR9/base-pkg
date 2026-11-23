package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
	"unicode/utf8"

	"github.com/K7mQ2vR9/base-pkg/es/esv7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// ErrNotFound 表示文档或索引不存在，适用于 errors.Is，与 OpenSearch 错误独立。
var ErrNotFound = errors.New("elasticsearch7 document or index not found")

// ErrInvalidResponse 表示缺失响应、正文或必要字段，不能将异常响应当作成功或零条数据。
var ErrInvalidResponse = errors.New("elasticsearch7 invalid response")

// ErrIncompleteSearch 表示搜索超时、分片失败、提前终止或总数仅为下界，不能作为完整结果使用。
// 使用 errors.Is 判断；此错误不会回显原始查询或分片错误，且在写入目标对象前返回。
var ErrIncompleteSearch = errors.New("elasticsearch7 search result is incomplete")

// defaultMaxResponseBytes 为未配置客户端或独立解析响应的默认成功正文上限，单位为字节，固定 16 MiB。
const defaultMaxResponseBytes int64 = 16 << 20

// ErrResponseTooLarge 表示成功响应正文超过读取上限，使用 errors.Is 判断；不包含原始正文或查询。
// ReadResponse 返回此错误时不提供截断字节，模型解析不会因此修改业务目标。
var ErrResponseTooLarge = errors.New("elasticsearch7 response body exceeds limit")

// limitedResponseBody 按实例配置限制实际读取字节，不信任 Content-Length。
// 单个读取者消费；Close 可并发或重复调用，超过上限只多读一个判定字节并立即关闭底层正文。
type limitedResponseBody struct {
	source    io.ReadCloser
	remaining int64
	exceeded  bool
	closeOnce sync.Once
	closeErr  error
}

// newLimitedResponseBody 绑定已校验的实例上限，0 使用独立解析器的默认值，不允许无限读取。
func newLimitedResponseBody(source io.ReadCloser, maxBytes int64) *limitedResponseBody {
	if maxBytes == 0 {
		maxBytes = defaultMaxResponseBytes
	}
	return &limitedResponseBody{source: source, remaining: maxBytes}
}

func (body *limitedResponseBody) Read(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if body.exceeded {
		return 0, ErrResponseTooLarge
	}
	if body.remaining == 0 {
		// 上限内正文仍可能恰好结束；只用一个额外字节区分 EOF 与实际超限，不把它交给解析器。
		var probe [1]byte
		count, err := body.source.Read(probe[:])
		if count > 0 {
			body.exceeded = true
			_ = body.Close()
			return 0, ErrResponseTooLarge
		}
		return 0, err
	}
	if int64(len(data)) > body.remaining {
		data = data[:body.remaining]
	}
	count, err := body.source.Read(data)
	body.remaining -= int64(count)
	return count, err
}

func (body *limitedResponseBody) Close() error {
	body.closeOnce.Do(func() {
		body.closeErr = body.source.Close()
	})
	return body.closeErr
}

// ResponseError 保留 HTTP 状态和最多 128 字节的结构化错误类型，不暴露原始正文或 error.reason。
type ResponseError struct {
	// StatusCode 为服务端 HTTP 状态码，404 可通过 errors.Is 匹配 ErrNotFound。
	StatusCode int
	// ErrorType 只保留结构化 error.type，最多 128 字节；缺失时使用 HTTP 状态文本。
	ErrorType string
}

// Error 返回状态与错误类型摘要，不包含原始正文或 error.reason。
func (e *ResponseError) Error() string {
	return fmt.Sprintf("elasticsearch7 HTTP %d (%s)", e.StatusCode, e.ErrorType)
}

// Unwrap 使 HTTP 404 可通过 errors.Is(err, ErrNotFound) 识别，其他状态不映射业务错误。
func (e *ResponseError) Unwrap() error {
	if e.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	return nil
}

// ReadResponse 接管并关闭有效响应体，成功正文受客户端 MaxResponseBytes 限制，独立响应默认 16 MiB。
// 超限返回可用 errors.Is 判断的 ErrResponseTooLarge，不提供截断正文；错误响应仍最多读取 4 KiB。
// 所有提前返回均关闭正文，不无限 drain；ctx 须与请求相同以中断网络读取。
func ReadResponse(ctx context.Context, response *esapi.Response) (body []byte, err error) {
	if response == nil {
		return nil, ErrInvalidResponse
	}
	if response.Body != nil {
		defer func() {
			if closeErr := response.Body.Close(); err == nil && closeErr != nil {
				err = &operationError{"close response", closeErr}
			}
		}()
	}
	if ctx == nil {
		return nil, errors.New("elasticsearch7 response context is required")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if response.Body == nil || response.StatusCode < 200 {
		return nil, ErrInvalidResponse
	}
	var reader io.Reader = response.Body
	if response.StatusCode >= 300 {
		reader = io.LimitReader(reader, 4096)
	} else if _, bounded := response.Body.(*limitedResponseBody); !bounded {
		// 非本客户端创建的响应也有默认保护；已有实例包装不再套默认值，允许显式调大上限。
		response.Body = newLimitedResponseBody(response.Body, defaultMaxResponseBytes)
		reader = response.Body
	}
	body, err = io.ReadAll(reader)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if err != nil {
		return nil, &operationError{"read response", err}
	}
	if response.StatusCode >= 300 {
		var payload struct {
			Error struct {
				Type string `json:"type"`
			} `json:"error"`
		}
		_ = json.Unmarshal(body, &payload)
		kind := payload.Error.Type
		if kind == "" {
			kind = http.StatusText(response.StatusCode)
		}
		if len(kind) > 128 {
			kind = kind[:128]
			for !utf8.ValidString(kind) && len(kind) > 0 {
				kind = kind[:len(kind)-1]
			}
		}
		return nil, &ResponseError{StatusCode: response.StatusCode, ErrorType: kind}
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, ErrInvalidResponse
	}
	return body, nil
}

// CheckResponse 用于探活和写操作，校验 HTTP 状态与 JSON 正文并关闭响应体。
func CheckResponse(ctx context.Context, response *esapi.Response) error {
	body, err := ReadResponse(ctx, response)
	if err != nil {
		return err
	}
	var result json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		return &operationError{"decode response", err}
	}
	if raw := bytes.TrimSpace(result); len(raw) == 0 || raw[0] != '{' {
		return ErrInvalidResponse
	}
	return ctx.Err()
}

// ParseGetResponse 关闭响应体并直接解码到调用方目标；404 或 found=false 返回 ErrNotFound。
// 目标须为非 nil 指针，成功时向实现 Schema 的目标设置文档 ID。
func ParseGetResponse(ctx context.Context, target any, response *esapi.Response) error {
	body, err := ReadResponse(ctx, response)
	if err != nil {
		return err
	}
	var document struct {
		ID     string          `json:"_id"`
		Found  *bool           `json:"found"`
		Source json.RawMessage `json:"_source"`
	}
	if err := json.Unmarshal(body, &document); err != nil {
		return &operationError{"decode get response", err}
	}
	if document.Found != nil && !*document.Found {
		return ErrNotFound
	}
	if document.ID == "" || len(document.Source) == 0 || bytes.Equal(bytes.TrimSpace(document.Source), []byte("null")) {
		return ErrInvalidResponse
	}
	if err := json.Unmarshal(document.Source, target); err != nil {
		return &operationError{"decode document", err}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if schema, ok := target.(esv7.Schema); ok {
		schema.SetId(document.ID)
	}
	return nil
}

// ParseSearchResponse 关闭响应体并将 _source 数组解码到目标切片，保留 int64 精度，空命中返回空数组。
// target 必须为非 nil 切片指针；支持 []T 和 []*T，实现 Schema 的元素按命中顺序回填原生 _id。
// 需回填 ID 的命中缺少非空 _id 时返回 ErrInvalidResponse；普通非 Schema 元素仍仅解码 _source。
// 必须提供 timed_out 和 _shards.failed；超时、失败分片或 terminated_early=true 返回 ErrIncompleteSearch。
// total 支持数字或 relation=eq 的对象；gte 返回 ErrIncompleteSearch，缺失/未知关系或非法字段返回 ErrInvalidResponse。
// 完整性检查先于目标解码；临时切片解码及 ID 回填成功后才替换 target，不回滚自定义方法的外部副作用。
func ParseSearchResponse(ctx context.Context, target any, response *esapi.Response) (int64, error) {
	body, err := ReadResponse(ctx, response)
	if err != nil {
		return 0, err
	}
	var result struct {
		TimedOut        *bool `json:"timed_out"`
		TerminatedEarly *bool `json:"terminated_early"`
		Shards          *struct {
			Failed *int64 `json:"failed"`
		} `json:"_shards"`
		Hits *struct {
			Total json.RawMessage `json:"total"`
			Hits  []struct {
				ID     string          `json:"_id"`
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, &operationError{"decode search response", err}
	}
	// 不将缺失或 null 的完整性字段按零值处理，避免把删减过的响应误判为完整结果。
	if result.TimedOut == nil || result.Shards == nil || result.Shards.Failed == nil || *result.Shards.Failed < 0 {
		return 0, ErrInvalidResponse
	}
	if *result.TimedOut || *result.Shards.Failed > 0 || (result.TerminatedEarly != nil && *result.TerminatedEarly) {
		return 0, ErrIncompleteSearch
	}
	if result.Hits == nil || result.Hits.Hits == nil {
		return 0, ErrInvalidResponse
	}
	var total int64
	if raw := bytes.TrimSpace(result.Hits.Total); len(raw) > 0 && raw[0] == '{' {
		var object struct {
			Value    *int64 `json:"value"`
			Relation string `json:"relation"`
		}
		if err := json.Unmarshal(raw, &object); err != nil {
			return 0, &operationError{"decode search total", err}
		}
		if object.Value == nil || *object.Value < 0 {
			return 0, ErrInvalidResponse
		}
		// gte 仅表示总数下界，不满足本接口返回准确总数的约定。
		switch object.Relation {
		case "eq":
		case "gte":
			return 0, ErrIncompleteSearch
		default:
			return 0, ErrInvalidResponse
		}
		total = *object.Value
	} else if err := json.Unmarshal(result.Hits.Total, &total); err != nil {
		return 0, &operationError{"decode search total", err}
	}
	if total < 0 || bytes.Equal(bytes.TrimSpace(result.Hits.Total), []byte("null")) {
		return 0, ErrInvalidResponse
	}
	sources := make([]json.RawMessage, 0, len(result.Hits.Hits))
	ids := make([]string, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		if len(hit.Source) == 0 || bytes.Equal(bytes.TrimSpace(hit.Source), []byte("null")) {
			return 0, ErrInvalidResponse
		}
		sources = append(sources, hit.Source)
		ids = append(ids, hit.ID)
	}
	data, err := json.Marshal(sources)
	if err != nil {
		return 0, &operationError{"encode search documents", err}
	}
	if err := decodeSearchDocuments(ctx, target, data, ids); err != nil {
		return 0, err
	}
	return total, nil
}

// decodeSearchDocuments 先解码到新的同类型切片，再按命中顺序设置文档 ID，最后替换目标切片。
// 普通解码失败、ID 缺失或提交前上下文已取消时不替换目标；自定义解码器/SetId 的外部副作用不回滚。
func decodeSearchDocuments(ctx context.Context, target any, data []byte, ids []string) error {
	value := reflect.ValueOf(target)
	if !value.IsValid() || value.Kind() != reflect.Pointer || value.IsNil() || value.Elem().Kind() != reflect.Slice {
		return &operationError{"decode search documents", fmt.Errorf("search target must be a non-nil pointer to a slice")}
	}
	decoded := reflect.New(value.Elem().Type())
	if err := json.Unmarshal(data, decoded.Interface()); err != nil {
		return &operationError{"decode search documents", err}
	}
	documents := decoded.Elem()
	// 自定义切片解码器不能改变命中数量，否则无法按位置可靠对应原生 ID。
	if documents.Len() != len(ids) {
		return ErrInvalidResponse
	}
	for index := 0; index < documents.Len(); index++ {
		document := documents.Index(index)
		// []T 使用元素地址调用指针接收者；[]*T 直接使用元素，兼容嵌入 EsSchema 的模型。
		schema, ok := document.Addr().Interface().(esv7.Schema)
		if !ok {
			schema, ok = document.Interface().(esv7.Schema)
		}
		if !ok {
			continue
		}
		if ids[index] == "" {
			return ErrInvalidResponse
		}
		if receiver := reflect.ValueOf(schema); receiver.Kind() == reflect.Pointer && receiver.IsNil() {
			return ErrInvalidResponse
		}
		// 原生 _id 优先于 _source 内映射到同一字段的业务值。
		schema.SetId(ids[index])
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if documents.IsNil() {
		// 正常空命中始终交付非 nil 空切片，包括自定义解码器返回 nil 的情况。
		documents.Set(reflect.MakeSlice(documents.Type(), 0, 0))
	}
	value.Elem().Set(documents)
	return nil
}
