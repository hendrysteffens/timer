package util

import (
	"net/http"
	"bytes"
	"os"
	"path/filepath"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"strconv"
	"crypto/tls"
	"encoding/xml"
)

var FILE = "tsconfig.json"

func DoGet(endpoint string) []byte {

	client := &http.Client{}

	request, error := http.NewRequest("GET", endpoint,nil)

	request.Header.Add("assertion", GenerateToken())

	response, error := client.Do(request)

	if error != nil {

		println(error)

	}

	buffer := new(bytes.Buffer)

	buffer.ReadFrom(response.Body)

	return buffer.Bytes()

}



func LoadConfiguration() Config {

	var config Config
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executablePath := filepath.Dir(ex)
	configFile, err := os.Open(executablePath + "/" + FILE)

	jsonParser := json.NewDecoder(configFile)

	jsonParser.Decode(&config)

	configFile.Close()

	if err != nil {

		fmt.Println(err.Error())

	}

	return config

}



func GenerateToken() string{

	base64Text := base64.StdEncoding.EncodeToString([]byte(LoadConfiguration().Username+":"+LoadConfiguration().Password))

	client := &http.Client{}

	request, error := http.NewRequest("POST", LoadConfiguration().Url + "senior/auth?gestor=S",nil)

	if error != nil {

		println(error)

	}

	request.Header.Add("Authorization", "Basic " + base64Text)

	response, error := client.Do(request)

	buffer := new(bytes.Buffer)

	buffer.ReadFrom(response.Body)

	var token Token

	json.Unmarshal(buffer.Bytes(), &token)

	return token.Token

}

func GetSoapResponseFromSonata(date time.Time) []byte {
	soap := LoadConfiguration().SoapEnv
	soap = strings.Replace(soap, "$day", strconv.Itoa(date.Day()), 2)
	soap = strings.Replace(soap, "$month", strconv.Itoa(int(date.Month())), 2)
	soap = strings.Replace(soap, "$year", strconv.Itoa(date.Year()), 2)
	soap = strings.Replace(soap, "$user", LoadConfiguration().Username, 1)
	soap = strings.Replace (soap, "$password", LoadConfiguration().PortalPassword, 1)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport:tr}
	request, error := http.NewRequest("POST", "https://sonata:1818/ConsultaPontoWS/ConsultaPonto",strings.NewReader(soap))
	response, error := client.Do(request)
	if error != nil {
		println(error)
	}
	buffer := new(bytes.Buffer)
	if response != nil {
		buffer.ReadFrom(response.Body)
	}
	return buffer.Bytes();
}

func GetDateTimesFromXml(sonataXml []byte) Return {

	var r Return

	err := xml.Unmarshal(sonataXml, &r)

	if err != nil {

		println(err)

	}

	return r

}
func (t Time) String() string {

	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)

}

func (time* Time) ToString() string {
	hour := strconv.Itoa(time.Hour)
	minute := strconv.Itoa(time.Minute)
	if time.Hour < 10 {
		hour = "0" + hour
	}
	if time.Minute < 10 {
		minute = "0" + minute
	}
	return hour + ":" + minute
}


type Return struct {
	Times []Time `xml:"Body>consultaPontoResponse>return>clock>time"`
	WorkedTime Time `xml:"Body>consultaPontoResponse>return>workedTime"`
}

type Config struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	PortalPassword string `json:"portalPassword"`
	SoapEnv string `json:"soap"`
}

type Token struct {
	Token      string `json:"token"`
}

type Time struct {
	Hour   int `xml:"hour"`
	Minute int `xml:"minute"`
}