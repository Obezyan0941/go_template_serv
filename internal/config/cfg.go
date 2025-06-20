package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	db_manager "test_backend/internal/db"
)

func LoadDBConfig() (*db_manager.DBConfig, error) {
	cfg := &db_manager.DBConfig{
		Addr:     getEnv("POSTGRES_HOST"),
		User:     getEnv("POSTGRES_USER_RW"),
		Password: getEnv("POSTGRES_PASSWORD_RW"),
		Database: getEnv("POSTGRES_DB"),
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid database config: %w", err)
	}

	return cfg, nil
}

func validateConfig(cfg interface{}) error {
	v := reflect.ValueOf(cfg)

	if v.Kind() == reflect.Ptr { // if recieved a pointer, retreive value
		v = v.Elem()
	}
	structTypeName := reflect.TypeOf(cfg)

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", v.Kind())
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() != reflect.String {
			log.Printf("Field %s in struct %s is not a string. Skipping check", v.Type().Field(i).Name, structTypeName)
			continue
		}

		value := field.Interface()
		if value == "" {
			return fmt.Errorf("Empty field: %s", value)
		}
	}
	return nil
}

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
