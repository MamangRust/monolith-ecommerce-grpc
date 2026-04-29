package merchantawardapimapper

type MerchantAwardResponseMapper interface {
	QueryMapper() MerchantAwardQueryResponseMapper
	CommandMapper() MerchantAwardCommandResponseMapper
}

type merchantAwardResponseMapper struct {
	queryMapper   MerchantAwardQueryResponseMapper
	commandMapper MerchantAwardCommandResponseMapper
}

func NewMerchantAwardResponseMapper() MerchantAwardResponseMapper {
	return &merchantAwardResponseMapper{
		queryMapper:   NewMerchantAwardQueryResponseMapper(),
		commandMapper: NewMerchantAwardCommandResponseMapper(),
	}
}

func (m *merchantAwardResponseMapper) QueryMapper() MerchantAwardQueryResponseMapper {
	return m.queryMapper
}

func (m *merchantAwardResponseMapper) CommandMapper() MerchantAwardCommandResponseMapper {
	return m.commandMapper
}
