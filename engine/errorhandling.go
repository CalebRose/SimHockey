package engine

import "log"

func HandleMissingPlayer(p GamePlayer, instance string) {
	if p.ID == 0 {
		log.Panicln("ERROR! Could not retrieve player in following instance: ", instance)
	}
}
