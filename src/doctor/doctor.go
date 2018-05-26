package doctor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var doctorRepo *gorm.DB

type Doctor struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	HourRate float64 `json:"hourRate" gorm:"column:hour_rate"`
	DtRate   float64 `json:"dtRate" gorm:"column:dt_rate"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	Info     string  `json:"info"`
}

type DoctorModule struct {
	Create  func(w http.ResponseWriter, r *http.Request)
	Update  func(w http.ResponseWriter, r *http.Request)
	Delete  func(w http.ResponseWriter, r *http.Request)
	List    func(w http.ResponseWriter, r *http.Request)
	GetById func(w http.ResponseWriter, r *http.Request)
}

func NewDoctorModule(db *gorm.DB) *DoctorModule {
	doctorRepo = db
	return &DoctorModule{
		Create:  createDoctor,
		Update:  updateDoctor,
		Delete:  deleteDoctor,
		List:    listDoctor,
		GetById: getDoctorById,
	}
}
func getDoctorById(w http.ResponseWriter, r *http.Request) {
	var (
		doctor         Doctor
		querySyringUrl = mux.Vars(r)
	)
	doctorRepo.First(&doctor, querySyringUrl["id"])
	json.NewEncoder(w).Encode(doctor)
}
func createDoctor(w http.ResponseWriter, r *http.Request) {
	var (
		doctor Doctor
	)
	err := json.NewDecoder(r.Body).Decode(&doctor)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	doctorRepo.Create(&doctor)
	json.NewEncoder(w).Encode(doctor)
}
func updateDoctor(w http.ResponseWriter, r *http.Request) {
	var (
		doctor Doctor
	)
	err := json.NewDecoder(r.Body).Decode(&doctor)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	doctorRepo.Model(&doctor).Updates(doctor)
	json.NewEncoder(w).Encode(doctor)
}
func deleteDoctor(w http.ResponseWriter, r *http.Request) {
	var (
		doctor         Doctor
		querySyringUrl = mux.Vars(r)
	)
	doctor.Id, _ = strconv.Atoi(querySyringUrl["id"])
	doctorRepo.Delete(&doctor)
	json.NewEncoder(w).Encode(doctor)
}
func listDoctor(w http.ResponseWriter, r *http.Request) {
	var doctors []*Doctor
	params := r.URL.Query()
	if params.Get("_start") != "" && params.Get("_end") != "" && params.Get("_order") != "" && params.Get("_sort") != "" {
		start, _ := strconv.Atoi(params.Get("_start"))
		end, _ := strconv.Atoi(params.Get("_end"))
		limit := end - start
		order := fmt.Sprintf("%s %s", params.Get("_sort"), params.Get("_order"))
		if params.Get("q") != "" {
			filter := fmt.Sprintf("%s%s%s", "%", params.Get("q"), "%")
			doctorRepo.Where("name LIKE ?", filter).Or("phone LIKE ?", filter).Order(order).Offset(start).Limit(limit).Find(&doctors)
		} else {
			doctorRepo.Order(order).Offset(start).Limit(limit).Find(&doctors)
		}
		total := strconv.Itoa(len(doctors))
		contentRange := fmt.Sprintf("doctors %d-%d/%d", start, end, len(doctors))
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("X-Total-Count", total)
		json.NewEncoder(w).Encode(doctors)
		return
	}
	doctorRepo.Find(&doctors)
	w.Header().Set("X-Total-Count", strconv.Itoa(len(doctors)))
	json.NewEncoder(w).Encode(doctors)
	return
}
