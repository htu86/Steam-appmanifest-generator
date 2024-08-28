package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"io"
	"encoding/json"
	"path/filepath"
	"runtime"
)

type Data struct {
	GameName string `json:"name"`
}

type AppDetails struct {
  Success bool `json:"success"`
	Data Data `json:"data"`
}

func getOS() string{
	currentOs := runtime.GOOS
	fmt.Println("Operating system:", currentOs)
	return currentOs
}

func getAppName(appID int) string{
	steamURL := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%d", appID)

	response, err := http.Get(steamURL)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	responseText, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	var result map[string]AppDetails
	err = json.Unmarshal(responseText, &result)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	appDetails := result[fmt.Sprintf("%d", appID)]
	return appDetails.Data.GameName
}

func createFile(fileName string, content string){
		// Creating the actual file
	
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintf(writer, content)
	fmt.Println("Appmanifest file created successfully!")

	// Flush the buffered writer

	if err := writer.Flush(); err != nil {
		fmt.Println("Error flushing writer:", err)
	}

}

func main(){
	var gameID int

	fmt.Print("Id for game: ")
	fmt.Scan(&gameID)

	if(gameID <= 0){
		fmt.Println("App ID cannot be 0 or less!")
		return
	}

	osType := getOS()

	// Get the user directory

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
// Library/Application Support/Steam/steamapps
	var directory string
	if osType == "darwin" {
		directory = filepath.Join(homeDir, "Library/Application Support/Steam/steamapps")
	} else if osType == "linux" {
		directory = filepath.Join(homeDir, ".steam/steam/SteamApps")
	} else {
		fmt.Println("OS is unsupported")
		return
	}

	fmt.Println(directory)

	gameName := getAppName(gameID)

	// Contents of the app-manifest file

	content := fmt.Sprintf(`"AppState"
{
	"AppID"  "%d"
	"Universe" "1"
	"installdir" "%v"
	"StateFlags" "1026"
}`, gameID, gameName)	

	// Creating the actual file

	file := filepath.Join(directory, fmt.Sprintf("appmanifest_%d.acf", gameID))
	createFile(file, content)
}