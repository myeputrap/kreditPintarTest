package helper

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/dongri/phonenumber"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// HTTPSimpleResponse is function for HTTPSimpleResponse
func HTTPSimpleResponse(c *fiber.Ctx, httpStatus int) error {
	return c.Status(httpStatus).SendString(fasthttp.StatusMessage(httpStatus))
}

func NumberUniformity(number string) (string, error) {
	uniformNumber := ""

	if number[0] == '0' {
		uniformNumber = phonenumber.Parse(number, "ID")
	} else {
		number = strings.Trim(number, "+")
		country := phonenumber.GetISO3166ByNumber(number, false)
		if country.Alpha2 != "" {
			uniformNumber = phonenumber.Parse(number, country.Alpha2)
		}
	}

	// log.Debug(uniformNumber)
	if uniformNumber != "" {
		uniformNumber = "+" + uniformNumber
		return uniformNumber, nil
	}

	return "", errors.New("Invalid mobile number")
}

func ValidateAge(birthdate string) (bool, error) {
	layout := "2006-01-02" // Format of the input birthdate
	currentTime := time.Now()

	// Parse the birthdate string
	birthDateTime, err := time.Parse(layout, birthdate)
	if err != nil {
		return false, err
	}

	// Calculate the age by subtracting birthdate from the current date
	age := currentTime.Year() - birthDateTime.Year()

	// Check if the birthdate hasn't occurred yet in the current year
	if currentTime.YearDay() < birthDateTime.YearDay() {
		age--
	}

	// Validate if age is less than 17
	if age < 17 {
		return false, nil
	}

	return true, nil
}

func ValidateDateFormat(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()

	// Check if the date string matches the desired format (YYYY-MM-DD)
	match, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateStr)
	if err != nil {
		// Return false if there was an error matching the regular expression
		return false
	}

	if !match {
		// Return false if the date string doesn't match the desired format
		return false
	}

	// Parse the date string to ensure it's a valid date
	_, err = time.Parse("2006-01-02", dateStr)
	// Return false if there was an error parsing the date string, return true if the date string is valid
	return err == nil
}

func CompareStructs(old interface{}, new interface{}) string {
	differences := make([]string, 0)

	oldVal := reflect.ValueOf(old)
	newVal := reflect.ValueOf(new)

	if oldVal.Kind() == reflect.Struct && newVal.Kind() == reflect.Struct {
		typeOfOld := reflect.TypeOf(old)
		for i := 0; i < typeOfOld.NumField(); i++ {
			field := typeOfOld.Field(i)
			oldField := oldVal.Field(i)

			// Check if the field is one of the desired variables
			if IsDesiredField(field.Name) {
				newField := newVal.FieldByName(field.Name)

				if newField.IsValid() {
					var oldValue, newValue interface{}
					if oldField.Kind() == reflect.Ptr {
						if oldField.IsNil() {
							oldValue = nil
						} else {
							oldValue = oldField.Elem().Interface()
						}
					} else {
						oldValue = oldField.Interface()
					}
					if newField.Kind() == reflect.Ptr {
						if newField.IsNil() {
							newValue = nil
						} else {
							newValue = newField.Elem().Interface()
						}
					} else {
						newValue = newField.Interface()
					}

					// Special handling for time.Time fields
					if oldField.Type() == reflect.TypeOf(time.Time{}) {
						oldTime := oldValue.(time.Time)

						newTime, err := time.Parse("2006-01-02", newValue.(string))
						if err != nil {
							fmt.Printf("Error parsing new time: %v\n", err)
							continue
						}

						oldDate := oldTime.Format("2006-01-02")
						newDate := newTime.Format("2006-01-02")

						if oldDate != newDate {
							diff := fmt.Sprintf("%s: %v into %v", field.Name, oldDate, newDate)
							differences = append(differences, diff)
						}
					} else {
						if !reflect.DeepEqual(oldValue, newValue) {
							diff := fmt.Sprintf("%s: %v into %v", field.Name, oldValue, newValue)
							differences = append(differences, diff)
						}
					}
				}
			}
		}
	}

	return strings.Join(differences, " ")
}

func IsDesiredField(fieldName string) bool {
	desiredFields := []string{
		"Nik",
		"Name",
		"BirthDate",
		"Address",
		"PhoneNumber",
		"Email",
		"HeirName",
		"HeirPhone",
	}

	for _, field := range desiredFields {
		if fieldName == field {
			return true
		}
	}

	return false
}

func ChangeFormatDateWithTimezone(input string) (output string, err error) {
	layout := "2006-01-02"

	birthDate, err := time.Parse(layout, input)
	if err != nil {
		fmt.Println("Error parsing birth date:", err)
		return
	}

	birthDateWithTimezone := birthDate.In(time.FixedZone("WIB", 7*60*60))
	birthDateWithTimezoneString := birthDateWithTimezone.Format("2006-01-02T15:04:05-07:00")
	return birthDateWithTimezoneString, nil
}

func ChangeFormatDateToISO(input string) (output string, err error) {

	// Parse the input date string to a time.Time object
	parsedTime, err := time.Parse(time.RFC3339, input)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	// Format the time in the desired output format
	output = parsedTime.Format("2006-01-02")
	return
}

func GetClientIP(c *fiber.Ctx) (out string, err error) {
	//log.Infof("X-Real-IP : %s", c.Get("X-Real-IP"))

	// log.Infof("X-Forwarded-For : %s", c.Get("X-Forwarded-For"))
	// if realIP := c.Get("X-Real-IP"); realIP != "" {
	// 	return realIP, nil
	// }

	//log.Infof("X-Forwarded-For : %s", c.Get("X-Forwarded-For"))
	if forwardedFor := c.Get("X-Forwarded-For"); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])

			return clientIP, nil
		}
		err = errors.New("malformed X-Forwarded-For header")
		return
	}

	clientIP := c.IP()
	out = clientIP
	return
}
