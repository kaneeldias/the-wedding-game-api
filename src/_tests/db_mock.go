package test

import (
	"errors"
	"reflect"
	"strings"
	"the-wedding-game-api/db"
	apperrors "the-wedding-game-api/errors"
	"the-wedding-game-api/types"
)

type MockDB struct {
	items []interface{}
	Error error
}

func (m *MockDB) GetSession() db.DatabaseInterface {
	return m
}

func (m *MockDB) Where(_ interface{}, _ ...interface{}) db.DatabaseInterface {
	return m
}

func (m *MockDB) First(dest interface{}, _ ...interface{}) db.DatabaseInterface {
	if len(m.items) == 0 {
		destValue := reflect.ValueOf(dest)
		m.Error = errors.New("record not found: " + destValue.Type().String())
		return m
	}

	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		m.Error = errors.New("dest must be a non-nil pointer")
		return m
	}

	itemValue := reflect.ValueOf(m.items[0])
	if itemValue.Kind() == reflect.Ptr {
		itemValue = itemValue.Elem()
	}

	destValue.Elem().Set(itemValue)
	m.items = m.items[1:]
	return m
}

func (m *MockDB) Create(value interface{}) db.DatabaseInterface {
	m.items = append(m.items, value)
	return m
}

func (m *MockDB) Find(dest interface{}, _ ...interface{}) db.DatabaseInterface {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		m.Error = errors.New("dest must be a non-nil pointer")
		return m
	}

	if destValue.Elem().Kind() != reflect.Slice {
		m.Error = errors.New("dest must be a pointer to a slice")
		return m
	}

	sliceType := destValue.Elem().Type()
	sliceValue := reflect.MakeSlice(sliceType, len(m.items), len(m.items))
	for i := 0; i < len(m.items); i++ {
		itemValue := reflect.ValueOf(m.items[i])
		if itemValue.Kind() == reflect.Ptr {
			itemValue = itemValue.Elem()
		}
		sliceValue.Index(i).Set(itemValue)
	}

	destValue.Elem().Set(sliceValue)
	return m
}

func (m *MockDB) GetPointsForUser(_ uint) (uint, error) {
	if m.Error != nil {
		return 0, apperrors.NewDatabaseError(m.Error.Error())
	}

	return 100, nil
}

func (m *MockDB) GetLeaderboard() ([]types.LeaderboardEntry, error) {
	if m.Error != nil {
		return nil, apperrors.NewDatabaseError(m.Error.Error())
	}

	return []types.LeaderboardEntry{
		{Username: "user1", Points: 100},
		{Username: "user2", Points: 200},
		{Username: "user3", Points: 300},
	}, nil
}

func (m *MockDB) GetError() error {
	if m.Error == nil {
		return nil
	}

	if strings.Contains(m.Error.Error(), "record not found") {
		return apperrors.NewRecordNotFoundError(m.Error.Error())
	}

	return apperrors.NewDatabaseError(m.Error.Error())
}
