package bannerapimapper

type BannerResponseMapper interface {
	QueryMapper() BannerQueryResponseMapper
	CommandMapper() BannerCommandResponseMapper
}

type bannerResponseMapper struct {
	queryMapper   BannerQueryResponseMapper
	commandMapper BannerCommandResponseMapper
}

func NewBannerResponseMapper() BannerResponseMapper {
	return &bannerResponseMapper{
		queryMapper:   NewBannerQueryResponseMapper(),
		commandMapper: NewBannerCommandResponseMapper(),
	}
}

func (m *bannerResponseMapper) QueryMapper() BannerQueryResponseMapper {
	return m.queryMapper
}

func (m *bannerResponseMapper) CommandMapper() BannerCommandResponseMapper {
	return m.commandMapper
}
