package player

type Input struct {
	Id    uint32
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

func (i *Input) HasInput() bool {
	return i.Up || i.Down || i.Left || i.Right
}
