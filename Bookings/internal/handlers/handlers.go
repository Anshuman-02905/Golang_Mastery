package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/config"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/driver"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/forms"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/models"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/render"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/repository"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/repository/dbrepo"
)

// Template data hold data  sent from handlers to Templates

// REpo the repositody used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct { //Repository patter

	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Creates the NewRpositoyr
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewpostgresDBRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// New Handlers sers the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	//send some data
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})

}
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Cannot get Reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName
	res.Room.ID = res.RoomID

	m.App.Session.Put(r.Context(), "reservatiion", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDAte.Format("2006-01-02")
	fmt.Println(sd)
	fmt.Println(ed)

	stringmap := make(map[string]string)
	stringmap["start_date"] = sd
	stringmap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringmap,
	})
}

// Post Reservation handles posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Cannot get Reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot Parse Form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// //01/02 03:04:05PM '06 -0700
	//layout := "2006-01-02 15:04:05 -0700 MST"
	layout := "2006-01-02"

	start_date, err := time.Parse(layout, r.Form.Get("start_date"))
	if err != nil {
		log.Println(err)
		log.Println(r.Form.Get("start_date"))
		m.App.Session.Put(r.Context(), "error", "Cannot Parese Start Date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	end_date, err := time.Parse(layout, r.Form.Get("end_date"))
	if err != nil {
		log.Println(err)
		log.Println(r.Form.Get("end_date"))
		m.App.Session.Put(r.Context(), "error", "Cannot Parese End Date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	log.Println("Date parsed")

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot process Room_id")
		fmt.Println(" ERROR ROOM ID")

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.InfoLog.Println(start_date, end_date)
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: start_date,
		EndDAte:   end_date,
		RoomID:    roomId,
		Room:      res.Room,
	}
	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		m.App.Session.Put(r.Context(), "error", "Cannot process Form")
		http.Error(w, "My Own Error", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	new_ReservationID, err := m.DB.InsertReservation(reservation)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot insert into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDAte:       reservation.EndDAte,
		RoomID:        reservation.RoomID,
		ReservationID: new_ReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot insert into Room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	Confirmation_msg := fmt.Sprintf(`
		<strong> Reservation Confirmation </strong><br>
		Dear %s:, <br>
		This is to Confirm your reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDAte.Format("2006-01-02"))

	Owner_msg := fmt.Sprintf(`
	<strong> Reservation Confirmation for Owner <strong><br>
	DEARn %s:<br>
	This is the OWNER confirmation for romom confirmation from  %s to %s.
 	`, "OWNER", reservation.StartDate.Format("2006-01-02"), reservation.EndDAte.Format("2006-01-02"),
	)

	m.SendMailData(reservation.Email, "me@here.com", "Reservation Confirmation", Confirmation_msg)
	m.SendMailData("owner@here.com", "me@here.com", "OWNER Confirmation", Owner_msg)

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) SendMailData(To, From, Subject, msg string) {

	msg_sent := models.MailData{
		To:       To,
		From:     From,
		Subject:  Subject,
		Content:  msg,
		Template: "basic.html",
	}
	m.App.MailChan <- msg_sent
	return

}

// renders the Generals pages
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// renders the major page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// renders the availability  page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// render the login Page
// renders the major page

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot Parse Form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	start_date, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot Parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	end_date, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot Parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(start_date, end_date)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot get rooms from search availibility")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if len(rooms) == 0 {
		//No Availability
		m.App.Session.Put(r.Context(), "error", "No availbility")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{

		StartDate: start_date,
		EndDAte:   end_date,
	}
	m.App.Session.Put(r.Context(), "reservation", res)
	w.WriteHeader(http.StatusOK)
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

type jsonResopnse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:room_id`
	StartDate string `json:start_date`
	EndDate   string `json:end_date`
}

// HANDLES Availability JSON for availability and send JSON Message
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	//build a json request
	err := r.ParseForm()
	if err != nil {

		resp := jsonResopnse{
			OK:      false,
			Message: "InternalServerError",
		}
		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return

	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"
	start_date, err := time.Parse(layout, sd)
	if err != nil {

		resp := jsonResopnse{
			OK:      false,
			Message: "Error Parsing StartDate",
		}
		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}
	end_date, err := time.Parse(layout, ed)
	if err != nil {
		resp := jsonResopnse{
			OK:      false,
			Message: "Error Parsing EndDate",
		}
		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}
	log.Println(start_date, end_date)

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))

	if err != nil {
		resp := jsonResopnse{
			OK:      false,
			Message: "Error Parsing RoomId",
		}
		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}
	available, err := m.DB.SearchAvailabilityByDatesByRoomID(start_date, end_date, roomID)
	if err != nil {
		resp := jsonResopnse{
			OK:      false,
			Message: "Error Connecting Database",
		}
		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}
	resp := jsonResopnse{
		OK:        available,
		Message:   "!",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "    ")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contacts.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {

		m.App.ErrorLog.Println("Cannot get reservation from Sessiion")
		m.App.Session.Put(r.Context(), "error", "Cant get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDAte.Format("2006-01-02")

	stringmap := make(map[string]string)
	stringmap["start_date"] = sd
	stringmap["end_date"] = ed

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringmap,
	})

}

// ChooseRoom Displaya list of avaialble rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.URL.Path, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot REservation from Sessiion")
		m.App.Session.Put(r.Context(), "error", "Cannot REservation  from Sessiion")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomID
	res.Room, err = m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.ErrorLog.Println("Cannot Room from DB")
		m.App.Session.Put(r.Context(), "error", "Cannot Room  from DB")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes usrl parameters  builds a session variable
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	//id s e
	roomid, err := strconv.Atoi(r.URL.Query().Get(("id")))
	m.App.ErrorLog.Printf("Cannot parse room_id %d", roomid)

	if err != nil {

		//m.App.ErrorLog.Printf("Cannot parse room_id %s",roomid)
		m.App.Session.Put(r.Context(), "error", "Cannot  parse room_id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}

	sd := r.URL.Query().Get(("s"))
	ed := r.URL.Query().Get(("e"))
	m.App.InfoLog.Println(roomid, sd, ed)

	var res models.Reservation

	layout := "2006-01-02"
	start_date, err := time.Parse(layout, sd)
	if err != nil {

		m.App.ErrorLog.Println("Cannot parse Start Date")
		m.App.Session.Put(r.Context(), "error", "Cannot parse Start Date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}
	end_date, err := time.Parse(layout, ed)
	if err != nil {
		m.App.ErrorLog.Printf("Cannot parse End Date %s\n", ed)
		m.App.Session.Put(r.Context(), "error", "Cannot parse End Date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(roomid)
	if err != nil {
		m.App.ErrorLog.Println("Cannot Room from DB")
		m.App.Session.Put(r.Context(), "error", "Cannot Room  from DB")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName
	res.Room.ID = res.RoomID
	res.RoomID = roomid
	res.StartDate = start_date
	res.EndDAte = end_date

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

//renders the availability  page
// func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
// 	render.Template(w, r,"search-availability.page.tmpl", &models.TemplateData{})
// }

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Handles loging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
		// TODO - take user back to page
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid Login Credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{})
}
func (m *Repository) AdminReservationCalender(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-reservations-calender.page.tmpl", &models.TemplateData{})
}