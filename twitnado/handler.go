package twitnado

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ThatCraws/twitnado-extractor/utils"
	"github.com/gin-gonic/gin"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Implementation of the nadoHandlerInterface
type nadoHandler struct {
	scraper    *twitterscraper.Scraper
	mongClient *mongo.Client
}

func NewNadoHandler(connUrl string) *nadoHandler {
	pScraper := twitterscraper.New()
	pScraper = pScraper.WithReplies(true)
	pScraper.SetSearchMode(twitterscraper.SearchLatest)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connUrl))
	if err != nil {
		log.Fatalf("Error connecting to database. Error: %s", err.Error())
	}

	ret := &nadoHandler{
		scraper:    pScraper,
		mongClient: client,
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

	// For better overview in output make an empty line after setup of GIN
	fmt.Println()

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
	fmt.Println("-----------------------------")
	fmt.Printf("Count: %d\n-----------------------------\n", count)

	ctx.JSON(http.StatusOK, allTweets)
}

// Handles /store-endpoint
func (handler *nadoHandler) store(ctx *gin.Context) {
	// read request-body
	var buf []byte

	if ctx.Request.Body == nil {
		log.Print("No body given? Aborting")
		return
	}

	buf, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Error reading request-body. Error: %s", err.Error())
		return
	}

	fmt.Printf("Body: %s\n", string(buf))

	var body []*twitterscraper.TweetResult
	err = json.Unmarshal(buf, &body)
	if err != nil {
		log.Printf("Error unmarshalling Request-Body. Error: %s", err.Error())
		return
	}

	collection := handler.mongClient.Database("scrape").Collection(utils.GetEnvVal("mong_collection", "tweets"))

	// convert to interface slice (to work with InsertMany...)
	toInsert := make([]interface{}, len(body))
	for i := 0; i < len(body); i++ {
		toInsert[i] = body[i]
	}

	_, err = collection.InsertMany(context.TODO(), toInsert)
	if err != nil {
		log.Printf("Unable to insert entries to DB. Error: %s", err.Error())
		ctx.Status(http.StatusInternalServerError)
	}

	ctx.Status(http.StatusOK)
}
