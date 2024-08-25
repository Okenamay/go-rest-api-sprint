package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Структура задачи:
type Task struct {
	ID          string   `json:"id"`          // ID задачи
	Description string   `json:"description"` // Заголовок
	Note        string   `json:"note"`        // Описание задачи
	Application []string `json:"application"` // Приложения, которыми будете пользоваться
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Application: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Application: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта

// 1. Обработчик для получения всех задач:
func getTasks(w http.ResponseWriter, r *http.Request) {
	// Сериализуем данные из слайса tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		// При ошибке возвращаем статус 500 Internal Server Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// В заголовке записываем тип контента – данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// При успешном запросе возвращаем статус 200 OK
	w.WriteHeader(http.StatusOK)
	// Записываем сериализованные данные в формате JSON в тело ответа
	w.Write(resp)
}

// 2. Обработчик для добавления задачи по телу запроса:
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer
	var checkId string
	// При ошибке возвращаем статус 400 Bad Request
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		// При ошибке возвращаем статус 400 Bad Request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вводим проверку на наличие задачи с указанным в теле запроса ID:
	checkId = task.ID
	_, check := tasks[checkId]
	if check {
		// В случае наличия задачи с указанным в теле запроса ID в мапе возвращаем статус 400 Bad Request
		http.Error(w, "Такое задание уже есть", http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	// При успешном запросе возвращаем статус 201 Created
	w.WriteHeader(http.StatusCreated)
}

// 3. Обработчик для получения задачи по ID:
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		// В случае отсутствия задачи в мапе возвращаем статус 400 Bad Request
		http.Error(w, "Задание не найдено", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		// В случае ошибки возвращаем статус 400 Bad Request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// При успешном выполнении запроса возвращаем статус 200 OK
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		// При наличии ошибки записи выводим ошибку:
		fmt.Println("При записи возникла ошибка:", err)
	}
}

// 4. Обработчик удаления задачи по ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		// В случае отсутствия задачи в мапе возвращаем статус 400 Bad Request
		http.Error(w, "Задание не найдено", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	// При успешном выполнении запроса возвращаем статус 200 OK
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// 1. Регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используем обработчик `getTasks`
	r.Get("/tasks", getTasks)
	// 2. Регистрируем в роутере эндпойнт `/tasks` с методом POST, для которого используем обработчик `postTask`
	r.Post("/tasks", postTask)
	// 3. Регистрируем в роутере эндпойнт `/task/{id}` с методом GET, для которого используем обработчик `getTask`
	r.Get("/task/{id}", getTask)
	// 4. Регистрируем в роутере эндпойнт `/task/{id}` с методом DELETE, для которого используем обработчик `deleteTask`
	r.Delete("/task/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
