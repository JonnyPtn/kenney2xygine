package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Sprite - each sprite in the atlas
type Sprite struct {
	Name string `xml:"name,attr"`
	X    string `xml:"x,attr"`
	Y    string `xml:"y,attr"`
	W    string `xml:"width,attr"`
	H    string `xml:"height,attr"`
}

// TextureAtlas - the texture atlas
type TextureAtlas struct {
	Path    string   `xml:"imagePath,attr"`
	Sprites []Sprite `xml:"SubTexture"`
}

// Error checking function
func errCheck(e error) {
	if e != nil {
		panic(e)
	}
}

// Write a string to a file
func writeLine(f *os.File, s string) {
	_, err := f.WriteString(s + "\n")
	errCheck(err)
}

// Converts the xml spritesheets provided with kenney assets to the format
// required for xygine spritesheet usage
//
// Kenney - kenny.nl
// xygine - github.com/fallahn/xygine
//
func main() {

	// Read the file passed in on the command line
	// Defaults to finding all local xml files
	fileName := flag.String("file", "", "File to convert")
	flag.Parse()

	// Make sure a file has been given
	if len(*fileName) == 0 {
		flag.Usage()
		return
	}

	// Open the file
	xmlFile, err := os.Open(*fileName)
	errCheck(err)

	defer xmlFile.Close()

	fmt.Println("Successfully opened " + *fileName)

	// Read as a byte array
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// Unmarshal the xml data
	var textureAtlas TextureAtlas

	xml.Unmarshal(byteValue, &textureAtlas)

	fmt.Println("Using image " + textureAtlas.Path)

	// Strip extension from filename and add xygine spritesheet extension
	extension := filepath.Ext(*fileName)
	outputFile := (*fileName)[0 : len(*fileName)-len(extension)]
	outputFile += ".spt"

	fmt.Println("Creating " + outputFile)

	f, err := os.Create(outputFile)
	errCheck(err)

	defer f.Close()

	// Opening empyty line and bracket, not sure why...
	writeLine(f, " ")
	writeLine(f, "{")

	// Write src parameter
	writeLine(f, "	src = "+textureAtlas.Path)

	// And set smooth to false by default
	writeLine(f, "	smooth = false")

	// Now output each sprite
	for _, sprite := range textureAtlas.Sprites {

		// Sprite name
		writeLine(f, "	sprite "+sprite.Name)

		// Opening bracket for sprite deets
		writeLine(f, "	{")

		// Sprite rect
		writeLine(f, "		bounds = "+sprite.X+","+sprite.Y+","+sprite.W+","+sprite.H)

		// Default White colour
		writeLine(f, "		colour = 255,255,255,255")

		// Close bracket for sprite deets
		writeLine(f, "	}")
	}

	// Final close bracker
	writeLine(f, "}")

	fmt.Println("Successfully converted " + *fileName + " to " + outputFile)
}
