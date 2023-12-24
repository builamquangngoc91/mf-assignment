package enums

type TransactionStatus int64

const (
	Completed TransactionStatus = iota + 1
)

var TransactionStatusMap = map[TransactionStatus]string{
	Completed: "Completed",
}

func (tt TransactionStatus) String() string {
	return TransactionStatusMap[tt]
}
