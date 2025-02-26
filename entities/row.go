package entities

type Row struct {
	RowIndex int
	ColIndex int
	Type     ITypeSystem
	Value    any
}

func NewRow(rowIndex int, colIndex int, t ITypeSystem, data any) *Row {
	return &Row{
		RowIndex: rowIndex,
		ColIndex: colIndex,
		Type:     t,
		Value:    data,
	}
}
