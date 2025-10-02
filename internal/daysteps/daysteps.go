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

var ErrCritical = errors.New("critical error")
var ErrConversion = errors.New("conversion error")

func parsePackage(data string) (int, time.Duration, error) {
	slice := strings.Split(data, ",")
	if len(slice) != 2 {
		return 0, time.Duration(0), fmt.Errorf("%w: Не переданы шаги или продолжительность прогулки", ErrCritical)
	}

	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, time.Duration(0), fmt.Errorf("%w: %w", ErrConversion, err)
	}
	if steps <= 0 {
		return 0, time.Duration(0), fmt.Errorf("%w: Количество шагов должно быть больше 0", ErrCritical)
	}
	timeWalk, err := time.ParseDuration(slice[1])
	if err != nil {
		return 0, time.Duration(0), fmt.Errorf("%w: %w", ErrConversion, err)
	}
	if timeWalk <= 0 {
		return 0, time.Duration(0), fmt.Errorf("%w: Продолжительность прогулки должно быть больше 0", ErrCritical)
	}
	return steps, timeWalk, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, timeWalk, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	// пройденное растояние в КМ
	distanceTraveledKm := (float64(steps) * stepLength) / mInKm
	spentCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, timeWalk)
	if err != nil {
		return ""
	}
	// Нужно вычеслить колории через WalkingSpentCalories() из пакета spentcalories
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceTraveledKm, spentCalories)
}
