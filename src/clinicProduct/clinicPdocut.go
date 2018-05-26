package clinicProduct

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var clinicProductRepo *gorm.DB

type ClinicProduct struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`
	RelatedProducts string  `json:"relatedProducts,omitempty" gorm:"column:related_products"`
	Price           float64 `json:"price"`
	Status          string  `json:"status"`
	Info            string  `json:"info"`
}

type RelatedPdocuct struct {
	ProductId int     `json:"product_id"`
	Amount    float64 `json:"amount"`
}

type ClinicProductModule struct {
	Create  func(w http.ResponseWriter, r *http.Request)
	Update  func(w http.ResponseWriter, r *http.Request)
	Delete  func(w http.ResponseWriter, r *http.Request)
	List    func(w http.ResponseWriter, r *http.Request)
	GetById func(w http.ResponseWriter, r *http.Request)
}

func NewClinicProductModule(db *gorm.DB) *ClinicProductModule {
	clinicProductRepo = db
	return &ClinicProductModule{
		Create:  createClinicProduct,
		Update:  updateClinicProduct,
		Delete:  deleteClinicProduct,
		List:    listClinicProduct,
		GetById: getClinicProductById,
	}
}
func getClinicProductById(w http.ResponseWriter, r *http.Request) {
	var (
		clinicProduct  ClinicProduct
		querySyringUrl = mux.Vars(r)
	)
	clinicProductRepo.First(&clinicProduct, querySyringUrl["id"])
	json.NewEncoder(w).Encode(clinicProduct)
}
func createClinicProduct(w http.ResponseWriter, r *http.Request) {
	var (
		clinicProduct ClinicProduct
	)
	err := json.NewDecoder(r.Body).Decode(&clinicProduct)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	clinicProductRepo.Create(&clinicProduct)
	json.NewEncoder(w).Encode(clinicProduct)
}
func updateClinicProduct(w http.ResponseWriter, r *http.Request) {
	var (
		clinicProduct ClinicProduct
	)
	err := json.NewDecoder(r.Body).Decode(&clinicProduct)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	clinicProductRepo.Model(&clinicProduct).Updates(clinicProduct)
	json.NewEncoder(w).Encode(clinicProduct)
}
func deleteClinicProduct(w http.ResponseWriter, r *http.Request) {
	var (
		clinicProduct  ClinicProduct
		querySyringUrl = mux.Vars(r)
	)
	clinicProduct.Id, _ = strconv.Atoi(querySyringUrl["id"])
	clinicProductRepo.Delete(&clinicProduct)
	json.NewEncoder(w).Encode(clinicProduct)
}
func listClinicProduct(w http.ResponseWriter, r *http.Request) {
	var clinicProducts []*ClinicProduct
	params := r.URL.Query()
	if params.Get("_start") != "" && params.Get("_end") != "" && params.Get("_order") != "" && params.Get("_sort") != "" {
		start, _ := strconv.Atoi(params.Get("_start"))
		end, _ := strconv.Atoi(params.Get("_end"))
		limit := end - start
		order := fmt.Sprintf("%s %s", params.Get("_sort"), params.Get("_order"))
		if params.Get("q") != "" {
			filter := fmt.Sprintf("%s%s%s", "%", params.Get("q"), "%")
			clinicProductRepo.Where("name LIKE ?", filter).Or("info LIKE ?", filter).Order(order).Offset(start).Limit(limit).Find(&clinicProducts)
		} else {
			clinicProductRepo.Order(order).Offset(start).Limit(limit).Find(&clinicProducts)
		}
		total := strconv.Itoa(len(clinicProducts))
		contentRange := fmt.Sprintf("clinicProducts %d-%d/%d", start, end, len(clinicProducts))
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("X-Total-Count", total)
		json.NewEncoder(w).Encode(clinicProducts)
		return
	}
	clinicProductRepo.Find(&clinicProducts)
	w.Header().Set("X-Total-Count", strconv.Itoa(len(clinicProducts)))
	json.NewEncoder(w).Encode(clinicProducts)
	return
}
