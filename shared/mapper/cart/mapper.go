package cartapimapper

type CartResponseMapper interface {
	QueryMapper() CartQueryResponseMapper
	CommandMapper() CartCommandResponseMapper
}

type cartResponseMapper struct {
	queryMapper   CartQueryResponseMapper
	commandMapper CartCommandResponseMapper
}

func NewCartResponseMapper() CartResponseMapper {
	return &cartResponseMapper{
		queryMapper:   NewCartQueryResponseMapper(),
		commandMapper: NewCartCommandResponseMapper(),
	}
}

func (t *cartResponseMapper) QueryMapper() CartQueryResponseMapper {
	return t.queryMapper
}

func (t *cartResponseMapper) CommandMapper() CartCommandResponseMapper {
	return t.commandMapper
}
