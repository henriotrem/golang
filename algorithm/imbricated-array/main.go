package main

import "fmt"

type ImbricatedArray []interface{}

func (array *ImbricatedArray) Display() {
	for _, element := range *array {
		switch value := element.(type) {
		case int:
			fmt.Println(value)
		case ImbricatedArray:
			value.Display()
		}
	}
}

func main() {
	imbricatedArray := ImbricatedArray{2, 4, ImbricatedArray{1, 4}, ImbricatedArray{1, ImbricatedArray{1, 4}}}
	imbricatedArray.Display()
}
