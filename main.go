package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const privatekey = "9a093885ed19ba7b36011112d4c50f6b1bb991f6"
const publickey = "11cfd16bad4271b95c2537a6d78f9ca8"
const gateway = "https://gateway.marvel.com/v1/public/characters?"

func numbertoString(number int) string {
	return strconv.Itoa(number)
}

//retorna la url de petición
func getRequest(character map[string]string) string {
	has := md5.New()
	timestamp := time.Now()
	var request bytes.Buffer
	io.WriteString(has, strconv.Itoa(timestamp.Second())) //timestamp by second
	io.WriteString(has, privatekey)                       //private key
	io.WriteString(has, publickey)                        //public key
	request.WriteString(gateway)
	_, ok := character["name"]
	if ok {
		request.WriteString("name=" + character["name"] + "&")
	}
	_, ok = character["limit"]
	if ok {
		request.WriteString("limit=" + character["limit"] + "&")
	}
	request.WriteString("ts=" + numbertoString(timestamp.Second()) + "&apikey=" + publickey + "&hash=" + hex.EncodeToString(has.Sum(nil)))
	return request.String()
}

//retorna los datos obtenidos de la petición
func getCharacters(character map[string]string) string {
	req := getRequest(character)
	fmt.Println(req)
	resp, err := http.Get(req)
	if err != nil {
		fmt.Println("No se pudo obtener los personajes!")
		os.Exit(-1)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return string(data)
}

func menu() {
	fmt.Println()
	fmt.Println("Opciones:")
	fmt.Println("1. Buscar por nombre.")
	fmt.Println("2. Listar todos los personajes (limite de 20).")
	fmt.Println("3. Salir.")
	fmt.Print("Opción: ")
}

func main() {
	var value int
	param := make(map[string]string)
	for true {
		menu()
		_, err := fmt.Scanln(&value)
		if err != nil {
			fmt.Println("Error de opción ingresda, intentalo de nuevo!")
			err = nil
		}
		switch value {
		case 1:
			{
				var name string
				fmt.Print("Nombre: ")
				fmt.Scanln(&name)
				param["name"] = name
				fmt.Println(getCharacters(param))
				delete(param, "name")
				fmt.Println("Enter para seguir...")
				fmt.Scanln()
			}

		case 2:
			{
				param["limit"] = "20"
				fmt.Println(getCharacters(param))
				delete(param, "limit")
				fmt.Println("Enter para seguir...")
				fmt.Scanln()
			}
		case 3:
			{
				os.Exit(2)
			}
		}
	}

}
