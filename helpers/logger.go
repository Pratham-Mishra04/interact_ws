package helpers

import "github.com/Pratham-Mishra04/interactWS/initializers"

func LogInfo(customString string) {
	initializers.Logger.Info(customString)
}

func LogWarn(customString string, err error) {
	initializers.Logger.Warnw(customString, "Error", err)
}

func LogError(customString string, err error) {
	initializers.Logger.Errorw(customString, "Error", err)
}

func LogFatal(customString string, err error) {
	initializers.Logger.Infow(customString, "Error", err)
}
