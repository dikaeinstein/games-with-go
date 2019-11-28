package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dikaeinstein/games-with-go/textadventure/story"
)

func main() {
	starText := `
You are in a large chamber, deep underground.
You see three passages leading out. A north passage leads into darkness.
To the south, a passage appears to head upward. The eastern passage
appears flat and well travelled`

	start := story.NewNode(starText, nil)
	darkRoom := story.NewNode("It is pitch black. You cannot see a thing", nil)
	darkRoomLit := story.NewNode("The dark passage is now lit by your lantern. You can continue north or head back south.", nil)
	grue := story.NewNode("While stumbling around in the darkness, you are eaten by a grue.", nil)
	trap := story.NewNode("You head down the well travelled path when suddenly a trap door opens and you fall into a pit.", nil)
	treasure := story.NewNode("You arrive at a small chamber filled with treasure!", nil)

	// add choices to story nodes
	start.AddChoice("N", "Go North", darkRoom)
	start.AddChoice("S", "Go South", darkRoom)
	start.AddChoice("E", "Go East", trap)

	darkRoom.AddChoice("S", "Try to go back south", grue)
	darkRoom.AddChoice("O", "Turn on lantern", darkRoomLit)

	darkRoomLit.AddChoice("N", "Go North", treasure)
	darkRoomLit.AddChoice("S", "Go South", start)

	// Play the game
	play(start)
}

func play(start *story.Node) {
	currentNode := start
	for {
		currentNode.Render()

		if currentNode.Choices == nil {
			// End game when no node to advance/traverse to
			break
		}

		r := bufio.NewReader(os.Stdin)
		cmd, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}

		cmd = strings.TrimSpace(cmd)
		currentNode = currentNode.ExecuteCmd(cmd)
	}

	fmt.Println("The End!")
}
