package styles

type DocSheet interface {
	Sheet
}

type DocSheetData struct {
	SheetData
}

func NewDocSheet() *DocSheetData {
	return &DocSheetData{newSheetData()}
}
