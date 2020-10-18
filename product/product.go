package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
)

// definir o tipo de um unico produto 
type Product struct {
	// serializar campos `json:xxx`
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float64 `json:"price,string"`
}

// definir o tipo de varios produtos 
type Products struct {
	Products []Product
}

func loadData() []byte {
	// carregar o arquivo json 
	// receber resultados da função os.Open em duas variaveis!!! success no jsonFile e erros no err 	
	jsonFile, err := os.Open("products.json")

	// se existir algum erro (tratamento de erros)
	if err != nil {
		// fmt.Println("ola") modulo nativo do go 
		fmt.Println("erro: ",err.Error())
	}

	// fechar o arquivo 
	defer jsonFile.Close()

	// ler e transformar o dado do arquivo em conteudo em memoria 
	data, err := ioutil.ReadAll(jsonFile)
	return data
}

// request e respose nos argumentos 
func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := loadData()
	
	// resposta retornada
	w.Write([]byte(products))
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	// pegar argumento do paramento para pegar apenas um produto 
	vars := mux.Vars(r)
	data := loadData()

	// coleção de produtos do tipo Products
	var products Products
	// converter json em estrutura da linguagem 
	json.Unmarshal(data, &products) //usar &products como referencia

	// percorrer array para verificar se existe um item com id passado
	for _, v := range products.Products {
		if v.Uuid == vars["id"] {
			product, _ := json.Marshal(v)
			w.Write([]byte(product))
		}
	}
}

// função inicial do arquivo 
func main() {
	// criar roteamento com o gotilla;mux
	r := mux.NewRouter()
	
	// rotas do app
	r.HandleFunc("/products", ListProducts)
	r.HandleFunc("/product/{id}", GetProductById)

	// servidor vai ouvir na mesma porta 
	http.ListenAndServe(":8081", r)
	
}

// comando para gerenciar os modulos do go
// $ go mod init product 
