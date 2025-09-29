package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
)


//================ Fields ===================================

type BudgetItem struct {
	Id int `json:"id"`
	ItemName string `json:"itemName"`
	ItemValue float64 `json:"itemValue"`
	Cols int `json:"cols"`
	Rows int `json:"rows"`
	Color string `json:"color"`
}

type BudgetSection struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Items []BudgetItem `json:"items"`
}

var items = []BudgetItem { 
	{Id: 2, ItemName: "Mortgage", ItemValue: 3000.69, Cols: 1, Rows: 1, Color: "lightblue"},
	{Id: 3, ItemName: "Internet", ItemValue: 70.99, Cols: 1, Rows: 1, Color: "lightgreen"},
}

var section = BudgetSection {
	Id: 1, 
	Title: "Housing",
	Items: items,
}

var incomeSection = BudgetSection {
	Id: 4,
	Title: "Income",
	Items: make([]BudgetItem, 0),
}

var sections = []BudgetSection {incomeSection, section}

var maxId int


//================ Functions ===================================

func getSections(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, World!")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sections)
}

func addNewSection(w http.ResponseWriter, r *http.Request) {
	var newSection BudgetSection
	err := json.NewDecoder(r.Body).Decode(&newSection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close() 

	maxId++
	newSection.Id = maxId
	sections = append(sections, newSection)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sections)
}

func editSection(w http.ResponseWriter, r *http.Request) {
	sectionId := r.URL.Path[len("/plan/"):]

	if sectionId == "" {
		http.Error(w, "Section ID is required", http.StatusBadRequest)
		return
	}

	id_int, _ := strconv.Atoi(sectionId)
	index := findSectionIndexById(id_int)

	if index == -1 {
		http.Error(w, "Section ID not found", http.StatusNotFound)
		return
	}

	var sectionToReplace BudgetSection
	err := json.NewDecoder(r.Body).Decode(&sectionToReplace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close() 

	sections[index] = sectionToReplace

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sections)
}

func deleteSection(w http.ResponseWriter, r *http.Request) {
	sectionId := r.URL.Path[len("/plan/"):]

	if sectionId == "" {
		http.Error(w, "Section ID is required", http.StatusBadRequest)
		return
	}

	id_int, _ := strconv.Atoi(sectionId)
	indexToRemove := findSectionIndexById(id_int)

	if indexToRemove == -1 {
		http.Error(w, "Section ID not found", http.StatusNotFound)
		return
	}

	if sections[indexToRemove].Title == "Income" {
		http.Error(w, "Cannot delete Income section", http.StatusBadRequest)
		return
	}

	sections = append(sections[:indexToRemove], sections[indexToRemove+1:]...)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sections)
}

func findSectionIndexById(sectionId int) int {
	for i, section := range sections {
		if section.Id == sectionId {
			return i
		}
	}

	return -1
}


func addNewItem(w http.ResponseWriter, r *http.Request) {
	sectionId := r.URL.Path[len("/item/"):]

	if sectionId == "" {
		http.Error(w, "Section ID is required", http.StatusBadRequest)
		return
	}

	id_int, _ := strconv.Atoi(sectionId)
	index := findSectionIndexById(id_int)

	if index == -1 {
		http.Error(w, "Section ID not found", http.StatusNotFound)
		return
	}

	var newItem BudgetItem
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close() 

	maxId++
	newItem.Id = maxId

	sections[index].Items = append(sections[index].Items, newItem)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sections)
}

func getIncome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incomeSection)
}


//================ Handlers ===================================


func handlePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	switch r.Method {
		case "GET":
			getSections(w, r)
		case "POST":
			addNewSection(w, r)
		case "PUT":
			editSection(w, r)
		case "DELETE":
			deleteSection(w, r)
		case "OPTIONS":
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func handleItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	switch r.Method {
	case "POST":
		addNewItem(w,r)
	case "OPTIONS":
			w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func handleIncome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	switch r.Method {
	case "GET":
		getIncome(w,r)
	case "OPTIONS":
			w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}


//================ Main ===================================

func main() {
	maxId = 3

	http.HandleFunc("/plan/", handlePlan)
	http.HandleFunc("/plan", handlePlan)

	http.HandleFunc("/item/", handleItem)

	http.HandleFunc("/income", handleIncome)

	fmt.Println("Server starting on port 4321...")

	if err := http.ListenAndServe(":4321", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}