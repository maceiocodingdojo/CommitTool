package gct

//import "github.com/google/go-github/github"
import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"strconv"
	"encoding/json"
	"os/exec"
)

type User struct {
	Login		string
	Id			int
	Avatar_url	string
	Name		string
	Email		string
}

func main() {
	userName := os.Args[1]

	user, err := findUser(userName)
	if err != nil {
		fmt.Println("Bad user")
		return
	}

	setGitConfig("user.name", user.Name)
	setGitConfig("user.email", user.Email)
}

func findUser(userName string) (user *User, err error){
	request := "https://api.github.com/users/"+userName
	fmt.Println("Request "+ request)
	resp, err := http.Get(request)
	if err != nil {
		fmt.Println("Can't request the user. Maybe without internet.")
		return nil, err
	}

	// calls 'resp.Body.Close()' after the ending of this funcion
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("Can't read the json")
		return nil, err
	}
	return user, nil
}

func setGitConfig(config string, value string){
	if len(value) > 0 {
		fmt.Println(value)
		exec.Command("git", "config", config, strconv.Quote(value)).Start()
	} else {
		fmt.Printf("Can't set %s with a empty value.\n", config)
	}
}

