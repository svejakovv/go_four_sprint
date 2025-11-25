package daysteps

import (
	"fmt"
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
		return 0, 0, fmt.Errorf("ожидается два значения, которые разделены запятой")
	}

	count, err := strconv.Atoi(strings.TrimSpace(stepTime[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("не получилось прочитать количество шагов")
	}

	if count <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}

	duration, err := time.ParseDuration(strings.TrimSpace(stepTime[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("не получилось прочитать время")
	}

	return count, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	count, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if count <= 0 {
		fmt.Println(err)
		return ""
	}

	metricWalk := float64(count) * stepLength
	kilometrWalk := metricWalk / mInKm
	calories, err := spentcalories.WalkingSpentCalories(count, weight, height, duration)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", count, kilometrWalk, calories)
}
