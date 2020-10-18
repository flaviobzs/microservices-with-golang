package main

import (
	"checkout/queue"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float32 `json:"price,string"`
}

// formato da informação de uma order 
type Order struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	ProductId string `json:"product_id"`
}

var productsUrl string

func init() {
	productsUrl = os.Getenv("PRODUCT_URL")
}

func displayCheckout(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	response, err := http.Get(productsUrl + "/product/" + vars["id"])
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))

	var product Product
	// converter dados de json para esse tipo de struct 
	json.Unmarshal(data, &product)

	t := template.Must(template.ParseFiles("templates/checkout.html"))
	t.Execute(w, product)
}

func finish(w http.ResponseWriter, r *http.Request)  {

	var order Order
	// informaçãoes que vem do formulario 
	order.Name = r.FormValue("name")
	order.Email = r.FormValue("email")
	order.Phone = r.FormValue("phone")
	order.ProductId = r.FormValue("product_id")

	// dado do erro vai ser ignorado se existir 
	// converter em json 
	data, _ := json.Marshal(order)
	fmt.Println(string(data))

	// conectar com que fila 
	connection := queue.Connect()
	// passar dados, e informações da fila 
	queue.Notify(data, "checkout_ex", "", connection)

	w.Write([]byte("Processou!"))
}

func main() {
	r := mux.NewRouter()

	// rotas 
	r.HandleFunc("/finish", finish)
	r.HandleFunc("/{id}", displayCheckout)

	http.ListenAndServe(":8082", r)
}
