package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type location struct {
	door      int
	infoForGo string              // i.e ничего интересного
	whereToGo string              // i.e можно пройти - коридор
	link      map[string]bool     // i.e [комната] = true
	surfaces  map[string][]string // i.e на столе
	init      string
	tasks     []string
	emptySurf string
}
type command func(string) string

const UNKNOWN_COMAND = "unknown command"
const NO_PATH = "no path to "
const INIT_KITCHEN = "you're in the kitchen, "
const INIT_CORRIDOR = "you're in the corridor. "
const NO_STORAGE = "nothere to put"
const ITEM_INSIDE = "item being add to inventory: "
const NO_SUCH = "there's no such"
const EMPTY_ROOM = "empty room. "
const AND = " and"
const NEED = "need to "
const NO_INV = "no such item in the inventory - "
const KEYS = "keys"
const DOOR = "door"
const CANNOT_USE = "no item to use with"
const GO = "go"
const TAKE = "take"
const LOOK = "look"
const USE = "use"
const WEAR = "wear"
const BPON = "you've put a backpack on"
const HOME = "home"

var places map[string]*location
var curPlace string
var surfacesList = []string{"on the table", "on the chair"}
var KITCHEN = []string{
	"kitchen",
	"can pass to - corridor",
	"kitchen, nothing interesting",
}
var CORRIDOR = []string{
	"corridor",
	"can pass to - kitchen, room, outdoors",
	"nothing interesting",
}
var ROOM = []string{
	"room",
	"can pass to - corridor",
	"you're in your room",
}
var STREET = []string{
	"outdoors",
	"can pass to - home",
	"it's spring outdoors",
}
var TABLE_KITCHEN = []string{
	"tea",
}
var TABLE_ROOM = []string{
	"keys",
	"lectures",
}
var CHAIR_ROOM = []string{
	"backpack",
}
var TASKS_KITCH = []string{
	"pack a backpack",
	"go to uni",
}
var doorSt = []string{"", "door shut", "door open"}
var surfaces = []string{"on the table", "on the chair"}
var backpack bool
var inventory map[string]bool
var interfaces map[string][]string
var action map[string]map[string]func() string
var commands map[string]func([]string) string

var keyDoorInteraction = func() string {
	if places[STREET[0]].door == 1 {
		places[STREET[0]].door = 2
	} else if places[STREET[0]].door == 2 {
		places[STREET[0]].door = 1
	}

	return doorSt[places[STREET[0]].door]
}

func main() {
	initGame()
	
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]

		if text[len(text)-1] == '\r' {
			text = text[:len(text)-1]
		}

		if text == "exit" {
			break
		} else {
			fmt.Println(handleCommand(text))
		}

	}
}

func createLocation(door int, init string, emptySurf string, tasks []string, loc []string) {
	places[loc[0]] = &location{door, loc[2], loc[1], make(map[string]bool), make(map[string][]string), init, tasks, emptySurf}
}

func createPlaces() {
	places = make(map[string]*location)

	createLocation(0, INIT_KITCHEN, INIT_KITCHEN, TASKS_KITCH, KITCHEN)
	createLocation(0, INIT_CORRIDOR, "", make([]string, 0), CORRIDOR)
	createLocation(0, "", EMPTY_ROOM, make([]string, 0), ROOM)
	createLocation(1, "", "", make([]string, 0), STREET)

	setStep(ROOM[0], CORRIDOR[0])
	setStep(CORRIDOR[0], KITCHEN[0])
	setStep(CORRIDOR[0], STREET[0])

	setSurfaces(ROOM[0], surfaces[0], TABLE_ROOM)
	setSurfaces(ROOM[0], surfaces[1], CHAIR_ROOM)
	setSurfaces(KITCHEN[0], surfaces[0], TABLE_KITCHEN)
}

func setStep(place1 string, place2 string) {
	places[place1].link[place2] = true
	places[place2].link[place1] = true
}

func setSurfaces(place string, surface string, objects []string) {
	places[place].surfaces[surface] = make([]string, len(objects), len(objects))
	copy(places[place].surfaces[surface], objects)
}

func setActions(what string, where string, f func() string) {
	if len(action[what]) == 0 {
		action[what] = make(map[string]func() string)
	}

	action[what][where] = f
}

func setInterfaces(what string, where string) {
	if len(interfaces[what]) == 0 {
		interfaces[what] = make([]string, 0)
	}

	interfaces[what] = append(interfaces[what], where)
}

func setCommand(command string, f func([]string) string) {
	commands[command] = f
}

func initGame() {
	createPlaces()

	backpack = false

	curPlace = KITCHEN[0]

	action = make(map[string]map[string]func() string)
	setActions(KEYS, DOOR, keyDoorInteraction)

	interfaces = make(map[string][]string)
	setInterfaces(KEYS, DOOR)

	commands = make(map[string]func([]string) string)
	setCommand(GO, commandGo)
	setCommand(LOOK, commandLookAround)
	setCommand(TAKE, commandTake)
	setCommand(WEAR, commandPutOnBackpack)
	setCommand(USE, commandUse)

	inventory = make(map[string]bool)
}

func handleCommand(command string) string {
	words := strings.Split(command, " ")
	flags := false
	for key := range commands {
		if key == words[0] {
			flags = true
			break
		}
	}
	if flags {
		return commands[words[0]](words)
	}
	return UNKNOWN_COMAND
}

func commandTake(words []string) string {
	if !backpack {
		return NO_STORAGE
	}
	item := words[1]
	for idx := range places[curPlace].surfaces {
		for thing := range places[curPlace].surfaces[idx] {
			if places[curPlace].surfaces[idx][thing] == item {
				places[curPlace].surfaces[idx] = deleteStuff(places[curPlace].surfaces[idx], places[curPlace].surfaces[idx][thing])
				inventory[item] = true
				return ITEM_INSIDE + item
			}
		}
	}
	return NO_SUCH
}

func commandPutOnBackpack(words []string) string {
	flag := false
	for idx := range places[curPlace].surfaces {
		for thing := range places[curPlace].surfaces[idx] {
			if places[curPlace].surfaces[idx][thing] == words[1] {
				flag = true
				break
			}
		}
	}
	if !flag {
		return NO_SUCH
	}
	backpack = true
	places[ROOM[0]].surfaces[surfaces[1]] = deleteStuff(places[ROOM[0]].surfaces[surfaces[1]], "рюкзак")
	return BPON
}

func deleteStuff(something []string, object string) []string {
	var x int
	for ind := range something {
		if something[ind] == object {
			x = ind
			break
		}
	}
	newSlice := something[:x]
	newSlice = append(newSlice, something[x+1:]...)
	return newSlice
}

func voidMapSurf(check map[string][]string) bool {
	if len(check) == 0 {
		return false
	}

	return true
}

func voidSlices(check map[string][]string) bool {
	x := true
	for idx := range check {
		x = x && (len(check[idx]) == 0)
	}
	return x
}

func emptySurface(loc string) string {
	return places[loc].emptySurf
}

func printTasks(tasks []string) string {
	strTask := NEED

	for task := range tasks {
		if backpack && task == 0 {
			continue
		}
		strTask += tasks[task] + AND + " "
	}
	strTask = strTask[:len(strTask)-2]
	return strTask
}

func commandLookAround(words []string) string {
	if voidMapSurf(places[curPlace].surfaces) {
		str := places[curPlace].init
		if voidSlices(places[curPlace].surfaces) {
			return emptySurface(curPlace) + places[curPlace].whereToGo
		}
		for key := range surfacesList {
			valExist := places[curPlace].surfaces[surfacesList[key]]
			if len(valExist) != 0 {
				str += surfacesList[key] + ": "
				for surf := range places[curPlace].surfaces[surfacesList[key]] {
					str += places[curPlace].surfaces[surfacesList[key]][surf] + ", "
				}
			}
		}

		if len(places[curPlace].tasks) != 0 {
			str += printTasks(places[curPlace].tasks)
		}

		str = str[:len(str)-2]
		str += ". " + places[curPlace].whereToGo

		return str
	} else {
		return places[curPlace].init + places[curPlace].whereToGo
	}
}

func commandGo(words []string) string {
	if words[1] == HOME && curPlace == STREET[0] {
		words[1] = CORRIDOR[0]
	}
	flag := false
	for idx := range places {
		if words[1] == idx {
			flag = true
			break
		}
	}
	if !flag {
		return NO_PATH + words[1]
	}
	if places[words[1]].door == 1 {
		return doorSt[1]
	}
	if places[curPlace].link[words[1]] {
		curPlace = words[1]
		return places[words[1]].infoForGo + ". " + places[words[1]].whereToGo
	}
	return NO_PATH + words[1]
}

func use(what string, where string) string {
	return action[what][where]()
}

func commandUse(words []string) string {
	if !inventory[words[1]] {
		return NO_INV + words[1]
	}
	for idx := range interfaces[words[1]] {
		if interfaces[words[1]][idx] == words[2] {
			return use(words[1], words[2])
		}
	}
	return CANNOT_USE
}
