package shippingaddressapimapper

type ShippingAddressResponseMapper interface {
	QueryMapper() ShippingAddressQueryResponseMapper
	CommandMapper() ShippingAddressCommandResponseMapper
}

type shippingAddressResponseMapper struct {
	queryMapper   ShippingAddressQueryResponseMapper
	commandMapper ShippingAddressCommandResponseMapper
}

func NewShippingAddressResponseMapper() ShippingAddressResponseMapper {
	return &shippingAddressResponseMapper{
		queryMapper:   NewShippingAddressQueryResponseMapper(),
		commandMapper: NewShippingAddressCommandResponseMapper(),
	}
}

func (m *shippingAddressResponseMapper) QueryMapper() ShippingAddressQueryResponseMapper {
	return m.queryMapper
}

func (m *shippingAddressResponseMapper) CommandMapper() ShippingAddressCommandResponseMapper {
	return m.commandMapper
}
