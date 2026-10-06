package handler

import (
	"math"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func createPaginationMeta(page, pageSize, totalRecords int) *pbcommon.PaginationMeta {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func formatTimestamp(v interface{}) string {
	switch t := v.(type) {
	case pgtype.Timestamptz:
		if t.Valid {
			return t.Time.Format("2006-01-02 15:04:05.000")
		}
	case pgtype.Timestamp:
		if t.Valid {
			return t.Time.Format("2006-01-02 15:04:05.000")
		}
	}
	return ""
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func int32PtrToInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func float64PtrToFloat32(f *float64) float32 {
	if f == nil {
		return 0
	}
	return float32(*f)
}

func weightPtrToInt32(w interface{}) int32 {
	switch v := w.(type) {
	case *int32:
		if v != nil {
			return *v
		}
	case int32:
		return v
	}
	return 0
}

func mapToProtoProductResponse(item interface{}) *pbproduct.ProductResponse {
	switch v := item.(type) {
	case *db.Product:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetProductsRow:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetProductsByMerchantRow:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetProductsByCategoryNameRow:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetProductByIDRow:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateProductRow:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateProductRow:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.TrashProductRow:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.RestoreProductRow:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateProductCountStockRow:
		return &pbproduct.ProductResponse{
			Id:           int32(v.ProductID),
			CountInStock: int32(v.CountInStock),
		}
	default:
		return nil
	}
}

func mapToProtoProductResponseDeleteAt(item interface{}) *pbproduct.ProductResponseDeleteAt {
	var res *pbproduct.ProductResponseDeleteAt
	var deletedAt interface{}

	switch v := item.(type) {
	case *db.Product:
		res = &pbproduct.ProductResponseDeleteAt{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetProductsActiveRow:
		res = &pbproduct.ProductResponseDeleteAt{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetProductsTrashedRow:
		res = &pbproduct.ProductResponseDeleteAt{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.TrashProductRow:
		res = &pbproduct.ProductResponseDeleteAt{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.RestoreProductRow:
		res = &pbproduct.ProductResponseDeleteAt{
			Id:           int32(v.ProductID),
			MerchantId:   int32(v.MerchantID),
			CategoryId:   int32(v.CategoryID),
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        int32(v.Price),
			CountInStock: int32(v.CountInStock),
			Brand:        stringPtrToString(v.Brand),
			Weight:       weightPtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}

	if val := formatTimestamp(deletedAt); val != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: val}
	}

	return res
}
