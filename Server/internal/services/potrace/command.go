package potrace

import (
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/jigsawcorp/log3900/internal/datastore"
	"gitlab.com/jigsawcorp/log3900/internal/services/potrace/model"
	"gitlab.com/jigsawcorp/log3900/internal/svgparser"
	"gitlab.com/jigsawcorp/log3900/pkg/strbuilder"
	"io/ioutil"
	"os/exec"
)

const defaultDraw = 0.5

func checkCommand(command string) bool {
	_, err := exec.LookPath(command)
	if err != nil {
		return false
	}
	return true
}

//Trace trace the image to the svg returns the key of the svg file generated
func Trace(imageKey string, blacklevel float64) (string, error) {
	svgKey := datastore.GenerateFileKey()
	svgPath := datastore.GetPath(svgKey)
	imagePath := datastore.GetPath(imageKey)

	cmd := fmt.Sprintf("convert -flatten -resize 1065x690 -background white -gravity center -extent 1125x750 %s bmp:- | potrace - -s -k %f -o %s", imagePath, blacklevel, svgPath)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("Failed to execute the conversion")
	}

	return svgKey, nil
}

//Translate changes the potrace svg to be compatible with our custom svg format
func Translate(svgKey string, brushSize int, mode int, update bool) error {
	file, err := datastore.GetFile(svgKey)
	if err != nil {
		return err
	}
	byteValue, _ := ioutil.ReadAll(file)
	var xmlSvg model.XMLSvg

	err = xml.Unmarshal(byteValue, &xmlSvg)
	if err != nil {
		return err
	}

	xmlSvg.G.Style = "fill:none;stroke:black"

	//Check all the paths
	for i := range xmlSvg.G.XMLPaths {
		var path *model.XMLPath
		path = &xmlSvg.G.XMLPaths[i]
		if path.ID == uuid.Nil {
			path.ID = uuid.New()
			path.Brush = "circle"
		}
		if brushSize > 0 {
			path.BrushSize = brushSize
		}

		if update {
			commands := svgparser.ParseD(path.D, []svgparser.TransformCommand{})
			if len(commands) > 0 {
				path.FirstCommand = commands[0]
			}
		}
	}
	if !update {
		splitPath(&xmlSvg, xmlSvg.G.Transform)
	}

	ChangeOrder(&xmlSvg.G.XMLPaths, mode)
	if xmlSvg.G.Transform != "" {
		xmlSvg.G.Transform = ""
	}
	data, err := xml.Marshal(&xmlSvg)
	err = datastore.PutFile(&data, svgKey)
	if err != nil {
		return err
	}
	return nil
}

//splitPath export paths
func splitPath(svg *model.XMLSvg, transform string) {
	builder := strbuilder.StrBuilder{}
	builder.Grow(128)
	newPaths := make([]model.XMLPath, 0, 20)
	transformCommands := svgparser.TransformParse(transform, svg.Width, svg.Height)

	for i := range svg.G.XMLPaths {
		commands := svgparser.ParseD(svg.G.XMLPaths[i].D, transformCommands)

		lastChunk := 0
		subPathCount := 0

		for j := range commands {
			char := commands[j].Command
			if (char == 'm' || char == 'M') && j != 0 {
				singleCommands := commands[lastChunk:j]
				processChunk(&singleCommands, &(svg.G.XMLPaths)[i], &newPaths, &builder, lastChunk)

				lastChunk = j
				subPathCount++
			} else if j == 0 {
				if char == 'm' || char == 'M' {
					subPathCount++
				}
			}
		}

		if subPathCount >= 1 {
			singleCommands := commands[lastChunk:]
			processChunk(&singleCommands, &(svg.G.XMLPaths)[i], &newPaths, &builder, lastChunk)
		}
	}
	svg.G.XMLPaths = newPaths
}

func processChunk(commands *[]svgparser.Command, currentPath *model.XMLPath, paths *[]model.XMLPath, builder *strbuilder.StrBuilder, j int) {
	builder.Reset()

	singlePath := model.XMLPath{
		ID:        uuid.New(),
		Eraser:    currentPath.Eraser,
		Brush:     currentPath.Brush,
		BrushSize: currentPath.BrushSize,
		Color:     currentPath.Color,
	}
	length := 0.0
	for i := range *commands {
		if i == 0 {
			singlePath.FirstCommand = (*commands)[i]
		}
		builder.WriteString((*commands)[i].Encode())
		length += (*commands)[i].ComputeLength()
	}
	singlePath.Length = length
	singlePath.Time = int(defaultDraw * length)
	singlePath.D = builder.StringVal()
	*paths = append(*paths, singlePath)
}
