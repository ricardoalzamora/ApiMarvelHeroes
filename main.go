package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Datos struct {
	Code            int    `json:"code"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
	Etag            string `json:"etag"`
	Data            struct {
		Offset  int `json:"offset"`
		Limit   int `json:"limit"`
		Total   int `json:"total"`
		Count   int `json:"count"`
		Results []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Modified    string `json:"modified"`
			Thumbnail   struct {
				Path      string `json:"path"`
				Extension string `json:"extension"`
			} `json:"thumbnail"`
			ResourceURI string `json:"resourceURI"`
			Comics      struct {
				Available     int    `json:"available"`
				CollectionURI string `json:"collectionURI"`
				Items         []struct {
					ResourceURI string `json:"resourceURI"`
					Name        string `json:"name"`
				} `json:"items"`
				Returned int `json:"returned"`
			} `json:"comics"`
			Series struct {
				Available     int    `json:"available"`
				CollectionURI string `json:"collectionURI"`
				Items         []struct {
					ResourceURI string `json:"resourceURI"`
					Name        string `json:"name"`
				} `json:"items"`
				Returned int `json:"returned"`
			} `json:"series"`
			Stories struct {
				Available     int    `json:"available"`
				CollectionURI string `json:"collectionURI"`
				Items         []struct {
					ResourceURI string `json:"resourceURI"`
					Name        string `json:"name"`
					Type        string `json:"type"`
				} `json:"items"`
				Returned int `json:"returned"`
			} `json:"stories"`
			Events struct {
				Available     int    `json:"available"`
				CollectionURI string `json:"collectionURI"`
				Items         []struct {
					ResourceURI string `json:"resourceURI"`
					Name        string `json:"name"`
				} `json:"items"`
				Returned int `json:"returned"`
			} `json:"events"`
			Urls []struct {
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"urls"`
		} `json:"results"`
	} `json:"data"`
}

const (
	privatekey = "9a093885ed19ba7b36011112d4c50f6b1bb991f6"
	publickey  = "11cfd16bad4271b95c2537a6d78f9ca8"
	gateway    = "https://gateway.marvel.com/v1/public/characters?"
)

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
	_, ok = character["orderBy"]
	if ok {
		request.WriteString("orderBy=" + character["orderBy"] + "&")
	}
	_, ok = character["limit"]
	if ok {
		request.WriteString("limit=" + character["limit"] + "&")
	}
	request.WriteString("ts=" + strconv.Itoa(timestamp.Second()) + "&apikey=" + publickey + "&hash=" + hex.EncodeToString(has.Sum(nil)))
	return request.String()
}

//retorna los datos obtenidos de la petición
func getCharacters(character map[string]string) Datos {
	req := getRequest(character)
	resp, err := http.Get(req)
	if err != nil {
		fmt.Println("No se pudo obtener los personajes!")
		os.Exit(-1)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var datosStruct Datos
	err = json.Unmarshal(data, &datosStruct)
	if err != nil {
		fmt.Println("No se puedo procesar")
	}
	return datosStruct
}

func mostrarResultados(datos Datos) {
	for _, results := range datos.Data.Results {
		fmt.Println("-------------------------------")
		fmt.Println("Nombre: ", results.Name)
		fmt.Println("ID: ", results.ID)
		fmt.Println("Descripción: ", results.Description)
		fmt.Println("Modificado: ", results.Modified)

		fmt.Println()
		fmt.Println("Comics: ")
		if len(results.Comics.Items) == 0 {
			println("Ninguno")
		}
		for _, comics := range results.Comics.Items {
			fmt.Println(comics.Name)
		}

		fmt.Println()
		fmt.Println("Series: ")
		if len(results.Series.Items) == 0 {
			println("Ninguno")
		}
		for _, series := range results.Series.Items {
			fmt.Println(series.Name)
		}

		fmt.Println()
		fmt.Println("Historias: ")
		if len(results.Stories.Items) == 0 {
			println("Ninguno")
		}
		for _, historias := range results.Stories.Items {
			fmt.Println(historias.Name)
		}

		fmt.Println()
		fmt.Println("Eventos: ")
		if len(results.Events.Items) == 0 {
			println("Ninguno")
		}
		for _, eventos := range results.Events.Items {
			fmt.Println(eventos.Name)
		}
		fmt.Println()
	}
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
	for {
		menu()
		_, err := fmt.Scanln(&value)
		if err != nil {
			fmt.Println("Error de opción ingresda, intentalo de nuevo!")
			err = nil
		}
		switch value {
		case 1:
			var name string
			fmt.Print("Nombre: ")
			fmt.Scanln(&name)
			param["name"] = name
			mostrarResultados(getCharacters(param))
			delete(param, "name")
			fmt.Println("Enter para seguir...")
			fmt.Scanln()
		case 2:
			param["limit"] = "20"
			param["orderBy"] = "name"
			mostrarResultados(getCharacters(param))
			delete(param, "limit")
			delete(param, "orderBy")
			fmt.Println("Enter para seguir...")
			fmt.Scanln()
		case 3:
			os.Exit(1)
		}
	}

}
