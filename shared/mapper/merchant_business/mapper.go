package merchantbusinessapimapper

type MerchantBusinessResponseMapper interface {
	QueryMapper() MerchantBusinessQueryResponseMapper
	CommandMapper() MerchantBusinessCommandResponseMapper
}

type merchantBusinessResponseMapper struct {
	queryMapper   MerchantBusinessQueryResponseMapper
	commandMapper MerchantBusinessCommandResponseMapper
}

func NewMerchantBusinessResponseMapper() MerchantBusinessResponseMapper {
	return &merchantBusinessResponseMapper{
		queryMapper:   NewMerchantBusinessQueryResponseMapper(),
		commandMapper: NewMerchantBusinessCommandResponseMapper(),
	}
}

func (m *merchantBusinessResponseMapper) QueryMapper() MerchantBusinessQueryResponseMapper {
	return m.queryMapper
}

func (m *merchantBusinessResponseMapper) CommandMapper() MerchantBusinessCommandResponseMapper {
	return m.commandMapper
}
