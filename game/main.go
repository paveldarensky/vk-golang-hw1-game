package main

import (
	"fmt"
	"strings"
)

/*
	код писать в этом файле
	наверняка у вас будут какие-то структуры с методами, глобальные переменные ( тут можно ), функции
*/

// ТИПЫ/АЛИАСЫ
type Item string

// СТРУКТУРЫ МОДЕЛИРОВАНИЯ
type Furniture struct {
	Name  string   // например: "на столе", "на стуле"
	Items []string // слайс предметов
}

type Room struct {
	Name        string           // имя комнаты
	Description string           // описание комнаты
	Exits       map[string]*Room // карта переходов: "коридор" -> &Room{Name: "коридор"}
	ExitsOrder  []string         // фиксированный порядок выходов для тестов
	Furniture   []Furniture      // список мебели в этой комнате
	TaskPrompt  string           // доп. текст задачи
}

type Player struct {
	CurrentRoom *Room    // текущая комната, где мы
	HasBackpack bool     // надет ли рюкзак?
	Inventory   []string // инвентарь с предметами
}

// ГЛОБАЛЬНОЕ СОСТОЯНИЕ
var (
	player   Player // игрок текущий
	doorOpen bool   // дверь открыта?
)

func main() {
	/*
		в этой функции можно ничего не писать,
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/

	initGame()

	_ = handleCommand("some command")
}

// ИНИЦИАЛИЗАЦИЯ
func initGame() {
	/*
		эта функция инициализирует игровой мир - все комнаты
		если что-то было - оно корректно перезатирается
	*/

	doorOpen = false

	// 1. создать комнаты
	kitchenRoom := &Room{
		Name:        "кухня",
		Description: "ты находишься на кухне",
		TaskPrompt:  "надо собрать рюкзак и идти в универ",
		Exits:       make(map[string]*Room),
		ExitsOrder:  []string{"коридор"},
		Furniture: []Furniture{
			{
				Name:  "на столе",
				Items: []string{"чай"},
			},
		},
	}

	lobbyRoom := &Room{
		Name:        "коридор",
		Description: "ничего интересного",
		Exits:       make(map[string]*Room),
		ExitsOrder:  []string{"кухня", "комната", "улица"},
	}

	// пусть будет спальня = комната игрока, чтоб не путаться с переменными
	bedroom := &Room{
		Name:        "комната",
		Description: "ты в своей комнате",
		Exits:       make(map[string]*Room),
		ExitsOrder:  []string{"коридор"},
		Furniture: []Furniture{
			{
				Name:  "на столе",
				Items: []string{"ключи", "конспекты"},
			},
			{
				Name:  "на стуле",
				Items: []string{"рюкзак"},
			},
		},
	}

	streetRoom := &Room{
		Name:        "улица",
		Description: "на улице весна",
		Exits:       make(map[string]*Room),
		ExitsOrder:  []string{"домой"},
	}

	// 2. связать комнаты переходами
	kitchenRoom.Exits["коридор"] = lobbyRoom

	lobbyRoom.Exits["кухня"] = kitchenRoom
	lobbyRoom.Exits["комната"] = bedroom
	lobbyRoom.Exits["улица"] = streetRoom

	bedroom.Exits["коридор"] = lobbyRoom

	streetRoom.Exits["домой"] = lobbyRoom

	// 3. игрок
	player = Player{
		CurrentRoom: kitchenRoom,
		HasBackpack: false,
		Inventory:   []string{},
	}

}

// ОБРАБОТКА КОМАНД
func handleCommand(command string) string {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/
	parts := strings.Split(command, " ")
	action := parts[0]

	switch action {
	case "осмотреться":
		return player.CurrentRoom.LookAround()

	case "идти":
		if len(parts) < 2 {
			return "куда идти?"
		}
		return goTo(parts[1])

	case "надеть":
		if len(parts) < 2 {
			return "что надеть?"
		}
		return putOn(parts[1])

	case "взять":
		if len(parts) < 2 {
			return "что взять?"
		}
		return takeItem(parts[1])

	case "применить":
		if len(parts) < 3 {
			return "недостаточно аргументов"
		}
		return applyItem(parts[1], parts[2])

	default:
		return "неизвестная команда"
	}
}

func goTo(target string) string {
	// проверка закрытой двери
	if player.CurrentRoom.Name == "коридор" && target == "улица" && !doorOpen {
		return "дверь закрыта"
	}

	// проверка на возможность перехода из текущей комнаты в целевую
	nextRoom, exists := player.CurrentRoom.Exits[target]
	if !exists {
		return fmt.Sprintf("нет пути в %s", target)
	}

	// иначе, если имеется переход, то перемещаем игрока
	player.CurrentRoom = nextRoom

	if nextRoom.Name == "кухня" {
		return "кухня, ничего интересного. можно пройти - коридор"
	}

	// формируем ответ согласно тестам
	return formatRoomMessage(nextRoom)
}

func formatRoomMessage(r *Room) string {
	// список доступных выходов
	exitsList := getExitsString(r)
	return fmt.Sprintf("%s. можно пройти - %s", r.Description, exitsList)
}

func getExitsString(r *Room) string {
	// собираем выходы через запятую
	return strings.Join(r.ExitsOrder, ", ")
}

func (r *Room) LookAround() string {
	var furnitureParts []string

	for _, f := range r.Furniture {
		// если на этой мебели есть хотя бы один предмет
		if len(f.Items) > 0 {
			// объединяем предметы через запятую: "ключи, конспекты"
			itemsStr := strings.Join(f.Items, ", ")
			// формируем кусок: "на столе: ключи, конспекты"
			furnitureParts = append(furnitureParts, fmt.Sprintf("%s: %s", f.Name, itemsStr))
		}
	}

	exitsStr := fmt.Sprintf("можно пройти - %s", getExitsString(r))

	// если коридор
	if r.Name == "коридор" {
		return fmt.Sprintf("%s. %s", r.Description, exitsStr)
	}

	// если всё пусто (нет предметов)
	if len(furnitureParts) == 0 {
		return fmt.Sprintf("пустая комната. %s", exitsStr)
	}

	itemsJoined := strings.Join(furnitureParts, ", ")

	// динамическая фраза для кухни (зависит от рюкзака)
	if r.Name == "кухня" {
		task := "надо идти в универ"
		if !player.HasBackpack {
			task = "надо собрать рюкзак и идти в универ"
		}
		return fmt.Sprintf("%s, %s, %s. %s", r.Description, itemsJoined, task, exitsStr)
	}

	// обычная комната с предметами
	return fmt.Sprintf("%s. %s", itemsJoined, exitsStr)
}

func putOn(item string) string {
	if item == "рюкзак" {
		// проверяем, есть ли рюкзак в комнате
		if player.CurrentRoom.TakeItem("рюкзак") {
			player.HasBackpack = true
			return "вы надели: рюкзак"
		}
	}
	return "нет такого"
}

// TakeItem проверяет, есть ли предмет в комнате, и убирает его, если found == true
func (r *Room) TakeItem(itemName string) bool {
	for i := range r.Furniture {
		for j, item := range r.Furniture[i].Items {
			if item == itemName {
				// удаляем элемент из слайса
				r.Furniture[i].Items = append(r.Furniture[i].Items[:j], r.Furniture[i].Items[j+1:]...)
				return true
			}
		}
	}
	return false
}

func takeItem(item string) string {
	// проверка на рюкзак
	if !player.HasBackpack {
		return "некуда класть"
	}

	// пытаемся забрать из комнаты
	if player.CurrentRoom.TakeItem(item) {
		player.Inventory = append(player.Inventory, item)
		return fmt.Sprintf("предмет добавлен в инвентарь: %s", item)
	}

	return "нет такого"
}

func applyItem(item, object string) string {
	// проверяем, есть ли предмет в инвентаре
	hasItem := false
	for _, invItem := range player.Inventory {
		if invItem == item {
			hasItem = true
			break
		}
	}

	if !hasItem {
		return fmt.Sprintf("нет предмета в инвентаре - %s", item)
	}

	// применение ключей на дверь
	if item == "ключи" && object == "дверь" {
		doorOpen = true
		return "дверь открыта"
	}

	// применение к нерасчетному объекту
	return "не к чему применить"
}
