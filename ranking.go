package vvb

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Endpoint for ranking requests, %s is replaced with the series code.
const epRankings = "http://www.volleyvvb.be/Competitie/rangschikking_xml.php?reeks=%s"

// rankingQuery represents the result of a ranking query.
type rankingQuery struct {
	Results []*Ranking `xml:"rangschikking"`
}

// Ranking represents a single team in a RankingQuery result.
type Ranking struct {
	Series      string  `xml:"reeks"`                      // Series shortcode
	Type        string  `xml:"wedstrijdtype"`              // Either "Hoofd" or "Reserven"
	Position    float32 `xml:"volgorde"`                   // Position in ranking
	Team        string  `xml:"ploegnaam"`                  // Team name
	GamesPlayed uint16  `xml:"aantalGespeeldeWedstrijden"` // Number of games played
	MayorWins   uint8   `xml:"aantalGewonnen30_31"`        // Number of games won with 3-2
	MinorWins   uint8   `xml:"aantalGewonnen32"`           // Number of games won with 3-0 or 3-1
	MayorLosses uint8   `xml:"aantalVerloren30_31"`        // Number of games lost with 0-3 or 0-1
	MinorLosses uint8   `xml:"aantalVerloren32"`           // Number of games lost with 2-3
	SetsWon     uint16  `xml:"aantalGewonnenSets"`         // Number of sets won
	SetsLost    uint16  `xml:"aantalVerlorenSets"`         // Number of sets lost
	Score       uint32  `xml:"puntentotaal"`               // Total score
}

// GetRankings fetches rankings from volleyvvb.be
func GetRankings(series string) (r []*Ranking, err error) {

	// Fast path: Series codes are always at least 3 characters in length.
	if len(series) < 3 {
		return
	}

	// Contact XML api at volleyvvb.be
	body, err := get(fmt.Sprintf(epRankings, series))
	if err != nil || body == nil {
		return
	}
	defer body.Close()

	query := new(rankingQuery)
	dec := xml.NewDecoder(body)

	// Ignore charsets for now
	dec.CharsetReader = func(c string, in io.Reader) (io.Reader, error) {
		return in, nil
	}

	// Attempt to decode ranking response
	if err = dec.Decode(query); err != nil {
		return
	}

	return query.Results, nil

}
