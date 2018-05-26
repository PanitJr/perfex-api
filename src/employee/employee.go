package employee

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var employeeRepo *gorm.DB

type Employee struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	HourRate float64 `json:"hourRate" gorm:"column:hour_rate"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	Info     string  `json:"info"`
}

type EmployeeModule struct {
	Create  func(w http.ResponseWriter, r *http.Request)
	Update  func(w http.ResponseWriter, r *http.Request)
	Delete  func(w http.ResponseWriter, r *http.Request)
	List    func(w http.ResponseWriter, r *http.Request)
	GetById func(w http.ResponseWriter, r *http.Request)
}

func NewEmployeeModule(db *gorm.DB) *EmployeeModule {
	employeeRepo = db
	return &EmployeeModule{
		Create:  creatEemployee,
		Update:  updatEemployee,
		Delete:  deletEemployee,
		List:    listEmployee,
		GetById: getEmployeeById,
	}
}
func getEmployeeById(w http.ResponseWriter, r *http.Request) {
	var (
		employee       Employee
		querySyringUrl = mux.Vars(r)
	)
	employeeRepo.First(&employee, querySyringUrl["id"])
	json.NewEncoder(w).Encode(employee)
}
func creatEemployee(w http.ResponseWriter, r *http.Request) {
	var (
		employee Employee
	)
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	employeeRepo.Create(&employee)
	json.NewEncoder(w).Encode(employee)
}
func updatEemployee(w http.ResponseWriter, r *http.Request) {
	var (
		employee Employee
	)
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	employeeRepo.Model(&employee).Updates(employee)
	json.NewEncoder(w).Encode(employee)
}
func deletEemployee(w http.ResponseWriter, r *http.Request) {
	var (
		employee       Employee
		querySyringUrl = mux.Vars(r)
	)
	employee.Id, _ = strconv.Atoi(querySyringUrl["id"])
	employeeRepo.Delete(&employee)
	json.NewEncoder(w).Encode(employee)
}
func listEmployee(w http.ResponseWriter, r *http.Request) {
	var employees []*Employee
	params := r.URL.Query()
	if params.Get("_start") != "" && params.Get("_end") != "" && params.Get("_order") != "" && params.Get("_sort") != "" {
		start, _ := strconv.Atoi(params.Get("_start"))
		end, _ := strconv.Atoi(params.Get("_end"))
		limit := end - start
		order := fmt.Sprintf("%s %s", params.Get("_sort"), params.Get("_order"))
		if params.Get("q") != "" {
			filter := fmt.Sprintf("%s%s%s", "%", params.Get("q"), "%")
			employeeRepo.Where("name LIKE ?", filter).Or("phone LIKE ?", filter).Order(order).Offset(start).Limit(limit).Find(&employees)
		} else {
			employeeRepo.Order(order).Offset(start).Limit(limit).Find(&employees)
		}
		total := strconv.Itoa(len(employees))
		contentRange := fmt.Sprintf("employees %d-%d/%d", start, end, len(employees))
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("X-Total-Count", total)
		json.NewEncoder(w).Encode(employees)
		return
	}
	employeeRepo.Find(&employees)
	w.Header().Set("X-Total-Count", strconv.Itoa(len(employees)))
	json.NewEncoder(w).Encode(employees)
	return
}
