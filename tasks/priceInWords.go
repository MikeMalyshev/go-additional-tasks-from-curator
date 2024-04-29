package tasks

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

func parseBeforeThousand(value int, femaleUnits bool) string {
	hundredsInWords := []string{"", "сто ", "двести ", "триста ", "четыреста ", "пятьсот ", "шестьсот ", "семьсот ", "восемьсот ", "девятьсот "}
	tensInWords := []string{"", "десять ", "двадцать ", "тридцать ", "сорок ", "пятьдесят ", "шестьдесят ", "семьдесят ", "восемьдесят ", "девяносто "}
	onesInWords := []string{"", "один ", "два ", "три ", "четыре ", "пять ", "шесть ", "семь ", "восемь ", "девять "}
	onesInWordsFemale := []string{"", "одна ", "две ", "три ", "четыре ", "пять ", "шесть ", "семь ", "восемь ", "девять "}
	if femaleUnits {
		onesInWords = onesInWordsFemale
	}
	exceptions := []string{"одиннадцать ", "двенадцать ", "тринадцать ", "четырнадцать ", "пятнадцать ", "шестнадцать ", "семнадцать ", "восемнадцать ", "девятнадцать "}

	hundreds := value / 100
	if value%100 > 10 && value%100 < 20 {
		return hundredsInWords[hundreds] + exceptions[value%100-11]
	}

	tens := value % 100 / 10
	ones := value % 10
	return hundredsInWords[hundreds] + tensInWords[tens] + onesInWords[ones]
}

func getCorrectTextNumeral(value int) int {
	if value%100 > 10 && value%100 < 20 {
		return 2
	} else if value%10 == 1 {
		return 0
	} else if value%10 > 1 && value%10 < 5 {
		return 1
	} else {
		return 2
	}
}

func parseInteger(integerValue int) string {
	thousandDimensions := [...]string{"тысяча", "тысячи", "тысяч"}
	millionDimensions := [...]string{"миллион", "миллиона", "миллионов"}
	billionDimensions := [...]string{"миллиард", "миллиарда", "миллиардов"}
	rublesDimensions := [...]string{"рубль", "рубля", "рублей"}

	billions := integerValue / 1e9
	millions := integerValue % 1e9 / 1e6
	thousands := integerValue % 1e6 / 1e3
	other := integerValue % 1e3

	var textValue string

	if billions != 0 {
		textValue += parseBeforeThousand(billions, false) + billionDimensions[getCorrectTextNumeral(billions)] + " "
	}
	if millions != 0 {
		textValue += parseBeforeThousand(millions, false) + millionDimensions[getCorrectTextNumeral(millions)] + " "
	}
	if thousands != 0 {
		textValue += parseBeforeThousand(thousands, true) + thousandDimensions[getCorrectTextNumeral(thousands)] + " "
	}
	if integerValue != 0 {
		textValue += parseBeforeThousand(other, false) + rublesDimensions[getCorrectTextNumeral(integerValue)] + " "
	} else {
		textValue = "Ноль " + rublesDimensions[getCorrectTextNumeral(integerValue)] + " "
	}
	return textValue
}

func parseFraction(fraction float64) string {
	value := int(math.Round((fraction - math.Floor(fraction)) * 100))

	kopeckDimensions := [...]string{"копейка", "копейки", "копеек"}
	var textValue string
	if value == 0 {
		textValue = "0 " + kopeckDimensions[getCorrectTextNumeral(value)]
	} else {
		textValue = parseBeforeThousand(value, false) + kopeckDimensions[getCorrectTextNumeral(value)]
	}
	return textValue
}

func ValueInTextFormat(value float64) (string, error) {
	if value >= 1e12 {
		return "", fmt.Errorf("value is too big")
	}
	if value < 0 {
		return "", fmt.Errorf("incorrect negative value")
	}
	result := parseInteger(int(value)) + parseFraction(value)

	return result, nil
}

func TestValueInTextFormat() {
	fmt.Println("\nThis utility converts the numeral price to price in words, for example:")

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	price1 := float64(r1.Intn(1e3-1)) + float64(r1.Intn(99))/100.0
	price2 := float64(r1.Intn(1e6-1)) + float64(r1.Intn(99))/100.0
	price3 := float64(r1.Intn(1e12 - 1))
	result1, _ := ValueInTextFormat(price1)
	fmt.Printf("\n\t%.2f : %s\n", price1, result1)
	result2, _ := ValueInTextFormat(price2)
	fmt.Printf("\n\t%.2f : %s\n", price2, result2)
	result3, _ := ValueInTextFormat(price3)
	fmt.Printf("\n\t%.2f : %s\n", price3, result3)

	fmt.Printf("\nWhat's your price? Price: ")
	var price float64
	_, err := fmt.Fscan(os.Stdin, &price)
	if err != nil || price < 0 {
		fmt.Println("\nError: Incorrect input\n")
	}
	result, _ := ValueInTextFormat(price)
	fmt.Printf("\n\t%.2f : %s\n", price, result)
}
