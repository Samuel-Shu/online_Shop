package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"online_Shop_api/llms_web/global"
	"online_Shop_api/llms_web/proto"
)

func HealthCheck(c *gin.Context) {
	check, err := global.LlmsSrvClient.HealthCheck(context.Background(), &proto.EmptyWithLlms{})
	if err != nil {
		zap.S().Errorw("健康检查出错", err)
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": check.Success,
		"msg":    check.Message,
	})
}

func Chat(c *gin.Context) {
	type chat struct {
		Sender   string
		MetaData map[string]string
		Content  string
	}

	var request chat
	if err := c.Bind(&request); err != nil {
		zap.S().Errorw("【Chat-api】解析请求参数错误", err)
		HandleValidatorError(c, err)
		return
	}

	rsp, err := global.LlmsSrvClient.SendMessage(context.Background(), &proto.ChatMessageRequest{
		Sender:   request.Sender,
		Metadata: request.MetaData,
		Content:  request.Content,
	})
	if err != nil {
		zap.S().Errorw("发送聊天消息失败", err)
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": rsp.Success,
		"msg":    rsp.Message,
	})

}

func UploadFile(c *gin.Context) {
	fileName, _ := c.GetPostForm("file_name")

	formFile, err := c.FormFile("file_content")
	if err != nil {
		zap.S().Errorw("获取文件数据失败", err)
		HandleValidatorError(c, err)
		return
	}

	open, err := formFile.Open()
	if err != nil {
		HandleValidatorError(c, err)
		return
	}

	defer func(open multipart.File) {
		err := open.Close()
		if err != nil {
			HandleValidatorError(c, err)
			return
		}
	}(open)

	all, err := io.ReadAll(open)
	if err != nil {
		HandleValidatorError(c, err)
		return
	}
	rsp, err := global.LlmsSrvClient.UploadFile(context.Background(), &proto.UploadFileRequest{
		Filename:    fileName,
		FileContent: all,
	})
	if err != nil {
		zap.S().Errorw("【UploadFile-srv】上传文件失败", err)
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": rsp.Success,
		"msg":    rsp.Message,
	})
}
