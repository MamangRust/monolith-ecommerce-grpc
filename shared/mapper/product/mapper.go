package productapimapper

type ProductResponseMapper interface {
	QueryMapper() ProductQueryResponseMapper
	CommandMapper() ProductCommandResponseMapper
}

type productResponseMapper struct {
	queryMapper   ProductQueryResponseMapper
	commandMapper ProductCommandResponseMapper
}

func NewProductResponseMapper() ProductResponseMapper {
	return &productResponseMapper{
		queryMapper:   NewProductQueryResponseMapper(),
		commandMapper: NewProductCommandResponseMapper(),
	}
}

func (m *productResponseMapper) QueryMapper() ProductQueryResponseMapper {
	return m.queryMapper
}

func (m *productResponseMapper) CommandMapper() ProductCommandResponseMapper {
	return m.commandMapper
}
