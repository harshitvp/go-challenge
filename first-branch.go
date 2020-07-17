package main

cdcdimport (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var responseObject Response

type Response struct {
	Characters []Characters `json:"character"`
	mux        sync.Mutex
}

type Characters struct {
	Name     string `json:"name"`
	MaxPower int    `json:"max_power"`
	GetCount int    `json:"get_count"`
}

func increasePower(responseObject *Response) {
	wg.Add(1)
	for true {
		time.Sleep(time.Duration(10) * time.Second)
		responseObject.mux.Lock()
		for i := 0; i < len(responseObject.Characters); i++ {
			responseObject.Characters[i].MaxPower += 2
		}
		responseObject.mux.Unlock()
	}
	wg.Done()
}

func getPower(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	charName := vars["char_name"]

	flag := false

	for i := 0; i < len(responseObject.Characters); i++ {
		if strings.EqualFold(charName, responseObject.Characters[i].Name) {
			fmt.Fprintf(w, "Character Name: "+responseObject.Characters[i].Name)
			fmt.Fprintf(w, "\nPower : %d", responseObject.Characters[i].MaxPower)
			flag = true
		}
	}

	if flag == false {
		removeLeastPowered()
		getCharacter(charName)

		fmt.Fprintf(w, "Character Name : "+responseObject.Characters[len(responseObject.Characters)-1].Name)
		fmt.Fprintf(w, "\nPower : %d", responseObject.Characters[len(responseObject.Characters)-1].MaxPower)
	}
}

func removeLeastPowered() {
	leastPower := responseObject.Characters[0].MaxPower
	index := 0

	for i := 1; i < len(responseObject.Characters); i++ {
		if leastPower > responseObject.Characters[i].MaxPower {
			leastPower = responseObject.Characters[i].MaxPower
			index = i
		}
	}
	responseObject.Characters = append(responseObject.Characters[:index], responseObject.Characters[index+1:]...)
}

func getCharacter(charName string) {

	response1, err := http.Get("https://run.mocky.io/v3/1e2a8dda-8310-4acc-83d6-ad0f1c30b3f3")
	if err != nil {
		fmt.Print("Error while processing the link(s)")
		return
	}

	responseData1, _ := ioutil.ReadAll(response1.Body)

	var tmpRecords Response
	json.Unmarshal(responseData1, &tmpRecords)

	for i := 0; i < len(tmpRecords.Characters); i++ {
		if strings.EqualFold(charName, tmpRecords.Characters[i].Name) {
			responseObject.Characters = append(responseObject.Characters, tmpRecords.Characters[i])
			fmt.Println("Character Retrievied : " + charName)
			return
		}
	}

	response2, err2 := http.Get("http://www.mocky.io/v2/5ecfd5dc3200006200e3d64b")
	if err2 != nil {
		fmt.Print("Error while processing the link(s)")
		return
	}
	responseData2, _ := ioutil.ReadAll(response2.Body)
	json.Unmarshal(responseData2, &tmpRecords)

	for i := 0; i < len(tmpRecords.Characters); i++ {
		if strings.EqualFold(charName, tmpRecords.Characters[i].Name) {
			responseObject.Characters = append(responseObject.Characters, tmpRecords.Characters[i])
			fmt.Println("Character Retrievied : " + charName)
			return
		}
	}

	response3, err3 := http.Get("http://www.mocky.io/v2/5ecfd6473200009dc1e3d64e")
	if err3 != nil {
		fmt.Print("Error while processing the link(s)")
		return
	}
	responseData3, _ := ioutil.ReadAll(response3.Body)
	json.Unmarshal(responseData3, &tmpRecords)

	for i := 0; i < len(tmpRecords.Characters); i++ {
		if strings.EqualFold(charName, tmpRecords.Characters[i].Name) {
			responseObject.Characters = append(responseObject.Characters, tmpRecords.Characters[i])
			fmt.Println("Character Retrievied : " + charName)
			return
		}
	}
}

func getCharacters() {
	fmt.Println("Retrieving character data from the API...")
	response1, err1 := http.Get("https://run.mocky.io/v3/1e2a8dda-8310-4acc-83d6-ad0f1c30b3f3")
	response2, err2 := http.Get("http://www.mocky.io/v2/5ecfd5dc3200006200e3d64b")
	response3, err3 := http.Get("http://www.mocky.io/v2/5ecfd6473200009dc1e3d64e")

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Print("Error while processing the link(s)")
		os.Exit(1)
	}

	responseData1, err1 := ioutil.ReadAll(response1.Body)
	responseData2, err2 := ioutil.ReadAll(response2.Body)
	responseData3, err3 := ioutil.ReadAll(response3.Body)
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Print("Error encountered. Exiting...")
		os.Exit(1)
	}

	var tmpRecords Response

	json.Unmarshal(responseData1, &responseObject)
	json.Unmarshal(responseData2, &tmpRecords)

	for i := 0; i < len(tmpRecords.Characters) && len(responseObject.Characters) < 16; i++ {
		responseObject.Characters = append(responseObject.Characters, tmpRecords.Characters[i])
	}

	json.Unmarshal(responseData3, &tmpRecords)
	for i := 0; i < len(tmpRecords.Characters) && len(responseObject.Characters) < 16; i++ {
		responseObject.Characters = append(responseObject.Characters, tmpRecords.Characters[i])
	}

	fmt.Println("Character data retrieved from the API...")

}

func main() {
	getCharacters()

	go increasePower(&responseObject)

	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/getPower/{char_name}", getPower).Methods(http.MethodGet)
	http.ListenAndServe(":8080", myRouter)

	wg.Wait()
}
