package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Item представляет товар на складе
type Item struct {
	ID       int
	Name     string
	Category string
	Price    float64
	Quantity int
}

// Inventory хранит список товаров и предустановленные категории
type Inventory struct {
	items      []Item
	categories []string
	nextID     int
}

func NewInventory() *Inventory {
	return &Inventory{
		items: []Item{
			{ID: 1, Name: "Ноутбук ASUS", Category: "Электроника", Price: 350000, Quantity: 8},
			{ID: 2, Name: "Беспроводная мышь", Category: "Аксессуары", Price: 12000, Quantity: 25},
			{ID: 3, Name: "Механическая клавиатура", Category: "Аксессуары", Price: 28000, Quantity: 14},
			{ID: 4, Name: "Монитор 27 дюймов", Category: "Электроника", Price: 145000, Quantity: 3},
		},
		categories: []string{
			"Электроника",
			"Аксессуары",
			"Бытовая техника",
		},
		nextID: 5,
	}
}

// GetCategories возвращает список доступных категорий
func (inv *Inventory) GetCategories() []string {
	return inv.categories
}

// AddCategory добавляет новую категорию, если её ещё нет в списке
func (inv *Inventory) AddCategory(newCat string) {
	newCat = strings.TrimSpace(newCat)
	if newCat == "" {
		return
	}
	for _, cat := range inv.categories {
		if strings.EqualFold(cat, newCat) {
			return
		}
	}
	inv.categories = append(inv.categories, newCat)
}

// AddItem добавляет новый товар с валидацией
func (inv *Inventory) AddItem(name, category string, price float64, qty int) {
	item := Item{
		ID:       inv.nextID,
		Name:     strings.TrimSpace(name),
		Category: strings.TrimSpace(category),
		Price:    price,
		Quantity: qty,
	}
	inv.items = append(inv.items, item)
	inv.nextID++
	fmt.Printf("Товар #%d ('%s') успешно добавлен на склад!\n", item.ID, item.Name)
}

// PrintAll выводит все товары склада
func (inv *Inventory) PrintAll() {
	if len(inv.items) == 0 {
		fmt.Println("Склад пуст.")
		return
	}
	fmt.Println("\n СПИСОК ВСЕХ ТОВАРОВ НА СКЛАДЕ ")
	fmt.Printf("%-4s | %-24s | %-14s | %-12s | %-8s | %s\n", "ID", "Название", "Категория", "Цена (тг)", "Остаток", "Сумма (тг)")
	fmt.Println(strings.Repeat("-", 80))
	for _, item := range inv.items {
		total := item.Price * float64(item.Quantity)
		fmt.Printf("#%-3d | %-24s | %-14s | %10.2f тг | %5d шт | %10.2f тг\n",
			item.ID, item.Name, item.Category, item.Price, item.Quantity, total)
	}
	fmt.Println(strings.Repeat("-", 80))
}

// UpdateQuantity обновляет количество товара по ID
func (inv *Inventory) UpdateQuantity(id int, newQty int) bool {
	for i := range inv.items {
		if inv.items[i].ID == id {
			inv.items[i].Quantity = newQty
			return true
		}
	}
	return false
}

// CalculateTotalValue вычисляет общую стоимость всех товаров
func (inv *Inventory) CalculateTotalValue() (totalCost float64, totalItems int) {
	for _, item := range inv.items {
		totalCost += item.Price * float64(item.Quantity)
		totalItems += item.Quantity
	}
	return
}

// FilterByCategory выводит товары конкретной категории
func (inv *Inventory) FilterByCategory(cat string) {
	normCat := strings.ToLower(strings.TrimSpace(cat))
	found := false
	fmt.Printf("\n ТОВАРЫ КАТЕГОРИИ '%s' \n", cat)
	for _, item := range inv.items {
		if strings.ToLower(item.Category) == normCat {
			fmt.Printf("#%-3d | %-24s | %10.2f тг | %d шт\n", item.ID, item.Name, item.Price, item.Quantity)
			found = true
		}
	}
	if !found {
		fmt.Printf("В категории '%s' товаров не найдено.\n", cat)
	}
}

// FindLowStock выводит товары с остатком меньше порогового значения
func (inv *Inventory) FindLowStock(threshold int) {
	found := false
	fmt.Printf("\n ТОВАРЫ С НИЗКИМ ОСТАТКОМ (< %d шт) \n", threshold)
	for _, item := range inv.items {
		if item.Quantity < threshold {
			fmt.Printf("Внимание! #%-3d %-24s — осталось всего %d шт (Категория: %s)\n",
				item.ID, item.Name, item.Quantity, item.Category)
			found = true
		}
	}
	if !found {
		fmt.Printf("Товаров с остатком менее %d шт нет. Склад укомплектован!\n", threshold)
	}
}

func main() {
	inv := NewInventory()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n INVENTORY MANAGER ")
		fmt.Println("1. Добавить товар")
		fmt.Println("2. Список всех товаров")
		fmt.Println("3. Изменить количество товара")
		fmt.Println("4. Общая стоимость склада")
		fmt.Println("5. Поиск по категории")
		fmt.Println("6. Товары с низким остатком")
		fmt.Println("0. Выход")
		fmt.Print("Выберите действие: ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("Название товара: ")
			scanner.Scan()
			name := scanner.Text()

			// Выбор категории из предустановленного списка
			cats := inv.GetCategories()
			fmt.Println("\nВыберите категорию товара:")
			for i, c := range cats {
				fmt.Printf("%d. %s\n", i+1, c)
			}
			fmt.Printf("%d. [+ Добавить новую категорию]\n", len(cats)+1)
			fmt.Print("Введите номер категории: ")

			scanner.Scan()
			catChoice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			var category string

			if err != nil || catChoice < 1 || catChoice > len(cats)+1 {
				fmt.Println("Ошибка: неверный выбор категории!")
				continue
			}

			if catChoice == len(cats)+1 {
				fmt.Print("Введите название новой категории: ")
				scanner.Scan()
				category = strings.TrimSpace(scanner.Text())
				if category == "" {
					fmt.Println("Ошибка: название категории не может быть пустым!")
					continue
				}
				inv.AddCategory(category)
			} else {
				category = cats[catChoice-1]
			}

			fmt.Print("Цена за единицу (тг): ")
			scanner.Scan()
			price, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)
			if err != nil || price <= 0 {
				fmt.Println("Ошибка: некорректная цена!")
				continue
			}

			fmt.Print("Количество на складе: ")
			scanner.Scan()
			qty, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil || qty < 0 {
				fmt.Println("Ошибка: некорректное количество!")
				continue
			}

			inv.AddItem(name, category, price, qty)

		case "2":
			inv.PrintAll()

		case "3":
			fmt.Print("Введите ID товара для изменения: ")
			scanner.Scan()
			id, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				fmt.Println("Ошибка: некорректный ID!")
				continue
			}

			fmt.Print("Введите новое количество: ")
			scanner.Scan()
			newQty, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil || newQty < 0 {
				fmt.Println("Ошибка: количество не может быть отрицательным!")
				continue
			}

			if inv.UpdateQuantity(id, newQty) {
				fmt.Printf("Количество товара #%d успешно обновлено на %d шт.\n", id, newQty)
			} else {
				fmt.Printf("Товар с ID %d не найден!\n", id)
			}

		case "4":
			cost, count := inv.CalculateTotalValue()
			fmt.Println("\n СВОДКА ПО СКЛАДУ ")
			fmt.Printf("Всего наименований: %d\n", len(inv.items))
			fmt.Printf("Всего единиц продукции: %d шт\n", count)
			fmt.Printf("Общая стоимость склада: %.2f тг\n", cost)

		case "5":
			cats := inv.GetCategories()
			if len(cats) == 0 {
				fmt.Println("Список категорий пуст.")
				continue
			}

			fmt.Println("\nВыберите категорию для поиска:")
			for i, c := range cats {
				fmt.Printf("%d. %s\n", i+1, c)
			}
			fmt.Print("Введите номер категории: ")

			scanner.Scan()
			catChoice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil || catChoice < 1 || catChoice > len(cats) {
				fmt.Println("Ошибка: неверный выбор категории!")
				continue
			}

			selectedCategory := cats[catChoice-1]
			inv.FilterByCategory(selectedCategory)

		case "6":
			fmt.Print("Введите пороговое значение (например, 5): ")
			scanner.Scan()
			threshold, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil || threshold <= 0 {
				threshold = 5 // значение по умолчанию
			}
			inv.FindLowStock(threshold)

		case "0":
			fmt.Println("Завершение работы программы Inventory Manager.")
			return

		default:
			fmt.Println("Неверный выбор. Введите пункт меню от 0 до 6.")
		}
	}
}
