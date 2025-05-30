package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quaters", "/generals-quaters", "GET", http.StatusOK},
	{"SA", "/search-availability", "GET", http.StatusOK},
	{"RS", "/reservation-summary", "GET", http.StatusOK},
	{"C", "/contact", "GET", http.StatusOK},
	{"MS", "/majors-suite", "GET", http.StatusOK},

	// {"post-searchAvail", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},

	// {"post-searchAvalJson", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},

}

func Test_handlers(t *testing.T) {

	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s  expeceted %d but fot %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{ID: 1,
			RoomName: "General's Quaters"},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCTX(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code got %d , wanter %d", rr.Code, http.StatusOK)
	}
	//test case when reservation is not in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code got %d , wanter %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//Test for GetRoomByID
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code got %d , wanter %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCTX(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestRepository_PostReservation(t *testing.T) {

	reserve := models.Reservation{
		Room: models.Room{
			ID:       1,
			RoomName: "Generatl's Qaters",
		},
	}

	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Anshuman")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Mandal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12312321213")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCTX(req)
	session.Put(ctx, "reservation", reserve)

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Error(rr.Code, http.StatusSeeOther, rr.Code == http.StatusSeeOther)

		t.Errorf("PostResercation handler returned wrong response code got %d -%T, wanter %d-%T,", rr.Code, http.StatusTemporaryRedirect, rr.Code, http.StatusTemporaryRedirect)
	}
	//test for missing post body

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCTX(req)
	session.Put(ctx, "reservation", reserve)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("PostResercation has missing post body handler returned wrong response code got %d -%T, wanter %d-%T,", rr.Code, http.StatusTemporaryRedirect, rr.Code, http.StatusTemporaryRedirect)
	}

	//TEST FOR INVALID START DATE
	reqBody = "start_date=2050-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Anshuman")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Mandal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12312321213")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	session.Put(ctx, "reservation", reserve)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("PostResercation has missing post body handler returned wrong response code got %d -%T, wanter %d-%T,", rr.Code, http.StatusTemporaryRedirect, rr.Code, http.StatusTemporaryRedirect)
	}

	// //TEST FOR INVALID END DATE
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Anshuman")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Mandal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12312321213")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	session.Put(ctx, "reservation", reserve)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("Should be nvalid has Start date is wrong handler returned wrong response code got %d -%T, wanter %d-%T,", rr.Code, http.StatusTemporaryRedirect, rr.Code, http.StatusTemporaryRedirect)
	}
	// //test for invalid RoonID

	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-10")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Anshuman")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Mandal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12312321213")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=ss")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	session.Put(ctx, "reservation", reserve)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("Should be nvalid has End date is wrong returned wrong response code got %d -%T, wanter %d-%T,", rr.Code, http.StatusTemporaryRedirect, rr.Code, http.StatusTemporaryRedirect)
	}
	// //test for invalid Form

	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-10")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=''")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Mandal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12312321213")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	session.Put(ctx, "reservation", reserve)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {

		t.Errorf("PostResercation Cannot processs form returned wrong response code got %d -%T, wanter %d-%T,", rr.Code, http.StatusTemporaryRedirect, rr.Code, http.StatusTemporaryRedirect)
	}
	//Testing Insert Reservaton

	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-10")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Anshuman")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Mandal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12312321213")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	session.Put(ctx, "reservation", reserve)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("PostResercation Reservation should send error returned wrong response code got %d -%T, wanter %d-%T,", rr.Code, http.StatusTemporaryRedirect, rr.Code, http.StatusTemporaryRedirect)
	}

	//Testing and Insert RoomRestriction
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-10")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Anshuman")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Mandal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12312321213")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	session.Put(ctx, "reservation", reserve)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("PostResercation Restriction Should send error wrong response code got %d -%T, wanter %d-%T,", rr.Code, http.StatusTemporaryRedirect, rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostAvailability(t *testing.T) {

	tableTest := []struct {
		startdate      string
		enddate        string
		expectedStatus int
	}{
		{startdate: "2050-01-01", enddate: "2050-01-10", expectedStatus: http.StatusOK},
		{startdate: "2050-01", enddate: "2050-01-10", expectedStatus: http.StatusTemporaryRedirect},
		{startdate: "2050-01-01", enddate: "2050-01", expectedStatus: http.StatusTemporaryRedirect},
		{startdate: "1990-04-30", enddate: "2050-01-10", expectedStatus: http.StatusTemporaryRedirect},
		{startdate: "2050-01-02", enddate: "2050-01-10", expectedStatus: http.StatusOK},
		{startdate: "2050-01-01", enddate: "1990-04-30", expectedStatus: http.StatusSeeOther},
	}

	for _, i := range tableTest {

		reqBody := fmt.Sprintf("start=%v", i.startdate)
		reqBody = fmt.Sprintf("%s&%s", reqBody, fmt.Sprintf("end=%s", i.enddate))

		req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
		ctx := getCTX(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.PostAvailability)
		handler.ServeHTTP(rr, req)

		if rr.Code != i.expectedStatus {

			t.Errorf("PostAvailability Should send error wrong response code got %d, wanter %d,", rr.Code, i.expectedStatus)

		}
	}

	reqBody := strings.NewReader("name=%ZZ")
	req, _ := http.NewRequest("POST", "/", reqBody)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {

		t.Errorf("PostAvailability Should send error wrong response code got %d, wanter %d,", rr.Code, http.StatusTemporaryRedirect)

	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {

	tableTest := []struct {
		startdate      string
		enddate        string
		room_id        string
		expectedStatus int
	}{
		{startdate: "2050-01-01", enddate: "2050-01-10", room_id: "1", expectedStatus: http.StatusOK},
		{startdate: "2050-01", enddate: "2050-01-10", room_id: "1", expectedStatus: http.StatusInternalServerError},
		{startdate: "2050-01-01", enddate: "2050-01", room_id: "1", expectedStatus: http.StatusInternalServerError},
		{startdate: "2050-01-01", enddate: "2050-01-10", room_id: "**", expectedStatus: http.StatusInternalServerError},
		{startdate: "2050-01-01", enddate: "2050-01-10", room_id: "2", expectedStatus: http.StatusInternalServerError},

		// {startdate: "1990-04-30", enddate: "2050-01-10", expectedStatus: http.StatusTemporaryRedirect},
		// {startdate: "2050-01-02", enddate: "2050-01-10", expectedStatus: http.StatusOK},
		// {startdate: "2050-01-01", enddate: "1990-04-30", expectedStatus: http.StatusSeeOther},
	}
	for _, i := range tableTest {

		reqBody := fmt.Sprintf("start=%v", i.startdate)
		reqBody = fmt.Sprintf("%s&%s", reqBody, fmt.Sprintf("end=%s", i.enddate))
		reqBody = fmt.Sprintf("%s&%s", reqBody, fmt.Sprintf("room_id=%s", i.room_id))

		req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
		ctx := getCTX(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AvailabilityJSON)
		handler.ServeHTTP(rr, req)

		if rr.Code != i.expectedStatus {

			t.Errorf("PostAvailability Should send error wrong response code got %d, wanter %d,", rr.Code, i.expectedStatus)

		}
	}

	reqBody := strings.NewReader("name=%ZZ")
	req, _ := http.NewRequest("POST", "/", reqBody)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {

		t.Errorf("PostAvailability Should send error wrong response code got %d, wanter %d,", rr.Code, http.StatusInternalServerError)

	}
}

func TestRepository_ReservationSummary(t *testing.T) {

	tableTest := []struct {
		startdate      string
		enddate        string
		room_id        string
		expectedStatus int
	}{
		{startdate: "2050-01-01", enddate: "2050-01-10", room_id: "1", expectedStatus: http.StatusOK},
	}
	for _, i := range tableTest {
		layout := "2006-01-02"
		st, _ := time.Parse(layout, i.startdate)
		et, _ := time.Parse(layout, i.enddate)

		reservation := models.Reservation{
			RoomID: 1,
			Room: models.Room{ID: 1,
				RoomName: "General's Quaters"},
			StartDate: st,
			EndDAte:   et,
		}

		req, _ := http.NewRequest("GET", "/reservation-summary", nil)
		ctx := getCTX(req)
		session.Put(ctx, "reservation", reservation)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.ReservationSummary)
		handler.ServeHTTP(rr, req)

		if rr.Code != i.expectedStatus {

			t.Errorf("PostAvailability Should send error wrong response code got %d, wanter %d,", rr.Code, i.expectedStatus)

		}
	}

}

func TestRepository_ChooseRoom(t *testing.T) {
	tableTest := []struct {
		name           string
		roomID         string
		expectedStatus int
		keepReservaion bool
	}{
		{"Normal", "1", http.StatusSeeOther, true},
		{"UrlintConversionError", "ss", http.StatusTemporaryRedirect, true},
		{"NoReservation", "1", http.StatusTemporaryRedirect, false},
		{"DatabaseError", "3", http.StatusTemporaryRedirect, true},
	}

	for _, i := range tableTest {

		reservation := models.Reservation{}
		urlString := fmt.Sprintf("/choose-room/%s", i.roomID)
		req, _ := http.NewRequest("GET", urlString, nil)

		ctx := getCTX(req)
		if i.keepReservaion {
			session.Put(ctx, "reservation", reservation)
		}
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.ChooseRoom)
		handler.ServeHTTP(rr, req)

		if rr.Code != i.expectedStatus {
			t.Errorf("FOR %s Expected the status code to be  %v but got %v", i.name, i.expectedStatus, rr.Code)
			t.Error("*****", req.URL.Path)
		}
	}
}

func TestRepository_BookRoom(t *testing.T) {
	tableTest := []struct {
		name           string
		startdate      string
		enddate        string
		roomid         string
		expectedStatus int
	}{
		{"normal", "2025-06-30", "2025-07-05", "1", http.StatusSeeOther},
		{"endDate", "2025-06-30", "2025-07", "1", http.StatusTemporaryRedirect},
		{"startdate", "2025-06", "2025-07-05", "1", http.StatusTemporaryRedirect},
		{"RoomDB", "2025-06-30", "2025-07-05", "3", http.StatusTemporaryRedirect},
		{"RoomId Parse", "2025-06-30", "2025-07-05", "dd", http.StatusTemporaryRedirect},
	}

	for _, i := range tableTest {

		queryParams := url.Values{}
		queryParams.Add("id", i.roomid)
		queryParams.Add("s", i.startdate)
		queryParams.Add("e", i.enddate)

		urlString := fmt.Sprintf("/book-room")
		parsedUrl, _ := url.Parse(urlString)

		parsedUrl.RawQuery = queryParams.Encode()

		req, _ := http.NewRequest("GET", parsedUrl.String(), nil)
		ctx := getCTX(req)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.BookRoom)
		handler.ServeHTTP(rr, req)
		if rr.Code != i.expectedStatus {
			t.Errorf("FOR %s Expected the status code to be  %v but got %v", i.name, i.expectedStatus, rr.Code)
		}

	}
}
