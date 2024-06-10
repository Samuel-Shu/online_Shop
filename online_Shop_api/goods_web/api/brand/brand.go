package brand

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	pn := c.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := c.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.GoodsSrvClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       int32(pnInt),
		PagePerNums: int32(pSizeInt),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	result := make([]interface{}, 0)
	reMap := make(map[string]interface{})
	reMap["total"] = rsp.Total
	for _, value := range rsp.Data[pnInt:pnInt*pSizeInt+pSizeInt] {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	reMap["data"] = result

	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context)  {
	brandForm := forms.BrandForm{}
	if err := c.ShouldBindJSON(&brandForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["logo"] = rsp.Logo

	c.JSON(http.StatusOK, request)
}

func Delete(c *gin.Context)  {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteBrand(context.Background(), &proto.BrandRequest{
		Id: int32(i),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func Update(c *gin.Context)  {
	brandForm := forms.BrandForm{}
	if err := c.ShouldBindJSON(&brandForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateBrand(context.Background(), &proto.BrandRequest{
		Id: int32(i),
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func GetCategoryBrand(c *gin.Context)  {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: int32(i),
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, rsp.Data)
}

func GetCategoryBrandList(c *gin.Context)  {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	c.JSON(http.StatusOK, result)
}

func CategoryBrandList(c *gin.Context)  {
	rsp, err := global.GoodsSrvClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["category"] = map[string]interface{}{
			"id": value.Category.Id,
			"name": value.Category.Name,
		}
		reMap["brand"] = map[string]interface{}{
			"id": value.Brand.Id,
			"name": value.Brand.Name,
			"logo": value.Brand.Logo,
		}

		result = append(result, reMap)
	}

	reMap["data"] = result

	c.JSON(http.StatusOK, reMap)
}

func NewCategoryBrand(c *gin.Context)  {
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := c.ShouldBindJSON(&categoryBrandForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId: int32(categoryBrandForm.BrandId),
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id

	c.JSON(http.StatusOK, response)
}

func UpdateCategoryBrand(c *gin.Context)  {
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := c.ShouldBindJSON(&categoryBrandForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id:         int32(i),
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId: int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}

func DeleteCategoryBrand(c *gin.Context)  {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{Id: int32(i)})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.Status(http.StatusOK)
}
