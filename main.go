package main

import (
    "fmt"
    "net/http"
    //"io/ioutil"
    "encoding/json"
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "time"
    "os"
    "bytes"
    "strings"
)

func check(e error) {
    if e != nil {
        //fmt.Println(e)
        panic(e)
    }
}

type MMSPayload struct {
	From string `json:"from"`
	To string `json:"to"`
	Text string `json:"text"`
	Media string `json:"media"`
}


func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func getLink(token, secret, userId, uri string) (link string, next bool) {
    req, err := http.NewRequest("GET", uri, nil)
    check(err)
    req.SetBasicAuth(token, secret)
    resp, err := http.DefaultClient.Do(req)
    check(err)
    defer resp.Body.Close()
    if resp.Header.Get("link") != "" {
        link = resp.Header.Get("link")
        link = strings.Split(link, "<")[1]
        link = strings.Split(link, ">")[0]
        next = true
    } else {
        next = false
    }

    return

}

func getAllPages(token, secret, userId, initialuri string) (uris []string) {
    link, next := getLink(token, secret, userId, initialuri)
    for next {
        uris = append(uris, link)
        link, next = getLink(token, secret, userId, link)
    }

    if len(uris) > 1 {
        // Pop the last one as its blank
        uris = uris[:len(uris)-1]
    }

    return uris
}

func grabMessages(db *sql.DB, token, secret, uri, reportFile, userId, surveyNumber string, admin []string) {
    req, err := http.NewRequest("GET", uri, nil)
    check(err)
    req.SetBasicAuth(token, secret)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    check(err)
    defer resp.Body.Close()
    
    var record Messages
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		fmt.Println(err)
	}

    for _, message := range record {
        messagetext := strings.Replace(message.Text, "\n","",-1)
        if stringInSlice(message.From, admin) && strings.EqualFold(messagetext, "report") {
            // Check if we already should have sent this, if so do nothing, if not, send a mms
            notindb := checkReportRequests(db, message.MessageID)
            if notindb == true {
                fmt.Println("I got a report request")
                sendReport(db, message.MessageID, reportFile, userId, token, secret, message.From, surveyNumber)
            }
        } else {
            store_message(db, message.MessageID, message.From, messagetext, message.Time)
        }
    }
}

func createReport(db *sql.DB, reportFile string) {
    var messageid string
    var messagefrom string
    var messagetext string
    var messagetime string
    var s []string

    rows := ReadItems(db)
    for rows.Next() {
        err := rows.Scan(&messageid, &messagefrom, &messagetext, &messagetime)
        check(err)
        line := fmt.Sprintf("%v,%v,%v,%v", messageid, messagefrom, messagetext, messagetime)
        s = append(s, line)
    }
    rows.Close()

    file, err := os.Create(reportFile)
    check(err)
    for _, v := range s {
        fmt.Fprintln(file, v)
    }
    file.Close()
}
        

func store_message(db *sql.DB, messageid, from, messagetext, messagetime string) {
    // Store current, and average
    items := []TestItem{
        TestItem{messageid, from, messagetext, messagetime},
    }

    StoreItem(db, items)
}

func checkReportRequests(db *sql.DB, reportmessageid string) (notindb bool) {
    var s []string
    var messageid string

    rows := ReadItemsReport(db)
    for rows.Next() {
        err := rows.Scan(&messageid)
        check(err)
        s = append(s, messageid)
    }
    rows.Close()

    if stringInSlice(reportmessageid, s) != true {
        notindb = true
    }
    return
}


func sendReport(db *sql.DB, reportmessageid, reportFile, userId, token, secret, admin, surveyNumber string) {

    // Store in db so we dont send this again
    items := []ReportItem{
        ReportItem{reportmessageid},
    }

    // Write report
    createReport(db, reportFile)
    // upload report to catapult
    uploadMedia(reportFile, userId, token, secret)
    // send report as a mms to the admin who requested it
    sendMMS(reportFile, userId, token, secret, admin, surveyNumber)
    // Record that we sent a response
    StoreItemReport(db, items)

}

func uploadMedia(reportFile, userId, token, secret string) {
    f, err := os.Open(reportFile)
    check(err)
    defer f.Close()

    uri := fmt.Sprintf("https://api.catapult.inetwork.com/v1/users/%v/media/report.csv", userId)
    req, err := http.NewRequest("PUT", uri, f)
    check(err)
    req.SetBasicAuth(token, secret)
    //req.Header.Set("Content-Type", "text/csv")
    req.Header.Set("Content-Type", "text/comma-separated-values")

    resp, err := http.DefaultClient.Do(req)
    check(err)
    defer resp.Body.Close()
    fmt.Println("I just uploaded a file I think")
    if resp.StatusCode != 200 {
        fmt.Println("Bad upload, response code :", resp.StatusCode)
    }
}

func sendMMS(reportFile, userId, token, secret, admin, surveyNumber string) {
    media := fmt.Sprintf("https://api.catapult.inetwork.com/v1/users/%v/media/report.csv", userId)
    data := MMSPayload {
        From: surveyNumber,
        To: admin,
        Text: "Here is the report you requested",
        Media: media,
    }
    fmt.Println(data)
    payloadBytes, err := json.Marshal(data)
    check(err)
    body := bytes.NewReader(payloadBytes)
    uri := fmt.Sprintf("https://api.catapult.inetwork.com/v1/users/%v/messages", userId)
    req, err := http.NewRequest("POST", uri, body)
    check(err)
    req.SetBasicAuth(token, secret)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    check(err)
    defer resp.Body.Close()
    if resp.StatusCode != 201 {
        fmt.Println("Bad upload, response code :", resp.StatusCode)
    }

}


func main() {
    sqldb, userId, token, secret, reportFile, surveyNumber, admin := config()
    db := InitDB(sqldb)
    CreateTable(db)
    CreateTableReport(db)
    
    // Grab all pages
    initialuri := fmt.Sprintf("https://api.catapult.inetwork.com/v1/users/%v/messages?direction=in&size=100&to=%v", userId, surveyNumber)
    uris := getAllPages(token, secret, userId, initialuri)
    for {
        if len(uris) > 1 {
            for _, uri := range uris {
                fmt.Println(uri)
                grabMessages(db, token, secret, uri, reportFile, userId, surveyNumber, admin)
            }
        } else {
            grabMessages(db, token, secret, initialuri, reportFile, userId, surveyNumber, admin)
        }
        fmt.Println("Sleeping for 30 now")
        time.Sleep(30 * time.Second)
    }
}
