package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/perimeterx/marshmallow"
)

// the only known property in the input json ahead of time
type userID struct {
	ID int `json:"id"`
}

// return a json string with scrubbed sensitive information from the provided input and sensitive fields files
func ScrubPersonalInfo(inputPath, sensitiveFieldsPath string) (string, error) {
	if inputPath == "" {
		return "", errors.New("missing required input file path")
	}
	if sensitiveFieldsPath == "" {
		return "", errors.New("missing required sensitiveFields file path")
	}

	inputFile, _ := ioutil.ReadFile(inputPath)
	sensitiveFields, _ := getSensitiveFields(sensitiveFieldsPath)

	uid := userID{}
	input, err := marshmallow.Unmarshal([]byte(inputFile), &uid)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// keep track of original field values before scrub
	savedValues := make([]any, 0)

	// recursively find and scrub fields from input at any level
	scrubRecursive(&input, "", sensitiveFields, &savedValues, true /* mask */, false /* doScrub */)

	// get json from the scrub string to return
	var b []byte
	b, _ = json.Marshal(input)

	// reset all scrubbed values back to their original values
	scrubRecursive(input, "", sensitiveFields, &savedValues, false /* unmask */, false /* doScrub */)

	// return the scrubbed string
	return string(b), nil
}

// recursively work through the unstructured json to scrub sensitive fields
func scrubRecursive(field interface{}, fieldName string, sensitiveFields map[string]bool, savedValues *[]any, mask bool, doScrub bool) {
	// field must be a pointer and addressable
	addrValue := reflect.ValueOf(field)
	if addrValue.Kind() != reflect.Ptr {
		return
	}
	fieldValue := addrValue.Elem()
	if !fieldValue.IsValid() {
		return
	}
	fieldType := fieldValue.Type()

	if fieldType.Kind() == reflect.Map {
		// got an object, loop through each property
		for _, fKey := range fieldValue.MapKeys() {
			fValue := fieldValue.MapIndex(fKey).Interface()
			scrubRecursive(&fValue, fKey.String(), sensitiveFields, savedValues, mask, doScrub)
			fieldValue.SetMapIndex(fKey, reflect.ValueOf(fValue))
		}
		return
	}

	// skip if no field names
	if fieldName == "" {
		return
	}
	// skip these types
	if !fieldValue.CanSet() || fieldValue.IsZero() {
		return
	}

	if fieldType.Kind() == reflect.Interface {
		_, doFieldScrub := sensitiveFields[strings.ToLower(fieldName)]
		// if parent field is sensitive field, scrub all children (sensitive or not), track in doScrub
		ok := doScrub || doFieldScrub

		// scrub leaf nodes; otherwise, continue recursing
		switch fValue := fieldValue.Interface().(type) {
		case string:
			if !ok {
				return
			}
			if mask {
				*savedValues = append(*savedValues, fValue)
				sampleRegexp := regexp.MustCompile(`[A-Za-z0-9]`)
				result := sampleRegexp.ReplaceAllString(fValue, "*")
				fieldValue.Set(reflect.ValueOf(result))
			} else {
				fieldValue.Set(reflect.ValueOf((*savedValues)[0]))
				*savedValues = (*savedValues)[1:]
			}
		case bool:
			if !ok {
				return
			}
			if mask {
				*savedValues = append(*savedValues, fValue)
				fieldValue.Set(reflect.ValueOf("-"))
			} else {
				fieldValue.Set(reflect.ValueOf((*savedValues)[0]))
				*savedValues = (*savedValues)[1:]
			}
		case int:
			if !ok {
				return
			}
		case float64:
			if !ok {
				return
			}
			if mask {
				format := "%g"
				if fValue == math.Trunc(fValue) {
					format = "%.0f"
				}
				*savedValues = append(*savedValues, fValue)
				str := fmt.Sprintf(format, fValue)
				exp := regexp.MustCompile(`[0-9]`)
				result := exp.ReplaceAllString(str, "*")
				fieldValue.Set(reflect.ValueOf(result))
			} else {
				fieldValue.Set(reflect.ValueOf((*savedValues)[0]))
				*savedValues = (*savedValues)[1:]
			}
		case uint64:
			if !ok {
				return
			}
			if mask {
				*savedValues = append(*savedValues, fieldValue.Uint())
				exp := regexp.MustCompile(`[0-9]`)
				str := fmt.Sprintf("%d", fieldValue.Uint())
				result := exp.ReplaceAllString(str, "*")
				fieldValue.Set(reflect.ValueOf(result))
			} else {
				fieldValue.Set(reflect.ValueOf((*savedValues)[0]))
				*savedValues = (*savedValues)[1:]
			}
		case []interface{}:
			// landed in an array
			for key, value := range fValue {
				scrubRecursive(&value, fieldName, sensitiveFields, savedValues, mask, false /* doScrub */)
				fValue[key] = value
			}
		default:
			// dealing with an object
			value := fieldValue.Elem().Interface()
			m := value.(map[string]interface{})
			_, doScrub := sensitiveFields[strings.ToLower(fieldName)]
			scrubRecursive(&m, fieldName, sensitiveFields, savedValues, mask, doScrub)
		}
	}
}

// return map of sensitive fields to scrub
func getSensitiveFields(sensitiveFieldsPath string) (map[string]bool, error) {
	sensitiveFields := make(map[string]bool)

	readFile, err := os.Open(sensitiveFieldsPath)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		sensitiveFields[fileScanner.Text()] = true
	}
	readFile.Close()

	return sensitiveFields, nil
}
