package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type User struct {
	Login      string
	Id         int
	Avatar_url string
	Name       string
	Email      string
}

const fileName = "users.json"

var users []User

func main() {
	load()

	login := os.Args[1]

	user := findUser(login)
	if user == nil {
		fmt.Println("Bad user")
		return
	}

	setGitConfig("user.name", user.Name)
	setGitConfig("user.email", user.Email)
}

func load() {
	file, _ := os.Open(fileName)
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&users)
}

func findUser(login string) *User {
	user := lUser(login)
	if user == nil {
		fmt.Println("Can't find user on local data.")
		user = dUser(login)
	}
	return user
}

func lUser(login string) *User {
	for i := 0; i < len(users); i++ {
		if strings.EqualFold(login, users[i].Login) {
			return &users[i]
		}
	}
	return nil
}

func dUser(login string) (user *User) {
	request := "https://api.github.com/users/" + login
	fmt.Println("Request " + request)
	resp, err := http.Get(request)
	if err != nil {
		fmt.Println("Can't request the user. Maybe without internet.")
		return nil
	}

	// calls 'resp.Body.Close()' after the ending of this funcion
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("Can't read the json")
		return nil
	}

	go save(user)
	return user
}

func save(user *User) {
	users = append(users, *user)
	file, _ := os.Create(fileName)
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.Encode(&users)
}

func setGitConfig(config string, value string) {
	if len(value) > 0 {
		fmt.Println(value)
		err := exec.Command("git", "config", config, strconv.Quote(value)).Run()
		if err != nil {
			fmt.Printf("Can't change property %s\n", config)
			fmt.Println(err)
		}
	} else {
		exec.Command("git", "config", "--unset", config).Run()
		fmt.Printf("Can't set %s with a empty value.\n", config)
	}
}
