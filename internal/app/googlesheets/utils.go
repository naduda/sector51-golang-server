package googlesheets

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// PageFile ...
const (
	credentialsFileName = "configs/credentials.json"
	tokenFile           = "./configs/token.json"
	PageFile            = "./configs/page-id.json"
)

// GooglePageSettings ...
type GooglePageSettings struct {
	PageID string `json:"pageId"`
}

// GetClient ...
func GetClient() (*http.Client, error) {
	config, err := GetGoogleConfig()
	if err != nil {
		return nil, err
	}

	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

// GetSheetService ...
func GetSheetService() (*sheets.Service, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	return sheets.New(client)
}

// GetGoogleConfig ...
func GetGoogleConfig() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(credentialsFileName)
	if err != nil {
		return nil, err
	}

	return google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
}

// GetPageID ...
func GetPageID() (*GooglePageSettings, error) {
	res := &GooglePageSettings{}
	f, err := os.Open(PageFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)

	if err = json.Unmarshal(byteValue, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// SaveToken ...
func SaveToken(token *oauth2.Token) error {
	f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

// func main() {
// 	b, err := ioutil.ReadFile("./configs/credentials.json")
// 	if err != nil {
// 		log.Fatalf("Unable to read client secret file: %v", err)
// 	}

// 	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
// 	if err != nil {
// 		log.Fatalf("Unable to parse client secret file to config: %v", err)
// 	}
// 	client := getClient(config)

// 	srv, err := sheets.New(client)
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve Sheets client: %v", err)
// 	}

// 	// Prints the names and majors of students in a sample spreadsheet:
// 	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
// 	spreadsheetId := "1iZjgXCIecU0M9JOGCq-Zk0s67-nTbqXZXOEBbkojO7s"
// 	readRange := "Sheet1!A16:C20"
// 	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve data from sheet: %v", err)
// 	}

// 	if len(resp.Values) == 0 {
// 		fmt.Println("No data found.")
// 	} else {
// 		fmt.Println("Name, Major:")
// 		for _, row := range resp.Values {
// 			// Print columns A and E, which correspond to indices 0 and 4.
// 			fmt.Printf("%s, %s\n", row[0], row[1])
// 		}
// 	}
// }
