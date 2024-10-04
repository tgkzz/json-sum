package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

const (
	objectNumber = 1000000
	filename     = "./test/zxc.json"
)

type Object struct {
	A int `json:"a"`
	B int `json:"b"`
}

func generateObj(numObjects int) []Object {

	objects := make([]Object, numObjects)
	for i := 0; i < numObjects; i++ {
		objects[i] = Object{
			A: rand.Intn(21) - 10,
			B: rand.Intn(21) - 10,
		}
	}
	return objects
}

func saveObj(objects []Object) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(objects); err != nil {
		return err
	}

	return nil
}

func countSum(objs []Object) int {
	s := 0
	for _, obj := range objs {
		s += obj.A + obj.B
	}
	return s
}

// here we generate tests
func main() {
	objects := generateObj(objectNumber)

	if err := saveObj(objects); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(countSum(objects))
}

// notes
// asd.json sum = 212
// qwe.json sum = -352
// zxc.json sum = -4482
