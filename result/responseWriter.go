package result

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/K7mQ2vR9/base-pkg/language"
	"github.com/zeromicro/go-zero/core/logc"
)

// writeJSONResponse 先完整编码再提交响应；用户自定义 MarshalJSON 只随本次编码执行，不做预编码再重试。
// 编码失败后只编码固定的内部错误结构，不再次执行失败数据的序列化。
func writeJSONResponse(ctx context.Context, w http.ResponseWriter, status int, body any) {
	payload, err := json.Marshal(body)
	if err != nil {
		status, fallback := internalResponseError(ctx, fmt.Errorf("encode JSON response: %w", err))
		// fallback 由本包创建，只包含数字、字符串和空对象，不含可失败的自定义序列化器。
		fallbackPayload, _ := json.Marshal(fallback)
		writeResponseBytes(ctx, w, status, fallbackPayload)
		return
	}
	writeResponseBytes(ctx, w, status, payload)
}

// writeInternalResponseError 处理 AES 明文编码或加密失败，不沿错误链采用调用方携带的业务码。
func writeInternalResponseError(ctx context.Context, w http.ResponseWriter, err error) {
	status, body := internalResponseError(ctx, err)
	writeJSONResponse(ctx, w, status, body)
}

// internalResponseError 保持 HTTP 200/code500 的既有内部错误协议及当前请求语言、Trace。
// 底层诊断仅写入服务端日志，公开响应不包含错误文本、原始数据或密文。
func internalResponseError(ctx context.Context, err error) (int, *ResponseBean) {
	logc.Errorf(ctx, "HTTP response preparation failed: %v", err)
	message := language.TranslateMessage(http.StatusText(http.StatusInternalServerError), headInfo.GetContentLanguage(ctx))
	return http.StatusOK, Error(http.StatusInternalServerError, message, headInfo.GetTrace(ctx))
}

// writeResponseBytes 提交已经编码的 JSON；任何写入错误或短写仅记录一次，不重发或拼接另一份响应。
// 调用方须在写入其他正文之前选择一个响应出口；本函数不接管任意上游已经提交的响应。
func writeResponseBytes(ctx context.Context, w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	n, err := w.Write(payload)
	if err != nil {
		logc.Errorf(ctx, "HTTP response write failed: %v", err)
		return
	}
	if n < len(payload) {
		logc.Errorf(ctx, "HTTP response short write: expected %d bytes, wrote %d", len(payload), n)
	}
}
