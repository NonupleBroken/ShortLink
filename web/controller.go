package web

import (
	"ShortLink/config"
	"ShortLink/db/redis"
	"ShortLink/logger"
	"ShortLink/util"
	"github.com/gin-gonic/gin"
	"time"
)

func adminAuth(serverConfig *config.ServerConfig) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		serverConfig.Username: serverConfig.Password,
	})
}

func success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func fail(c *gin.Context, statusCode int, code int, msg string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func systemError(c *gin.Context) {
	c.JSON(500, gin.H{
		"code": -1,
		"msg":  "系统错误",
		"data": nil,
	})
}

func index(c *gin.Context) {
	success(c, nil)
}

func getShortLinkByLinkID(c *gin.Context) {
	var link Link
	err := c.ShouldBindUri(&link)
	if err != nil {
		fail(c, 400, 1, "参数错误", nil)
		return
	}

	isExists, err := redis.Client.HExists(redis.Ctx, "short_link", link.LinkID).Result()
	if err != nil {
		logger.S.Error("delete link redis error: ", err)
		systemError(c)
		return
	}
	if !isExists {
		fail(c, 200, 1, "短链接不存在", nil)
		return
	}
	result, err := redis.Client.HGet(redis.Ctx, "short_link", link.LinkID).Result()
	if err != nil {
		logger.S.Error("get link redis error: ", err)
		systemError(c)
		return
	}
	c.Redirect(302, result)

}

func addLink(c *gin.Context) {
	var link AddLinkParams
	err := c.BindJSON(&link)
	if err != nil {
		fail(c, 400, 1, "参数错误", nil)
		return
	}
	if link.ShortLinkLength < 1 || link.ShortLinkLength > 32 {
		fail(c, 400, 1, "短链接长度有误", nil)
		return
	}
	if len(link.LinkContent) > 4096 {
		fail(c, 400, 1, "链接长度过长", nil)
		return
	}

	shortLink := util.GetRandomStr(link.ShortLinkLength)

	mutex := redis.NewMutex(shortLink, time.Second * 5)
	mutex.Lock()
	defer mutex.UnLock()

	isExists, err := redis.Client.HExists(redis.Ctx, "short_link", shortLink).Result()
	if err != nil {
		logger.S.Error("add link redis error: ", err)
		systemError(c)
		return
	}
	if isExists {
		fail(c, 200, 1, "短链接已满", nil)
		return
	}
	_, err = redis.Client.HSet(redis.Ctx, "short_link", shortLink, link.LinkContent).Result()
	if err != nil {
		logger.S.Error("add link redis error: ", err)
		systemError(c)
		return
	}

	success(c, gin.H{"link_id": shortLink})
}

func deleteLink(c *gin.Context) {
	var link DeleteLinkParams
	err := c.BindJSON(&link)
	if err != nil {
		fail(c, 400, 1, "参数错误", nil)
		return
	}
	if len(link.LinkID) > 32 {
		fail(c, 400, 1, "短链接长度有误", nil)
		return
	}

	mutex := redis.NewMutex(link.LinkID, time.Second * 5)
	mutex.Lock()
	defer mutex.UnLock()

	isExists, err := redis.Client.HExists(redis.Ctx, "short_link", link.LinkID).Result()
	if err != nil {
		logger.S.Error("delete link redis error: ", err)
		systemError(c)
		return
	}
	if !isExists {
		fail(c, 200, 1, "短链接不存在", nil)
		return
	}
	_, err = redis.Client.HDel(redis.Ctx, "short_link", link.LinkID).Result()
	if err != nil {
		logger.S.Error("delete link redis error: ", err)
		systemError(c)
		return
	}

	success(c, nil)
}

func checkLink(c *gin.Context) {
	result, err := redis.Client.HGetAll(redis.Ctx, "short_link").Result()
	if err != nil {
		logger.S.Error("delete link redis error: ", err)
		systemError(c)
		return
	}
	var shortLinks []CheckShortLink
	for k, v := range result {
		link := CheckShortLink{LinkID: k, Content: v}
		shortLinks = append(shortLinks, link)
	}
	success(c, gin.H{"short_links": shortLinks})
}
