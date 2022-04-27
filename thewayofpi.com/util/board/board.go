package board

type (
	Board struct {
		height  uint
		width	uint
		board	[][]rune
	}
)

func New(h,w int) *Board {
	return &Board{}
}

func (b *Board) String() string {
	return ""
}