package main

/*
	код писать в этом файле
	наверняка у вас будут какие-то структуры с методами, глобальные перменные ( тут можно ), функции
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/sergeymirzoev/lectures-2021-2/1/99_hw/game/inventory"
	"gitlab.com/sergeymirzoev/lectures-2021-2/1/99_hw/game/room"
)

var rooms map[string]room.RoomIFace
var commands map[string]string
var currentRoom room.RoomIFace
var myInventory *inventory.InventoryImpl

func main() {
	/*
		в этой функции можно ничего не писать
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/
	initGame()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result := handleCommand(scanner.Text())
		fmt.Printf("%s\n", result)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error occured while reading from os.Stdin: %v\n", err)
		return
	}
}

//nolint:deadcode,unused
func initGame() {
	/*
		эта функция инициализирует игровой мир - все команты
		если что-то было - оно корректно перезатирается
	*/

	rooms = map[string]room.RoomIFace{
		"комната": &room.LivingRoom{
			RoomBase: room.RoomBase{
				Connections: []string{
					"коридор",
				},
				Items: map[string]inventory.Item{
					"ключи": {
						Location: "на столе",
						Name:     "ключи",
						Applications: []string{
							"дверь",
						},
					},
					"конспекты": {
						Location: "на столе",
						Name:     "конспекты",
					},
					"рюкзак": {
						Location: "на стуле",
						Name:     "рюкзак",
					},
				},
				Name: "комната",
			},
		},
		"коридор": &room.Corridor{
			RoomBase: room.RoomBase{
				Connections: []string{
					"кухня",
					"комната",
					"улица",
				},
				Name:       "коридор",
				DoorClosed: true,
			},
		},
		"кухня": &room.Kitchen{
			RoomBase: room.RoomBase{
				Connections: []string{
					"коридор",
				},
				Items: map[string]inventory.Item{
					"чай": {
						Location: "на столе",
						Name:     "чай",
					},
				},
				Name: "кухня",
			},
			Tasks: []string{
				"собрать рюкзак",
				"идти в универ",
			},
		},
		"улица": &room.Street{
			RoomBase: room.RoomBase{
				Connections: []string{
					"коридор",
				},
				Name: "улица",
			},
		},
	}

	commands = map[string]string{
		"осмотреться": "lookAround",
		"идти":        "walkTo",
		"надеть":      "take",
		"взять":       "take",
		"применить":   "use",
	}

	currentRoom = rooms["кухня"]

	myInventory = &inventory.InventoryImpl{
		Items: make(map[string]inventory.Item, 10),
	}
}

func handleCommand(command string) string {
	result, err := handleCommandKeepError(command)
	if err != nil {
		result = fmt.Sprint(err)
	}

	return result
}

func walkTo(command string, parameters []string) (string, error) {
	if len(parameters) == 1 {
		hasPath, err := currentRoom.HasPath(parameters[0])
		if hasPath {
			currentRoom = rooms[parameters[0]]
			return currentRoom.Greet(), nil
		} else {
			return "", err
		}
	} else {
		return "", fmt.Errorf("WrongNumberOfArgsError: command \"%s\" supports 1 param, %d given", command, len(parameters))
	}
}

func use(command string, parameters []string) (string, error) {
	if len(parameters) == 2 {
		if myInventory.Check(parameters[0]) {
			item, err := myInventory.Get(parameters[0])
			if err != nil {
				return "", err
			}

			hasItem := currentRoom.HasItem(parameters[1]) || parameters[1] == "дверь"
			if !hasItem {
				return "", fmt.Errorf("не к чему применить")
			}

			err = item.Use(parameters[1])
			if err != nil {
				return "", err
			}

			if parameters[1] == "дверь" {
				var usedDoorStatus string
				usedDoorStatus, err = currentRoom.UseDoor()
				return usedDoorStatus, err
			}
		} else {
			return "", fmt.Errorf("нет предмета в инвентаре - %s", parameters[0])
		}
	} else {
		return "", fmt.Errorf("WrongNumberOfArgsError: command \"%s\" supports 2 param, %d given", command, len(parameters))
	}

	return "", fmt.Errorf("reached unimplemeted")
}

func take(command string, parameters []string) (string, error) {
	if len(parameters) == 1 {
		hasItem := currentRoom.HasItem(parameters[0])
		if !hasItem {
			return "", fmt.Errorf("нет такого")
		}

		item, err := currentRoom.TakeItem(parameters[0])
		if err != nil {
			return "", err
		}

		if parameters[0] == "рюкзак" {
			err = rooms["кухня"].RemoveTask("собрать рюкзак")
			if err != nil {
				currentRoom.PutItem(item)
				return "", err
			}
		}

		s, err := myInventory.Add(item)
		if err != nil {
			currentRoom.PutItem(item)
		}
		return s, err
	} else {
		return "", fmt.Errorf("WrongNumberOfArgsError: command \"%s\" supports 1 param, %d given", command, len(parameters))
	}
}

func handleCommandKeepError(command string) (string, error) {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/
	args := strings.Split(command, " ")
	command = commands[args[0]]
	parameters := args[1:]

	if command == "" {
		return "", fmt.Errorf("неизвестная команда")
	}

	switch command {
	case "lookAround":
		return currentRoom.LookAround(), nil

	case "walkTo":
		return walkTo(command, parameters)

	case "use":
		return use(command, parameters)

	case "take":
		return take(command, parameters)
	}

	return "", fmt.Errorf("reached unimplemeted")
}
