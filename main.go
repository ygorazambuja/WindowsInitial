package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func installChocolatey() {
	chocolateyCommand := "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))"

	if c, err := exec.Command("powershell", chocolateyCommand).CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", c)
	}
}

func runScripts() {

	apps := fetchAppList()

	for element := apps.Front(); element != nil; element = element.Next() {
		runScript(*element)
	}
}

func fetchAppList() list.List {
	url := "https://raw.githubusercontent.com/ygorazambuja/WindowsInitial/main/apps.json"

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	var result Result

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("can not unmarshal JSON")
	}

	appList := list.New()
	for _, app := range result {
		appList.PushBack(app.Appname)
	}

	return *appList
}

func runScript(app list.Element) {

	appValue := fmt.Sprintf("%v", app.Value)
	fmt.Println("Instalando:", appValue)

	if c, err := exec.Command("cmd", "/c", "choco", "install", appValue).CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", c)
	}

}

type Result []struct {
	Appname string `json:"appname"`
}

func main() {
	// fmt.Println("Executando scripts")

	installChocolatey()
	runScripts()
}
