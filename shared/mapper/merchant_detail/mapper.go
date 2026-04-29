package merchantdetailapimapper

type MerchantDetailResponseMapper interface {
	QueryMapper() MerchantDetailQueryResponseMapper
	CommandMapper() MerchantDetailCommandResponseMapper
}

type merchantDetailResponseMapper struct {
	queryMapper   MerchantDetailQueryResponseMapper
	commandMapper MerchantDetailCommandResponseMapper
}

func NewMerchantDetailResponseMapper() MerchantDetailResponseMapper {
	return &merchantDetailResponseMapper{
		queryMapper:   NewMerchantDetailQueryResponseMapper(),
		commandMapper: NewMerchantDetailCommandResponseMapper(),
	}
}

func (m *merchantDetailResponseMapper) QueryMapper() MerchantDetailQueryResponseMapper {
	return m.queryMapper
}

func (m *merchantDetailResponseMapper) CommandMapper() MerchantDetailCommandResponseMapper {
	return m.commandMapper
}
