package game

type ObjSubj struct {
	Name  string
	Key   string
	Exist bool
	Lock  bool
}
type Link struct {
	Rfrom *Room
	Rto   *Room
	Name  string
	Lock  bool
}

type Back struct {
	Things map[string]bool
	Full   bool
}

type Game struct {
	Rooms   map[string]*Room
	Links   []Link
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
	// Watch    bool
}
