package portal

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

var rootTemplate *template.Template

func ImportTemplates() error {
	var err error
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	studentsFiles := []string{"../../portal/students.html", "./portal/students.html"}
	studentFiles := []string{"../../portal/student.html", "./portal/student.html"}
	var studentsFile string
	var studentFile string
	for _, file := range studentsFiles {
		file = filepath.Join(currentDir, file)
		_, err := os.Stat(file)
		if err == nil {
			studentsFile = file
			break
		}
	}
	for _, file := range studentFiles {
		file = filepath.Join(currentDir, file)
		_, err := os.Stat(file)
		if err == nil {
			fmt.Println(file)
			studentFile = file
			break
		}
	}
	if studentsFile == "" || studentFile == "" {
		panic("file does not exist")
	}
	rootTemplate, err = template.ParseFiles(
		studentsFile,
		studentFile)

	if err != nil {
		return err
	}

	return nil
}
