package main

import (
	// "http/net"
	"bufio"
	"fmt"
	"os"
)

func main(){
	var fileID int

	fmt.Print("Id for game: ")
	fmt.Scan(&fileID)
	fileName := fmt.Sprintf("appmanifest_%d.acf", fileID)

	if(fileID <= 0){
		fmt.Println("App ID cannot be 0 or less!")
		return
	}

	content := fmt.Sprintf(`"AppState"
	{
		"AppID"  "%d"
		"Universe" "1"
		"installdir" "APPNAME"
		"StateFlags" "1026"
	}`, fileID)	

	
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintf(writer, content)
	fmt.Println("appmanifest file created successfully!")

	// Flush the buffered writer
	if err := writer.Flush(); err != nil {
			fmt.Println("Error flushing writer:", err)
	}
}