package category

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"online_Shop_api/goods_web/forms"
	"online_Shop_api/goods_web/global"
	"online_Shop_api/goods_web/proto"
	"strconv"
	"strings"
)

func RemoveTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转化为http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "内部错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func List(c *gin.Context)  {
	r, err := global.GoodsSrvClient.GetAllCategoryList(context.Background(), &proto.MyEmpty{})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	data := make([]interface{}, 0)
	err = json.Unmarshal([]byte(r.JsonData), &data)
	if err != nil {
		zap.S().Errorw("【List】查询【分类列表】失败：", err.Error())
	}

	c.JSON(http.StatusOK, data)
}

func Detail(c *gin.Context)  {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	reMap := make(map[string]interface{})
	subCategory := make([]interface{}, 0)
	if r, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: int32(i),
	}); err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}else {
		for _, value := range r.SubCategorys {
			subCategory = append(subCategory, map[string]interface{}{
				"id": value.Id,
				"name": value.Name,
				"level": value.Level,
				"parent_category": value.ParentCategory,
				"is_tab": value.IsTab,
			})
		}
		reMap["id"] = r.Info.Id
		reMap["name"] = r.Info.Name
		reMap["level"] = r.Info.Level
		reMap["parent_category"] = r.Info.ParentCategory
		reMap["is_tab"] = r.Info.IsTab
		reMap["sub_categorys"] = subCategory

		c.JSON(http.StatusOK, reMap)
	}
	return
}

func New(c *gin.Context)  {
	categoryForm := forms.CategoryForm{}
	if err := c.ShouldBindJSON(&categoryForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name: categoryForm.Name,
		IsTab: *categoryForm.IsTab,
		Level: categoryForm.Level,
		ParentCategory: categoryForm.ParentCategory,
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["parent"] = rsp.ParentCategory
	request["level"] = rsp.Level
	request["is_tab"] = rsp.IsTab

	c.JSON(http.StatusOK, request)
}

func Delete(c *gin.Context)  {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if _, err := global.GoodsSrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{
		Id: int32(i),
	}); err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func Update(c *gin.Context) {
	categoryForm := forms.UpdateCategoryForm{}
	if err := c.ShouldBindJSON(&categoryForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id: int32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}

	_, err = global.GoodsSrvClient.UpdateCategory(context.Background(), request)
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}
