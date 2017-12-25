package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TimeLayout :
const TimeLayout = "2006-01-02 15:04:05"

// DateLayout :
const DateLayout = "2006-01-02"

// Atot : string to time.Time
func Atot(s string) (time.Time, error) {
	t, err := time.Parse(TimeLayout, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// Ttoa : time.Time to string
func Ttoa(t time.Time) string {
	return t.Format(TimeLayout)
}

// Stoj : struct to json
func Stoj(st interface{}) string {
	s, _ := json.Marshal(st)
	return string(s)
}

// Stom : struct to map
func Stom(st interface{}) map[string]interface{} {
	elem := reflect.ValueOf(&st).Elem()
	typet := elem.Type()
	mapt := map[string]interface{}{}
	for i := 0; i < typet.NumField(); i++ {
		mapt[typet.Field(i).Name] = elem.Field(i).Interface()
	}
	return mapt
}

// Now :
func Now() time.Time {
	now := time.Now()
	_, s := now.Zone()
	now = now.Add(time.Second * time.Duration(s))
	return now
}

// LocalTime :
func LocalTime(t time.Time) time.Time {
	_, s := t.Zone()
	local := t.Add(time.Second * time.Duration(s))
	return local
}

// IntSliceToString :
func IntSliceToString(slice []int64) (str string) {
	for _, v := range slice {
		str = str + strconv.FormatInt(v, 10) + ","
	}
	if len(str) > 0 {
		str = str[0 : len(str)-1]
	}
	return str
}

// StringSliceToString :
func StringSliceToString(slice []string) (str string) {
	for _, v := range slice {
		str = str + v + ","
	}
	if len(str) > 0 {
		str = str[0 : len(str)-1]
	}
	return str
}

// StringSliceToString4SQL :
func StringSliceToString4SQL(slice []string) (str string) {
	for _, v := range slice {
		str = str + "'" + v + "'" + ","
	}
	if len(str) > 0 {
		str = str[0 : len(str)-1]
	}
	return str
}

// GetMonthStringList : Get month string list, "2016-01", ...
func GetMonthStringList(year int) (monthStringList []string) {
	now := time.Now()
	if now.Year() == year {
		for m := time.January; m <= now.Month(); m++ {
			var monthString string
			if m < time.October {
				monthString = strconv.Itoa(year) + "-0" + strconv.Itoa(int(m))
			} else {
				monthString = strconv.Itoa(year) + "-" + strconv.Itoa(int(m))
			}

			monthStringList = append(monthStringList, monthString)
		}
	} else {
		for m := time.January; m <= time.December; m++ {
			var monthString string
			if m < time.October {
				monthString = strconv.Itoa(year) + "-0" + strconv.Itoa(int(m))
			} else {
				monthString = strconv.Itoa(year) + "-" + strconv.Itoa(int(m))
			}

			monthStringList = append(monthStringList, monthString)
		}
	}
	return monthStringList
}

// GetRandomString : Get random string
func GetRandomString(length int) string {
	var randomString string
	array := "abcdefghijklmnopqrstuvwxyz0123456789"
	var randSource = rand.NewSource(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		randNumber := randSource.Int63() % 36
		randomString = randomString + string(array[randNumber])
	}
	return randomString
}

// EncryptByAES : Encrypt by AES
func EncryptByAES(input, secret, iv []byte) (output []byte, err error) {
	//fmt.Printf("input len: %v", len(input))
	//	fmt.Printf("encrypt input len: %v, input: %v\n", len(input), input)
	if len(input)%aes.BlockSize != 0 {
		//return nil, errors.New("input size is not multiple of aes block size")
		for i := 0; i < len(input)%aes.BlockSize; i++ {
			input = append(input, 0x00)
		}
	}
	//	fmt.Printf("encrypt input len: %v, input: %v\n", len(input), input)

	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, errors.New("new aes cipher error")
	}

	output = make([]byte, len(input))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(output, input)

	return output, nil
}

// DecryptByAES : Decrypt by AES
func DecryptByAES(input, secret []byte, iv []byte) (output []byte, err error) {
	//	fmt.Printf("decrypt input len: %v, input: %v\n", len(input), input)
	// CBC mode always works in whole blocks.
	if len(input)%aes.BlockSize != 0 {
		return nil, errors.New("input size is not multiple of the aes block size")
	}
	//	fmt.Printf("decrypt input len: %v, input: %v\n", len(input), input)

	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, errors.New("new aes cipher error")
	}

	output = make([]byte, len(input))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(output, input)

	decryptStr := strings.TrimRight(string(output), string(0x00))
	decryptStr = strings.TrimRight(decryptStr, string(" "))

	return []byte(decryptStr), nil
}

// DeleteTokenParaFromURLParas :
func DeleteTokenParaFromURLParas(urlParas string) (urlParasNoToken string) {
	parasSlice := strings.Split(urlParas, "&")
	parasSliceNoToken := []string{}
	for _, v := range parasSlice {
		if len(v) > 5 && v[:5] == "token" {
			continue
		} else {
			parasSliceNoToken = append(parasSliceNoToken, v)
		}
	}
	urlParasNoToken = strings.Join(parasSliceNoToken, "&")
	return urlParasNoToken
}

// IsPhone : Is phone no or not
func IsPhone(input string) bool {
	var reg = regexp.MustCompile("^((13[0-9])|(15[^4,\\D])|(18[0,2,3,5-9]))\\d{8}$")
	return reg.MatchString(input)
}

// IsEmail : Is email or not
func IsEmail(input string) bool {
	var reg = regexp.MustCompile("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+")
	return reg.MatchString(input)
}

// MD5 : MD5
func MD5(input string) string {
	dataMD5 := md5.Sum([]byte(input))
	dataMD5Slice := []byte{}
	for _, v := range dataMD5 {
		dataMD5Slice = append(dataMD5Slice, v)
	}
	output := hex.EncodeToString(dataMD5Slice)
	return output
}
