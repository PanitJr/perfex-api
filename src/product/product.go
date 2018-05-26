package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var productRepo *gorm.DB

type fillter struct {
	q string `json:"q"`
}

type Product struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	Brand      string     `json:"brand"`
	Category   string     `json:"category"`
	Amount     float64    `json:"amount"`
	Price      float64    `json:"price"`
	Info       string     `json:"info"`
	Unit       string     `json:"unit"`
	LastUpdate *time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

type ProductModule struct {
	Create  func(w http.ResponseWriter, r *http.Request)
	Update  func(w http.ResponseWriter, r *http.Request)
	Delete  func(w http.ResponseWriter, r *http.Request)
	List    func(w http.ResponseWriter, r *http.Request)
	GetById func(w http.ResponseWriter, r *http.Request)
}

func NewProductModule(db *gorm.DB) *ProductModule {
	productRepo = db
	return &ProductModule{
		Create:  createProduct,
		Update:  updateProduct,
		Delete:  deleteProduct,
		List:    listProduct,
		GetById: getProductById,
	}
}
func getProductById(w http.ResponseWriter, r *http.Request) {
	var (
		product        Product
		querySyringUrl = mux.Vars(r)
	)
	productRepo.First(&product, querySyringUrl["id"])
	json.NewEncoder(w).Encode(product)
	return
}
func createProduct(w http.ResponseWriter, r *http.Request) {
	var (
		product Product
	)
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	t := time.Now()
	product.LastUpdate = &t
	productRepo.Create(&product)
	json.NewEncoder(w).Encode(product)
	return
}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	var (
		product Product
	)
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	t := time.Now()
	product.LastUpdate = &t
	productRepo.Model(&product).Updates(product)
	json.NewEncoder(w).Encode(product)
	return
}
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	var (
		product        Product
		querySyringUrl = mux.Vars(r)
	)
	product.Id, _ = strconv.Atoi(querySyringUrl["id"])
	productRepo.Delete(&product)
	json.NewEncoder(w).Encode(product)
	return
}
func listProduct(w http.ResponseWriter, r *http.Request) {
	var (
		products []*Product
		//fil      fillter
	)
	params := r.URL.Query()
	if params.Get("_start") != "" && params.Get("_end") != "" && params.Get("_order") != "" && params.Get("_sort") != "" {
		start, _ := strconv.Atoi(params.Get("_start"))
		end, _ := strconv.Atoi(params.Get("_end"))
		limit := end - start
		order := fmt.Sprintf("%s %s", params.Get("_sort"), params.Get("_order"))
		// if params.Get("filter") != "" {
		// 	err := json.Unmarshal([]byte(params.Get("filter")), &fil)
		// 	if err == nil {
		// 		nameFilter := fmt.Sprintf("%s%s%s", "%", fil.q, "%")
		// 		productRepo.Where("name LIKE ?", nameFilter).Order(order).Offset(start).Limit(limit).Find(&products)
		// 	}
		// } else
		if params.Get("q") != "" {
			filter := fmt.Sprintf("%s%s%s", "%", params.Get("q"), "%")
			productRepo.Where("name LIKE ?", filter).Or("category LIKE ?", filter).Or("brand LIKE ?", filter).Order(order).Offset(start).Limit(limit).Find(&products)
		} else {
			productRepo.Order(order).Offset(start).Limit(limit).Find(&products)
		}
		total := strconv.Itoa(len(products))
		contentRange := fmt.Sprintf("products %d-%d/%d", start, end, len(products))
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("X-Total-Count", total)
		json.NewEncoder(w).Encode(products)
		return
	}
	productRepo.Find(&products)
	w.Header().Set("X-Total-Count", strconv.Itoa(len(products)))
	json.NewEncoder(w).Encode(products)
	return
}
