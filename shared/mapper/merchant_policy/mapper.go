package merchantpolicyapimapper

type MerchantPolicyResponseMapper interface {
	QueryMapper() MerchantPolicyQueryResponseMapper
	CommandMapper() MerchantPolicyCommandResponseMapper
}

type merchantPolicyResponseMapper struct {
	queryMapper   MerchantPolicyQueryResponseMapper
	commandMapper MerchantPolicyCommandResponseMapper
}

func NewMerchantPolicyResponseMapper() MerchantPolicyResponseMapper {
	return &merchantPolicyResponseMapper{
		queryMapper:   NewMerchantPolicyQueryResponseMapper(),
		commandMapper: NewMerchantPolicyCommandResponseMapper(),
	}
}

func (m *merchantPolicyResponseMapper) QueryMapper() MerchantPolicyQueryResponseMapper {
	return m.queryMapper
}

func (m *merchantPolicyResponseMapper) CommandMapper() MerchantPolicyCommandResponseMapper {
	return m.commandMapper
}
