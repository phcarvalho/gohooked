package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TaskRequest struct {
	Id          string   `json:"id"`
	Url         string   `json:"url"`
	Payload     string   `json:"payload"`
	Headers     []Header `json:"headers"`
	CallbackUrl string   `json:"callbackUrl"`
}

type TaskResult struct {
	Id          string
	Payload     string
	CallbackUrl string
}

type TaskResponse struct {
	Id      string `json:"id"`
	Payload string `json:"payload"`
}

type TaskItem struct {
	Id          string `json:"id"`
	StartedAt   string `json:"startedAt"`
	RunningTime int    `json:"runningTime"`
}

type TaskListResponse struct {
	Count int        `json:"count"`
	Tasks []TaskItem `json:"tasks"`
}

func (app *Application) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	var task TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	app.inputChannel <- task

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Task scheduled")
	return
}

func (app *Application) handleTaskList(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()

	tasks := []TaskItem{}

	for id, startedAt := range app.tasks {
		task := TaskItem{
			Id:          id,
			StartedAt:   startedAt.Format(time.RFC3339),
			RunningTime: int(currentTime.Sub(startedAt) / 1000000),
		}

		tasks = append(tasks, task)
	}

	response := TaskListResponse{
		Count: len(tasks),
		Tasks: tasks,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (app *Application) taskWorker(inputChannel <-chan TaskRequest, outputChannel chan<- TaskResult, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}

	for task := range inputChannel {
		request, err := http.NewRequest(http.MethodPost, task.Url, bytes.NewReader([]byte(task.Payload)))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating the request:", err)
			continue
		}

    request.Header.Set("Content-Type", "application/json")

		for _, header := range task.Headers {
			request.Header.Set(header.Key, header.Value)
		}

		res, err := client.Do(request)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error running task:", err)
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading response body:", err)
			continue
		}
		defer res.Body.Close()

		outputChannel <- TaskResult{
			Id:          task.Id,
			Payload:     string(body),
			CallbackUrl: task.CallbackUrl,
		}
	}
}

func (app *Application) taskResponder(outputChannel <-chan TaskResult) {
	client := &http.Client{}

	for task := range outputChannel {
		response := TaskResponse{
			Id:      task.Id,
			Payload: task.Payload,
		}

		responseJson, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error enconding JSON:", err)
			continue
		}

		request, err := http.NewRequest("POST", task.CallbackUrl, bytes.NewReader(responseJson))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating request:", err)
			continue
		}

    request.Header.Set("Content-Type", "application/json")

		res, err := client.Do(request)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error sending the response:", err)
			continue
		}
		defer res.Body.Close()

		delete(app.tasks, task.Id)
	}
}
