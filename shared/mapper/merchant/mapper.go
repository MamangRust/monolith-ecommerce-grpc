package merchantapimapper

type MerchantResponseMapper interface {
	QueryMapper() MerchantQueryResponseMapper
	CommandMapper() MerchantCommandResponseMapper
}

type merchantResponseMapper struct {
	queryMapper   MerchantQueryResponseMapper
	commandMapper MerchantCommandResponseMapper
}

func NewMerchantResponseMapper() MerchantResponseMapper {
	return &merchantResponseMapper{
		queryMapper:   NewMerchantQueryResponseMapper(),
		commandMapper: NewMerchantCommandResponseMapper(),
	}
}

func (m *merchantResponseMapper) QueryMapper() MerchantQueryResponseMapper {
	return m.queryMapper
}

func (m *merchantResponseMapper) CommandMapper() MerchantCommandResponseMapper {
	return m.commandMapper
}
