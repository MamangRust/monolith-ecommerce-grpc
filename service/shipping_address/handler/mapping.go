package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbshipping_address "github.com/MamangRust/monolith-ecommerce-pb/shipping_address"
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
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoShippingResponse(shipping interface{}) *pbshipping_address.ShippingResponse {
	switch s := shipping.(type) {
	case *db.ShippingAddress:
		return &pbshipping_address.ShippingResponse{
			Id:             int32(s.ShippingAddressID),
			OrderId:        int32(s.OrderID),
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      s.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      s.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetShippingAddressRow:
		return &pbshipping_address.ShippingResponse{
			Id:             int32(s.ShippingAddressID),
			OrderId:        int32(s.OrderID),
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      s.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      s.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetShippingByIDRow:
		return &pbshipping_address.ShippingResponse{
			Id:             int32(s.ShippingAddressID),
			OrderId:        int32(s.OrderID),
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      s.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      s.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetShippingAddressByOrderIDRow:
		return &pbshipping_address.ShippingResponse{
			Id:             int32(s.ShippingAddressID),
			OrderId:        int32(s.OrderID),
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      s.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      s.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.CreateShippingAddressRow:
		return &pbshipping_address.ShippingResponse{
			Id:             int32(s.ShippingAddressID),
			OrderId:        int32(s.OrderID),
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      s.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      s.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.UpdateShippingAddressRow:
		return &pbshipping_address.ShippingResponse{
			Id:             int32(s.ShippingAddressID),
			OrderId:        int32(s.OrderID),
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      s.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      s.UpdatedAt.Time.Format("2006-01-02"),
		}
	default:
		return nil
	}
}

func mapToProtoShippingResponseDeleteAt(shipping interface{}) *pbshipping_address.ShippingResponseDeleteAt {
	switch s := shipping.(type) {
	case *db.ShippingAddress:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt.Valid {
			deletedAt = wrapperspb.String(s.DeletedAt.Time.Format("2006-01-02"))
		}
		return &pbshipping_address.ShippingResponseDeleteAt{
			Id:             int32(s.ShippingAddressID),
			OrderId:        int32(s.OrderID),
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      s.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      s.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:      deletedAt,
		}
	case *db.GetShippingAddressActiveRow:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt.Valid {
			deletedAt = wrapperspb.String(s.DeletedAt.Time.Format("2006-01-02"))
		}
		return &pbshipping_address.ShippingResponseDeleteAt{
			Id:             int32(s.ShippingAddressID),
			OrderId:        int32(s.OrderID),
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      s.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      s.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:      deletedAt,
		}
	case *db.GetShippingAddressTrashedRow:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt.Valid {
			deletedAt = wrapperspb.String(s.DeletedAt.Time.Format("2006-01-02"))
		}
		return &pbshipping_address.ShippingResponseDeleteAt{
			Id:             int32(s.ShippingAddressID),
			OrderId:        int32(s.OrderID),
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      s.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:      s.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:      deletedAt,
		}
	default:
		return nil
	}
}
