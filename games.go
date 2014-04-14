package vvb

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Endpoint for games requests, %s is replaced with the series code.
const epGames = "http://www.volleyvvb.be/Competitie/wedstrijden_xml.php?reeks=%s"

// gamesQuery represents the result of a games query.
type gamesQuery struct {
	Results []*Game `xml:"wedstrijd"`
}

// Game represents a single game in a GamesQuery result.
type Game struct {
	Nr       uint16 `xml:"nr"`
	Date     string `xml:"datum"`          // Date (yyyy-mm-dd)
	Time     string `xml:"aanvangsuur"`    // Start time (hh:mm)
	Series   string `xml:"reeks"`          // Series shortcode
	Home     string `xml:"thuisploeg"`     // Name of home team
	Visitors string `xml:"bezoekersploeg"` // Name of visiting team
	Venue    string `xml:"sporthal"`       // Venue name
	Score    string `xml:"uitslag"`        // Total score
	Set1     string `xml:"uitslag_set_1"`  // Scores for set 1, if available.
	Set2     string `xml:"uitslag_set_2"`  // Scores for set 2, if available.
	Set3     string `xml:"uitslag_set_3"`  // Scores for set 3, if available.
	Set4     string `xml:"uitslag_set_4"`  // Scores for set 4, if available.
	Set5     string `xml:"uitslag_set_5"`  // Scores for set 5, if available.
}

// GetRankings fetches rankings from volleyvvb.be
func GetGames(series string) (r []*Game, err error) {

	// Fast path: Series codes are always at least 3 characters in length.
	if len(series) < 3 {
		return
	}

	// Contact XML api at volleyvvb.be
	body, err := get(fmt.Sprintf(epGames, series))
	if err != nil || body == nil {
		return
	}
	defer body.Close()

	query := new(gamesQuery)
	dec := xml.NewDecoder(body)

	// Ignore charsets for now
	dec.CharsetReader = func(c string, in io.Reader) (io.Reader, error) {
		return in, nil
	}

	// Attempt to decode games response
	if err = dec.Decode(query); err != nil {
		return
	}

	return query.Results, nil

}
