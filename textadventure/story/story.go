package story

import (
	"fmt"
	"strings"
)

// Node represents a story node
type Node struct {
	Text    string
	Choices []Choice
}

// Choice represents a single story node choice
type Choice struct {
	cmd         string
	description string
	nextStory   *Node
}

// NewNode initializes a new Node with the given text and predefined
// list of choice. You can pass in nil for choices and use the
// AddChoice method of the story Node to add a single choice
func NewNode(text string, choices []Choice) *Node {
	return &Node{Text: text, Choices: choices}
}

// AddChoice adds a new choice to the story node choices
func (s *Node) AddChoice(cmd, description string, nextStory *Node) {
	c := Choice{cmd, description, nextStory}
	s.Choices = append(s.Choices, c)
}

// Render prints the content of this story node
func (s *Node) Render() {
	fmt.Println(s.Text)

	for _, choice := range s.Choices {
		fmt.Println(choice.cmd, ":", choice.description)
	}
}

// ExecuteCmd executes the input cmd against the list of choices
// and returns the appropriate next node in the graph
func (s *Node) ExecuteCmd(cmd string) *Node {
	if s.Choices == nil {
		fmt.Println("Add choices first to execute a command")
	}

	for _, choice := range s.Choices {
		if strings.ToLower(choice.cmd) == strings.ToLower(cmd) {
			return choice.nextStory
		}
	}

	fmt.Println("Sorry, I didn't understand that.")
	// return current storyNode, so user starts the node again
	return s
}
