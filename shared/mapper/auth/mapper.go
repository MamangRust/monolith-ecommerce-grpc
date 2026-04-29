package authapimapper

type AuthResponseMapper interface {
	QueryMapper() AuthQueryResponseMapper
	CommandMapper() AuthCommandResponseMapper
}

type authResponseMapper struct {
	queryMapper   AuthQueryResponseMapper
	commandMapper AuthCommandResponseMapper
}

func NewAuthResponseMapper() AuthResponseMapper {
	return &authResponseMapper{
		queryMapper:   NewAuthQueryResponseMapper(),
		commandMapper: NewAuthCommandResponseMapper(),
	}
}

func (m *authResponseMapper) QueryMapper() AuthQueryResponseMapper {
	return m.queryMapper
}

func (m *authResponseMapper) CommandMapper() AuthCommandResponseMapper {
	return m.commandMapper
}
