package spentcalories

import (
	"errors"
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
		err := errors.New("data parsing error")
		return 0, "", 0, err
	}

	//count, err := strconv.Atoi(strings.TrimSpace(countTrainTime[0]))
	count, err := strconv.Atoi(countTrainTime[0])

	if err != nil {
		err := errors.New("couldn't convert steps")
		return 0, "", 0, err
	}

	if count <= 0 {
		err := errors.New("the number of steps must be greater than zero")
		return 0, "", 0, err
	}

	//duration, err := time.ParseDuration(strings.TrimSpace(countTrainTime[2]))
	duration, err := time.ParseDuration(countTrainTime[2])

	if err != nil {
		return 0, "", 0, err
	}

	if duration <= time.Duration(0) {
		err := errors.New("activity is too short")
		return 0, "", 0, err
	}

	return count, countTrainTime[1], duration, nil
}

func distance(steps int, height float64) float64 {
	step := height * stepLengthCoefficient
	distance := (float64(steps) * step) / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= time.Duration(0) {
		return 0.0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	count, typeTrain, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if height <= 0 || weight <= 0 {
		err := errors.New("height and weight should be positive")
		log.Println(err)
		return "", err
	}

	//dist := distance(count, height)
	//speed := meanSpeed(count, height, duration)

	switch typeTrain {
	case "Ходьба":

		dist := distance(count, height)

		averageSpeed := meanSpeed(count, height, duration)

		calorise, err := WalkingSpentCalories(count, weight, height, duration)

		if err != nil {
			return "", err
		}

		result := fmt.Sprintf("Тип тренировки: Ходьба\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			duration.Hours(), dist, averageSpeed, calorise)

		return result, nil
	case "Бег":

		dist := distance(count, height)

		averageSpeed := meanSpeed(count, height, duration)

		calorise, err := RunningSpentCalories(count, weight, height, duration)

		if err != nil {
			return "", err
		}

		result := fmt.Sprintf("Тип тренировки: Бег\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			duration.Hours(), dist, averageSpeed, calorise)

		return result, nil
	default:

		err := errors.New("неизвестный тип тренировки")
		return "", err
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		err := errors.New("incorrect number of steps when counting calories while running")
		return 0.0, err
	}

	if duration <= time.Duration(0) {
		err := errors.New("incorrect time when counting calories while running")
		return 0.0, err
	}

	if weight <= 0.0 {
		err := errors.New("incorrect weight when counting calories while running")
		return 0.0, err
	}

	if height <= 0.0 {
		err := errors.New("incorrect height when counting calories while running")
		return 0.0, err
	}

	return (meanSpeed(steps, height, duration) * weight * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		err := errors.New("incorrect number of steps when counting calories while walking")
		return 0.0, err
	}

	if duration <= time.Duration(0) {
		err := errors.New("incorrect time when counting calories while walking")
		return 0.0, err
	}

	if weight <= 0.0 {
		err := errors.New("incorrect weight when counting calories while walking")
		return 0.0, err
	}

	if height <= 0.0 {
		err := errors.New("incorrect height when counting calories while walking")
		return 0.0, err
	}

	calories := (meanSpeed(steps, height, duration) * weight * duration.Minutes()) / minInH
	return calories * walkingCaloriesCoefficient, nil
}
