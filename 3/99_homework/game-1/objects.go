package game

type Command struct {
	player  *Player
	command string
}

type ObjSubj struct {
	Key      string
	Exist    bool
	Lock     bool
	NameRoom string
}
type Player struct {
	InRoom  *Room
	RefBack *Back
	Name    string
	msg     chan string
}
type Back struct {
	Things map[string]bool
	Full   bool
}
type Game struct {
	msgin chan *Command

	Priory  []string
	Rooms   map[string]*Room
	Players map[string]*Player
	Aliases map[string]string
}

type Room struct {
	Name     string
	Msg      map[string]string
	Things   map[string]bool
	Subjects map[string]ObjSubj
	Act      string
	LinkRoom map[string]string
}

func (r *Room) Label() string {
	return r.Name
}

func (r *Room) Type() string {
	return "комнаты"
}
