package transactionapimapper

type TransactionResponseMapper interface {
	QueryMapper() TransactionQueryResponseMapper
	CommandMapper() TransactionCommandResponseMapper
	StatsMapper() TransactionStatsResponseMapper
}

type transactionResponseMapper struct {
	queryMapper   TransactionQueryResponseMapper
	commandMapper TransactionCommandResponseMapper
	statsMapper   TransactionStatsResponseMapper
}

func NewTransactionResponseMapper() TransactionResponseMapper {
	return &transactionResponseMapper{
		queryMapper:   NewTransactionQueryResponseMapper(),
		commandMapper: NewTransactionCommandResponseMapper(),
		statsMapper:   NewTransactionStatsResponseMapper(),
	}
}

func (t *transactionResponseMapper) QueryMapper() TransactionQueryResponseMapper {
	return t.queryMapper
}

func (t *transactionResponseMapper) CommandMapper() TransactionCommandResponseMapper {
	return t.commandMapper
}

func (t *transactionResponseMapper) StatsMapper() TransactionStatsResponseMapper {
	return t.statsMapper
}
