package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func createPaginationMeta(page, pageSize, totalRecords int) *pb.PaginationMeta {
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoShippingResponse(shipping interface{}) *pb.ShippingResponse {
	switch s := shipping.(type) {
	case *db.ShippingAddress:
		return &pb.ShippingResponse{
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
			Courier:        s.Courier,
		}
	case *db.GetShippingAddressRow:
		return &pb.ShippingResponse{
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
			Courier:        s.Courier,
		}
	case *db.GetShippingByIDRow:
		return &pb.ShippingResponse{
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
			Courier:        s.Courier,
		}
	case *db.GetShippingAddressByOrderIDRow:
		return &pb.ShippingResponse{
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
			Courier:        s.Courier,
		}
	case *db.CreateShippingAddressRow:
		return &pb.ShippingResponse{
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
			Courier:        s.Courier,
		}
	case *db.UpdateShippingAddressRow:
		return &pb.ShippingResponse{
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
			Courier:        s.Courier,
		}
	default:
		return nil
	}
}

func mapToProtoShippingResponseDeleteAt(shipping interface{}) *pb.ShippingResponseDeleteAt {
	switch s := shipping.(type) {
	case *db.ShippingAddress:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt.Valid {
			deletedAt = wrapperspb.String(s.DeletedAt.Time.Format("2006-01-02"))
		}
		return &pb.ShippingResponseDeleteAt{
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
			Courier:        s.Courier,
		}
	case *db.GetShippingAddressActiveRow:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt.Valid {
			deletedAt = wrapperspb.String(s.DeletedAt.Time.Format("2006-01-02"))
		}
		return &pb.ShippingResponseDeleteAt{
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
			Courier:        s.Courier,
		}
	case *db.GetShippingAddressTrashedRow:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt.Valid {
			deletedAt = wrapperspb.String(s.DeletedAt.Time.Format("2006-01-02"))
		}
		return &pb.ShippingResponseDeleteAt{
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
			Courier:        s.Courier,
		}
	default:
		return nil
	}
}
