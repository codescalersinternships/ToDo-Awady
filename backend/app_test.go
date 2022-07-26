package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

func MakeTempFile(t testing.TB) string {
	f, err := os.CreateTemp("", "go-sqlite-test")
	defer f.Close()
	if err != nil {
		t.Fatalf("Error making temp file: %q", err.Error())
	}
	return f.Name()
}

func TestGetToDoHandler(t *testing.T) {
	file := MakeTempFile(t)
	defer os.Remove(file)
	a, _ := NewApp(file, LISTENURL, mux.NewRouter())
	a.db.AddToDo("first todo")
	a.db.AddToDo("second todo")
	t.Run("Getting existing ID", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "localhost:5000/todo/1", nil)
		response := httptest.NewRecorder()
		request = mux.SetURLVars(request, map[string]string{"id": "1"})
		a.GetToDoHandler(response, request)
		got := response.Body.String()
		want := "{\"id\":1,\"text\":\"first todo\",\"status\":false}\n"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("Getting non-existing ID", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "localhost:5000/todo/1000", nil)
		response := httptest.NewRecorder()
		request = mux.SetURLVars(request, map[string]string{"id": "1000"})
		a.GetToDoHandler(response, request)
		got := response.Code
		want := http.StatusNotFound
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}

	})
}

func TestGetAllToDosHandler(t *testing.T) {
	file := MakeTempFile(t)
	defer os.Remove(file)
	a, _ := NewApp(file, LISTENURL, mux.NewRouter())
	a.db.AddToDo("first todo")
	a.db.AddToDo("second todo")
	t.Run("Getting all todos", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "localhost:5000/todo", nil)
		response := httptest.NewRecorder()

		a.GetAllToDosHandler(response, request)

		got := response.Body.String()
		want := "[{\"id\":1,\"text\":\"first todo\",\"status\":false},{\"id\":2,\"text\":\"second todo\",\"status\":false}]\n"
		if got != want {
			t.Errorf("got %q want %q", response.Body.String(), want)
		}

	})

}

func TestAddToDoHandler(t *testing.T) {
	file := MakeTempFile(t)
	defer os.Remove(file)
	a, _ := NewApp(file, LISTENURL, mux.NewRouter())
	a.db.AddToDo("first todo")
	a.db.AddToDo("second todo")
	var todoJSON = []byte(`{"text":"new todo"}`)
	t.Run("adding new todo", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "localhost:5000/todo", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()

		a.AddToDoHandler(response, request)

		got := response.Body.String()
		want := "{\"id\":3,\"text\":\"new todo\",\"status\":false}\n"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("sending empty body", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "localhost:5000/todo", nil)
		response := httptest.NewRecorder()

		a.AddToDoHandler(response, request)

		got := response.Code
		want := 500
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}

func TestUpdateToDoHandler(t *testing.T) {
	file := MakeTempFile(t)
	defer os.Remove(file)
	a, _ := NewApp(file, LISTENURL, mux.NewRouter())
	a.db.AddToDo("first todo")
	a.db.AddToDo("second todo")
	var todoJSON = []byte(`{"text":"updating second todo", "status": true}`)
	t.Run("updating second todo with empty body", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "localhost:5000/todo/2", nil)
		response := httptest.NewRecorder()
		request = mux.SetURLVars(request, map[string]string{"id": "2"})

		a.UpdateToDoHandler(response, request)

		got := response.Code
		want := 500
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
	t.Run("updating second todo", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "localhost:5000/todo/2", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		request = mux.SetURLVars(request, map[string]string{"id": "2"})

		a.UpdateToDoHandler(response, request)

		got := response.Body.String()
		want := "{\"id\":2,\"text\":\"updating second todo\",\"status\":true}\n"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("updating non-existing todo", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "localhost:5000/todo/1000", bytes.NewBuffer(todoJSON))
		response := httptest.NewRecorder()
		request = mux.SetURLVars(request, map[string]string{"id": "1000"})

		a.UpdateToDoHandler(response, request)

		got := response.Code
		want := 404
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}

func TestDeleteToDoHandler(t *testing.T) {
	file := MakeTempFile(t)
	defer os.Remove(file)
	a, _ := NewApp(file, LISTENURL, mux.NewRouter())
	a.db.AddToDo("first todo")
	a.db.AddToDo("second todo")
	t.Run("deleting existing todo", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "localhost:5000/todo/2", nil)
		response := httptest.NewRecorder()
		request = mux.SetURLVars(request, map[string]string{"id": "2"})

		a.DeleteToDoHandler(response, request)

		got := response.Code
		want := http.StatusNoContent
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}

	})
	t.Run("deleting non-existing todo", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "localhost:5000/todo/1000", nil)
		response := httptest.NewRecorder()
		request = mux.SetURLVars(request, map[string]string{"id": "1000"})

		a.DeleteToDoHandler(response, request)

		got := response.Code
		want := 404
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}

	})
}
