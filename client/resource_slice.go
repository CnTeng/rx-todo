package client

type ResourceStatus int

const (
	None ResourceStatus = iota
	Add
	Change
	Delete
)

type ResourceSlice interface {
	GetIDByIndex(idx int) (int64, error)
	GetIDsByIndexs(idxs []int) ([]int64, error)
}
