package broadcast

import (
	"log"
	"strings"
)

type Station struct {
	Name string
	Rank int
}

func (stn Station) String() string {
	return stn.Name
}

func (stn Station) ClassName() string {
	return strings.ToLower(
		strings.ReplaceAll(stn.Name, " ", "-"),
	)
}

var (
	SkySports    = Station{"Sky Sports", 0}
	TNTSports    = Station{"TNT Sports", 10}
	AmazonPrime = Station{"Amazon Prime Video", 12}
	BBC          = Station{"BBC", 15}
	BBCOne       = Station{"BBC One", 20}
	BBCTwo       = Station{"BBC Two", 30}
	ITV1         = Station{"ITV1", 40}
	ITV4         = Station{"ITV4", 44}
	ChannelFour  = Station{"Channel 4", 50}
	Talksport    = Station{"talkSPORT", 60}
	Talksport2   = Station{"talkSPORT2", 70}
	Radio5       = Station{"Radio 5 Live", 80}
	Radio5Extra  = Station{"Radio 5 Sports Extra", 90}
	BlankStation = Station{"", 9999}
)

func StationFromString(name string) Station {
	for _, station := range []Station{
		SkySports,
		TNTSports,
		BBCOne,
		BBCTwo,
		ITV1,
		ITV4,
		ChannelFour,
		Talksport,
		Talksport2,
		Radio5,
		Radio5Extra,
		AmazonPrime,
	} {
		if name == station.Name {
			return station
		}
	}
	log.Printf("station not found: '%s'\n", name)
	return Station{name, 9999}
}
