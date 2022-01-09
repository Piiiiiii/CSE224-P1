package main

import (
	"fmt"
	"log"
	"os"
	"sort"
    "bytes"
)

const keySize int64 = 10
const valueSize int64 = 90
const kvSize int64 = keySize + valueSize

type KeyValue struct {
	key   []byte
	value []byte
}

func read(path string) ([]byte, int64, error) {
	stats, err := os.Stat(path)
	if err == nil {
		log.Printf("found file: %s \n", path)
	} else if os.IsNotExist(err) {
		err = fmt.Errorf("file: %s does not exist \n", path)
		return []byte{}, 0, err
	} else {
		err = fmt.Errorf("get stats of file %s err: %v \n", path, err)
		return []byte{}, 0, err
	}
	size := stats.Size()
	// size = (size/kvSize)*kvSize + kvSize

	file, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("os open %s err: %v \n", path, err)
		return []byte{}, 0, err
	}

	buffer := make([]byte, size)
	num, err := file.Read(buffer)
	if err != nil {
		err = fmt.Errorf("read %s err: %v \n", path, err)
		return []byte{}, 0, err
	}

	return buffer, int64(num), nil
}

func getKVSlice(file []byte, size int64) (KVList, error) {
	//if size%kvSize != 0 {
	//	err := fmt.Errorf("size % 100 = %d != 0", size)
	//	return []KeyValue{}, err
	//}
	var index int64 = 0
	var KVSlice []KeyValue
	for index+kvSize <= size {
		key := file[index : index+keySize]
		value := file[index+keySize : index+kvSize]
		kv := KeyValue{key, value}
		KVSlice = append(KVSlice, kv)

		index += kvSize
	}
	return KVSlice, nil
}

func getByteSlice(kvList KVList) []byte {
    var byteSlice []byte
    for _, kv := range kvList {
        byteSlice = append(byteSlice, append(kv.key, kv.value...)...)
    }
    return byteSlice
}

type KVList []KeyValue

func (s KVList) Len() int {
	return len(s)
}

func (s KVList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s KVList) Less(i, j int) bool {
	return bytes.Compare(s[i].key, s[j].key) < 0
	// key1, key2 := s[i].key, s[j].key
	// for index := int64(0); index < keySize; index++ {
	// 	if key1[index] < key2[index] {
	// 		return true
	// 	} else if key1[index] > key2[index] {
	// 		return false
	// 	}
	// }
	// return false
}

func (s KVList) PrintKeys() {
	for i := 0; i < len(s); i++ {
		log.Println(s[i].key)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	log.Printf("Sorting %s to %s\n", inputPath, outputPath)

	file, size, err := read(inputPath)
	if err != nil {
		panic(err)
	}
	log.Printf("read file: %s success. size is %v \n", inputPath, size)

	kvList, err := getKVSlice(file, size)
	if err != nil {
		panic(err)
	}

	kvList.PrintKeys()
	sort.Sort(kvList)
	log.Println()
	kvList.PrintKeys()

    byteSlice := getByteSlice(kvList)
    
    outputFile, err := os.OpenFile(
        outputPath,
        os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
        0666,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer outputFile.Close()

    bytesWritten, err := outputFile.Write(byteSlice)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Wrote %d bytes into %s.\n", bytesWritten, outputPath)

}
