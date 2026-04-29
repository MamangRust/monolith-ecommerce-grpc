package merchantsociallinkapimapper

type MerchantSocialLinkResponseMapper interface {
	QueryMapper() MerchantSocialLinkQueryResponseMapper
	CommandMapper() MerchantSocialLinkCommandResponseMapper
}

type merchantSocialLinkResponseMapper struct {
	queryMapper   MerchantSocialLinkQueryResponseMapper
	commandMapper MerchantSocialLinkCommandResponseMapper
}

func NewMerchantSocialLinkResponseMapper() MerchantSocialLinkResponseMapper {
	return &merchantSocialLinkResponseMapper{
		queryMapper:   NewMerchantSocialLinkQueryResponseMapper(),
		commandMapper: NewMerchantSocialLinkCommandResponseMapper(),
	}
}

func (m *merchantSocialLinkResponseMapper) QueryMapper() MerchantSocialLinkQueryResponseMapper {
	return m.queryMapper
}

func (m *merchantSocialLinkResponseMapper) CommandMapper() MerchantSocialLinkCommandResponseMapper {
	return m.commandMapper
}
