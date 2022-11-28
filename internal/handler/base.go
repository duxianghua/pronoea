package handler

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

/**
 * 基础的返回结构，有点类似java定义的统一返回Result对象。
 * Json格式是：
 * {
 *   "code":"200",
 *   "msg":"OK",
 *   "data":...,
 * }
 **/
type Result struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

/**
 * @description:  获取返回结构对象，不对外暴露
 * @param {string} code 返回码
 * @param {string} msg 返回消息
 * @param {interface{}} data 返回数据
 * @return {*}
 */
func getResult(code string, msg string, data interface{}) Result {
	var result = &Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	return *result
}

/**
 * @description: 获取返回结构Map，不对外暴露
 * @param {string} code 返回码
 * @param {interface{}} msg 返回消息
 * @param {interface{}} data 返回数据
 * @return {*}
 */
// func getResultMap(code string, msg interface{}, data interface{}) map[string]interface{} {
// 	result := make(map[string]interface{})
// 	result["code"] = code
// 	result["msg"] = msg
// 	result["data"] = data
// 	return result
// }

func error(c *gin.Context, errcode string, msg string, data interface{}) {
	c.JSON(http.StatusInternalServerError, getResult(errcode, msg, data))
}

func ErrorHandler(c *gin.Context) {
	c.Next()
	for _, err := range c.Errors {
		// log, handle, etc.
		log.Error().Msg(err.Error())
		if c.Writer.Status() == http.StatusOK {
			c.Status(http.StatusInternalServerError)
		}
		c.AbortWithStatusJSON(
			-1,
			Result{
				Code: "-1",
				Msg:  err.Error(),
				Data: nil,
			})
	}

}

func RedirectHomeHandler(c *gin.Context) {
	c.Redirect(302, "/dashboard")
}

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			//打印错误堆栈信息
			log.Error().Interface("error", err)
			debug.PrintStack()
			//封装通用json返回
			error(c, "-1", "call api error!", err)
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
