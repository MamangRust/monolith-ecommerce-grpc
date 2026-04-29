package merchantdocumentsapimapper

type MerchantDocumentResponseMapper interface {
	QueryMapper() MerchantDocumentQueryResponseMapper
	CommandMapper() MerchantDocumentCommandResponseMapper
}

type merchantDocumentResponseMapper struct {
	queryMapper   MerchantDocumentQueryResponseMapper
	commandMapper MerchantDocumentCommandResponseMapper
}

func NewMerchantDocumentResponseMapper() MerchantDocumentResponseMapper {
	return &merchantDocumentResponseMapper{
		queryMapper:   NewMerchantDocumentQueryResponseMapper(),
		commandMapper: NewMerchantDocumentCommandResponseMapper(),
	}
}

func (m *merchantDocumentResponseMapper) QueryMapper() MerchantDocumentQueryResponseMapper {
	return m.queryMapper
}

func (m *merchantDocumentResponseMapper) CommandMapper() MerchantDocumentCommandResponseMapper {
	return m.commandMapper
}
