package enums

type TransactionType int64

const (
	Deposit TransactionType = iota + 1
	Withdrawal
	Transfer
)

var TransactionTypeMap = map[TransactionType]string{
	Deposit:    "Deposit",
	Withdrawal: "Withdrawal",
	Transfer:   "Transfer",
}

func (tt TransactionType) String() string {
	return TransactionTypeMap[tt]
}
