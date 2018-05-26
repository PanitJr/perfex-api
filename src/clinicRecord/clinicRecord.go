package clinicRecord

import (
	"encoding/json"
	"fmt"
	"net/http"
	"perfex-api/src/clinicProduct"
	"perfex-api/src/customer"
	"perfex-api/src/product"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var clinicRecordRepo *gorm.DB

type ClinicRecord struct {
	Id              int        `json:"id"`
	CustomerId      int        `json:"customerId" gorm:"column:customer_id"`
	ClinicProductId int        `json:"clinicProductId" gorm:"column:clinic_product_id"`
	EmployeeId      int        `json:"employeeId" gorm:"column:employee_id"`
	DoctorId        int        `json:"doctorId" gorm:"column:doctor_id"`
	Price           float64    `json:"price"`
	Discount        float64    `json:"discount"`
	Paid            float64    `json:"paid"`
	Left            float64    `json:"left"`
	CreateAt        *time.Time `json:"createAt" gorm:"column:create_at"`
	LastUpdate      *time.Time `json:"lastUpdate" gorm:"column:last_update"`
}

type ClinicRecordModule struct {
	Create  func(w http.ResponseWriter, r *http.Request)
	Update  func(w http.ResponseWriter, r *http.Request)
	Delete  func(w http.ResponseWriter, r *http.Request)
	List    func(w http.ResponseWriter, r *http.Request)
	GetById func(w http.ResponseWriter, r *http.Request)
}

func NewClinicRecordModule(db *gorm.DB) *ClinicRecordModule {
	clinicRecordRepo = db
	return &ClinicRecordModule{
		Create:  createClinicRecord,
		Update:  updateClinicRecord,
		Delete:  deleteClinicRecord,
		List:    listClinicRecord,
		GetById: getClinicRecordById,
	}
}
func getClinicRecordById(w http.ResponseWriter, r *http.Request) {
	var (
		clinicRecord   ClinicRecord
		querySyringUrl = mux.Vars(r)
	)
	clinicRecordRepo.First(&clinicRecord, querySyringUrl["id"])
	json.NewEncoder(w).Encode(clinicRecord)
}
func createClinicRecord(w http.ResponseWriter, r *http.Request) {
	var (
		clinicRecord ClinicRecord
		cust         customer.Customer
		prod         product.Product
		clinicProd   clinicProduct.ClinicProduct
		relatedProds []*clinicProduct.RelatedPdocuct
	)
	err := json.NewDecoder(r.Body).Decode(&clinicRecord)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	t := time.Now()
	//update product
	clinicRecordRepo.First(&clinicProd, clinicRecord.ClinicProductId)
	err = json.Unmarshal([]byte(clinicProd.RelatedProducts), &relatedProds)
	if err == nil {
		for _, relatedProd := range relatedProds {
			clinicRecordRepo.First(&prod, relatedProd.ProductId)
			leftAmount := prod.Amount - relatedProd.Amount
			if leftAmount < 10 {
				clinicRecordRepo.Model(&clinicProd).Where("id = ?", clinicRecord.ClinicProductId).Update("status", "RUNNING LOW")
			}
			clinicRecordRepo.Model(&prod).Where("id = ?", relatedProd.ProductId).Update("amount", leftAmount)
			clinicRecordRepo.Model(&prod).Where("id = ?", relatedProd.ProductId).Update("last_update", &t)
		}
	}
	clinicRecord.CreateAt = &t
	clinicRecord.LastUpdate = &t
	//update cust
	clinicRecordRepo.Model(&cust).Where("id = ?", clinicRecord.CustomerId).Update("last_update", &t)
	clinicRecordRepo.Create(&clinicRecord)
	json.NewEncoder(w).Encode(clinicRecord)
}
func updateClinicRecord(w http.ResponseWriter, r *http.Request) {
	var (
		clinicRecord ClinicRecord
	)
	err := json.NewDecoder(r.Body).Decode(&clinicRecord)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	t := time.Now()
	clinicRecord.LastUpdate = &t
	clinicRecordRepo.Model(&clinicRecord).Updates(clinicRecord)
	json.NewEncoder(w).Encode(clinicRecord)
}
func deleteClinicRecord(w http.ResponseWriter, r *http.Request) {
	var (
		clinicRecord   ClinicRecord
		querySyringUrl = mux.Vars(r)
	)
	clinicRecord.Id, _ = strconv.Atoi(querySyringUrl["id"])
	clinicRecordRepo.Delete(&clinicRecord)
	json.NewEncoder(w).Encode(clinicRecord)
}
func listClinicRecord(w http.ResponseWriter, r *http.Request) {
	var (
		clinicRecords []*ClinicRecord
		query         map[string]interface{}
	)
	params := r.URL.Query()
	query = make(map[string]interface{})
	if params.Get("_start") != "" && params.Get("_end") != "" && params.Get("_order") != "" && params.Get("_sort") != "" {
		start, _ := strconv.Atoi(params.Get("_start"))
		end, _ := strconv.Atoi(params.Get("_end"))
		limit := end - start
		order := fmt.Sprintf("%s %s", params.Get("_sort"), params.Get("_order"))
		if params.Get("customerId") != "" {
			query["customer_id"] = params.Get("customerId")
		}
		if params.Get("clinicProductId") != "" {
			query["clinic_product_id"] = params.Get("clinicProductId")
		}
		if params.Get("employeeId") != "" {
			query["employee_id"] = params.Get("employeeId")
		}
		if params.Get("doctorId") != "" {
			query["doctor_id"] = params.Get("doctorId")
		}
		if params.Get("customerId") != "" || params.Get("clinicProductId") != "" || params.Get("employeeId") != "" || params.Get("doctorId") != "" {
			clinicRecordRepo.Where(query).Order(order).Offset(start).Limit(limit).Find(&clinicRecords)
		} else {
			clinicRecordRepo.Order(order).Offset(start).Limit(limit).Find(&clinicRecords)
		}
		total := strconv.Itoa(len(clinicRecords))
		contentRange := fmt.Sprintf("clinicRecords %d-%d/%d", start, end, len(clinicRecords))
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("X-Total-Count", total)
		json.NewEncoder(w).Encode(clinicRecords)
		return
	}
	clinicRecordRepo.Find(&clinicRecords)
	w.Header().Set("X-Total-Count", strconv.Itoa(len(clinicRecords)))
	json.NewEncoder(w).Encode(clinicRecords)
	return
}
