package help

import (
	"fmt"
)

func Run() {
	fmt.Printf("\n\nCommands : dnote\n")
	fmt.Printf("       --> use [book name]              * Change the book to write your note in * \n")
	fmt.Printf("                                        * alias : u * \n\n")
	fmt.Printf("       --> new [content]                * Write a new note under the current book * \n")
	fmt.Printf("                                        * alias : n * \n\n")
	fmt.Printf("       --> edit [note index] [content]  * Overwrite a note under the current book * \n")
	fmt.Printf("                                        * alias : e * \n\n")
	fmt.Printf("                                        * option : -b *Specify the name of the book to read from * \n\n")
	fmt.Printf("       --> delete [item]                * Delete a note in the current book * \n")
	fmt.Printf("                                        * alias : d * \n")
	fmt.Printf("                                        * option : -b *Specify the name of the book to delete from * \n\n")
	fmt.Printf("                                                   --book *Specify the name of the book to be deleted * \n\n")
	fmt.Printf("       --> books                        * List all the books that you created * \n")
	fmt.Printf("                                        * alias : b * \n\n")
	fmt.Printf("       --> notes                        * List all the notes in the current book * \n")
	fmt.Printf("                                        * option: -b *Specify the name of the book to read from * \n\n")
	fmt.Printf("       --> sync                         * Sync notes with Dnote server * \n\n")
	fmt.Printf("       --> login                        * Start a login procedure which will store the APIKey to communicate with the server * \n\n")
	fmt.Printf("       --> help                         * Print a list of commands and what they do * \n\n\n")
}