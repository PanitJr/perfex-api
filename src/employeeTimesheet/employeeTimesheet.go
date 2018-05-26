package employeeTimesheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"perfex-api/src/employee"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var employeeTimesheetRepo *gorm.DB

type EmployeeTimesheet struct {
	Id             int        `json:"id"`
	EmployeeId     int        `json:"employeeId" gorm:"column:employee_id"`
	CheckIn        *time.Time `json:"checkIn" gorm:"column:check_in"`
	CheckOut       *time.Time `json:"checkOut" gorm:"column:check_out"`
	CalculatedHour float64    `json:"calculatedHour" gorm:"column:calculated_hour"`
	Total          float64    `json:"total"`
	Aditional      float64    `json:"aditional"`
	Status         string     `json:"status"`
	Info           string     `json:"info"`
}

type EmployeeTimesheetModule struct {
	Create  func(w http.ResponseWriter, r *http.Request)
	Update  func(w http.ResponseWriter, r *http.Request)
	Delete  func(w http.ResponseWriter, r *http.Request)
	List    func(w http.ResponseWriter, r *http.Request)
	GetById func(w http.ResponseWriter, r *http.Request)
}

func NewEmployeeTimesheetModule(db *gorm.DB) *EmployeeTimesheetModule {
	employeeTimesheetRepo = db
	return &EmployeeTimesheetModule{
		Create:  createEmployeeTimesheet,
		Update:  updateEmployeeTimesheet,
		Delete:  deleteEmployeeTimesheet,
		List:    listEmployeeTimesheet,
		GetById: getEmployeeTimesheetById,
	}
}
func getEmployeeTimesheetById(w http.ResponseWriter, r *http.Request) {
	var (
		employeeTimesheet EmployeeTimesheet
		querySyringUrl    = mux.Vars(r)
	)
	employeeTimesheetRepo.First(&employeeTimesheet, querySyringUrl["id"])
	json.NewEncoder(w).Encode(employeeTimesheet)
}
func createEmployeeTimesheet(w http.ResponseWriter, r *http.Request) {
	var (
		employeeTimesheet EmployeeTimesheet
	)
	err := json.NewDecoder(r.Body).Decode(&employeeTimesheet)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	t := time.Now()
	employeeTimesheet.CheckIn = &t
	employeeTimesheetRepo.Create(&employeeTimesheet)
	json.NewEncoder(w).Encode(employeeTimesheet)
}
func updateEmployeeTimesheet(w http.ResponseWriter, r *http.Request) {
	var (
		employeeTimesheet EmployeeTimesheet
		empl              employee.Employee
	)
	err := json.NewDecoder(r.Body).Decode(&employeeTimesheet)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	employeeTimesheet.Total = 0.00
	if employeeTimesheet.Status == "CHECKOUT" {
		t := time.Now()
		employeeTimesheet.CheckOut = &t
		diff := employeeTimesheet.CheckOut.Sub(*employeeTimesheet.CheckIn)
		employeeTimesheet.CalculatedHour = diff.Hours()
		empl = employee.Employee{Id: employeeTimesheet.EmployeeId}
		employeeTimesheetRepo.Find(&empl)
		employeeTimesheet.Total = (employeeTimesheet.CalculatedHour * empl.HourRate) + employeeTimesheet.Aditional
		//fmt.Printf("%+v", employeeTimesheet)
	}
	employeeTimesheetRepo.Model(&employeeTimesheet).Updates(employeeTimesheet)
	json.NewEncoder(w).Encode(employeeTimesheet)
}
func deleteEmployeeTimesheet(w http.ResponseWriter, r *http.Request) {
	var (
		employeeTimesheet EmployeeTimesheet
		querySyringUrl    = mux.Vars(r)
	)
	employeeTimesheet.Id, _ = strconv.Atoi(querySyringUrl["id"])
	employeeTimesheetRepo.Delete(&employeeTimesheet)
	json.NewEncoder(w).Encode(employeeTimesheet)
}
func listEmployeeTimesheet(w http.ResponseWriter, r *http.Request) {
	var employeeTimesheets []*EmployeeTimesheet
	params := r.URL.Query()
	if params.Get("_start") != "" && params.Get("_end") != "" && params.Get("_order") != "" && params.Get("_sort") != "" {
		start, _ := strconv.Atoi(params.Get("_start"))
		end, _ := strconv.Atoi(params.Get("_end"))
		limit := end - start
		order := fmt.Sprintf("%s %s", params.Get("_sort"), params.Get("_order"))
		if params.Get("employeeId") != "" {
			employeeTimesheetRepo.Where("employee_id = ?", params.Get("employeeId")).Order(order).Offset(start).Limit(limit).Find(&employeeTimesheets)
		} else {
			employeeTimesheetRepo.Order(order).Offset(start).Limit(limit).Find(&employeeTimesheets)
		}
		total := strconv.Itoa(len(employeeTimesheets))
		contentRange := fmt.Sprintf("employeeTimesheets %d-%d/%d", start, end, len(employeeTimesheets))
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("X-Total-Count", total)
		json.NewEncoder(w).Encode(employeeTimesheets)
		return
	}
	employeeTimesheetRepo.Find(&employeeTimesheets)
	w.Header().Set("X-Total-Count", strconv.Itoa(len(employeeTimesheets)))
	json.NewEncoder(w).Encode(employeeTimesheets)
	return
}
