package googlesheets

import "fmt"

// SheetExprorter ...
type SheetExprorter struct {
	clientID     string
	clientSecret string
}

// NewSheetExporter ...
func NewSheetExporter() *SheetExprorter {
	return &SheetExprorter{
		clientID:     "",
		clientSecret: "",
	}
}

// Backup ...
func (e *SheetExprorter) Backup(sheetID string) error {
	srv, err := GetSheetService()
	if err != nil {
		return err
	}

	readRange := "Sheet1!A1"
	resp, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		return err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			fmt.Printf("%s\n", row[0])
		}
	}

	return nil
}
