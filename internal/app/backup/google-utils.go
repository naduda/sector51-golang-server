package backup

import (
	"context"
	"encoding/json"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	credentialsFileName = "configs/credentials.json"
	tokenFile           = "./configs/token.json"
)

// GetClient ...
func GetClientOptions() (option.ClientOption, error) {
	config, err := GetGoogleConfig()
	if err != nil {
		return nil, err
	}

	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		return nil, err
	}

	return option.WithTokenSource(config.TokenSource(context.Background(), tok)), nil
}

// GetDriveService ...
func GetDriveService() (*drive.Service, error) {
	clientOptions, err := GetClientOptions()
	if err != nil {
		return nil, err
	}
	return drive.NewService(context.Background(), clientOptions)
}

// GetGoogleConfig ...
func GetGoogleConfig() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(credentialsFileName)
	if err != nil {
		return nil, err
	}

	return google.ConfigFromJSON(b,
		drive.DriveMetadataReadonlyScope,
		drive.DriveScope,
		drive.DriveFileScope,
		drive.DriveReadonlyScope)
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
