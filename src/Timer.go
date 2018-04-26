package src

import (
	"fmt"
	"os"
	"encoding/json"
	"time"
	"util"
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 {

		switch args[0] {

		case "today" : informationOfDay(time.Now().Local())
		case "yesterday" : informationOfDay(time.Now().Local().AddDate(0, 0, -1))
		case "week" : informationOfWeek()
		case "add" : informationOfWeek()

		}

	} else {
		informationOfDay(time.Now().Local())
	}
}

func informationOfWeek() {
	day := time.Now().Local()
	day = day.AddDate(0,0,-1)
	response := util.DoGet(util.LoadConfiguration().Url + "colaboradores/1-1-3219/acerto/" + day.Format("2006-01-02") + "/marcacoes-originais")
	var timesInStringArray []string
	json.Unmarshal(response, &timesInStringArray)
	responseText := string(response)

	times := toTimeArray(timesInStringArray)

	fmt.Println(times)
	fmt.Println(responseText)
}



func toTimeArray(arrayOfString []string) (timesArray []time.Time) {
	for _, element := range arrayOfString  {
		t, _ := time.Parse("15:04", element)
		timesArray = append(timesArray, t)
	}
	return timesArray
}
func informationOfDay(date time.Time) {
	returnXml := util.GetDateTimesFromXml(util.GetSoapResponseFromSonata(date))

	for _, t := range returnXml.Times {
		fmt.Println(fmt.Sprintf("%02d:%02d", t.Hour, t.Minute))
	}

	fmt.Print("Tempo trabalhado:")

	fmt.Println(fmt.Sprintf("%02d:%02d", returnXml.WorkedTime.Hour, returnXml.WorkedTime.Minute))
	workedTime, _ := time.Parse("15:04", returnXml.WorkedTime.ToString())
	timeToLeave, howMuchTimeToLeave := calculateHowMuchTimeToLeave(workedTime)

	fmt.Println("Você deve sair: " + timeToLeave.Format(time.Kitchen))

	fmt.Print("Faltam : ")

	fmt.Print(howMuchTimeToLeave.Format("15:04"))

	fmt.Println(" para você sair.")
}

func calculateHowMuchTimeToLeave(timeWorked time.Time)(timeToLeave time.Time, howMuchTimeToLeave time.Time) {
	howMuchTimeToLeave = howMuchTimeToleave(timeWorked)
	timeToLeave = timeToleave(howMuchTimeToLeave)
	return
}



func timeToleave(timeWorked time.Time) time.Time {
	timeToLeave := time.Now().Local().Add(time.Hour * time.Duration(timeWorked.Hour()) +
		time.Minute * time.Duration(timeWorked.Minute()))
	return timeToLeave
}

func timeDtoToTime(timeDto util.Time) time.Time {
	timeReturn := time.Now()
	timeReturn.Hour()
	timeReturn.Minute()
	return timeReturn
}

func howMuchTimeToleave(timeWorked time.Time) time.Time {
	normalWorkTime,_ := time.Parse("15:04", "08:30")
	timeWorked = normalWorkTime.Add(time.Hour * time.Duration(-timeWorked.Hour()) +
		time.Minute * time.Duration(-timeWorked.Minute()))
	if timeWorked.Hour() > 8 {
		timeWorked,_ = time.Parse("15:04", "00:00")
	}
	return timeWorked
}


