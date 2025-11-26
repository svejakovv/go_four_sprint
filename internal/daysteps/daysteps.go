package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	stepTime := strings.Split(data, ",")

	if len(stepTime) != 2 {
		return 0, 0, errors.New("data parsing error")
	}

	//count, err := strconv.Atoi(strings.TrimSpace(stepTime[0]))
	count, err := strconv.Atoi(stepTime[0])

	if err != nil {
		return 0, 0, errors.New("couldn't parse the number of steps")
	}

	if count <= 0 {
		return 0, 0, errors.New("the number of steps must be greater than zero")
	}

	//duration, err := time.ParseDuration(strings.TrimSpace(stepTime[1]))
	duration, err := time.ParseDuration(stepTime[1])

	if err != nil {
		return 0, 0, errors.New("couldn't parse the time")
	}

	if duration <= time.Duration(0) {
		return 0, 0, errors.New("duration activity is too short")
	}

	return count, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	count, duration, err := parsePackage(data)

	if err != nil {
		log.Print(err)
		return ""
	}

	metricWalk := float64(count) * stepLength
	kmWalk := metricWalk / mInKm
	calories, err := spentcalories.WalkingSpentCalories(count, weight, height, duration)

	if err != nil {
		log.Print(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", count, kmWalk, calories)
}
