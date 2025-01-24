package environment

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Remember firt you have to load the enc file after that we can get the value of the env data.
func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}
func GetString(key, defaultValue string) string {

	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return val
}

func GetIntegerValue(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return valAsInt
}
