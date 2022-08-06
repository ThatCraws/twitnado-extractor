package twitnado

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

// Implementation of the nadoHandlerInterface
type nadoHandler struct {
	scraper *twitterscraper.Scraper
}

func NewNadoHandler() *nadoHandler {
	pScraper := twitterscraper.New()
	pScraper = pScraper.WithReplies(true)
	pScraper.SetSearchMode(twitterscraper.SearchTop)

	ret := &nadoHandler{
		scraper: pScraper,
	}

	return ret
}

// Handles /search-endpoint
func (handler *nadoHandler) searchQuery(ctx *gin.Context) {
	// Query-Parameter for search-query
	searchQuery := ctx.Query("q")
	if searchQuery == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Query-Parameter for number of tweets to scrape
	numberRetStr := ctx.DefaultQuery("n", "-1")
	fmt.Printf("You are searching for: %s\n", searchQuery)

	numberRet, err := strconv.Atoi(numberRetStr)
	if err != nil {
		log.Printf("Couldn't convert Tweet limit / query-parameter 'p'. Error: %s\nUsing default", err.Error())
		numberRet = -1
	}

	count := 0
	var allTweets []*twitterscraper.TweetResult
	for tweet := range handler.scraper.SearchTweets(context.Background(), searchQuery, numberRet) {
		if tweet.Error != nil {
			log.Print(tweet.Error.Error())
		}
		allTweets = append(allTweets, tweet)
		count++

		fmt.Println("-----------------------------")
		fmt.Printf("--- Author: %s ---\n", tweet.Username)
		fmt.Println(tweet.Text)
	}
	fmt.Println("-----------------------------\n-----------------------------")
	fmt.Printf("Count: %d", count)

	ctx.JSON(http.StatusOK, allTweets)
}

// Handles /store-endpoint
func (handler *nadoHandler) store(ctx *gin.Context) {
	fmt.Println("Triggered store endpoint and me... REEEEEEEEEEEEEEEEEEEE")

	// read request-body
	var buf []byte
	n, err := ctx.Request.Body.Read(buf)

	if err != nil {
		log.Printf("Error reading request-body. Error: %s", err.Error())
	}

	fmt.Printf("Read %d chars. Body: %s\n", n, string(buf))

	var body *queryBody = &queryBody{}
	err = json.Unmarshal(buf, body)

	if err != nil {
		log.Printf("Error unmarshalling Request-Body. Error: %s", err.Error())
	}

	ctx.String(http.StatusOK, body.Query)
}

type queryBody struct {
	Query string
}
