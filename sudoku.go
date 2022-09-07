package main

type SudokuCell struct {
	value    byte
	readonly bool
}

// NewSudokuCell returns a new, modifiable SudokuCell.
func NewSudokuCell(val byte) *SudokuCell {
	return &SudokuCell{
		value: val,
	}
}

// NewSudokuCell returns a new, readonly SudokuCell.
func NewReadonlySudokuCell(val byte) *SudokuCell {
	return NewSudokuCell(val).SetReadonly(true)
}

func (c *SudokuCell) Value() byte {
	return c.value
}

func (c *SudokuCell) SetValue(v byte) *SudokuCell {
	c.value = v
	if v == '0' {
		c.value = ' '
	}
	return c
}

func (c *SudokuCell) Readonly() bool {
	return c.readonly
}

func (c *SudokuCell) SetReadonly(v bool) *SudokuCell {
	c.readonly = v
	return c
}

func (c SudokuCell) Rune() rune {
	return rune(c.value)
}
