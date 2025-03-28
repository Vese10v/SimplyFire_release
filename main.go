package main

import (
	"database/sql"
	"image/color"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lukasjarosch/go-docx"
	"github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
)

const (
	ISOformat = "2006-01-02"
)

// *функция для проверки имени и фамилии на числовые значения
func IsEntryTextInt(entry_text string) bool {
	for _, elem_int := range entry_text {
		if unicode.IsNumber(elem_int) {
			logrus.Errorf("01.0: 1-й этап проверки первичных данных (имя / фамилия) на числовые значения: UNSUCCESS (числовые значения обнаружены) → %v", entry_text)
			return true
		}
	}
	logrus.Infof("01.0: 1-й этап проверки первичных данных (имя / фамилия) на числовые значения: SUCCESS (числовые значения не обнаружены) → %v", entry_text)
	return false
}

// **функция для проверки имени и фамилии на знаки, отличные от кириллицы
func IsEntryTextCyrillic(entry_text string) bool {
	for _, elem_cyrillic := range entry_text {
		if elem_cyrillic < unicode.MaxASCII {
			logrus.Errorf("01.1: 1-й этап проверки первичных данных (имя / фамилия) на знаки, отличные от кириллицы: UNSUCCESS (обнаружены отличные знаки, отличные от кириллицы) → %v", entry_text)
			return true
		}
	}
	logrus.Infof("01.1: 1-й этап проверки первичных данных (имя / фамилия) на знаки, отличные от кириллицы: SUCCESS (знаки, отличные от кириллицы, не обнаружены) → %v", entry_text)
	return false
}

// ***функция для проверки даты на соответствие с форматом даты ISOformat
func EntryTextConvertToDate(birth_date_entry_text string) bool {
	birth_date_entry_text_convert_to_date, err := time.Parse(ISOformat, birth_date_entry_text)
	if err == nil {
		logrus.Infof("01.2: 1-й этап проверки первичных данных (дата рождения) на соответствие формату даты ISOformat: SUCCESS → %v", birth_date_entry_text_convert_to_date)
		return true
	} else {
		logrus.Errorf("01.2: 1-й этап проверки первичных данных (дата рождения) на соответствие формату даты ISOformat: UNSUCCESS → %v", birth_date_entry_text)
		return false
	}
}

// ****функция для преобразования должностей в родительный падеж
func ConvertEntryDataPositionToGenitiveCase(entry_data string) string {
	entry_data_fields_array := strings.Fields(entry_data)
	entry_data_result_array := []string{}
	// условный перебор элементов массива согласно длине массива
	if len(entry_data_fields_array) == 1 {
		for i := 0; i <= len(entry_data_fields_array[0]); i++ {
			to_rune_0 := []rune(entry_data_fields_array[0])
			to_rune_0_string := string(to_rune_0)
			to_rune_0_string_1_letter_end := string(to_rune_0[len(to_rune_0)-1:])
			if to_rune_0_string_1_letter_end == "ь" {
				delete_1_letter := string(to_rune_0[:len(to_rune_0)-1])
				result := delete_1_letter + "я"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				if to_rune_0_string_1_letter_end == "б" || to_rune_0_string_1_letter_end == "в" || to_rune_0_string_1_letter_end == "г" || to_rune_0_string_1_letter_end == "д" || to_rune_0_string_1_letter_end == "ж" || to_rune_0_string_1_letter_end == "з" || to_rune_0_string_1_letter_end == "к" || to_rune_0_string_1_letter_end == "л" || to_rune_0_string_1_letter_end == "м" || to_rune_0_string_1_letter_end == "н" || to_rune_0_string_1_letter_end == "п" || to_rune_0_string_1_letter_end == "р" || to_rune_0_string_1_letter_end == "с" || to_rune_0_string_1_letter_end == "т" || to_rune_0_string_1_letter_end == "ф" || to_rune_0_string_1_letter_end == "х" || to_rune_0_string_1_letter_end == "ц" || to_rune_0_string_1_letter_end == "ч" || to_rune_0_string_1_letter_end == "ш" || to_rune_0_string_1_letter_end == "щ" {
					result := to_rune_0_string + "а"
					entry_data_result_array = append(entry_data_result_array, result)
					break
				}
			}
		}
	} else if len(entry_data_fields_array) == 2 {
		for i := 0; i <= len(entry_data_fields_array[0]); i++ {
			to_rune_0 := []rune(entry_data_fields_array[0])
			to_rune_0_string := string(to_rune_0)
			to_rune_0_string_1_letter_end := string(to_rune_0[len(to_rune_0)-1:])
			to_rune_0_string_2_letter_end := string(to_rune_0[len(to_rune_0)-2:])
			to_rune_0_string_2_letter_end_cut := string(to_rune_0[:len(to_rune_0)-2])
			to_rune_0_2_letter_end_cut := []rune(to_rune_0_string_2_letter_end_cut)
			to_rune_0_2_letter_end_cut_1_letter_end := string(to_rune_0_2_letter_end_cut[len(to_rune_0_2_letter_end_cut)-1:])
			if to_rune_0_string_1_letter_end == "ь" {
				delete_1_letter := string(to_rune_0[:len(to_rune_0)-1])
				result := delete_1_letter + "я"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_0_string_1_letter_end == "б" || to_rune_0_string_1_letter_end == "в" || to_rune_0_string_1_letter_end == "г" || to_rune_0_string_1_letter_end == "д" || to_rune_0_string_1_letter_end == "ж" || to_rune_0_string_1_letter_end == "з" || to_rune_0_string_1_letter_end == "к" || to_rune_0_string_1_letter_end == "л" || to_rune_0_string_1_letter_end == "м" || to_rune_0_string_1_letter_end == "н" || to_rune_0_string_1_letter_end == "п" || to_rune_0_string_1_letter_end == "р" || to_rune_0_string_1_letter_end == "с" || to_rune_0_string_1_letter_end == "т" || to_rune_0_string_1_letter_end == "ф" || to_rune_0_string_1_letter_end == "х" || to_rune_0_string_1_letter_end == "ц" || to_rune_0_string_1_letter_end == "ч" || to_rune_0_string_1_letter_end == "ш" || to_rune_0_string_1_letter_end == "щ" {
				result := to_rune_0_string + "а"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				if to_rune_0_string_2_letter_end == "ый" {
					delete_2_letters := string(to_rune_0[:len(to_rune_0)-2])
					result := delete_2_letters + "ого"
					entry_data_result_array = append(entry_data_result_array, result)
					break
				} else if to_rune_0_string_2_letter_end == "ий" && (to_rune_0_2_letter_end_cut_1_letter_end == "ж" || to_rune_0_2_letter_end_cut_1_letter_end == "ш" || to_rune_0_2_letter_end_cut_1_letter_end == "ч" || to_rune_0_2_letter_end_cut_1_letter_end == "щ") {
					delete_2_letters := string(to_rune_0[:len(to_rune_0)-2])
					result := delete_2_letters + "его"
					entry_data_result_array = append(entry_data_result_array, result)
					break
				} else {
					delete_2_letters := string(to_rune_0[:len(to_rune_0)-2])
					result := delete_2_letters + "ого"
					entry_data_result_array = append(entry_data_result_array, result)
					break
				}
			}
		}
		for i := 0; i <= len(entry_data_fields_array[1]); i++ {
			to_rune_1 := []rune(entry_data_fields_array[1])
			to_rune_1_string := string(to_rune_1)
			to_rune_1_string_1_letter_end := string(to_rune_1[len(to_rune_1)-1:])
			if to_rune_1_string_1_letter_end == "в" {
				result := to_rune_1
				entry_data_result_array = append(entry_data_result_array, string(result))
				break
			} else if to_rune_1_string_1_letter_end == "а" || to_rune_1_string_1_letter_end == "о" || to_rune_1_string_1_letter_end == "у" || to_rune_1_string_1_letter_end == "ы" || to_rune_1_string_1_letter_end == "э" || to_rune_1_string_1_letter_end == "я" || to_rune_1_string_1_letter_end == "е" || to_rune_1_string_1_letter_end == "ё" || to_rune_1_string_1_letter_end == "ю" || to_rune_1_string_1_letter_end == "и" {
				result := to_rune_1_string
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_1_string_1_letter_end == "в" {
				result := to_rune_1_string
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				result := to_rune_1_string + "а"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			}
		}
	} else if len(entry_data_fields_array) == 3 {
		for i := 0; i <= len(entry_data_fields_array[0]); i++ {
			to_rune_0 := []rune(entry_data_fields_array[0])
			to_rune_0_string := string(to_rune_0)
			to_rune_0_string_1_letter_end := string(to_rune_0[len(to_rune_0)-1:])
			to_rune_0_string_1_letter_end_cut := string(to_rune_0[:len(to_rune_0)-1])
			to_rune_0_string_2_letter_end := string(to_rune_0[len(to_rune_0)-2:])
			to_rune_0_string_2_letter_end_cut := string(to_rune_0[:len(to_rune_0)-2])
			if to_rune_0_string_1_letter_end == "ь" {
				delete_1_letter := to_rune_0_string_1_letter_end_cut
				result := delete_1_letter + "я"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_0_string_1_letter_end == "б" || to_rune_0_string_1_letter_end == "в" || to_rune_0_string_1_letter_end == "г" || to_rune_0_string_1_letter_end == "д" || to_rune_0_string_1_letter_end == "ж" || to_rune_0_string_1_letter_end == "з" || to_rune_0_string_1_letter_end == "к" || to_rune_0_string_1_letter_end == "л" || to_rune_0_string_1_letter_end == "м" || to_rune_0_string_1_letter_end == "н" || to_rune_0_string_1_letter_end == "п" || to_rune_0_string_1_letter_end == "р" || to_rune_0_string_1_letter_end == "с" || to_rune_0_string_1_letter_end == "т" || to_rune_0_string_1_letter_end == "ф" || to_rune_0_string_1_letter_end == "х" || to_rune_0_string_1_letter_end == "ц" || to_rune_0_string_1_letter_end == "ч" || to_rune_0_string_1_letter_end == "ш" || to_rune_0_string_1_letter_end == "щ" {
				result := to_rune_0_string + "а"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				if to_rune_0_string_2_letter_end == "ый" {
					delete_2_letters := to_rune_0_string_2_letter_end_cut
					result := delete_2_letters + "ого"
					entry_data_result_array = append(entry_data_result_array, result)
					break
				} else if to_rune_0_string_2_letter_end == "ий" {
					delete_2_letters := to_rune_0_string_2_letter_end_cut
					result := delete_2_letters + "его"
					entry_data_result_array = append(entry_data_result_array, result)
					break
				}
			}
		}
		for i := 0; i <= len(entry_data_fields_array[1]); i++ {
			to_rune_1 := []rune(entry_data_fields_array[1])
			to_rune_1_string := string(to_rune_1)
			to_rune_1_string_1_letter_end := string(to_rune_1[len(to_rune_1)-1:])
			to_rune_1_string_1_letter_end_cut := string(to_rune_1[:len(to_rune_1)-1])
			if to_rune_1_string_1_letter_end == "б" || to_rune_1_string_1_letter_end == "в" || to_rune_1_string_1_letter_end == "г" || to_rune_1_string_1_letter_end == "д" || to_rune_1_string_1_letter_end == "ж" || to_rune_1_string_1_letter_end == "з" || to_rune_1_string_1_letter_end == "к" || to_rune_1_string_1_letter_end == "л" || to_rune_1_string_1_letter_end == "м" || to_rune_1_string_1_letter_end == "н" || to_rune_1_string_1_letter_end == "п" || to_rune_1_string_1_letter_end == "р" || to_rune_1_string_1_letter_end == "с" || to_rune_1_string_1_letter_end == "т" || to_rune_1_string_1_letter_end == "ф" || to_rune_1_string_1_letter_end == "х" || to_rune_1_string_1_letter_end == "ц" || to_rune_1_string_1_letter_end == "ч" || to_rune_1_string_1_letter_end == "ш" || to_rune_1_string_1_letter_end == "щ" {
				result := to_rune_1_string + "а"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_1_string_1_letter_end == "ь" {
				delete_1_letter := to_rune_1_string_1_letter_end_cut
				result := delete_1_letter + "я"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				result := to_rune_1_string
				entry_data_result_array = append(entry_data_result_array, result)
				break
			}
		}
		for i := 0; i <= len(entry_data_fields_array[2]); i++ {
			to_rune_2 := []rune(entry_data_fields_array[2])
			to_rune_2_string := string(to_rune_2)
			result := to_rune_2_string
			entry_data_result_array = append(entry_data_result_array, result)
			break
		}
	} else if len(entry_data_fields_array) >= 4 {
		for i := 0; i <= len(entry_data_fields_array[0]); i++ {
			to_rune_0 := []rune(entry_data_fields_array[0])
			to_rune_0_string := string(to_rune_0)
			to_rune_0_string_1_letter_end := string(to_rune_0[len(to_rune_0)-1:])
			to_rune_0_string_1_letter_end_cut := string(to_rune_0[:len(to_rune_0)-1])
			to_rune_0_string_2_letter_end := string(to_rune_0[len(to_rune_0)-2:])
			to_rune_0_string_2_letter_end_cut := string(to_rune_0[:len(to_rune_0)-2])
			if to_rune_0_string_1_letter_end == "ь" {
				delete_1_letter := to_rune_0_string_1_letter_end_cut
				result := delete_1_letter + "я"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_0_string_1_letter_end == "б" || to_rune_0_string_1_letter_end == "в" || to_rune_0_string_1_letter_end == "г" || to_rune_0_string_1_letter_end == "д" || to_rune_0_string_1_letter_end == "ж" || to_rune_0_string_1_letter_end == "з" || to_rune_0_string_1_letter_end == "к" || to_rune_0_string_1_letter_end == "л" || to_rune_0_string_1_letter_end == "м" || to_rune_0_string_1_letter_end == "н" || to_rune_0_string_1_letter_end == "п" || to_rune_0_string_1_letter_end == "р" || to_rune_0_string_1_letter_end == "с" || to_rune_0_string_1_letter_end == "т" || to_rune_0_string_1_letter_end == "ф" || to_rune_0_string_1_letter_end == "х" || to_rune_0_string_1_letter_end == "ц" || to_rune_0_string_1_letter_end == "ч" || to_rune_0_string_1_letter_end == "ш" || to_rune_0_string_1_letter_end == "щ" {
				result := to_rune_0_string + "а"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				if to_rune_0_string_2_letter_end == "ый" {
					delete_2_letters := to_rune_0_string_2_letter_end_cut
					result := delete_2_letters + "ого"
					entry_data_result_array = append(entry_data_result_array, result)
					break
				} else if to_rune_0_string_2_letter_end == "ий" {
					delete_2_letters := to_rune_0_string_2_letter_end_cut
					result := delete_2_letters + "его"
					entry_data_result_array = append(entry_data_result_array, result)
					break
				}
			}
		}
		for i := 0; i <= len(entry_data_fields_array[1]); i++ {
			to_rune_1 := []rune(entry_data_fields_array[1])
			to_rune_1_string := string(to_rune_1)
			to_rune_1_string_1_letter_end := string(to_rune_1[len(to_rune_1)-1:])
			to_rune_1_string_1_letter_end_cut := string(to_rune_1[:len(to_rune_1)-1])
			if to_rune_1_string_1_letter_end == "б" || to_rune_1_string_1_letter_end == "в" || to_rune_1_string_1_letter_end == "г" || to_rune_1_string_1_letter_end == "д" || to_rune_1_string_1_letter_end == "ж" || to_rune_1_string_1_letter_end == "з" || to_rune_1_string_1_letter_end == "к" || to_rune_1_string_1_letter_end == "л" || to_rune_1_string_1_letter_end == "м" || to_rune_1_string_1_letter_end == "н" || to_rune_1_string_1_letter_end == "п" || to_rune_1_string_1_letter_end == "р" || to_rune_1_string_1_letter_end == "с" || to_rune_1_string_1_letter_end == "т" || to_rune_1_string_1_letter_end == "ф" || to_rune_1_string_1_letter_end == "х" || to_rune_1_string_1_letter_end == "ц" || to_rune_1_string_1_letter_end == "ч" || to_rune_1_string_1_letter_end == "ш" || to_rune_1_string_1_letter_end == "щ" {
				result := to_rune_1_string + "а"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_1_string_1_letter_end == "ь" {
				delete_1_letter := to_rune_1_string_1_letter_end_cut
				result := delete_1_letter + "я"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				result := to_rune_1_string
				entry_data_result_array = append(entry_data_result_array, result)
				break
			}
		}
		for i := 0; i <= len(entry_data_fields_array[2]); i++ {
			to_rune_2 := []rune(entry_data_fields_array[2])
			to_rune_2_string := string(to_rune_2)
			result := to_rune_2_string
			entry_data_result_array = append(entry_data_result_array, result)
			break
		}
		for i := 0; i <= len(entry_data_fields_array[3]); i++ {
			to_rune_3 := []rune(entry_data_fields_array[3])
			to_rune_3_string := string(to_rune_3)
			to_rune_3_string_1_letter_end := string(to_rune_3[len(to_rune_3)-1:])
			if to_rune_3_string_1_letter_end == "т" {
				result := to_rune_3_string + "а"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				result := to_rune_3_string
				entry_data_result_array = append(entry_data_result_array, result)
				break
			}
		}
		for _, elem := range entry_data_fields_array[4:] {
			result := string(elem)
			entry_data_result_array = append(entry_data_result_array, result)
		}
	}
	entry_data_result_array_string := strings.Join(entry_data_result_array, " ")
	logrus.Infof("07.0: Преобразование должностей в р.п.: %v → %v", entry_data, entry_data_result_array_string)
	return entry_data_result_array_string
}

// ****функция для преобразования должности руководителя компании в творительный падеж
func ConvertEntryDataHeadCompanyPositionToInstrumentalCase(entry_data string) string {
	entry_data_fields_array := strings.Fields(entry_data)
	entry_data_result_array := []string{}
	// условный перебор элементов массива согласно длине массива
	if len(entry_data_fields_array) == 1 {
		for i := 0; i <= len(entry_data_fields_array[0]); i++ {
			to_rune_0 := []rune(entry_data_fields_array[0])
			to_rune_0_string := string(to_rune_0)
			to_rune_0_string_1_letter_end := string(to_rune_0[len(to_rune_0)-1:])
			if to_rune_0_string_1_letter_end == "б" || to_rune_0_string_1_letter_end == "в" || to_rune_0_string_1_letter_end == "г" || to_rune_0_string_1_letter_end == "д" || to_rune_0_string_1_letter_end == "ж" || to_rune_0_string_1_letter_end == "з" || to_rune_0_string_1_letter_end == "к" || to_rune_0_string_1_letter_end == "л" || to_rune_0_string_1_letter_end == "м" || to_rune_0_string_1_letter_end == "н" || to_rune_0_string_1_letter_end == "п" || to_rune_0_string_1_letter_end == "р" || to_rune_0_string_1_letter_end == "с" || to_rune_0_string_1_letter_end == "т" || to_rune_0_string_1_letter_end == "ф" || to_rune_0_string_1_letter_end == "х" || to_rune_0_string_1_letter_end == "ц" || to_rune_0_string_1_letter_end == "ч" || to_rune_0_string_1_letter_end == "ш" || to_rune_0_string_1_letter_end == "щ" {
				result := to_rune_0_string + "ом"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			}
		}
	} else if len(entry_data_fields_array) >= 2 {
		for i := 0; i <= len(entry_data_fields_array[0]); i++ {
			to_rune_0 := []rune(entry_data_fields_array[0])
			// to_rune_0_string := string(to_rune_0)
			to_rune_0_string_2_letter_end := string(to_rune_0[len(to_rune_0)-2:])
			to_rune_0_string_2_letter_cut := string(to_rune_0[:len(to_rune_0)-2])
			if to_rune_0_string_2_letter_end == "ый" {
				result := to_rune_0_string_2_letter_cut + "ым"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_0_string_2_letter_end == "ий" {
				result := to_rune_0_string_2_letter_cut + "им"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			}
		}
		for i := 0; i <= len(entry_data_fields_array[1]); i++ {
			to_rune_1 := []rune(entry_data_fields_array[1])
			to_rune_1_string := string(to_rune_1)
			to_rune_1_string_1_letter_end := string(to_rune_1[len(to_rune_1)-1:])
			if to_rune_1_string_1_letter_end == "б" || to_rune_1_string_1_letter_end == "в" || to_rune_1_string_1_letter_end == "г" || to_rune_1_string_1_letter_end == "д" || to_rune_1_string_1_letter_end == "ж" || to_rune_1_string_1_letter_end == "з" || to_rune_1_string_1_letter_end == "к" || to_rune_1_string_1_letter_end == "л" || to_rune_1_string_1_letter_end == "м" || to_rune_1_string_1_letter_end == "н" || to_rune_1_string_1_letter_end == "п" || to_rune_1_string_1_letter_end == "р" || to_rune_1_string_1_letter_end == "с" || to_rune_1_string_1_letter_end == "т" || to_rune_1_string_1_letter_end == "ф" || to_rune_1_string_1_letter_end == "х" || to_rune_1_string_1_letter_end == "ц" || to_rune_1_string_1_letter_end == "ч" || to_rune_1_string_1_letter_end == "ш" || to_rune_1_string_1_letter_end == "щ" {
				result := to_rune_1_string + "ом"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			}
		}
	}
	entry_data_result_array_string := strings.Join(entry_data_result_array, " ")
	logrus.Infof("07.0: Преобразование должности руководителя компании в т.п.: %v → %v", entry_data, entry_data_result_array_string)
	return entry_data_result_array_string
}

// ****функция для преобразования наименований подразделений в родительный падеж
func ConvertEntryDataDepartmentToGenitiveCase(entry_data string) string {
	entry_data_fields_array := strings.Fields(entry_data)
	entry_data_result_array := []string{}
	// условный перебор элементов массива согласно длине массива
	if len(entry_data_fields_array) >= 1 {
		for i := 0; i <= len(entry_data_fields_array[0]); i++ {
			to_rune_0 := []rune(entry_data_fields_array[0])
			to_rune_0_string := string(to_rune_0)
			to_rune_0_string_1_letter_end := string(to_rune_0[len(to_rune_0)-1:])
			to_rune_0_string_1_letter_cut := string(to_rune_0[:len(to_rune_0)-1])
			to_rune_0_string_2_letter_end := string(to_rune_0[len(to_rune_0)-2:])
			to_rune_0_string_2_letter_cut := string(to_rune_0[:len(to_rune_0)-2])
			if to_rune_0_string_1_letter_end == "б" || to_rune_0_string_1_letter_end == "в" || to_rune_0_string_1_letter_end == "г" || to_rune_0_string_1_letter_end == "д" || to_rune_0_string_1_letter_end == "ж" || to_rune_0_string_1_letter_end == "з" || to_rune_0_string_1_letter_end == "к" || to_rune_0_string_1_letter_end == "л" || to_rune_0_string_1_letter_end == "м" || to_rune_0_string_1_letter_end == "н" || to_rune_0_string_1_letter_end == "п" || to_rune_0_string_1_letter_end == "р" || to_rune_0_string_1_letter_end == "с" || to_rune_0_string_1_letter_end == "т" || to_rune_0_string_1_letter_end == "ф" || to_rune_0_string_1_letter_end == "х" || to_rune_0_string_1_letter_end == "ц" || to_rune_0_string_1_letter_end == "ч" || to_rune_0_string_1_letter_end == "ш" || to_rune_0_string_1_letter_end == "щ" {
				plus_1_letter := to_rune_0_string
				result := plus_1_letter + "а"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_0_string_2_letter_end == "ие" {
				delete_2_letter := to_rune_0_string_2_letter_cut
				result := delete_2_letter + "ия"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_0_string_1_letter_end == "а" {
				delete_1_letter := to_rune_0_string_1_letter_cut
				result := delete_1_letter + "ы"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_0_string_2_letter_end == "ия" {
				delete_2_letter := to_rune_0_string_2_letter_cut
				result := delete_2_letter + "ии"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				if to_rune_0_string == "-" || to_rune_0_string == "NULL" || to_rune_0_string == "0" {
					result := " "
					entry_data_result_array = append(entry_data_result_array, result)
					break
				}
			}
		}
		for _, elem := range entry_data_fields_array[1:] {
			result := string(elem)
			entry_data_result_array = append(entry_data_result_array, result)
			continue
		}
	}
	entry_data_result_array_string := strings.Join(entry_data_result_array, " ")
	logrus.Infof("07.0: Преобразование наименований подразделений в р.п.: %v → %v", entry_data, entry_data_result_array_string)
	return entry_data_result_array_string
}

// ****функция для преобразования наименований подразделений в предложный падеж
func ConvertEntryDataDepartmentToPrepositionalCase(entry_data string) string {
	entry_data_fields_array := strings.Fields(entry_data)
	entry_data_result_array := []string{}
	// условный перебор элементов массива согласно длине массива
	if len(entry_data_fields_array) >= 1 {
		for i := 0; i <= len(entry_data_fields_array[0]); i++ {
			to_rune_0 := []rune(entry_data_fields_array[0])
			to_rune_0_string := string(to_rune_0)
			to_rune_0_string_1_letter_end := string(to_rune_0[len(to_rune_0)-1:])
			to_rune_0_string_1_letter_cut := string(to_rune_0[:len(to_rune_0)-1])
			to_rune_0_string_2_letter_end := string(to_rune_0[len(to_rune_0)-2:])
			to_rune_0_string_2_letter_cut := string(to_rune_0[:len(to_rune_0)-2])
			if to_rune_0_string_1_letter_end == "б" || to_rune_0_string_1_letter_end == "в" || to_rune_0_string_1_letter_end == "г" || to_rune_0_string_1_letter_end == "д" || to_rune_0_string_1_letter_end == "ж" || to_rune_0_string_1_letter_end == "з" || to_rune_0_string_1_letter_end == "к" || to_rune_0_string_1_letter_end == "л" || to_rune_0_string_1_letter_end == "м" || to_rune_0_string_1_letter_end == "н" || to_rune_0_string_1_letter_end == "п" || to_rune_0_string_1_letter_end == "р" || to_rune_0_string_1_letter_end == "с" || to_rune_0_string_1_letter_end == "т" || to_rune_0_string_1_letter_end == "ф" || to_rune_0_string_1_letter_end == "х" || to_rune_0_string_1_letter_end == "ц" || to_rune_0_string_1_letter_end == "ч" || to_rune_0_string_1_letter_end == "ш" || to_rune_0_string_1_letter_end == "щ" {
				plus_1_letter := to_rune_0_string
				result := plus_1_letter + "е"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_0_string_2_letter_end == "ие" || to_rune_0_string_2_letter_end == "ия" {
				delete_2_letter := to_rune_0_string_2_letter_cut
				result := delete_2_letter + "ии"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else if to_rune_0_string_1_letter_end == "а" {
				delete_1_letter := to_rune_0_string_1_letter_cut
				result := delete_1_letter + "ы"
				entry_data_result_array = append(entry_data_result_array, result)
				break
			} else {
				if to_rune_0_string == "-" || to_rune_0_string == "NULL" || to_rune_0_string == "0" {
					result := " "
					entry_data_fields_array = append(entry_data_fields_array, result)
					break
				}
			}
		}
		for _, elem := range entry_data_fields_array[1:] {
			result := string(elem)
			entry_data_result_array = append(entry_data_result_array, result)
			continue
		}
	}
	entry_data_result_array_string := strings.Join(entry_data_result_array, " ")
	logrus.Infof("07.0: Преобразование наименований подразделений в п.п.: %v → %v", entry_data, entry_data_result_array_string)
	return entry_data_result_array_string
}

// *****функция для проверки первичных данных (имя, фамилия, дата рождения, компания) на соответствие с данными базы данных
func EntryTextMatchSQLDataTableStepOne(entry_text_array []string, SQL_data_table_array []string) bool {
	var status bool
	status = true
	for i := 0; i < len(entry_text_array)-2; i++ {
		if entry_text_array[i] == SQL_data_table_array[i] {
			logrus.Infof("06.0: 1-й этап проверки первичных данных (имя / фамилия / дата рождения + компания) на соответствие с данными из таблицы: SUCCESS → `%v` == `%v`", entry_text_array[i], SQL_data_table_array[i])
			status = true
			break
		} else {
			logrus.Errorf("06.0: 1-й этап проверки первичных данных (имя / фамилия / дата рождения + компания) на соответствие с данными из таблицы: UNSUCCESS → `%v` != `%v`", entry_text_array[i], SQL_data_table_array[i])
			status = false
			break
		}
	}
	return status
}

// *****функция для проверки вторичных данных (уровень функционала, статус разработчика ПО) на соответствие с данными базы данных
func SelectAnOptionMatchSQLDataTableStepTwo(entry_text_array []string, SQL_data_table_array []string) bool {
	var status bool
	status = true
	for i := 4; i < len(entry_text_array); i++ {
		if entry_text_array[i] == SQL_data_table_array[i] {
			logrus.Infof("08.0: 1-й этап проверки вторичных данных (уровень функционала / статус разработчика) на соответствие с данными из таблицы: SUCCESS → `%v` == `%v`", entry_text_array[i], SQL_data_table_array[i])
			status = true
		} else {
			logrus.Warningf("08.0: 1-й этап проверки вторичных данных (уровень функционала / статус разработчика) на соответствие с данными из таблицы: UNSUCCESS → `%v` != `%v`", entry_text_array[i], SQL_data_table_array[i])
			status = false
			break
		}
	}
	return status
}

// ******функция для отрытия папки с лог-файлом
func OpenLogFolder(url string) {
	err := open.Start(url)
	if err != nil {
		log.Printf("00.2: Открытие папки с лог-файлом программы: UNSUCCESS → %v", err)
	}
	log.Println("00.2: Открытие папки с лог-файлом программы: SUCCESS")
}

// *******создание структуры для базы данных
type Employee struct {
	// о сотруднике - вводная информация
	employee_name               string // имя сотрудника
	employee_surname            string // фамилия сотрудника
	employee_brief_name_surname string // сокращенное ФИО сотрудника (например: Н.Н. Петров)
	employee_birth_date         string // дата рождения сотрудника

	// о компании, её руководителе и подписанте должностной инструкции
	employee_company_name           string // название компании сотрудника
	company_head_position           string // должность руководителя компании
	jd_signatory_position           string // должность подписанта должностной инструкции
	jd_signatory_brief_name_surname string // сокращенное ФИО подписанта должностной инструкции (например: Н.Н. Петров)

	// о сотруднике - дополнительная информация
	employee_position string // должность сотрудника

	employee_management_level  string // уровень функционала сотрудника
	employee_programmer_status string // статус разработчика ПО

	employee_department_5_level        string // подразделение сотрудника 5 уровня (например, отдел) - данное направление подразделение может не использоваться (NULL)
	employee_department_4_level        string // подразделение сотрудника 4 уровня (например, отдел)
	employee_department_3_level        string // подразделение сотрудника 3 уровня (например, департамент)
	employee_department_2_level        string // подразделение сотрудника 2 уровня (например, дирекция)
	employee_department_1_level        string // подразделение сотрудника 1 уровня (например, блок программных продуктов и решений)
	employee_essential_education       string // требуемый уровень образования сотрудника (например, высшее)
	employee_essential_work_experience string // требуемый стаж работы (например, 1-3 лет)
	employee_job_duties                string // должностные обязанности сотрудника, то есть функционал (например: управление проектами; финансовое моделирование)

	// о прямом/административном и функциональном руководителе сотрудника
	employee_adm_head_position           string // должность прямого/административного руководителя сотрудника
	employee_adm_head_department_5_level string // подразделение прямого/административного руководителя сотрудника 5 уровня (например, отдел) - данное направление подразделение может не использоваться (NULL)
	employee_adm_head_department_4_level string // подразделение прямого/административного руководителя сотрудника 4 уровня (например, отдел)
	employee_adm_head_department_3_level string // подразделение прямого/административного руководителя сотрудника 3 уровня (например, департамент)
	employee_adm_head_department_2_level string // подразделение прямого/административного руководителя сотрудника 2 уровня (например, дирекция)
	employee_adm_head_department_1_level string // подразделение прямого/административного руководителя сотрудника 1 уровня (например, блок программных продуктов и решений)

	employee_fun_head_position           string // должность функционального руководителя сотрудника
	employee_fun_head_department_5_level string // подразделение функционального руководителя сотрудника 5 уровня (например, отдел) - данное направление подразделение может не использоваться (NULL)
	employee_fun_head_department_4_level string // подразделение функционального руководителя сотрудника 4 уровня (например, отдел)
	employee_fun_head_department_3_level string // подразделение функционального руководителя сотрудника 3 уровня (например, департамент)
	employee_fun_head_department_2_level string // подразделение функционального руководителя сотрудника 2 уровня (например, дирекция)
	employee_fun_head_department_1_level string // подразделение функционального руководителя сотрудника 1 уровня (например, блок программных продуктов и решений)
}

func main() {

	// -------- СОЗДАНИЕ ЛОГ-ФАЙЛА ДЛЯ ЗАПИСИ ПРОЦЕССОВ ВЫПОЛНЕНИЯ ПРОГРАММЫ
	log_file, err := os.OpenFile("Лог для SimplyFire", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logrus.Errorf("00.1: Создание лог-файла Simply.🚀.Fire: UNSUCCESS → %v", err)
		return
	}
	log.Printf("00.1: Создание лог-файла SimplyFire: SUCCESS")
	defer log_file.Close()

	logrus.SetOutput(log_file)

	// -------- СОЗДАНИЕ ОСНОВНОГО ОКНА WINDOW_MAJOR И ДВУХ ТЕМ
	// создание основного окна window_major с базовой темной темой
	application := app.New()
	window_major := application.NewWindow("SimplyFire")
	window_major.Resize(fyne.NewSize(420, 529))
	window_major.SetFixedSize(true)
	window_major.CenterOnScreen()
	application.Settings().SetTheme(theme.DarkTheme())

	// -------- СОЗДАНИЕ КНОПКИ "О программе"
	// создание кнопки "О программе"
	about_button := widget.NewButtonWithIcon(
		"",
		theme.InfoIcon(),
		func() {
			dialog.ShowCustom(
				"О программе SimplyFire",
				"ОК",
				widget.NewLabel("Программа предназначена для автоматизации\nсоставления должностных инструкций\nсотрудников компаний.\n\nВерсия программы: 1.0\n"),
				window_major,
			)
		},
	)

	about_button.Importance = widget.LowImportance
	about_button.Resize(fyne.NewSize(13, 13))
	about_button.Move(fyne.NewPos(309, 30))
	about_button_container := container.NewWithoutLayout(about_button)

	// -------- СОЗДАНИЕ ПОЛЕЙ ВВОДА ДАННЫХ, СЕЛЕКТИВНЫХ И ВСПОМОГАТЕЛЬНЫХ КНОПОК
	// цвет шрифта для текстовых полей
	window_major_color_for_text := color.NRGBA{R: 1, G: 130, B: 245, A: 255}

	// создание полей ввода данных о сотруднике
	entry_request := canvas.NewText("Введите Ваше имя", window_major_color_for_text)
	entry_request.Resize(fyne.NewSize(200, 40))
	entry_request.Move(fyne.NewPos(5, 30))

	entry_name := widget.NewEntry()
	entry_name.PlaceHolder = "..."
	entry_name.Resize(fyne.NewSize(402, 40))
	entry_name.Move(fyne.NewPos(5, 65))

	surname_request := canvas.NewText("Введите Вашу фамилию", window_major_color_for_text)
	surname_request.Resize(fyne.NewSize(200, 40))
	surname_request.Move(fyne.NewPos(5, 100))

	entry_surname := widget.NewEntry()
	entry_surname.PlaceHolder = "..."
	entry_surname.Resize(fyne.NewSize(402, 40))
	entry_surname.Move(fyne.NewPos(5, 135))

	birth_date_request := canvas.NewText("Введите Вашу дату рождения", window_major_color_for_text)
	birth_date_request.Resize(fyne.NewSize(200, 40))
	birth_date_request.Move(fyne.NewPos(5, 170))

	entry_birth_date := widget.NewEntry()
	entry_birth_date.PlaceHolder = "ГГГГ-ММ-ДД"
	entry_birth_date.Resize(fyne.NewSize(402, 40))
	entry_birth_date.Move(fyne.NewPos(5, 205))

	data_input_container := container.NewWithoutLayout(
		entry_name,
		entry_request,
		surname_request,
		entry_surname,
		birth_date_request,
		entry_birth_date,
	)

	// создание селективной кнопки выбора компании
	company_choose := widget.NewSelect(
		[]string{
			`ООО "Ромашка"`,
		},
		func(s string) {
		},
	)

	company_choose.PlaceHolder = "Выберете компанию, где Вы трудоустроены"
	company_choose.Resize(fyne.NewSize(400, 40))
	company_choose.Move(fyne.NewPos(5, 215))
	company_choose_container := container.NewWithoutLayout(company_choose)

	// создание селективной кнопки выбора функционального уровня
	management_level_choose := widget.NewSelect(
		[]string{
			"Руководящий",
			"Линейный",
		},
		func(s string) {},
	)

	management_level_choose.PlaceHolder = "Выберете Ваш функциональный уровень"
	management_level_choose.Resize(fyne.NewSize(400, 40))
	management_level_choose.Move(fyne.NewPos(5, 225))
	management_level_choose_container := container.NewWithoutLayout(management_level_choose)

	// создание кнопки-справки о функциональном уровне
	management_level_reference := widget.NewButtonWithIcon(
		"",
		theme.InfoIcon(),
		func() {
			dialog.ShowCustom(
				"Справка о функциональном уровне",
				"ОК",
				widget.NewLabel("К сотрудникам руководящего уровня относятся\nруководители подразделений и ТОП-менеджмент.\n\nК сотрудникам линейного уровня относятся\nспециалисты, которые выполняют операционные\nфункции.\n"),
				window_major,
			)
		},
	)

	management_level_reference.Importance = widget.LowImportance
	management_level_reference.Resize(fyne.NewSize(13, 13))
	management_level_reference.Move(fyne.NewPos(390, 235))
	management_level_reference_container := container.NewWithoutLayout(management_level_reference)

	// создание селективной кнопки выбора статуса разработчика ПО
	programmer_status := widget.NewSelect(
		[]string{
			"Разработчик ПО",
			"Не разработчик ПО",
		},
		func(s string) {
		},
	)

	programmer_status.PlaceHolder = "Выберете, являетесь ли Вы разработчиком ПО"
	programmer_status.Resize(fyne.NewSize(400, 40))
	programmer_status.Move(fyne.NewPos(5, 180))

	programmer_status_container := container.NewWithoutLayout(programmer_status)

	// -------- СОЗДАНИЕ КНОПКИ "Сформировать должностную инструкцию" И СВЕРКА УКАЗАННЫХ ПОЛЬЗОВАТЕЛЕМ ДАННЫХ С ТАБЛИЦЕЙ БАЗЫ ДАННЫХ
	// создание кнопки "Сформировать должностную инструкцию"
	check_button := widget.NewButton(
		"Сформировать должностную инструкцию",
		func() {
			// создание переменных из указанных пользователем данных
			name_output := entry_name.Text
			name_output_validate_int := IsEntryTextInt(name_output)
			name_output_validate_cyrillic := IsEntryTextCyrillic(name_output)
			surname_output := entry_surname.Text
			surname_output_validate_int := IsEntryTextInt(surname_output)
			surname_output_validate_cyrillic := IsEntryTextCyrillic(surname_output)

			birth_date_output := entry_birth_date.Text
			birth_date_output_validate_with_date_format := EntryTextConvertToDate(birth_date_output)
			company_choose_output := company_choose.Selected
			management_level_choose_output := management_level_choose.Selected
			programmer_status_output := programmer_status.Selected

			// вторичная проверка первичных данных (имя / фамилия / дата рождения) на 1) отсутствие числовых значений 2) знаков, отличных от кириллицы 3) формат даты рождения
			if name_output_validate_int == false && name_output_validate_cyrillic == false && surname_output_validate_int == false && surname_output_validate_cyrillic == false && birth_date_output_validate_with_date_format == true {
				logrus.Info("01.3: 2-й этап проверки первичных данных: SUCCESS")

				// внесение первичных (имя / фамилия / дата рождения + компания) и вторичных (функциональный уровень / статус разработчика ПО) данных в массив для их дальнейшей верификацией таблицей базы данных
				data_array := []string{}
				data_array = append(
					data_array,
					name_output,
					surname_output,
					birth_date_output,
					company_choose_output,
					management_level_choose_output,
					programmer_status_output,
				)
				for index, element := range data_array {
					logrus.Infof("02.0: Внесение указанных данных в массив: %v %v", index, element)
				}

				// установление соединения с базой данных
				employee_database, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/employee-database-test")
				if err != nil {
					logrus.Infof("03.0: Подключение к базе данных: UNSUCCESS → %v", err)
					dialog.ShowCustomConfirm(
						"Ошибка подключения к базе данных",
						"Да",
						"Нет",
						widget.NewLabel("Не удалось установить соединение с\nбазой данных.\n\nОзнакомиться с подробностями ошибки?\n"),
						func(b bool) {
							if b {
								OpenLogFolder("C:/Users/Public/Downloads/")
							}
						},
						window_major,
					)
					return
				}
				defer employee_database.Close()
				logrus.Info("03.0: Подключение к базе данных: SUCCESS")

				// запрос статуса у базы данных по наличию подходящих столбцов таблицы для извлечения из них нужных данных о сотруднике
				data_output, err := employee_database.Query("SELECT `employee_name`, `employee_surname`, `employee_brief_name_surname`, `employee_birth_date`, `employee_company_name`, `company_head_position`, `jd_signatory_position`, `jd_signatory_brief_name_surname`, `employee_position`, `employee_management_level`, `employee_programmer_status`, `employee_department_5_level`, `employee_department_4_level`, `employee_department_3_level`, `employee_department_2_level`, `employee_department_1_level`, `employee_essential_education`, `employee_essential_work_experience`, `employee_job_duties`, `employee_adm_head_position`, `employee_adm_head_department_5_level`, `employee_adm_head_department_4_level`, `employee_adm_head_department_3_level`, `employee_adm_head_department_2_level`, `employee_adm_head_department_1_level`, `employee_fun_head_position`, `employee_fun_head_department_5_level`, `employee_fun_head_department_4_level`, `employee_fun_head_department_3_level`, `employee_fun_head_department_2_level`, `employee_fun_head_department_1_level` FROM `employee-test-update`")
				if err != nil {
					logrus.Errorf("04.0: Выбор нужных столбцов таблицы базы данных для выборки данных: UNSUCCESS → %v", err)
					dialog.ShowCustomConfirm(
						"Ошибка запроса к таблице базы данных",
						"Да",
						"Нет",
						widget.NewLabel("Не удалось определить таблицу базы данных\nи / или столбцы таблицы, необходимые для\nвыборки нужной информации о\nсотруднике.\n\nОзнакомиться с подробностями ошибки?\n"),
						func(b bool) {
							if b {
								OpenLogFolder("C:/Users/Public/Downloads/")
							}
						},
						window_major,
					)
					return
				}
				logrus.Info("04.0: Выбор нужных столбцов таблицы базы данных для выборки данных: SUCCESS")

				// запуск цикла для сохранения содержимого столбцов таблицы в созданной структуре
				for data_output.Next() {
					var employee Employee
					err = data_output.Scan(

						// о сотруднике - вводная информация
						&employee.employee_name,               // имя сотрудника
						&employee.employee_surname,            // фамилия сотрудника
						&employee.employee_brief_name_surname, // сокращенное ФИО сотрудника (например: Н.Н. Петров)
						&employee.employee_birth_date,         // дата рождения сотрудника

						// о компании, её руководителе и подписанте должностной инструкции
						&employee.employee_company_name,           // название компании сотрудника
						&employee.company_head_position,           // должность руководителя компании
						&employee.jd_signatory_position,           // должность подписанта должностной инструкции
						&employee.jd_signatory_brief_name_surname, // сокращенное ФИО руководителя компании (например: Н.Н. Петров)

						// о сотруднике - дополнительная информация
						&employee.employee_position, // должность сотрудника

						&employee.employee_management_level,  // уровень функционала сотрудника
						&employee.employee_programmer_status, // статус разработчика ПО

						&employee.employee_department_5_level,        // подразделение сотрудника 5 уровня (например, отдел) - данное направление подразделение может не использоваться (NULL)
						&employee.employee_department_4_level,        // подразделение сотрудника 4 уровня (например, отдел) - данное направление подразделение может не использоваться (NULL)
						&employee.employee_department_3_level,        // подразделение сотрудника 3 уровня (например, департамент)
						&employee.employee_department_2_level,        // подразделение сотрудника 2 уровня (например, дирекция)
						&employee.employee_department_1_level,        // подразделение сотрудника 1 уровня (например, блок программных продуктов и решений)
						&employee.employee_essential_education,       // требуемый уровень образования (например, высшее)
						&employee.employee_essential_work_experience, // требуемый стаж работы (например, 2 года)
						&employee.employee_job_duties,                // должностные обязанность, то есть функционал (например: управление проектами; финансовое моделирование)

						// о прямом/административном и функциональном руководителе сотрудника
						&employee.employee_adm_head_position,           // должность прямого/административного руководителя сотрудника
						&employee.employee_adm_head_department_5_level, // подразделение прямого/административного руководителя сотрудника 5 уровня (например, отдел) - данное направление подразделение может не использоваться (NULL)
						&employee.employee_adm_head_department_4_level, // подразделение прямого/административного руководителя сотрудника 4 уровня (например, отдел) - данное направление подразделение может не использоваться (NULL)
						&employee.employee_adm_head_department_3_level, // подразделение прямого/административного руководителя сотрудника 3 уровня (например, департамент)
						&employee.employee_adm_head_department_2_level, // подразделение прямого/административного руководителя сотрудника 2 уровня (например, дирекция)
						&employee.employee_adm_head_department_1_level, // подразделение прямого/административного руководителя сотрудника 1 уровня (например, блок программных продуктов и решений)

						&employee.employee_fun_head_position,           // должность функционального руководителя сотрудника
						&employee.employee_fun_head_department_5_level, // подразделение функционального руководителя сотрудника 5 уровня (например, отдел) - данное направление подразделение может не использоваться (NULL)
						&employee.employee_fun_head_department_4_level, // подразделение функционального руководителя сотрудника 4 уровня (например, отдел) - данное направление подразделение может не использоваться (NULL)
						&employee.employee_fun_head_department_3_level, // подразделение функционального руководителя сотрудника 3 уровня (например, департамент)
						&employee.employee_fun_head_department_2_level, // подразделение функционального руководителя сотрудника 2 уровня (например, дирекция)
						&employee.employee_fun_head_department_1_level, // подразделение функционального руководителя сотрудника 1 уровня (например, блок программных продуктов и решений)
					)
					if err != nil {
						logrus.Errorf("05.0: Выборка данных из таблицы базы данных: UNSUCCESS → %v", err)
						dialog.ShowCustomConfirm(
							"Ошибка выборки из таблицы базы данных",
							"Да",
							"Нет",
							widget.NewLabel("Не удалось осуществить выборку данных о\nсотруднике из таблицы базы данных.\n\nОзнакомиться с подробностями ошибки?\n"),
							func(b bool) {
								if b {
									OpenLogFolder("C:/Users/Public/Downloads/")
								}
							},
							window_major,
						)
						return
					} else {
						logrus.Info("05.0: Выборка данных из таблицы базы данных: SUCCESS")

						// внесение первичных и вторичных сотрудника из таблицы базы данных в массив для последующей сверки двух массивов: "пользовательский" массив и массив данных из таблицы базы данных
						SQL_data_table_array := []string{employee.employee_name, employee.employee_surname, employee.employee_birth_date, employee.employee_company_name, employee.employee_management_level, employee.employee_programmer_status}

						// присвоение переменной entry_text_array значение переменной data_array (массив с указанными пользователем программы первичными и вторичными данными)
						entry_text_array := data_array

						// создание переменной для присвоения ей результата функции по проверке первичных данных
						entry_text_match_SQL_data_table_step_one_result := EntryTextMatchSQLDataTableStepOne(entry_text_array, SQL_data_table_array)

						// проверка двух массивов (массив с данными пользователя и массив с данными таблицы базы данных) с первичными и вторичными данными на соответствие
						if entry_text_match_SQL_data_table_step_one_result == true {
							logrus.Infof("06.1: 2-й этап проверки первичных данных (имя / фамилия / дата рождения + компания) на соответствие с данными таблицы базы данных: SUCCESS → %v", entry_text_match_SQL_data_table_step_one_result)
							dialog.ShowCustom(
								"Успешная верификация первичных данных",
								"ОК",
								widget.NewLabel("Указанные пользователем первичные данные\nопределены таблицей базы данных:\nсотрудник найден.\n\nВыполняется верификация вторичных данных\nпользователя: уровень функционала и\nстатус разработчика ПО.\n"),
								window_major,
							)
							time.Sleep(3 * time.Second)

							// при соответствии элементов первичных данных двух массивов создается карта с ключами шаблонизатора должностной инструкции и значениями ключей, взятыми из подходящих столбцов таблицы базы данных
							replaceMap := docx.PlaceholderMap{

								// о компании, её руководителе и подписанте должностной инструкции
								"_employee_company_name_":           employee.employee_company_name,                                                        // название компании сотрудника
								"_company_head_position_g_c_":       ConvertEntryDataPositionToGenitiveCase(employee.company_head_position),                // должность руководителя компании в родительном падеже
								"_company_head_position_i_c_":       ConvertEntryDataHeadCompanyPositionToInstrumentalCase(employee.company_head_position), // должность руководителя компании в творительном падеже
								"_jd_signatory_position_n_c_":       employee.jd_signatory_position,                                                        // должность подписанта должностной инструкции в именительном падежа
								"_jd_signatory_brief_name_surname_": employee.jd_signatory_brief_name_surname,                                              // сокращенное ФИО подписанта должностной инструкции (например: Н.Н. Петров)

								// о сотруднике
								"_employee_brief_name_surname_": employee.employee_brief_name_surname, // сокращенное ФИО сотрудника

								"_employee_position_g_c_": ConvertEntryDataPositionToGenitiveCase(employee.employee_position), // должность сотрудника в родительном падеже

								"_employee_department_5_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_department_5_level),      // подразделение сотрудника 5 уровня (например, отдел) в родительном падеже - данное подразделение может не иметь значений (NULL)
								"_employee_department_5_level_p_c_": ConvertEntryDataDepartmentToPrepositionalCase(employee.employee_department_5_level), // подразделение сотрудника 5 уровня (например, отдел) в предложном падеже - данное подразделение может не иметь значений (NULL)

								"_employee_department_4_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_department_4_level),      // подразделение сотрудника 4 уровня (например, отдел) в родительном падеже
								"_employee_department_4_level_p_c_": ConvertEntryDataDepartmentToPrepositionalCase(employee.employee_department_4_level), // подразделение сотрудника 4 уровня (например, отдел) в предложном падеже

								"_employee_department_3_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_department_3_level),      // подразделение сотрудника 3 уровня (например, департамент) в родительном падеже
								"_employee_department_3_level_p_c_": ConvertEntryDataDepartmentToPrepositionalCase(employee.employee_department_3_level), // подразделение сотрудника 3 уровня (например, департамент) в предложном падеже

								"_employee_department_2_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_department_2_level),      // подразделение сотрудника 2 уровня (например, дирекция) в родительном падеже
								"_employee_department_2_level_p_c_": ConvertEntryDataDepartmentToPrepositionalCase(employee.employee_department_2_level), // подразделение сотрудника 2 уровня (например, дирекция) в предложном падеже

								"_employee_department_1_level_n_c_": employee.employee_department_1_level,                                                // подразделение сотрудника 1 уровня (например, блок программных продуктов и решений) в именительном падеже
								"_employee_department_1_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_department_1_level),      // подразделение сотрудника 1 уровня (например, блок программных продуктов и решений) в родительном падеже
								"_employee_department_1_level_p_c_": ConvertEntryDataDepartmentToPrepositionalCase(employee.employee_department_1_level), // подразделение сотрудника 1 уровня (например, блок программных продуктов и решений) в предложном падеже

								"_employee_essential_education_":       employee.employee_essential_education,       // требуемый уровень образования сотрудника (например, высшее)
								"_employee_essential_work_experience_": employee.employee_essential_work_experience, // требуемый стаж работы (например, 1-3 лет)
								"_employee_job_duties_":                employee.employee_job_duties,                // должностные обязанности сотрудника, то есть функционал (например: управление проектами; финансовое моделирование)

								// о прямом/административном и функциональном руководителе сотрудника
								"_employee_adm_head_position_g_c_": ConvertEntryDataPositionToGenitiveCase(employee.employee_adm_head_position), // должность прямого/административного руководителя сотрудника в родительном падеже

								"_employee_adm_head_department_5_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_adm_head_department_5_level), // аналогично сотруднику
								"_employee_adm_head_department_4_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_adm_head_department_4_level), // аналогично сотруднику
								"_employee_adm_head_department_3_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_adm_head_department_3_level), // аналогично сотруднику
								"_employee_adm_head_department_2_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_adm_head_department_2_level), // аналогично сотруднику
								"_employee_adm_head_department_1_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_adm_head_department_1_level), // аналогично сотруднику

								"_employee_fun_head_position_g_c_": ConvertEntryDataPositionToGenitiveCase(employee.employee_fun_head_position), // должность функционального руководителя сотрудника в родительном падеже

								"_employee_fun_head_department_5_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_fun_head_department_5_level), // аналогично сотруднику
								"_employee_fun_head_department_4_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_fun_head_department_4_level), // аналогично сотруднику
								"_employee_fun_head_department_3_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_fun_head_department_3_level), // аналогично сотруднику
								"_employee_fun_head_department_2_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_fun_head_department_2_level), // аналогично сотруднику
								"_employee_fun_head_department_1_level_g_c_": ConvertEntryDataDepartmentToGenitiveCase(employee.employee_fun_head_department_1_level), // аналогично сотруднику
							}

							// создание переменной для присвоения ей результата функции по проверке вторичных данных (функциональный уровень / статус разработчика ПО)
							select_an_option_match_SQL_data_table_step_two_result := SelectAnOptionMatchSQLDataTableStepTwo(entry_text_array, SQL_data_table_array)

							// при соответствии элементов вторичных данных двух массивов начинает выполняться подготовка персонализированных шаблонизаторов должностной инструкции
							if select_an_option_match_SQL_data_table_step_two_result == true {
								logrus.Infof("08.1: 2-й этап проверки вторичных данных (уровень функционала / статус разработчика): SUCCESS → %v", select_an_option_match_SQL_data_table_step_two_result)
								logrus.Infof("09.0: Выгрузка персонализированного шаблонизатора для сотрудников определенного функционального уровня и статуса разработчика ПО: SUCCESS → %v", select_an_option_match_SQL_data_table_step_two_result)
								dialog.ShowCustom(
									"Успешная верификация вторичных данных",
									"ОК",
									widget.NewLabel("Указанные пользователем вторичные данные\nопределены таблицей базы данных.\n\nВыполняется подготовка профильной\nдолжностной инструкции.\n"),
									window_major,
								)
								time.Sleep(3 * time.Second)

								// создание условий для подтверждения функционального уровня и статуса разработчика ПО и дальнейшей подготовки персонализированного шаблонизатора должностной инструкции
								if data_array[4] == "Руководящий" && data_array[5] == "Разработчик ПО" {
									logrus.Info("09.1: Подготовка шаблонизатора должностной инструкции для сотрудника руководящего уровня со статусом разработчика ПО")

									// поиск шаблонизатора должностной инструкции для сотрудника руководящего уровня со статусом разработчика ПО
									template_doc, err := docx.Open("ДИ_руководящий уровень_разработчик ПО.docx")
									if err != nil {
										logrus.Errorf("10.0: Поиск шаблонизатора должностной инструкции: UNSUCCESS (шаблонизатор не обнаружен) → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка поиска шаблонизатора",
											"Да",
											"Нет",
											widget.NewLabel("Не удается найти именной файл с\nшаблонизатором должностной инструкции\nдля сотрудника руководящего уровня со\nстатусом разработчика ПО.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("10.0: Поиск шаблонизатора должностной инструкции: SUCCESS (шаблонизатор обнаружен)")
									}

									// замена ключей из шаблонизатора (ключи аналогичны ключам в карте ReplaceMap) значениями из таблицы базы данных
									err = template_doc.ReplaceAll(replaceMap)
									if err != nil {
										logrus.Errorf("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: UNSUCCESS → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка подстановки значений",
											"Да",
											"Нет",
											widget.NewLabel("Значения из карты ReplaceMap не могут\nбыть внесены в шаблонизатор должностной\nинструкции для сотрудника руководящего\nуровня со статусом разработчика ПО.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: SUCCESS")
									}

									// сохранение должностной инструкции с подтянутыми значениями из таблицы базы данных
									err = template_doc.WriteToFile("C:/Users/Public/Downloads/ДИ_руководящий уровень_разработчик ПО_replaced.docx")
									if err != nil {
										logrus.Errorf("12.0: Сохранение сформированной должностной инструкции: UNSUCCESS → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка сохранения",
											"Да",
											"Нет",
											widget.NewLabel("Сформированная должностная инструкция\nдля сотрудника руководящего уровня со\nстатусом разработчика ПО не может\nбыть успешно сохранена.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("12.0: Сохранение сформированной должностной инструкции: SUCCESS")
										logrus.Info("РЕЗУЛЬТАТ: Должностная инструкция выгружена в директорию: C:/Users/Public/Downloads")
										dialog.ShowCustomConfirm(
											"Преобразование выполнено",
											"Да",
											"Нет",
											widget.NewLabel("Должностная инструкция для сотрудника\nруководящего уровня со статусом разработчика\nПО сформирована и сохранена.\n\nОзнакомиться с должностной инструкцией?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									}
								} else if data_array[4] == "Руководящий" && data_array[5] == "Не разработчик ПО" {
									logrus.Info("09.1: Подготовка шаблонизатора должностной инструкции для сотрудника руководящего уровня без статуса разработчика ПО")

									// поиск шаблонизатора должностной инструкции для сотрудника руководящего уровня без статуса разработчика ПО
									template_doc, err := docx.Open("ДИ_руководящий уровень_не разработчик ПО.docx")
									if err != nil {
										logrus.Errorf("10.0: Поиск шаблонизатора должностной инструкции: UNSUCCESS (шаблонизатор не обнаружен) → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка поиска шаблонизатора",
											"Да",
											"Нет",
											widget.NewLabel("Не удается найти именной файл с\nшаблонизатором должностной инструкции\nдля сотрудника руководящего уровня без\nстатуса разработчика ПО.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("10.0: Поиск шаблонизатора должностной инструкции: SUCCESS (шаблонизатор обнаружен)")
									}

									// замена ключей из шаблонизатора (ключи аналогичны ключам в карте ReplaceMap) значениями из таблицы базы данных
									err = template_doc.ReplaceAll(replaceMap)
									if err != nil {
										logrus.Errorf("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: UNSUCCESS → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка подстановки значений",
											"Да",
											"Нет",
											widget.NewLabel("Значения из карты ReplaceMap не могут\nбыть внесены в шаблонизатор должностной\nинструкции для сотрудника руководящего\nуровня без статуса разработчика ПО.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: SUCCESS")
									}

									// сохранение должностной инструкции с подтянутыми значениями из таблицы базы данных
									err = template_doc.WriteToFile("C:/Users/Public/Downloads/ДИ_руководящий уровень_не разработчик ПО_replaced.docx")
									if err != nil {
										logrus.Errorf("12.0: Сохранение сформированной должностной инструкции: UNSUCCESS → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка сохранения",
											"Да",
											"Нет",
											widget.NewLabel("Сформированная должностная инструкция\nдля сотрудника руководящего уровня без\nстатуса разработчика ПО не может\nбыть успешно сохранена.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("12.0: Сохранение сформированной должностной инструкции: SUCCESS")
										logrus.Info("РЕЗУЛЬТАТ: Должностная инструкция выгружена в директорию: C:/Users/Public/Downloads")
										dialog.ShowCustomConfirm(
											"Преобразование выполнено",
											"Да",
											"Нет",
											widget.NewLabel("Должностная инструкция для сотрудника\nруководящего уровня без статуса разработчика\nПО сформирована и сохранена.\n\nОзнакомиться с должностной инструкцией?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									}
								} else if data_array[4] == "Линейный" && data_array[5] == "Разработчик ПО" {
									logrus.Info("09.1: Подготовка шаблонизатора должностной инструкции для сотрудника линейного уровня со статусом разработчика ПО")

									// поиск шаблонизатора должностной инструкции для сотрудника линейного уровня без статуса разработчика ПО
									template_doc, err := docx.Open("ДИ_линейный уровень_разработчик ПО.docx")
									if err != nil {
										logrus.Errorf("10.0: Поиск шаблонизатора должностной инструкции: UNSUCCESS (шаблонизатор не обнаружен) → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка поиска шаблонизатора",
											"Да",
											"Нет",
											widget.NewLabel("Не удается обнаружить именной файл с\nшаблонизатором должностной инструкции\nдля сотрудника линейного уровня со\nстатусом разработчика ПО.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("10.0: Поиск шаблонизатора должностной инструкции: SUCCESS (шаблонизатор обнаружен)")
									}

									// замена ключей из шаблонизатора (ключи аналогичны ключам в карте ReplaceMap) значениями из таблицы базы данных
									err = template_doc.ReplaceAll(replaceMap)
									if err != nil {
										logrus.Errorf("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: UNSUCCESS → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка подстановки значений",
											"Да",
											"Нет",
											widget.NewLabel("Значения из карты ReplaceMap не могут\nбыть внесены в шаблонизатор должностной\nинструкции для сотрудника линейного\nуровня со статусом разработчика ПО.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: SUCCESS")
									}

									// сохранение должностной инструкции с подтянутыми значениями из таблицы базы данных
									err = template_doc.WriteToFile("C:/Users/Public/Downloads/ДИ_линейный уровень_разработчик ПО_replaced.docx")
									if err != nil {
										logrus.Errorf("12.0: Сохранение сформированной должностной инструкции: UNSUCCESS → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка сохранения",
											"Да",
											"Нет",
											widget.NewLabel("Сформированная должностная инструкция\nдля сотрудника линейного уровня со\nстатусом разработчика ПО не может\nбыть успешно сохранена.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("12.0: Сохранение сформированной должностной инструкции: SUCCESS")
										logrus.Info("РЕЗУЛЬТАТ: Должностная инструкция выгружена в директорию: C:/Users/Public/Downloads")
										dialog.ShowCustomConfirm(
											"Преобразование выполнено",
											"Да",
											"Нет",
											widget.NewLabel("Должностная инструкция для сотрудника\nлинейного уровня со статусом разработчика\nПО сформирована и сохранена.\n\nОзнакомиться с должностной инструкцией?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									}
								} else if data_array[4] == "Линейный" && data_array[5] == "Не разработчик ПО" {
									logrus.Info("09.1: Подготовка шаблонизатора должностной инструкции для сотрудника линейного уровня без статуса разработчика ПО")

									// поиск шаблонизатора должностной инструкции для сотрудника линейного уровня без статуса разработчика ПО
									template_doc, err := docx.Open("ДИ_линейный уровень_не разработчик ПО.docx")
									if err != nil {
										logrus.Errorf("10.0: Поиск шаблонизатора должностной инструкции: UNSUCCESS (шаблонизатор не обнаружен) → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка поиска шаблонизатора",
											"Да",
											"Нет",
											widget.NewLabel("Не удается обнаружить именной файл с\nшаблонизатором должностной инструкции\nдля сотрудника линейного уровня без\nстатуса разработчика ПО.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("10.0: Поиск шаблонизатора должностной инструкции: SUCCESS (шаблонизатор обнаружен)")
									}

									// замена ключей из шаблонизатора (ключи аналогичны ключам в карте ReplaceMap) значениями из таблицы базы данных
									err = template_doc.ReplaceAll(replaceMap)
									if err != nil {
										logrus.Errorf("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: UNSUCCESS → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка подстановки значений",
											"Да",
											"Нет",
											widget.NewLabel("Значения из карты ReplaceMap не могут\nбыть внесены в шаблонизатор должностной\nинструкции для сотрудника линейного\nуровня без статуса разработчика ПО.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: SUCCESS")
									}

									// сохранение должностной инструкции с подтянутыми значениями из таблицы базы данных
									err = template_doc.WriteToFile("C:/Users/Public/Downloads/ДИ_линейный уровень_не разработчик ПО_replaced.docx")
									if err != nil {
										logrus.Errorf("12.0: Сохранение сформированной должностной инструкции: UNSUCCESS → %v", err)
										dialog.ShowCustomConfirm(
											"Ошибка сохранения",
											"Да",
											"Нет",
											widget.NewLabel("Сформированная должностная инструкция\nдля сотрудника линейного уровня без\nстатуса разработчика ПО не может\nбыть успешно сохранена.\n\nОзнакомиться с подробностями ошибки?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									} else {
										logrus.Info("12.0: Сохранение сформированной должностной инструкции: SUCCESS")
										logrus.Info("РЕЗУЛЬТАТ: Должностная инструкция выгружена в директорию: C:/Users/Public/Downloads")
										dialog.ShowCustomConfirm(
											"Преобразование выполнено",
											"Да",
											"Нет",
											widget.NewLabel("Должностная инструкция для сотрудника\nлинейного уровня без статуса разработчика\nПО сформирована и сохранена.\n\nОзнакомиться с должностной инструкцией?\n"),
											func(b bool) {
												if b {
													OpenLogFolder("C:/Users/Public/Downloads/")
												}
											},
											window_major,
										)
										return
									}
								}
							} else {
								// условия функциональгого уровня и статуса разработчика ПО не подтверждены. Выполняется подготовка шаблонизатора долностной инструкции общего вида
								logrus.Infof("08.1: 2-й этап проверки вторичных данных (уровень функционала / статус разработчика): UNSUCCESS → %v", select_an_option_match_SQL_data_table_step_two_result)
								logrus.Warningf("09.0: Выгрузка персонализированного шаблонизатора для сотрудников определенного функционального уровня и статуса разработчика ПО: UNSUCCESS → %v", select_an_option_match_SQL_data_table_step_two_result)
								logrus.Info("09.1: Подготовка должностной инструкции общего вида")
								dialog.ShowCustom(
									"Неуспешная верификация вторичных данных",
									"ОК",
									widget.NewLabel("Указанные пользователем вторичные данные\nне определены таблицей базы данных.\n\nВыполняется подготовка должностной\nинструкции общего вида.\n"),
									window_major,
								)
								time.Sleep(1 * time.Second)

								// поиск шаблонизатора должностной инструкции общего вида
								template_doc, err := docx.Open("ДИ_общий вид.docx")
								if err != nil {
									logrus.Errorf("10.0: Поиск шаблонизатора должностной инструкции: UNSUCCESS (шаблонизатор не обнаружен) → %v", err)
									dialog.ShowCustomConfirm(
										"Ошибка поиска шаблонизатора",
										"Да",
										"Нет",
										widget.NewLabel("Не удается обнаружить именной файл с\nшаблонизатором должностной инструкции\nобщего вида.\n\nОзнакомиться с подробностями ошибки?\n"),
										func(b bool) {
											if b {
												OpenLogFolder("C:/Users/Public/Downloads/")
											}
										},
										window_major,
									)
									return
								} else {
									logrus.Info("10.0: Поиск шаблонизатора должностной инструкции: SUCCESS (шаблонизатор обнаружен)")
								}

								// замена ключей из шаблонизатора (ключи аналогичны ключам в карте ReplaceMap) значениями из таблицы базы данных
								err = template_doc.ReplaceAll(replaceMap)
								if err != nil {
									logrus.Errorf("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: UNSUCCESS → %v", err)
									dialog.ShowCustomConfirm(
										"Ошибка подстановки значений",
										"Да",
										"Нет",
										widget.NewLabel("Значения из карты ReplaceMap не могут\nбыть внесены в шаблонизатор должностной\nинструкции общего вида.\n\nОзнакомиться с подробностями ошибки?\n"),
										func(b bool) {
											if b {
												OpenLogFolder("C:/Users/Public/Downloads/")
											}
										},
										window_major,
									)
									return
								} else {
									logrus.Info("11.0: Внесение значений в шаблонизатор должностной инструкции из карты ReplaceMap: SUCCESS")
								}

								// сохранение должностной инструкции с подтянутыми значениями из таблицы базы данных
								err = template_doc.WriteToFile("C:/Users/Public/Downloads/ДИ_общий вид_replaced.docx")
								if err != nil {
									logrus.Errorf("12.0: Сохранение сформированной должностной инструкции: UNSUCCESS → %v", err)
									dialog.ShowCustomConfirm(
										"Ошибка сохранения",
										"Да",
										"Нет",
										widget.NewLabel("Сформированная должностная инструкция\nобщего вида не может быть успешна\nсохранена.\n\nОзнакомиться с подробностями ошибки?\n"),
										func(b bool) {
											if b {
												OpenLogFolder("C:/Users/Public/Downloads/")
											}
										},
										window_major,
									)
									return
								} else {
									logrus.Info("12.0: Сохранение сформированной должностной инструкции: SUCCESS")
									logrus.Info("РЕЗУЛЬТАТ: Должностная инструкция выгружена в директорию: C:/Users/Public/Downloads")
									dialog.ShowCustomConfirm(
										"Преобразование выполнено",
										"Да",
										"Нет",
										widget.NewLabel("Должностная инструкция общего вида\nсформирована и сохранена.\n\nОзнакомиться с должностной инструкцией?\n"),
										func(b bool) {
											if b {
												OpenLogFolder("C:/Users/Public/Downloads/")
											}
										},
										window_major,
									)
								}
							}
						} else {
							logrus.Errorf("06.1: 2-й этап проверки первичных данных (имя / фамилия / дата рождения + компания) на соответствие с данными таблицы базы данных: UNSUCCESS → %v", entry_text_match_SQL_data_table_step_one_result)
							dialog.ShowCustomConfirm(
								"Ошибка верификации первичных данных",
								"Да",
								"Нет",
								widget.NewLabel("Указанные пользователем первичные данные\nне определены таблицей базы данных:\nсотрудник не найден.\n\nОзнакомиться с подробностями ошибки?\n"),
								func(b bool) {
									if b {
										OpenLogFolder("C:/Users/Public/Downloads/")
									}
								},
								window_major,
							)
						}
					}
				}
			} else {
				logrus.Error("01.3: 2-й этап проверки первичных данных: UNSUCCESS")
				dialog.ShowCustomConfirm(
					"Ошибка проверки вводных данных",
					"Да",
					"Нет",
					widget.NewLabel("Значения имени и фамилии могут содержать\nциферные символы и знаки, отличные\nот кириллицы.\n\nЗначение даты рождения может быть указано\nв неверном формате.\n\nОзнакомиться с подробностями ошибки?\n"),
					func(b bool) {
						if b {
							OpenLogFolder("C:/Users/Public/Downloads/")
						}
					},
					window_major,
				)
				return
			}
		},
	)

	check_button.Importance = widget.HighImportance
	check_button.Resize(fyne.NewSize(350, 40))
	check_button.Move(fyne.NewPos(5, 270))
	check_button_container := container.NewWithoutLayout(check_button)

	// -------- СОЗДАНИЕ КНОПКИ "Очистить указанные данные"

	// создание кнопки "Очистить указанные данные"
	entry_and_choose_data_reset := widget.NewButtonWithIcon(
		"",
		theme.DeleteIcon(),
		func() {
			entry_name.Text = ""
			entry_name.Refresh()
			entry_surname.Text = ""
			entry_surname.Refresh()
			entry_birth_date.Text = ""
			entry_birth_date.Refresh()
			company_choose.Selected = ""
			company_choose.Refresh()
			management_level_choose.Selected = ""
			management_level_choose.Refresh()
			programmer_status.Selected = ""
			programmer_status.Refresh()
		},
	)

	entry_and_choose_data_reset.Importance = widget.LowImportance
	entry_and_choose_data_reset.Resize(fyne.NewSize(40, 40.9))
	entry_and_choose_data_reset.Move(fyne.NewPos(365, 191))
	entry_and_choose_data_reset_container := container.NewWithoutLayout(entry_and_choose_data_reset)

	// -------- СОЗДАНИЕ ФОНА "SimplyFire"
	background := canvas.NewRectangle(color.RGBA{R: 6, G: 45, B: 65, A: 255})

	// -------- СОЗДАНИЕ ЛОГОТИПА "SimplyFire"
	logo_img := canvas.NewImageFromFile("SFire_Logo_1000_1000.png")
	logo_img_container := container.NewWithoutLayout(logo_img)
	logo_img.Resize(fyne.NewSize(223, 223))
	logo_img.Move(fyne.NewPos(98, -66))

	// -------- СОЗДАНИЕ ИКОНКИ "SimplyFire"
	icon, _ := fyne.LoadResourceFromPath("SFire_Icon_1000_1000.png")
	window_major.SetIcon(icon)

	// -------- ВЫВОД КНОПОК ОКНА WINDOW_MAJOR В ВЕРТИКАЛЬНОЕ ПОЛОЖЕНИЕ
	// вывод созданных контейнеров
	window_major.SetContent(
		container.NewStack(
			background,
			container.NewVBox(
				logo_img_container,
				about_button_container,
				data_input_container,
				company_choose_container,
				management_level_choose_container,
				management_level_reference_container,
				check_button_container,
				programmer_status_container,
				entry_and_choose_data_reset_container,
			),
		),
	)

	// -------- НАСТРОЙКИ ЗАПУСКА ПРОГРАММЫ "SimplyFire"
	// вывод окна window_major
	window_major.Show()

	// придание окну window_major статуса основого
	window_major.SetMaster()

	// запуск программы "SimplyFire"
	application.Run()
}
