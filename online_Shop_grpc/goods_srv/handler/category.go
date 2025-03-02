package handler

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"online_Shop/goods_srv/global"
	"online_Shop/goods_srv/model"
	"online_Shop/goods_srv/proto"
)

// 商品分类

func (g *GoodsServer) GetAllCategoryList(c context.Context, req *proto.MyEmpty) (*proto.CategoryListResponse, error) {
	var category []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&category)
	marshal, err := json.Marshal(&category)
	if err != nil {
		zap.S().Errorf("[GetAllCategoryList]获取所有商品分类时序列化失败！%s", err)
	}
	return &proto.CategoryListResponse{JsonData: string(marshal)}, nil
}

//  获取子分类

// GetSubCategory 根据商品分类id获取子分类
func (g *GoodsServer) GetSubCategory(c context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var category model.Category
	var categoryListResponse proto.SubCategoryListResponse
	if r := global.DB.First(&category, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             int32(category.ID),
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: int32(category.ParentCategoryID),
	}

	preloads := "SubCategory"
	if category.Level == 1 {
		preloads = "SubCategory.SubCategory"
	}

	var subCategorys []*model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	global.DB.Where(&model.Category{ParentCategoryID: uint(req.Id)}).Preload(preloads).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             int32(subCategory.ID),
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: int32(subCategory.ParentCategoryID),
		})
	}
	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse, nil
}

func (g *GoodsServer) CreateCategory(c context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}

	category.Name = req.Name
	category.Level = req.Level
	if req.Level != 1 {
		category.ParentCategoryID = uint(req.ParentCategory)
	}
	category.IsTab = req.IsTab
	global.DB.Save(&category)
	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}
func (g *GoodsServer) UpdateCategory(c context.Context, req *proto.CategoryInfoRequest) (*proto.MyEmpty, error) {
	var category model.Category

	if r := global.DB.First(&category, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = uint(req.ParentCategory)
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &proto.MyEmpty{}, nil
}

func (g *GoodsServer) DeleteCategory(c context.Context, req *proto.DeleteCategoryRequest) (*proto.MyEmpty, error) {
	if r := global.DB.Delete(&model.Category{}, req.Id); r.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &proto.MyEmpty{}, nil
}
