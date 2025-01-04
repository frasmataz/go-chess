package game

import "fmt"

type PieceClass int

const (
	Space PieceClass = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

type Piece struct {
	Class     PieceClass
	Colour    PlayerColour
	Symbol    rune
	FENSymbol rune
}

var Pieces = map[string]Piece{
	"space":       {Class: Space, Colour: Undefined, Symbol: ' '},
	"whitePawn":   {Class: Pawn, Colour: White, Symbol: '♟', FENSymbol: 'P'},
	"whiteKnight": {Class: Knight, Colour: White, Symbol: '♞', FENSymbol: 'N'},
	"whiteBishop": {Class: Bishop, Colour: White, Symbol: '♝', FENSymbol: 'B'},
	"whiteRook":   {Class: Rook, Colour: White, Symbol: '♜', FENSymbol: 'R'},
	"whiteQueen":  {Class: Queen, Colour: White, Symbol: '♛', FENSymbol: 'Q'},
	"whiteKing":   {Class: King, Colour: White, Symbol: '♚', FENSymbol: 'K'},
	"blackPawn":   {Class: Pawn, Colour: Black, Symbol: '♙', FENSymbol: 'p'},
	"blackKnight": {Class: Knight, Colour: Black, Symbol: '♘', FENSymbol: 'n'},
	"blackBishop": {Class: Bishop, Colour: Black, Symbol: '♗', FENSymbol: 'b'},
	"blackRook":   {Class: Rook, Colour: Black, Symbol: '♖', FENSymbol: 'r'},
	"blackQueen":  {Class: Queen, Colour: Black, Symbol: '♕', FENSymbol: 'q'},
	"blackKing":   {Class: King, Colour: Black, Symbol: '♔', FENSymbol: 'k'},
}

func GetPiece(name string) *Piece {
	piece := Pieces[name]
	return &piece
}

func FENToPiece(fenChar rune) (*Piece, error) {
	for _, piece := range Pieces {
		if piece.FENSymbol == fenChar {
			return &piece, nil
		}
	}

	return GetPiece("space"), fmt.Errorf("could not parse FEN char to piece: got '%c'", fenChar)
}
