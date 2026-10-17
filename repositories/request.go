package repositories

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Body adalah isi request (JSON atau form), setara `await request.req.json()` di Bun.
type Body map[string]any

func GetBody(request *gin.Context) Body {
	body := Body{}

	if strings.HasPrefix(request.ContentType(), "application/json") {
		_ = request.ShouldBindJSON(&body)
		return body
	}

	if err := request.Request.ParseForm(); err == nil {
		for key, values := range request.Request.PostForm {
			if len(values) > 0 {
				body[key] = values[0]
			}
		}
	}

	return body
}

// Has: field ada dan tidak kosong.
func (body Body) Has(key string) bool {
	return !isEmpty(body[key])
}

func (body Body) String(key string) string {
	switch value := body[key].(type) {
	case nil:
		return ""
	case string:
		return value
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return fmt.Sprint(value)
	}
}

// Nullable: nil jika field kosong (untuk kolom NULL).
func (body Body) Nullable(key string) *string {
	if !body.Has(key) {
		return nil
	}
	value := body.String(key)
	return &value
}

func (body Body) Int(key string) int {
	value, err := strconv.Atoi(body.String(key))
	if err != nil {
		return 0
	}
	return value
}

func isEmpty(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(typed) == ""
	case []any:
		return len(typed) == 0
	}
	return false
}
