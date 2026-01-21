package notification

import (
	"errors"
)

var ErrorWhileOpeningFile = errors.New("An error occurred while opening the file")
var ErrorWhileReadingFile = errors.New("An error occurred while reading the file")
var NotificationFileMalformed = errors.New("The notification file is malformed")
