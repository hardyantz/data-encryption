package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hardyantz/data-encryption/config"
	"github.com/hardyantz/data-encryption/conn"
	"github.com/hardyantz/data-encryption/pkg/member/domain"
	memberRepository "github.com/hardyantz/data-encryption/pkg/member/repository"
	memberUseCase "github.com/hardyantz/data-encryption/pkg/member/usecase"
	"github.com/labstack/echo"
)

func BenchmarkHTTPHandler_GetAll(b *testing.B) {
	db := conn.ConnectDB()
	conf := config.NewConfImpl()

	cacheExpire := 1 * time.Hour

	redis := config.NewCacheRedis(conn.RedisConnect(), cacheExpire)

	memberRepo := memberRepository.NewMemberRepoImplementation(db, conf, redis)
	memberUc := memberUseCase.NewMemberImplementation(memberRepo)
	h := NewHTTPHandler(memberUc)

	bodyData, _ := json.Marshal(domain.Member{})

	e := echo.New()
	req, _ := http.NewRequest(http.MethodGet, "member", strings.NewReader(string(bodyData)))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	b.ResetTimer()
	for i := 0; i < b.N; i ++ {
		_ = h.GetOne(c)
	}
}
