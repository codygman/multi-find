/*
	This program takes a list of search terms via stdin then concurrently searches them with find.
	Useful if you want a seperated list of results. If not use:
	find \( -iname "*py" -o -iname "*jpg" -o -iname "*txt" \)


*/
package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"flag"
	"runtime"
)

func process(directory string, searchTerm string, output chan []byte) {
	out, err := exec.Command("find", directory, "-type", "f", "-iname", searchTerm).Output()
	if err != nil {
		log.Fatal(fmt.Sprintf("Find error: %s", err))
	}
	output <- out
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	cwd, _ := os.Getwd()
	dirPtr := flag.String("dir", cwd, "Directory to search. Default is current working directory")
	flag.Parse()

	searchTerms := flag.Args()

	input := make(chan string, len(searchTerms))
	output := make(chan []byte, len(searchTerms))
	
	for _, term := range searchTerms {
		input <- term
	}
	// input is closed
	close(input)

	for searchTerm := range input {
		go process(*dirPtr, searchTerm, output)
		fmt.Println(string(<-output))
	}
}

