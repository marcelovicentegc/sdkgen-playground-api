package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"os"
	"os/exec"
)

// Code struct defines transactioned JSON
type Code struct {
	Sdkgen              string
	Target              string
	TargetFileExtension string
}

func main() {
	http.HandleFunc("/gen", gen)
	http.HandleFunc("/example", example)
	log.Println("Server is up and running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func example(responseWriter http.ResponseWriter, request *http.Request) {
	exampleCode := Code{Sdkgen: `type User {
	firstName: string
	lastName: string
	email: string
	password: string
	cpf: string
	birthDate: datetime
	gender: string
	status: string
	address: Address
	profile: Profile
}
	
type Message {
	date: date
	author: User
	mentions: User[]
	text: string?
}
	
type Address {
	countryCode: string
	stateCode: string
	city: string
	neighborhood: string
	street: string
	number: string
	complement: string
}
	
type Profile {
	username: string
	photoUrl: string
}`, Target: ``}

	encodedExampleCode, error := json.Marshal(exampleCode)

	if error != nil {
		http.Error(responseWriter, error.Error(), http.StatusInternalServerError)
		return
	}

	enableCors(&responseWriter)
	responseWriter.Header().Set("Content-Type", "application/json")

	responseWriter.Write(encodedExampleCode)
}

func gen(responseWriter http.ResponseWriter, request *http.Request) {
	enableCors(&responseWriter)
	requestBody := json.NewDecoder(request.Body)

	var code Code

	error := requestBody.Decode(&code)

	if error != nil {
		http.Error(responseWriter, error.Error(), http.StatusInternalServerError)
		return
	}

	sdkgen := code.Sdkgen
	target := code.Target
	targetFileExtension := code.TargetFileExtension

	sdkgenFile := createFile("playground.sdkgen")

	defer closeFile(sdkgenFile)

	writeFile(sdkgenFile, sdkgen)

	command := exec.Command("bash", "sdkgen.sh", targetFileExtension, target)
	currentDir := getCurrentDir()
	command.Dir = currentDir

	_, err := command.Output()

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	targetCodeFileBytes, err := ioutil.ReadFile("gen/playground." + targetFileExtension)
	target = string(targetCodeFileBytes)

	encodedTargetCode, err := json.Marshal(target)

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")

	responseWriter.Write(encodedTargetCode)
}

func createFile(path string) *os.File {
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	return file
}

func writeFile(file *os.File, content string) {
	fmt.Fprintln(file, content)
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
		return
	}
}

func enableCors(responseWriter *http.ResponseWriter) {
	(*responseWriter).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*responseWriter).Header().Set("Access-Control-Allow-Credentials", "true")
}

func getCurrentDir() string {
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	return dir
}
