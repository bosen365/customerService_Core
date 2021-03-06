package handle

import (
	"customerService_Core/common"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 管理员后台鉴权
func AdminOauthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err   error
			kfId  string
			token = c.Request.Header.Get("Authentication")
		)

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "API token required",
			})
			c.Abort()
			return
		}

		if kfId, err = AdminAuthToken2Model(c.Request.Header.Get("Authentication")); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		} else {
			c.Set("KFID", kfId)
		}

		c.Next()
	}
}

// 鉴权Token解码为模型
func AdminAuthToken2Model(token string) (kfId string, err error) {
	var (
		aes = common.AesEncrypt{}
	)
	decodeBytes, err := base64.StdEncoding.DecodeString(token)
	if bytes, err := aes.Decrypt(decodeBytes); err != nil {
		err = errors.New("API token Authentication failed")
	} else {
		kfId = string(bytes)
		if kfId == "" {
			err = errors.New("API token Authentication failed")
		}
	}
	return
}
