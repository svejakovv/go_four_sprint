package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	countTrainTime := strings.Split(data, ",")
	if len(countTrainTime) != 3 {
		return 0, "", 0, fmt.Errorf("длина слайса не совпадает")
	}

	count, err := strconv.Atoi(strings.TrimSpace(countTrainTime[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать количество шагов")
	}

	if count <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}

	trainingType := countTrainTime[1]

	duration, err := time.ParseDuration(strings.TrimSpace(countTrainTime[2]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать время")
	}

	return count, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	step := height * stepLengthCoefficient
	distance := (float64(steps) * step) / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	var calories float64

	count, typeTrain, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if height <= 0 || weight <= 0 {
		fmt.Errorf("рост и вес должны быть положительными")
		log.Println(err)
		return "", err
	}

	dist := distance(count, height)
	speed := meanSpeed(count, height, duration)

	switch typeTrain {
	case "Ходьба":
		c, err := WalkingSpentCalories(count, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
		calories = c
	case "Бег":
		c, err := RunningSpentCalories(count, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
		calories = c
	default:
		err := fmt.Errorf("неизвестный тип тренировки")
		log.Println(err)
		return "", err
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nСредняя скорость: %.2f км/ч.\nВы сожгли %.2f ккал.",
		count,
		dist,
		speed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps == 0 || duration == 0 || weight == 0.0 || height == 0.0 {
		return 0.0, fmt.Errorf("невозможно выполнить рассчет калорий при беге")
	}

	return (meanSpeed(steps, height, duration) * weight * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps == 0 || duration == 0 || weight == 0.0 || height == 0.0 {
		return 0.0, fmt.Errorf("невозможно выполнить рассчет калорий при ходьбе")
	}

	calories := (meanSpeed(steps, height, duration) * weight * duration.Minutes()) / minInH
	return calories * walkingCaloriesCoefficient, nil
}
