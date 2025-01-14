package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"loginTest/common"
	"loginTest/config"
	"loginTest/route"
	"net/http/httptest"
	"testing"
)

func ParseToStr(mp map[string]string) string {
	values := ""
	for key, val := range mp {
		values += "&" + key + "=" + val
	}
	temp := values[1:]
	values = "?" + temp
	return values
}

func PostForm(uri string, param map[string]string, router *gin.Engine) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", uri+ParseToStr(param), nil)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应handler接口
	router.ServeHTTP(w, req)
	return w
}

var r *gin.Engine

func Init() {
	config.InitConfig()
	r = gin.Default()
	r = route.CollectRoute(r)
}

func TestModifyPassword(t *testing.T) {
	Init()
	db := common.InitDB()
	rds := common.RedisInit()

	var w *httptest.ResponseRecorder

	param := make(map[string]string)
	param["Phone"] = "13822209685"
	param["Password"] = "123455"
	param["Password2"] = "123455"

	urlImport := "/api/auth/modifyPassword"
	w = PostForm(urlImport, param, r)
	assert.Equal(t, 200, w.Code)

	defer db.Close()
	defer rds.Close()
}
