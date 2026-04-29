package orderitemapimapper

type OrderItemResponseMapper interface {
	QueryMapper() OrderItemQueryResponseMapper
	CommandMapper() OrderItemCommandResponseMapper
}

type orderItemResponseMapper struct {
	queryMapper   OrderItemQueryResponseMapper
	commandMapper OrderItemCommandResponseMapper
}

func NewOrderItemResponseMapper() OrderItemResponseMapper {
	return &orderItemResponseMapper{
		queryMapper:   NewOrderItemQueryResponseMapper(),
		commandMapper: NewOrderItemCommandResponseMapper(),
	}
}

func (m *orderItemResponseMapper) QueryMapper() OrderItemQueryResponseMapper {
	return m.queryMapper
}

func (m *orderItemResponseMapper) CommandMapper() OrderItemCommandResponseMapper {
	return m.commandMapper
}
