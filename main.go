package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	// "models/neo4j_models"

	"github.com/gocolly/colly"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"SCPEntry
)

// FindObjectClass Find SCP Report Object class eg Safe, Euclid, Keter
func FindObjectClass(TextBody string) (string, error) {
	r, err := regexp.Compile("Object Class\\: (.+)")
	if err != nil {
		return "", err
	}
	ObjectClass := r.FindStringSubmatch(TextBody)[1]
	return ObjectClass, nil
}

func main() {
	var mutex = &sync.Mutex{}

	pageLister := colly.NewCollector()

	pageGetter := colly.NewCollector(
	// colly.Async(true),
	)

	pageLister.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	pageGetter.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	pageGetter.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 100,
		Delay:       1 * time.Second,
	})

	pageLister.OnHTML("li", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "SCP-") {
			link := e.ChildAttr("a", "href")
			if _, err := strconv.Atoi(link[5:]); err == nil {
				pageGetter.Visit(e.Request.AbsoluteURL(link))
			}

		}
	})
	pageGetter.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	pageGetter.OnHTML("html", func(e *colly.HTMLElement) {
		points := e.ChildText("span[class=rate-points]")
		if points == "" {
			return
			// log.Fatal("Error with finding Points")
		}
		intpoints, _ := strconv.Atoi(points[9:])
		scp := SCPEntry{
			Points: intpoints,
			Url:    e.Request.URL.Path,
		}
		log.Printf("Creating SCP for %s", e.Request.URL.Path)
		mutex.Lock()
		scp.CreateOrUpdate()
		mutex.Unlock()
		e.ForEach("div[class=page-tags] span a", func(_ int, el *colly.HTMLElement) {
			scpTag := SCPTag{Name: el.Attr("href")}
			log.Printf("Creating SCP Tag for %s", el.Attr("href"))
			mutex.Lock()
			scpTag.CreateOrUpdate()
			mutex.Unlock()
			scp.Tags = append(scp.Tags, scpTag)
		})
		log.Printf("Associating SCP Tags %s", e.Request.URL.Path)
		mutex.Lock()
		scp.CreateOrUpdate()
		mutex.Unlock()

		e.ForEach("div[id=main-content]", func(_ int, el *colly.HTMLElement) {
			objectclass, _ := FindObjectClass(el.Text)
			scp.ObjectClass = objectclass
		})
		mutex.Lock()
		scp.CreateOrUpdate()
		mutex.Unlock()
	})

	pageGetter.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	pageLister.Visit("http://www.scp-wiki.net/scp-series")
	pageGetter.Wait()
	pageLister.Visit("http://www.scp-wiki.net/scp-series-2")
	pageGetter.Wait()
	pageLister.Visit("http://www.scp-wiki.net/scp-series-3")
	pageGetter.Wait()
	pageLister.Visit("http://www.scp-wiki.net/scp-series-4")
	pageGetter.Wait()
	pageLister.Visit("http://www.scp-wiki.net/scp-series-5")
	pageGetter.Wait()
	pageLister.Visit("http://www.scp-wiki.net/scp-series-6")
	pageGetter.Wait()
}
