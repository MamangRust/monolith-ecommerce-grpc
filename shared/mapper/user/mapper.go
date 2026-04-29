package userapimapper

type UserResponseMapper interface {
	QueryMapper() UserQueryResponseMapper
	CommandMapper() UserCommandResponseMapper
}

type userResponseMapper struct {
	queryMapper   UserQueryResponseMapper
	commandMapper UserCommandResponseMapper
}

func NewUserResponseMapper() UserResponseMapper {
	return &userResponseMapper{
		queryMapper:   NewUserQueryResponseMapper(),
		commandMapper: NewUserCommandResponseMapper(),
	}
}

func (u *userResponseMapper) QueryMapper() UserQueryResponseMapper {
	return u.queryMapper
}

func (u *userResponseMapper) CommandMapper() UserCommandResponseMapper {
	return u.commandMapper
}
