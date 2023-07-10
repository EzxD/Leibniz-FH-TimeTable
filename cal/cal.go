package cal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	CalendarId string
)

func CreateEvent(title string, start string, end string) {
	srv, err := getCalClient()
	if err != nil {
		log.Fatalf("Unable to get CalClient: %v", err)
		return
	}
	event := &calendar.Event{
		Summary: title,
		Start: &calendar.EventDateTime{
			DateTime: start,
			TimeZone: "Europe/Berlin",
		},
		End: &calendar.EventDateTime{
			DateTime: end,
			TimeZone: "Europe/Berlin",
		},
	}

	event, err = srv.Events.Insert(CalendarId, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
	return nil
}

func getCalClient() (*calendar.Service, error) {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return srv, err
}

func DeleteAllUniEvents() {
	srv, err := getCalClient()
	if err != nil {
		log.Fatalf("Unable to get CalClient: %v", err)
		return
	}
	events, err := srv.Events.List(CalendarId).ShowDeleted(false).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve user's events: %v", err)
		return
	}

	for _, event := range events.Items {
		err := srv.Events.Delete(CalendarId, event.Id).Do()
		if err != nil {
			log.Fatalf("Unable to delete event: %v", err)
			return
		}
		log.Printf("Event deleted: %s\n", event.Summary)
	}
	return nil
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
