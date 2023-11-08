package room

type Street struct {
	RoomBase
}

func (s *Street) LookAround() string {
	return "на улице весна. можно пройти - домой"
}

func (s *Street) Greet() string {
	return "на улице весна. можно пройти - домой"
}
