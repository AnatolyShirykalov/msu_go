package game

type ObjSubj struct {
	Key      string
	Exist    bool
	Lock     bool
	NameRoom string
}

type Back struct {
	Things map[string]bool
	Full   bool
}

type Game struct {
	Priory  []string
	Rooms   map[string]*Room
	Players []Player
	Aliases map[string]string
}

type Room struct {
	Game     *Game
	Name     string
	Msg      map[string]string
	Things   map[string]bool
	Subjects map[string]ObjSubj
	Act      string
	LinkRoom map[string]string
}

type Player struct {
	InRoom  *Room
	RefBack *Back
}

func (r *Room) Label() string {
	return r.Name
}

func (r *Room) Type() string {
	return "комнаты"
}
