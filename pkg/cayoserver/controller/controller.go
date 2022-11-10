package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/model/alarm"
	"github.com/yunling101/cayo/pkg/types"
)

// BaseController
type BaseController struct{}

// render
func (w *BaseController) render(c *gin.Context, v interface{}) {
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.JSON(200, v)
}

// RenderDieIfError
func (w *BaseController) RenderDieError(c *gin.Context, err error, val interface{}) {
	w.RenderFail(c, val)
	log.Printf("logf path: %s error: %s", c.Request.RequestURI, err.Error())
}

// RenderFail
func (w *BaseController) RenderFail(c *gin.Context, val interface{}) {
	w.render(c, gin.H{"msg": val, "code": 0})
}

// RenderIfFail
func (w *BaseController) RenderIfFail(c *gin.Context, err error) bool {
	if err != nil {
		w.RenderFail(c, err.Error())
		return false
	}
	return true
}

// RenderSuccess
func (w *BaseController) RenderSuccess(c *gin.Context, data interface{}) {
	w.render(c, gin.H{"msg": "success", "code": 1, "data": data})
}

// SetQuery
// https://www.alexedwards.net/blog/change-url-query-params-in-go
// c.Params = []gin.Param{}
func (w *BaseController) SetQuery(c *gin.Context, k string, val string) {
	paramPairs := c.Request.URL.Query()
	paramPairs.Set(k, val)
	c.Request.URL.RawQuery = paramPairs.Encode()
}

// Params
func (w *BaseController) Params(c *gin.Context, args []types.Param) bool {
	// fmt.Println(fmt.Printf("%T", t))
	assert := types.Assert{Result: true}
	for i := 0; i < len(args); i++ {
		assert.Index = i
		switch args[i].Data.(type) {
		case string:
			if args[i].Data.(string) == "" {
				assert.Result = false
				break
			}
		case int:
			if args[i].Data.(int) == 0 {
				assert.Result = false
				break
			}
		case []int:
			if len(args[i].Data.([]int)) == 0 {
				assert.Result = false
				break
			}
		case []alarm.AlarmRule:
			if len(args[i].Data.([]alarm.AlarmRule)) == 0 {
				assert.Result = false
				break
			}
		}
		if !assert.Result {
			break
		}
	}
	if !assert.Result {
		w.RenderFail(c, args[assert.Index].Desc+"不能为空!")
	}
	return assert.Result
}

// BindJSON
func (w *BaseController) BindJSON(c *gin.Context, val interface{}) bool {
	if err := c.ShouldBindJSON(&val); err != nil {
		w.RenderDieError(c, err, "参数解析出错!")
		return false
	}
	return true
}
