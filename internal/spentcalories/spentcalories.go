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
	//lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

var ErrCritical = errors.New("critical error")
var ErrConversion = errors.New("conversion error")

func parseTraining(data string) (int, string, time.Duration, error) {
	slice := strings.Split(data, ",")
	if len(slice) != 3 {
		return 0, "", time.Duration(0), fmt.Errorf("%w: Не переданы шаги, вид активности или продолжительность активности", ErrCritical)
	}
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, "", time.Duration(0), fmt.Errorf("%w: %w", ErrConversion, err)
	}
	if steps <= 0 {
		return 0, "", time.Duration(0), fmt.Errorf("%w: Количество шагов должно быть больше 0", ErrCritical)
	}
	activity := slice[1]
	activityDuration, err := time.ParseDuration(slice[2])
	if err != nil {
		return 0, "", time.Duration(0), fmt.Errorf("%w: %w", ErrConversion, err)
	}
	if activityDuration <= 0 {
		return 0, "", time.Duration(0), fmt.Errorf("%w: Продолжительность активности должно быть больше 0", ErrCritical)
	}
	return steps, activity, activityDuration, nil
}

func distance(steps int, height float64) float64 {
	// длинна шага
	stepsLength := height * stepLengthCoefficient
	// пройденное растояние в КМ
	distanceTraveledKm := (stepsLength * float64(steps)) / mInKm
	return distanceTraveledKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}
	distanceTraveledKm := distance(steps, height)
	// Средняя скорость
	averageSpeed := distanceTraveledKm / duration.Hours()
	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, activityDuration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}
	distance := distance(steps, height)
	duration := activityDuration.Hours()
	averageSpeed := meanSpeed(steps, height, activityDuration)
	switch activity {
	case "Бег":
		spentCalories, err := RunningSpentCalories(steps, weight, height, activityDuration)
		if err != nil {
			return "", fmt.Errorf("ошибка в расчётах калорий")
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration, distance, averageSpeed, spentCalories), nil
	case "Ходьба":
		spentCalories, err := WalkingSpentCalories(steps, weight, height, activityDuration)
		if err != nil {
			return "", fmt.Errorf("ошибка в расчётах калорий")
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration, distance, averageSpeed, spentCalories), nil
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("%w: Количество шагов, продолжительность бега, вес и рост пользователя должно быть больше 0", ErrCritical)
	}
	// средняя скорость
	averageSpeed := meanSpeed(steps, height, duration)
	// продолжительность в минутах
	durationInMinutes := duration.Minutes()
	// потрачено калорий
	spentCalories := (weight * averageSpeed * durationInMinutes) / minInH
	return spentCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("%w: Количество шагов, продолжительность бега, вес и рост пользователя должно быть больше 0", ErrCritical)
	}
	// средняя скорость
	averageSpeed := meanSpeed(steps, height, duration)
	// продолжительность в минутах
	durationInMinutes := duration.Minutes()
	// потрачено калорий
	spentCalories := ((weight * averageSpeed * durationInMinutes) / minInH) * walkingCaloriesCoefficient
	return spentCalories, nil
}
