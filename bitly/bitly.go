package bitly

import (
	"github.com/thraxil/bitly"
	"log"
)

const token = "a8098938fc27ec4ff5f18225a74b6e65cdb4803a"

func GetShortnedLink(link string) string {

	c := bitly.NewConnection(token)
	shortLink, err := c.Shorten(link)
	if err != nil {
		log.Fatal(err)
	}
	return shortLink
}
