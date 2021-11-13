package morris

import (
	"fmt"
)

type edge struct {
	From int32
	To   int32
}

type IllegalMoveError struct {
	Description string
}

func (e IllegalMoveError) Error() string {
	return e.Description
}

var (
	legalMoves       *[]edge     = createMoveGraph()
	possibleMorrises *[][3]int32 = createMorrisSet()

	whitePiece string = "⚪"
	blackPiece string = "⚫"
)

func createMoveGraph() *[]edge {
	var graph []edge
	addEdge := func(a int32, b int32) {
		graph = append(graph, edge{From: a, To: b})
		graph = append(graph, edge{From: b, To: a})
	}
	addEdges := func(vertices ...int32) {
		for i := 1; i < len(vertices); i++ {
			addEdge(vertices[i-1], vertices[i])
		}
	}
	addEdges(1, 2, 14, 23, 22, 21, 9, 0, 1, 4, 3, 10, 18, 19, 20, 13, 5, 4, 7,
		8, 12, 17, 16, 15, 11, 6, 7)
	addEdges(9, 10, 11)
	addEdges(12, 13, 14)
	addEdges(16, 19, 22)
	return &graph
}

func createMorrisSet() *[][3]int32 {
	return &[][3]int32{
		// Horizontal
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {9, 19, 11}, {12, 13, 14},
		{15, 16, 17}, {18, 19, 20}, {21, 22, 23},
		// Vertical
		{0, 9, 21}, {3, 10, 18}, {6, 11, 15}, {1, 4, 7}, {16, 19, 22},
		{8, 12, 17}, {5, 13, 20}, {2, 14, 23},
	}
}

func IsBoardSpace(space int32) bool {
	if space < 0 || space >= 24 {
		return false
	} else {
		return true
	}
}

// Checks if a move is valid.
func (m *Move) IsValid() bool {
	for _, e := range *legalMoves {
		if e.From == m.From && e.To == m.To {
			return true
		}
	}
	return false
}

// Checks if a space is part of a morris.
func (bs *BoardState) HasMorrisAt(space int32) bool {
	for _, pm := range *possibleMorrises {
		if !morrisContainsSpace(pm, space) {
			continue
		}
		isMorris := checkMorris([3]BoardSpace{
			bs.Board[pm[0]],
			bs.Board[pm[1]],
			bs.Board[pm[2]],
		})
		if isMorris {
			return true
		}
	}
	return false
}

// Does this morris contain space?
func morrisContainsSpace(morris [3]int32, space int32) bool {
	for _, s := range morris {
		if s == space {
			return true
		}
	}
	return false
}

// Does this possible morris actually contain a morris?
func checkMorris(morris [3]BoardSpace) bool {
	color := morris[0]
	if color == BoardSpace_FREE {
		return false
	}
	for i := 1; i < len(morris); i++ {
		if morris[i] != color {
			return false
		}
	}
	return true
}

// Returns a string which can be printed to visually show the state of the
// board.
func (bs *BoardState) Visualize(showSpaceIds bool) string {
	board := bs.Board
	var boardString string
	if showSpaceIds {
		boardString = "Current Turn: %s \n\r" +
			" 0%s────────────1%s────────────2%s\n\r" +
			"  │              │              │ \n\r" +
			"  │   3%s───────4%s───────5%s   │ \n\r" +
			"  │    │         │         │    │ \n\r" +
			"  │    │   6%s──7%s──8%s   │    │ \n\r" +
			"  │    │    │         │    │    │ \n\r" +
			" 9%s─10%s─11%s      12%s─13%s─14%s\n\r" +
			"  │    │    │         │    │    │ \n\r" +
			"  │    │  15%s─16%s─17%s   │    │ \n\r" +
			"  │    │         │         │    │ \n\r" +
			"  │  18%s──────19%s──────20%s   │ \n\r" +
			"  │              │              │ \n\r" +
			"21%s───────────22%s───────────23%s\n\r" +
			"Graveyard:   %s%d   %s%d\n\r" +
			"%s\n\r"
	} else {
		boardString = "Current Turn: %s \n\r" +
			"  %s─────────────%s─────────────%s\n\r" +
			"  │              │              │ \n\r" +
			"  │    %s────────%s────────%s   │ \n\r" +
			"  │    │         │         │    │ \n\r" +
			"  │    │    %s───%s───%s   │    │ \n\r" +
			"  │    │    │         │    │    │ \n\r" +
			"  %s───%s───%s        %s───%s───%s\n\r" +
			"  │    │    │         │    │    │ \n\r" +
			"  │    │    %s───%s───%s   │    │ \n\r" +
			"  │    │         │         │    │ \n\r" +
			"  │    %s────────%s────────%s   │ \n\r" +
			"  │              │              │ \n\r" +
			"  %s─────────────%s─────────────%s\n\r" +
			"Graveyard:   %s%d   %s%d\n\r" +
			"%s\n\r"
	}
	return fmt.Sprintf(boardString,
		bs.Turn.Visualize("  "),
		board[0].Visualize("┌─"),
		board[1].Visualize("┬─"),
		board[2].Visualize("┐ "),
		board[3].Visualize("┌─"),
		board[4].Visualize("┼─"),
		board[5].Visualize("┐ "),
		board[6].Visualize("┌─"),
		board[7].Visualize("┴─"),
		board[8].Visualize("┐ "),
		board[9].Visualize("├─"),
		board[10].Visualize("┼─"),
		board[11].Visualize("┤ "),
		board[12].Visualize("├─"),
		board[13].Visualize("┼─"),
		board[14].Visualize("┤ "),
		board[15].Visualize("└─"),
		board[16].Visualize("┬─"),
		board[17].Visualize("┘ "),
		board[18].Visualize("└─"),
		board[19].Visualize("┼─"),
		board[20].Visualize("┘ "),
		board[21].Visualize("└─"),
		board[22].Visualize("┴─"),
		board[23].Visualize("┘ "),
		whitePiece, bs.WhiteGrave,
		blackPiece, bs.BlackGrave,
		bs.Phase.Visualize(bs.WhitePieces, bs.BlackPieces),
	)
}

func (p *Phase) Visualize(whiteLeft int32, blackLeft int32) string {
	switch *p {
	case Phase_PLACE:
		return fmt.Sprintf("Placeing - %s%d  %s%d", whitePiece, whiteLeft, blackPiece, blackLeft)
	case Phase_MOVE:
		return "Moving"
	case Phase_FLY:
		return "Flying"
	default:
		return ""
	}
}

// Returns a character representing the color of a BoardSpace. If it is empty,
// background will be used instead. Background should be two characters long.
func (bs *BoardSpace) Visualize(background string) string { // FIXME: Default terminal fonts cannot show emojis
	switch *bs {
	case BoardSpace_WHITE:
		return whitePiece
	case BoardSpace_BLACK:
		return blackPiece
	default:
		return background
	}
}

func SetSymbolMode(useFontSafeSymbols bool) {
	if useFontSafeSymbols {
		whitePiece = "■ "
		blackPiece = "□ "
	} else {
		whitePiece = "⚪"
		blackPiece = "⚫"
	}
}
