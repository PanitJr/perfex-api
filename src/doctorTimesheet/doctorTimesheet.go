package doctorTimesheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"perfex-api/src/clinicRecord"
	"perfex-api/src/doctor"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var doctorTimesheetRepo *gorm.DB

type DoctorTimesheet struct {
	Id             int        `json:"id"`
	DoctorId       int        `json:"doctorId" gorm:"column:doctor_id"`
	CheckIn        *time.Time `json:"checkIn" gorm:"column:check_in"`
	CheckOut       *time.Time `json:"checkOut" gorm:"column:check_out"`
	CalculatedDt   float64    `json:"calculatedDt" gorm:"column:calculated_dt"`
	CalculatedHour float64    `json:"calculatedHour" gorm:"column:calculated_hour"`
	Total          float64    `json:"total"`
	Aditional      float64    `json:"aditional"`
	Status         string     `json:"status"`
	Info           string     `json:"info"`
}

type DoctorTimesheetModule struct {
	Create  func(w http.ResponseWriter, r *http.Request)
	Update  func(w http.ResponseWriter, r *http.Request)
	Delete  func(w http.ResponseWriter, r *http.Request)
	List    func(w http.ResponseWriter, r *http.Request)
	GetById func(w http.ResponseWriter, r *http.Request)
}

func NewDoctorTimesheetModule(db *gorm.DB) *DoctorTimesheetModule {
	doctorTimesheetRepo = db
	return &DoctorTimesheetModule{
		Create:  createDoctorTimesheet,
		Update:  updateDoctorTimesheet,
		Delete:  deleteDoctorTimesheet,
		List:    listDoctorTimesheet,
		GetById: getDoctorTimesheetById,
	}
}
func getDoctorTimesheetById(w http.ResponseWriter, r *http.Request) {
	var (
		doctorTimesheet DoctorTimesheet
		querySyringUrl  = mux.Vars(r)
	)
	doctorTimesheetRepo.First(&doctorTimesheet, querySyringUrl["id"])
	json.NewEncoder(w).Encode(doctorTimesheet)
}
func createDoctorTimesheet(w http.ResponseWriter, r *http.Request) {
	var (
		doctorTimesheet DoctorTimesheet
	)
	err := json.NewDecoder(r.Body).Decode(&doctorTimesheet)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	t := time.Now()
	doctorTimesheet.CheckIn = &t
	doctorTimesheetRepo.Create(&doctorTimesheet)
	json.NewEncoder(w).Encode(doctorTimesheet)
}
func updateDoctorTimesheet(w http.ResponseWriter, r *http.Request) {
	var (
		doctorTimesheet DoctorTimesheet
		doc             doctor.Doctor
		records         []*clinicRecord.ClinicRecord
		totalDoctorWork float64
	)
	err := json.NewDecoder(r.Body).Decode(&doctorTimesheet)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	doctorTimesheet.Total = 0.00
	if doctorTimesheet.Status == "CHECKOUT" {
		t := time.Now()
		doctorTimesheet.CheckOut = &t
		diff := doctorTimesheet.CheckOut.Sub(*doctorTimesheet.CheckIn)
		doctorTimesheet.CalculatedHour = diff.Hours()
		doc = doctor.Doctor{Id: doctorTimesheet.DoctorId}
		doctorTimesheetRepo.Where("create_at BETWEEN ? AND ?", doctorTimesheet.CheckIn, doctorTimesheet.CheckOut).Where("doctor_id = ?", doctorTimesheet.DoctorId).Find(&records)
		doctorTimesheetRepo.Find(&doc)
		totalDoctorWork = 0.00
		for _, record := range records {
			totalDoctorWork = totalDoctorWork + record.Price
		}
		doctorTimesheet.CalculatedDt = (doc.DtRate / 100) * totalDoctorWork
		doctorTimesheet.Total = (doctorTimesheet.CalculatedHour * doc.HourRate) + doctorTimesheet.CalculatedDt + doctorTimesheet.Aditional
		//fmt.Printf("%+v", doctorTimesheet)
	}
	doctorTimesheetRepo.Model(&doctorTimesheet).Updates(doctorTimesheet)
	json.NewEncoder(w).Encode(doctorTimesheet)
}
func deleteDoctorTimesheet(w http.ResponseWriter, r *http.Request) {
	var (
		doctorTimesheet DoctorTimesheet
		querySyringUrl  = mux.Vars(r)
	)
	doctorTimesheet.Id, _ = strconv.Atoi(querySyringUrl["id"])
	doctorTimesheetRepo.Delete(&doctorTimesheet)
	json.NewEncoder(w).Encode(doctorTimesheet)
}
func listDoctorTimesheet(w http.ResponseWriter, r *http.Request) {
	var doctorTimesheets []*DoctorTimesheet
	params := r.URL.Query()
	if params.Get("_start") != "" && params.Get("_end") != "" && params.Get("_order") != "" && params.Get("_sort") != "" {
		start, _ := strconv.Atoi(params.Get("_start"))
		end, _ := strconv.Atoi(params.Get("_end"))
		limit := end - start
		order := fmt.Sprintf("%s %s", params.Get("_sort"), params.Get("_order"))
		if params.Get("doctorId") != "" {
			doctorTimesheetRepo.Where("doctor_id = ?", params.Get("doctorId")).Order(order).Offset(start).Limit(limit).Find(&doctorTimesheets)
		} else {
		doctorTimesheetRepo.Order(order).Offset(start).Limit(limit).Find(&doctorTimesheets)
		}
		total := strconv.Itoa(len(doctorTimesheets))
		contentRange := fmt.Sprintf("doctorTimesheets %d-%d/%d", start, end, len(doctorTimesheets))
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("X-Total-Count", total)
		json.NewEncoder(w).Encode(doctorTimesheets)
		return
	}
	doctorTimesheetRepo.Find(&doctorTimesheets)
	w.Header().Set("X-Total-Count", strconv.Itoa(len(doctorTimesheets)))
	json.NewEncoder(w).Encode(doctorTimesheets)
	return
}
