package api

type RegisterPlayer struct {
	Name string `json:"name" binding:"required" example:"my name"`
}

func (r *RegisterPlayer) Validation() error {

	return nil
}
