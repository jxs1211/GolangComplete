package httptest_demo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_helloHandler(t *testing.T) {
	r := SetupRouter()
	tests := []struct {
		name   string
		param  string
		expect string
	}{
		// TODO: Add test cases.
		{"base case", `{"name": "liwenzhou"}`, "hello liwenzhou"},
		{"bad case", "", "we need a name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				"POST",                      // 请求方法
				"/hello",                    // 请求URL
				strings.NewReader(tt.param), // 请求参数
			)
			// mock一个响应记录器
			w := httptest.NewRecorder()
			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)
			// 校验状态码是否符合预期
			assert.Equal(t, http.StatusOK, w.Code)
			// 解析并检验响应内容是否复合预期
			// var resp map[string]string
			// err := json.Unmarshal([]byte(w.Body.String()))
		})
	}
}
