package player

type MoveInput struct {
	Id    uint32
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

func (i *MoveInput) HasInput() bool {
	return i.Up || i.Down || i.Left || i.Right
}
