package user

type Status struct {
	username string
	locked   bool
}

func (status Status) isLocked() bool {
	return status.locked
}
