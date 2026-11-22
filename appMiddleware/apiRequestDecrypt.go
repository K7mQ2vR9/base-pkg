package appMiddleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/K7mQ2vR9/base-pkg/aesGCM"
	"github.com/K7mQ2vR9/base-pkg/consts"
	"github.com/K7mQ2vR9/base-pkg/result"
	"github.com/K7mQ2vR9/base-pkg/result/xcode"
)

var RequestDecryptError error = xcode.New(http.StatusBadRequest, "Request decryption failed. ")
var RequestBadError error = xcode.New(http.StatusBadRequest, "Request forbid. ")

// ApiRequestDecryptMiddleware 按当前进程的配置开关解密请求，并对独立使用时的读取施加上限。
type ApiRequestDecryptMiddleware struct {
	maxBodyBytes int64
}

// NewApiRequestDecryptMiddleware 校验并应用进程级配置，默认关闭 AES-GCM。
// 可选 bodyConfig 最多一项，省略时请求体上限默认 1 MiB；负数拒绝。
func NewApiRequestDecryptMiddleware(config aesGCM.Config, bodyConfig ...RequestBodyConfig) (*ApiRequestDecryptMiddleware, error) {
	if len(bodyConfig) > 1 {
		return nil, errors.New("only one request body configuration is allowed")
	}
	var body RequestBodyConfig
	if len(bodyConfig) == 1 {
		body = bodyConfig[0]
	}
	limiter, err := NewRequestBodyMiddleware(body)
	if err != nil {
		return nil, err
	}
	if err := aesGCM.Configure(config); err != nil {
		return nil, err
	}
	return &ApiRequestDecryptMiddleware{maxBodyBytes: limiter.maxBytes}, nil
}

// Handle 按配置开关解密请求，失败时返回统一错误响应。
func (m *ApiRequestDecryptMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := m.RequestDecrypt(r); err != nil {
			var readErr *requestBodyReadError
			if errors.As(err, &readErr) {
				WriteRequestBodyError(w, r, err)
				return
			}
			result.HttpErrorResult(r.Context(), w, err)
			return
		}
		next(w, r)
	}
}

// RequestDecrypt 处理 aes_data 封装，解密成功后重建请求体及其长度。
// 非空请求体带有 business-code 时要求提供 aes_data，其他普通请求可透传；旧头名不再读取。
// 复用前置缓存或按配置限量读取，读取失败时不向下游提供部分正文。
func (m *ApiRequestDecryptMiddleware) RequestDecrypt(r *http.Request) error {
	cipher := aesGCM.ConfiguredCipher()
	if cipher == nil || r.Body == nil {
		return nil
	}
	data, err := readRequestBody(nil, r, m.maxBodyBytes)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	var envelope struct {
		AESData json.RawMessage `json:"aes_data"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil || len(envelope.AESData) == 0 {
		if r.Header.Get(consts.BusinessCode) != "" {
			return RequestBadError
		}
		return nil
	}

	var ciphertext string
	if err := json.Unmarshal(envelope.AESData, &ciphertext); err != nil || ciphertext == "" {
		return RequestDecryptError
	}
	plaintext, err := cipher.Decrypt(ciphertext)
	if err != nil {
		return RequestDecryptError
	}
	replaceRequestBody(r, plaintext)
	return nil
}

// replaceRequestBody 为后续解析和请求重读提供相同的数据，更新长度并清除原传输编码。
func replaceRequestBody(r *http.Request, body []byte) {
	r.GetBody = func() (io.ReadCloser, error) {
		return &bufferedRequestBody{Reader: bytes.NewReader(body), data: body}, nil
	}
	r.Body, _ = r.GetBody()
	r.ContentLength = int64(len(body))
	r.TransferEncoding = nil
	if r.Header != nil {
		r.Header.Set("Content-Length", strconv.Itoa(len(body)))
		r.Header.Del("Transfer-Encoding")
	}
}
