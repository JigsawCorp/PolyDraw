package language

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gitlab.com/jigsawcorp/log3900/internal/context"
	"gitlab.com/jigsawcorp/log3900/internal/services/auth"
	"log"
	"net/http"
)

const (
	//FR represent language
	FR = 1
	//EN represent language
	EN = 0
)

var viperInstance *viper.Viper

//Init load all the strings in memory
func Init() {
	log.Printf("[18n] Loading all the strings translation")
	viperInstance = viper.New()
	viperInstance.SetConfigType("yaml")
	err := viperInstance.ReadConfig(bytes.NewBuffer(strings))
	if err != nil {
		panic(err)
	}
}

//MustGet returns the translated string for a key
func MustGet(key string, lang int) string {
	var realKey string
	switch lang {
	case 0:
		realKey = "en." + key
	case 1:
		realKey = "fr." + key
	default:
		panic("[i18n] The language is not valid")
	}
	if viperInstance.IsSet(realKey) {
		return viperInstance.GetString(realKey)
	}
	panic("[i18n] The key " + realKey + " cannot be found")
}

//MustGetRest returns the translated string with a request for a key
func MustGetRest(key string, r *http.Request) string {
	lang := r.Context().Value(context.CtxLang).(int)

	return MustGet(key, lang)
}

//MustGetSocket returns the translated string with a socketID for a key
func MustGetSocket(key string, socketID uuid.UUID) string {
	lang := auth.GetLang(socketID)
	return MustGet(key, lang)
}
