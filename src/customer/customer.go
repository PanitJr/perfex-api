package customer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var customerRepo *gorm.DB

type Customer struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Gender      string     `json:"gender"`
	Type        string     `json:"type"`
	National_id string     `json:"nationalId" gorm:"column:national_id"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	DateOfBirth *time.Time `json:"dateOfBirth" gorm:"column:date_of_birth"`
	Info        string     `json:"info"`
	JoinDate    *time.Time `json:"joinDate" gorm:"column:join_date"`
	LastUpdate  *time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

type CustomerModule struct {
	Create  func(w http.ResponseWriter, r *http.Request)
	Update  func(w http.ResponseWriter, r *http.Request)
	Delete  func(w http.ResponseWriter, r *http.Request)
	List    func(w http.ResponseWriter, r *http.Request)
	GetById func(w http.ResponseWriter, r *http.Request)
}

func NewCustomerModule(db *gorm.DB) *CustomerModule {
	customerRepo = db
	return &CustomerModule{
		Create:  createCustomer,
		Update:  updateCustomer,
		Delete:  deleteCustomer,
		List:    listCustomer,
		GetById: getCustomerById,
	}
}
func getCustomerById(w http.ResponseWriter, r *http.Request) {
	var (
		customer       Customer
		querySyringUrl = mux.Vars(r)
	)
	customerRepo.First(&customer, querySyringUrl["id"])
	json.NewEncoder(w).Encode(customer)
}
func createCustomer(w http.ResponseWriter, r *http.Request) {
	var (
		customer Customer
	)
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	t := time.Now()
	customer.JoinDate = &t
	customer.LastUpdate = &t
	customerRepo.Create(&customer)
	json.NewEncoder(w).Encode(customer)
}
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	var (
		customer Customer
	)
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	t := time.Now()
	customer.LastUpdate = &t
	customerRepo.Model(&customer).Updates(customer)
	json.NewEncoder(w).Encode(customer)
}
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	var (
		customer       Customer
		querySyringUrl = mux.Vars(r)
	)
	customer.Id, _ = strconv.Atoi(querySyringUrl["id"])
	customerRepo.Delete(&customer)
	json.NewEncoder(w).Encode(customer)
}
func listCustomer(w http.ResponseWriter, r *http.Request) {
	var customers []*Customer
	params := r.URL.Query()
	if params.Get("_start") != "" && params.Get("_end") != "" && params.Get("_order") != "" && params.Get("_sort") != "" {
		start, _ := strconv.Atoi(params.Get("_start"))
		end, _ := strconv.Atoi(params.Get("_end"))
		limit := end - start
		order := fmt.Sprintf("%s %s", params.Get("_sort"), params.Get("_order"))
		if params.Get("q") != "" {
			filter := fmt.Sprintf("%s%s%s", "%", params.Get("q"), "%")
			customerRepo.Where("name LIKE ?", filter).Or("phone LIKE ?", filter).Or("national_id LIKE ?", filter).Order(order).Offset(start).Limit(limit).Find(&customers)
		} else {
			customerRepo.Order(order).Offset(start).Limit(limit).Find(&customers)
		}
		total := strconv.Itoa(len(customers))
		contentRange := fmt.Sprintf("customers %d-%d/%d", start, end, len(customers))
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("X-Total-Count", total)
		json.NewEncoder(w).Encode(customers)
		return
	}
	customerRepo.Find(&customers)
	w.Header().Set("X-Total-Count", strconv.Itoa(len(customers)))
	json.NewEncoder(w).Encode(customers)
	return
}
