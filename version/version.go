package version

import (
	"fmt"
	"sync"

	"github.com/xackery/eqgamepatch/checksum"
)

var (
	versions sync.Map
)

func init() {
	versions.Store("85218FC053D8B367F2B704BAC5E30ACC", "sof")
	versions.Store("859E89987AA636D36B1007F11C2CD6E0", "uf")
	versions.Store("EF07EE6649C9A2BA2EFFC3F346388E1E78B44B48", "uf") ////one of the torrented uf clients, used by B&R too
	versions.Store("A9DE1B8CC5C451B32084656FCACF1103", "tit")        //p99 client
	versions.Store("BB42BC3870F59B6424A56FED3289C6D4", "tit")        //vanille
	versions.Store("368BB9F425C8A55030A63E606D184445", "rof")
	versions.Store("240C80800112ADA825C146D7349CE85B", "rof")
	versions.Store("A057A23F030BAA1C4910323B131407105ACAD14D", "rof")
	versions.Store("6BFAE252C1A64FE8A3E176CAEE7AAE60", "tbm") //This is one of the live EQ binaries.
	versions.Store("AD970AD6DB97E5BB21141C205CAD6E68", "tbm") //2016/08/27
}

// Detect returns the 3 letter shortname of a client
func Detect() (string, error) {
	path := "eqgame.exe"

	sum, err := checksum.Get(path)
	if err != nil {
		return "", err
	}

	ver, ok := versions.Load(sum)
	if !ok {
		return "unk", fmt.Errorf("unrecognized eqgame (checksum: %s)", sum)
	}
	return ver.(string), nil

}

// Name returns the full name of a provided version
func Name(version string) string {
	versions := map[string]string{
		"unk": "unknown",
		"rof": "Rain of Fear",
		"sof": "Secrets of Feydwer",
		"uf":  "Underfoot",
		"tbm": "The Broken Mirror",
	}
	for short, long := range versions {
		if version == short {
			return long
		}
	}
	return "Unknown"
}
