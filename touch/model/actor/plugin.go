package actor

import (
	"github.com/peers-touch/peers-touch-go/core/config"
)

func init() {
	config.RegisterOptions(&ymlOptions)
}

// see github.com/peers-touch/peers-touch-go/example/helloworld/conf/actor_*.yml
var ymlOptions struct {
	Peers struct {
		Actor struct {
			Person []struct {
				Name  string `pconf:"name"`
				Email string `pconf:"email"`
			} `pconf:"person"`
		} `pconf:"actor"`
	} `pconf:"peers"`
}
