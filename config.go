package main

import (
    "io/ioutil"
    "github.com/ghodss/yaml"
)

type Config struct {
    Sqldb string `json:"sqldb"`
    UserId string `json:"userId"`
    Token string `json:"token"`
    Secret string `json:"secret"`
    Admin []string `json:"admin"`
    ReportFile string `json:"reportFile"`
    SurveyNumber string `json:"surveyNumber"`
}

func config() (sqldb, userId, token, secret, reportFile, surveyNumber string, admin []string){
    var v Config
    config_file, err := ioutil.ReadFile("/etc/survey_hail/config.yaml")
    check(err)
    yaml.Unmarshal(config_file, &v)
    sqldb = v.Sqldb
    userId = v.UserId
    token = v.Token
    secret = v.Secret
    admin = v.Admin
    reportFile = v.ReportFile
    surveyNumber = v.SurveyNumber
    return
}
