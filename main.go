package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ProductId    int    `json:"productId"`
	Manufacturer string `json:"manufacturer"`
	Sku          string `json:"sku"`
	Upc          string `json:"upc"`
	PricePerUnit string `json:"pricePerunit"`
	Quantity     int    `json:"quantity"`
	ProductName  string `json:"productName"`
}

var ProductList []Product

func init() {
	productJson := `[
		{
			"productId":1,
			"manufacturer": "Sandeep",
			"sku": "POW738DNCIE",
			"upc": "1234567877",
			"pricePerunit": "45.80",
			"quantity": 1,
			"productName": "UPS"
		},
		{
			"productId":2,
			"manufacturer": "Sandeep1",
			"sku": "POW738DNCIE13",
			"upc": "E234567877",
			"pricePerunit": "00.80",
			"quantity": 12,
			"productName": "pspd"
		}		
	]`
	err := json.Unmarshal([]byte(productJson), &ProductList)
	if err != nil {
		log.Fatal(err)
	}
}

func getProductId() int {

	higtestId := -1

	for _, product := range ProductList {
		if product.ProductId > higtestId {
			higtestId = product.ProductId
		}
	}
	return higtestId + 1
}

func findProductById(productId int) (*Product, int) {
	for i, product := range ProductList {
		if product.ProductId == productId {
			return &product, i
		}
	}
	return nil, 0
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.Path, "products/")
	productId, err := strconv.Atoi(urlPath[len(urlPath)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	product, listItemIndex := findProductById(productId)
	if product == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		productJson, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(productJson)

	case http.MethodPut:
		var updateProduct Product
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &updateProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updateProduct.ProductId != productId {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		product = &updateProduct
		ProductList[listItemIndex] = *product
		w.WriteHeader(http.StatusOK)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productJson, err := json.Marshal(ProductList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJson)

	case http.MethodPost:
		var newProduct Product
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return

		}
		err = json.Unmarshal(body, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return

		}

		if newProduct.ProductId != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return

		}

		newProduct.ProductId = getProductId()
		ProductList = append(ProductList, newProduct)

		w.WriteHeader(http.StatusCreated)
		return

	}
}

type fooHandler struct {
	Message string
}

func (h *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.Message))
}

func Barhandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bar Handle"))
}

func main() {

	//Test

	http.Handle("/foo", &fooHandler{Message: "Foo Handle"})
	http.HandleFunc("/bar", Barhandler)

	// S
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productHandler)
	http.ListenAndServe(":5000", nil)
}
