package categoryapimapper

type CategoryResponseMapper interface {
	QueryMapper() CategoryQueryResponseMapper
	CommandMapper() CategoryCommandResponseMapper
	StatsMapper() CategoryStatsResponseMapper
}

type categoryResponseMapper struct {
	queryMapper   CategoryQueryResponseMapper
	commandMapper CategoryCommandResponseMapper
	statsMapper   CategoryStatsResponseMapper
}

func NewCategoryResponseMapper() CategoryResponseMapper {
	return &categoryResponseMapper{
		queryMapper:   NewCategoryQueryResponseMapper(),
		commandMapper: NewCategoryCommandResponseMapper(),
		statsMapper:   NewCategoryStatsResponseMapper(),
	}
}

func (m *categoryResponseMapper) QueryMapper() CategoryQueryResponseMapper {
	return m.queryMapper
}

func (m *categoryResponseMapper) CommandMapper() CategoryCommandResponseMapper {
	return m.commandMapper
}

func (m *categoryResponseMapper) StatsMapper() CategoryStatsResponseMapper {
	return m.statsMapper
}
