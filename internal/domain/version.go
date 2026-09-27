package domain

import (
	"errors"
)

type Version struct {
	currentNumber  int
	previousNumber int
}

func NewVersion(current, previous int) (Version, error) {
	if current < 1 || previous < 1 {
		return Version{}, errors.New("หมายเลข version ต้องมากกว่าหรือเท่ากับ 1 ครับ")
	}
	if current < previous {
		return Version{}, errors.New("current_number ต้องไม่น้อยกว่า previous_number ครับ")
	}
	return Version{
		currentNumber:  current,
		previousNumber: previous,
	}, nil
}

func InitialVersion() Version {
	return Version{
		currentNumber:  1,
		previousNumber: 1,
	}
}

func (v Version) Increment() Version {
	return Version{
		currentNumber:  v.currentNumber + 1,
		previousNumber: v.currentNumber,
	}
}

func (v Version) CurrentNumber() int {
	return v.currentNumber
}

func (v Version) PreviousNumber() int {
	return v.previousNumber
}
