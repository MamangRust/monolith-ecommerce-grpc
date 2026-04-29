package sliderapimapper

type SliderResponseMapper interface {
	QueryMapper() SliderQueryResponseMapper
	CommandMapper() SliderCommandResponseMapper
}

type sliderResponseMapper struct {
	queryMapper   SliderQueryResponseMapper
	commandMapper SliderCommandResponseMapper
}

func NewSliderResponseMapper() SliderResponseMapper {
	return &sliderResponseMapper{
		queryMapper:   NewSliderQueryResponseMapper(),
		commandMapper: NewSliderCommandResponseMapper(),
	}
}

func (m *sliderResponseMapper) QueryMapper() SliderQueryResponseMapper {
	return m.queryMapper
}

func (m *sliderResponseMapper) CommandMapper() SliderCommandResponseMapper {
	return m.commandMapper
}
