package task

type State string

const (
	New State = "new"
	Downloaded State = "downloaded"
	Empty State = "empty"
)

type Task struct {
	ReportID int64
	SellerID uint64
	SellerName string
	Status State
}