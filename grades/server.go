package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers() {
	handler := new(studentHandler)
	http.Handle("/students", handler)
	http.Handle("/students/", handler)
}

type studentHandler struct{}

// /students
// /students/{id}
// /students/{id}/grades
func (h studentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2:
		log.Println("get all")
		h.getAll(w, r)
	case 3:
		log.Println("get one")
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.getOne(w, r, id)
	case 4:
		log.Println("add grade")
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h studentHandler) getAll(w http.ResponseWriter, r *http.Request) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	data, err := h.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (h studentHandler) toJSON(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to encode")
	}
	return buf.Bytes(), nil
}

func (h studentHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("get student err: %q", err)
		return
	}
	data, err := h.toJSON(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Final to serialize student: %q", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (h studentHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Final to serialize student: %q", err)
		return
	}
	var g Grade
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to add grade: %q", err)
		return
	}
	student.Grades = append(student.Grades, g)
	w.WriteHeader(http.StatusCreated)
	data, err := h.toJSON(g)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Final to serialize grade: %q", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
