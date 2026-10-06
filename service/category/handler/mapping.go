package handler

import (
	"encoding/json"
	"log"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
)

func (h *Handler) mapToCategoryResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *db.GetCategoryByIDRow:
		return &pbcategory.CategoryResponse{
			Id:            int32(v.CategoryID),
			Name:          v.Name,
			Description:   *v.Description,
			SlugCategory:  *v.SlugCategory,
			ImageCategory: *v.ImageCategory,
			CreatedAt:     v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetCategoriesRow:
		return &pbcategory.CategoryResponse{
			Id:            int32(v.CategoryID),
			Name:          v.Name,
			Description:   *v.Description,
			SlugCategory:  *v.SlugCategory,
			ImageCategory: *v.ImageCategory,
			CreatedAt:     v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetCategoriesActiveRow:
		var deletedAt string
		if v.DeletedAt.Valid {
			deletedAt = v.DeletedAt.Time.Format("2006-01-02")
		}
		return &pbcategory.CategoryResponseDeleteAt{
			Id:            int32(v.CategoryID),
			Name:          v.Name,
			Description:   *v.Description,
			SlugCategory:  *v.SlugCategory,
			ImageCategory: *v.ImageCategory,
			CreatedAt:     v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
		}
	case *db.GetCategoriesTrashedRow:
		var deletedAt string
		if v.DeletedAt.Valid {
			deletedAt = v.DeletedAt.Time.Format("2006-01-02")
		}
		return &pbcategory.CategoryResponseDeleteAt{
			Id:            int32(v.CategoryID),
			Name:          v.Name,
			Description:   *v.Description,
			SlugCategory:  *v.SlugCategory,
			ImageCategory: *v.ImageCategory,
			CreatedAt:     v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
		}
	case *db.CreateCategoryRow:
		return &pbcategory.CategoryResponse{
			Id:            int32(v.CategoryID),
			Name:          v.Name,
			Description:   *v.Description,
			SlugCategory:  *v.SlugCategory,
			ImageCategory: *v.ImageCategory,
			CreatedAt:     v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.UpdateCategoryRow:
		return &pbcategory.CategoryResponse{
			Id:            int32(v.CategoryID),
			Name:          v.Name,
			Description:   *v.Description,
			SlugCategory:  *v.SlugCategory,
			ImageCategory: *v.ImageCategory,
			CreatedAt:     v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.Category:
		var deletedAt string
		if v.DeletedAt.Valid {
			deletedAt = v.DeletedAt.Time.Format("2006-01-02")
		}
		return &pbcategory.CategoryResponseDeleteAt{
			Id:            int32(v.CategoryID),
			Name:          v.Name,
			Description:   *v.Description,
			SlugCategory:  *v.SlugCategory,
			ImageCategory: *v.ImageCategory,
			CreatedAt:     v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     v.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
		}
	case *db.GetMonthlyTotalPriceRow:
		return &pbcategory.CategoriesMonthlyTotalPriceResponse{
			Year:         v.Year,
			Month:        v.Month,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetMonthlyTotalPriceByIdRow:
		return &pbcategory.CategoriesMonthlyTotalPriceResponse{
			Year:         v.Year,
			Month:        v.Month,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetMonthlyTotalPriceByMerchantRow:
		return &pbcategory.CategoriesMonthlyTotalPriceResponse{
			Year:         v.Year,
			Month:        v.Month,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetYearlyTotalPriceRow:
		return &pbcategory.CategoriesYearlyTotalPriceResponse{
			Year:         v.Year,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetYearlyTotalPriceByIdRow:
		return &pbcategory.CategoriesYearlyTotalPriceResponse{
			Year:         v.Year,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetYearlyTotalPriceByMerchantRow:
		return &pbcategory.CategoriesYearlyTotalPriceResponse{
			Year:         v.Year,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetMonthlyCategoryRow:
		return &pbcategory.CategoryMonthPriceResponse{
			Month:        v.Month,
			CategoryId:   int32(v.CategoryID),
			CategoryName: v.CategoryName,
			OrderCount:   int32(v.OrderCount),
			ItemsSold:    int32(v.ItemsSold),
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetMonthlyCategoryByMerchantRow:
		return &pbcategory.CategoryMonthPriceResponse{
			Month:        v.Month,
			CategoryId:   int32(v.CategoryID),
			CategoryName: v.CategoryName,
			OrderCount:   int32(v.OrderCount),
			ItemsSold:    int32(v.ItemsSold),
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetMonthlyCategoryByIdRow:
		return &pbcategory.CategoryMonthPriceResponse{
			Month:        v.Month,
			CategoryId:   int32(v.CategoryID),
			CategoryName: v.CategoryName,
			OrderCount:   int32(v.OrderCount),
			ItemsSold:    int32(v.ItemsSold),
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetYearlyCategoryRow:
		return &pbcategory.CategoryYearPriceResponse{
			Year:               v.Year,
			CategoryId:         int32(v.CategoryID),
			CategoryName:       v.CategoryName,
			OrderCount:         int32(v.OrderCount),
			ItemsSold:          int32(v.ItemsSold),
			TotalRevenue:       int32(v.TotalRevenue),
			UniqueProductsSold: int32(v.UniqueProductsSold),
		}
	case *db.GetYearlyCategoryByMerchantRow:
		return &pbcategory.CategoryYearPriceResponse{
			Year:               v.Year,
			CategoryId:         int32(v.CategoryID),
			CategoryName:       v.CategoryName,
			OrderCount:         int32(v.OrderCount),
			ItemsSold:          int32(v.ItemsSold),
			TotalRevenue:       int32(v.TotalRevenue),
			UniqueProductsSold: int32(v.UniqueProductsSold),
		}
	case *db.GetYearlyCategoryByIdRow:
		return &pbcategory.CategoryYearPriceResponse{
			Year:               v.Year,
			CategoryId:         int32(v.CategoryID),
			CategoryName:       v.CategoryName,
			OrderCount:         int32(v.OrderCount),
			ItemsSold:          int32(v.ItemsSold),
			TotalRevenue:       int32(v.TotalRevenue),
			UniqueProductsSold: int32(v.UniqueProductsSold),
		}
	default:
		log.Printf("Unknown type for mapping: %T", v)
		return nil
	}
}

func (h *Handler) mapToPayload(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
