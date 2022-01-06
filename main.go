package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"net/http"
	"time"
)

type Course struct {
	debut   string
	fin     string
	matiere string
	prof    string
	salle   string
}

type Day struct {
	date    string
	courses map[int]*Course
}

func main() {
	router := gin.Default()
	router.GET("/:classe/:groupe", getWeek)
	router.GET("/:classe/:groupe/:date", getWeek)
	router.POST("/:classe/:groupe/:date", refresh)

	router.Run("localhost:8080")
}

func getWeek(c *gin.Context) {
	class := c.Param("classe")
	group := c.Param("groupe")
	stringDate := c.Param("date")
	date := time.Now()
	if stringDate != "" {
		date, err := time.Parse(time.RFC3339, stringDate)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(date)
	}
	date = getStartDayOfWeek(date)
	fmt.Println(date)
	username := getStudentName(class, group)
	url := "https://edtmobiliteng.wigorservices.net/WebPsDyn.aspx?Action=posETUD&serverid=C&tel=" + username + "&date="

	week := setAsyncRequest(date, url)
	c.IndentedJSON(http.StatusOK, week)
}

func refresh(c *gin.Context) {

}

func setAsyncRequest(date time.Time, url string) map[int]*Day {
	c := colly.NewCollector()
	week := make(map[int]*Day)
	j := 0

	q, _ := queue.New(
		5, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		day := &Day{
			date: date.Format("01/02/2006"),
		}
		i := 0
		e.ForEach("div.Ligne", func(i1 int, divContainer *colly.HTMLElement) {
			day.courses[i] = &Course{
				divContainer.ChildText("div.Debut"),
				divContainer.ChildText("div.Fin"),
				divContainer.ChildText("div.Matiere"),
				divContainer.ChildText("div.Prof"),
				divContainer.ChildText("div.Salle"),
			}
			i++
		})

		week[j] = day
		j++
	})

	for i := 0; i < 5; i++ {
		q.AddURL(url + date.Format("01/02/2006") + "%208:00")
		date = date.AddDate(0, 0, 1)
	}

	q.Run(c)
	return week
}

func getStartDayOfWeek(tm time.Time) time.Time { //get monday 00:00:00
	weekday := time.Duration(tm.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := tm.Date()
	currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
}

func getStudentName(class string, group string) string {
	switch class {
	case "I2":
		switch group {
		case "G1":
			return "alexis.heroin"
		case "G2":
			return "benjamin.gonzalez"
		case "INFRA":
			return "alan.amoyel"
		}
	default:
		return "alexis.heroin"
	}
	return "alexis.heroin"
}
