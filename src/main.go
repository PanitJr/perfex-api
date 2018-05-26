package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"perfex-api/src/clinicProduct"
	"perfex-api/src/clinicRecord"
	"perfex-api/src/customer"
	"perfex-api/src/doctor"
	"perfex-api/src/doctorTimesheet"
	"perfex-api/src/employee"
	"perfex-api/src/employeeTimesheet"
	"perfex-api/src/product"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	addr                    = flag.String("addr", ":8080", "http service address")
	productModule           = &product.ProductModule{}
	customerModule          = &customer.CustomerModule{}
	clinicRecordModule      = &clinicRecord.ClinicRecordModule{}
	doctorModule            = &doctor.DoctorModule{}
	doctorTimesheetModule   = &doctorTimesheet.DoctorTimesheetModule{}
	employeeModule          = &employee.EmployeeModule{}
	employeeTimesheetModule = &employeeTimesheet.EmployeeTimesheetModule{}
	clinicProductModule     = &clinicProduct.ClinicProductModule{}
)

func main() {
	flag.Parse()
	db, err := gorm.Open("mysql", "root:password@/perfex?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	productModule = product.NewProductModule(db)
	clinicProductModule = clinicProduct.NewClinicProductModule(db)
	customerModule = customer.NewCustomerModule(db)
	clinicRecordModule = clinicRecord.NewClinicRecordModule(db)
	doctorModule = doctor.NewDoctorModule(db)
	doctorTimesheetModule = doctorTimesheet.NewDoctorTimesheetModule(db)
	employeeModule = employee.NewEmployeeModule(db)
	employeeTimesheetModule = employeeTimesheet.NewEmployeeTimesheetModule(db)
	defer db.Close()
	fmt.Println(http.ListenAndServe(*addr, GetRouter()))
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes array is a
type Routes []Route

func initRoutes() Routes {
	routes := Routes{
		//Route{"index", "GET", "/", http.FileServer(http.Dir("./public"))},
		// Products
		Route{"create-product", "POST", "/products", productModule.Create},
		Route{"update-product", "PUT", "/products/{id}", productModule.Update},
		Route{"delete-product", "DELETE", "/products/{id}", productModule.Delete},
		Route{"list-product", "GET", "/products", productModule.List},
		Route{"get-product", "GET", "/products/{id}", productModule.GetById},
		// Clinic-products
		Route{"create-clinic-product", "POST", "/clinic-products", clinicProductModule.Create},
		Route{"update-clinic-product", "PUT", "/clinic-products/{id}", clinicProductModule.Update},
		Route{"delete-clinic-product", "DELETE", "/clinic-products/{id}", clinicProductModule.Delete},
		Route{"list-clinic-product", "GET", "/clinic-products", clinicProductModule.List},
		Route{"get-clinic-product", "GET", "/clinic-products/{id}", clinicProductModule.GetById},
		// Doctors
		Route{"create-doctor", "POST", "/doctors", doctorModule.Create},
		Route{"update-doctor", "PUT", "/doctors/{id}", doctorModule.Update},
		Route{"delete-doctor", "DELETE", "/doctors/{id}", doctorModule.Delete},
		Route{"list-doctor", "GET", "/doctors", doctorModule.List},
		Route{"get-doctor", "GET", "/doctors/{id}", doctorModule.GetById},
		// Doctor-timesheets
		Route{"create-doctor-timesheet", "POST", "/doctor-timesheets", doctorTimesheetModule.Create},
		Route{"update-doctor-timesheet", "PUT", "/doctor-timesheets/{id}", doctorTimesheetModule.Update},
		Route{"delete-doctor-timesheet", "DELETE", "/doctor-timesheets/{id}", doctorTimesheetModule.Delete},
		Route{"list-doctor-timesheet", "GET", "/doctor-timesheets", doctorTimesheetModule.List},
		Route{"get-doctor-timesheet", "GET", "/doctor-timesheets/{id}", doctorTimesheetModule.GetById},
		// Employees
		Route{"create-employee", "POST", "/employees", employeeModule.Create},
		Route{"update-employee", "PUT", "/employees/{id}", employeeModule.Update},
		Route{"delete-employee", "DELETE", "/employees/{id}", employeeModule.Delete},
		Route{"list-employee", "GET", "/employees", employeeModule.List},
		Route{"get-employee", "GET", "/employees/{id}", employeeModule.GetById},
		// Employee-timesheets
		Route{"create-employee-timesheet", "POST", "/employee-timesheets", employeeTimesheetModule.Create},
		Route{"update-employee-timesheet", "PUT", "/employee-timesheets/{id}", employeeTimesheetModule.Update},
		Route{"delete-employee-timesheet", "DELETE", "/employee-timesheets/{id}", employeeTimesheetModule.Delete},
		Route{"list-employee-timesheet", "GET", "/employee-timesheets", employeeTimesheetModule.List},
		Route{"get-employee-timesheet", "GET", "/employee-timesheets/{id}", employeeTimesheetModule.GetById},
		// Clinic-records
		Route{"create-clinic-record", "POST", "/clinic-records", clinicRecordModule.Create},
		Route{"update-clinic-record", "PUT", "/clinic-records/{id}", clinicRecordModule.Update},
		Route{"delete-clinic-record", "DELETE", "/clinic-records/{id}", clinicRecordModule.Delete},
		Route{"list-clinic-record", "GET", "/clinic-records", clinicRecordModule.List},
		Route{"get-clinic-record", "GET", "/clinic-records/{id}", clinicRecordModule.GetById},
		// Customer
		Route{"create-customer", "POST", "/customers", customerModule.Create},
		Route{"update-customer", "PUT", "/customers/{id}", customerModule.Update},
		Route{"delete-customer", "DELETE", "/customers/{id}", customerModule.Delete},
		Route{"list-customer", "GET", "/customers", customerModule.List},
		Route{"get-customer", "GET", "/customers/{id}", customerModule.GetById},
		// Stop here if its Preflighted OPTIONS request
		Route{"Preflighted-products", "OPTIONS", "/products", Preflighted},
		Route{"Preflighted-product", "OPTIONS", "/products/{id}", Preflighted},
		Route{"Preflighted-clinic-products", "OPTIONS", "/clinic-products", Preflighted},
		Route{"Preflighted-clinic-product", "OPTIONS", "/clinic-products/{id}", Preflighted},
		Route{"Preflighted-doctors", "OPTIONS", "/doctors", Preflighted},
		Route{"Preflighted-doctor", "OPTIONS", "/doctors/{id}", Preflighted},
		Route{"Preflighted-doctor-timesheets", "OPTIONS", "/doctor-timesheets", Preflighted},
		Route{"Preflighted-doctor-timesheet", "OPTIONS", "/doctor-timesheets/{id}", Preflighted},
		Route{"Preflighted-employees", "OPTIONS", "/employees", Preflighted},
		Route{"Preflighted-employee", "OPTIONS", "/employees/{id}", Preflighted},
		Route{"Preflighted-employee-timesheets", "OPTIONS", "/employee-timesheets", Preflighted},
		Route{"Preflighted-employee-timesheet", "OPTIONS", "/employee-timesheets/{id}", Preflighted},
		Route{"Preflighted-clinic-records", "OPTIONS", "/clinic-records", Preflighted},
		Route{"Preflighted-clinic-record", "OPTIONS", "/clinic-records/{id}", Preflighted},
		Route{"Preflighted-customers", "OPTIONS", "/customers", Preflighted},
		Route{"Preflighted-customer", "OPTIONS", "/customers/{id}", Preflighted},
	}
	return routes
}

func GetRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range initRoutes() {
		var handler http.Handler
		handler = AccessControlAllowOrigin(route.HandlerFunc)
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func Logger(next http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func AccessControlAllowOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Total-Count, Content-Range")
		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count, Content-Range")
		// Lets Gorilla work
		next.ServeHTTP(w, r)
	})
}

func Preflighted(w http.ResponseWriter, r *http.Request) {
	return
}
