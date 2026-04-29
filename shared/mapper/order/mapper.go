package orderapimapper

type OrderResponseMapper interface {
	QueryMapper() OrderQueryResponseMapper
	CommandMapper() OrderCommandResponseMapper
	StatsMapper() OrderStatsResponseMapper
}

type orderResponseMapper struct {
	queryMapper   OrderQueryResponseMapper
	commandMapper OrderCommandResponseMapper
	statsMapper   OrderStatsResponseMapper
}

func NewOrderResponseMapper() OrderResponseMapper {
	return &orderResponseMapper{
		queryMapper:   NewOrderQueryResponseMapper(),
		commandMapper: NewOrderCommandResponseMapper(),
		statsMapper:   NewOrderStatsResponseMapper(),
	}
}

func (o *orderResponseMapper) QueryMapper() OrderQueryResponseMapper {
	return o.queryMapper
}

func (o *orderResponseMapper) CommandMapper() OrderCommandResponseMapper {
	return o.commandMapper
}

func (o *orderResponseMapper) StatsMapper() OrderStatsResponseMapper {
	return o.statsMapper
}
