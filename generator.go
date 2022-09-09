package main

// GeneratePuzzle is a mock function that returns a partial SudokuGrid.
func GeneratePuzzle() *SudokuGrid {
	g := NewSudokuGrid()
	for r, row := range [9][9]int{
		{0, 0, 5, 3, 0, 0, 0, 0, 0},
		{8, 0, 0, 0, 0, 0, 0, 2, 0},
		{0, 7, 0, 0, 1, 0, 5, 0, 0},
		{4, 0, 0, 0, 0, 5, 3, 0, 0},
		{0, 1, 0, 0, 7, 0, 0, 0, 6},
		{0, 0, 3, 2, 0, 0, 0, 8, 0},
		{0, 6, 0, 5, 0, 0, 0, 0, 9},
		{0, 0, 4, 0, 0, 0, 0, 3, 0},
		{0, 0, 0, 0, 0, 9, 7, 0, 0},
	} {
		for c, v := range row {
			cell := g.GetCell(r, c)
			cell.SetValue(v)
			if !cell.IsEmpty() {
				cell.SetReadonly(true)
			}
		}
	}
	return g
}
