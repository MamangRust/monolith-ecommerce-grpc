package reviewdetailapimapper

type ReviewDetailResponseMapper interface {
	QueryMapper() ReviewDetailQueryResponseMapper
	CommandMapper() ReviewDetailCommandResponseMapper
}

type reviewDetailResponseMapper struct {
	queryMapper   ReviewDetailQueryResponseMapper
	commandMapper ReviewDetailCommandResponseMapper
}

func NewReviewDetailResponseMapper() ReviewDetailResponseMapper {
	return &reviewDetailResponseMapper{
		queryMapper:   NewReviewDetailQueryResponseMapper(),
		commandMapper: NewReviewDetailCommandResponseMapper(),
	}
}

func (m *reviewDetailResponseMapper) QueryMapper() ReviewDetailQueryResponseMapper {
	return m.queryMapper
}

func (m *reviewDetailResponseMapper) CommandMapper() ReviewDetailCommandResponseMapper {
	return m.commandMapper
}
