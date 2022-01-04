package main

import (
    "os"
    "fmt"
    "log"
)

func read(path string) {
    file, err := os.Open(path)
    if err != nil {
        fmt.Println("os open err: ", err)
        panic("")
    }


    buffer := make([]byte, 100)
    data, err := file.Read(buffer)
    if err != nil {
        panic("read file err")
    }
    fmt.Printf("%d bytes: %s \n", data, string(buffer))
}

func main() {
    read("src/example1.dat")

    return

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}

	log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])
}
