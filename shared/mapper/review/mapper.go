package reviewapimapper

type ReviewResponseMapper interface {
	QueryMapper() ReviewQueryResponseMapper
	CommandMapper() ReviewCommandResponseMapper
}

type reviewResponseMapper struct {
	queryMapper   ReviewQueryResponseMapper
	commandMapper ReviewCommandResponseMapper
}

func NewReviewResponseMapper() ReviewResponseMapper {
	return &reviewResponseMapper{
		queryMapper:   NewReviewQueryResponseMapper(),
		commandMapper: NewReviewCommandResponseMapper(),
	}
}

func (m *reviewResponseMapper) QueryMapper() ReviewQueryResponseMapper {
	return m.queryMapper
}

func (m *reviewResponseMapper) CommandMapper() ReviewCommandResponseMapper {
	return m.commandMapper
}
