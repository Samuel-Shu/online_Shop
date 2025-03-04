package message

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"online_Shop_api/userop_web/api"
	"online_Shop_api/userop_web/forms"
	"online_Shop_api/userop_web/global"
	"online_Shop_api/userop_web/middleware"
	"online_Shop_api/userop_web/proto"
)

func New(c *gin.Context) {
	userId, _ := c.Get("userId")
	messageForm := forms.MessageForm{}
	if err := c.ShouldBindJSON(&messageForm); err != nil {
		api.HandleValidatorError(c, err)
		return
	}

	rsp, err := global.MessageSrvClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      userId.(int32),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id

	c.JSON(http.StatusOK, request)
}

func List(c *gin.Context) {
	//如果是管理员用户，则返回所有的订单
	request := &proto.MessageRequest{}
	userId, _ := c.Get("userId")
	claims, _ := c.Get("claims")
	model := claims.(*middleware.MyClaims)
	if model.AuthorityId == 1 {
		request.UserId = userId.(int32)
	}

	rsp, err := global.MessageSrvClient.MessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		api.HandleGrpcErrorToHttp(err, c)
		return
	}

	reMap := map[string]interface{}{
		"total": rsp.Total,
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["user_id"] = value.UserId
		reMap["type"] = value.MessageType
		reMap["subject"] = value.Subject
		reMap["message"] = value.Message
		reMap["file"] = value.File

		result = append(result, reMap)
	}

	reMap["data"] = result

	c.JSON(http.StatusOK, reMap)
}
