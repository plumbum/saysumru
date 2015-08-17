package saysumru

import (
	"errors"
	"strconv"
)

const (
	GenderMale = iota
	GenderFemale
	GenderNeuter
)

type sayWords struct {
	FirstDecade       [20]string
	FirstDecadeFemale [2]string
	FirstDecadeNeuter [2]string
	Decades           [8]string
	Hundreds          [9]string
	Thousands         [3]string
	Millions          [3]string
	Billions          [3]string
	Rubles            [3]string
	Kopek             [3]string
	Pieces            [3]string
}

var defaultSayWords sayWords = sayWords{
	FirstDecade: [...]string{
		"ноль",
		"один",
		"два",
		"три",
		"четыре",
		"пять",
		"шесть",
		"семь",
		"восемь",
		"девять",
		"десять",
		"одиннадцать",
		"двенадцать",
		"тринадцать",
		"четырнадцать",
		"пятнадцать",
		"шестнадцать",
		"семнадцать",
		"восемнадцать",
		"девятнадцать",
	},
	FirstDecadeFemale: [...]string{
		"одна",
		"две",
	},
	FirstDecadeNeuter: [...]string{
		"одно",
		"два",
	},
	Decades: [...]string{
		"двадцать",
		"тридцать",
		"сорок",
		"пятьдесят",
		"шестьдесят",
		"семьдесят",
		"восемьдесят",
		"девяносто",
	},
	Hundreds: [...]string{
		"сто",
		"двести",
		"триста",
		"четыреста",
		"пятьсот",
		"шестьсот",
		"семьсот",
		"восемьсот",
		"девятьсот",
	},
	Thousands: [...]string{
		"тысяча",
		"тысячи",
		"тысяч",
	},
	Millions: [...]string{
		"миллион",
		"миллиона",
		"миллионов",
	},
	Billions: [...]string{
		"миллиард",
		"миллиадра",
		"миллиардов",
	},
	Rubles: [...]string{
		"рубль",
		"рубля",
		"рублей",
	},
	Kopek: [...]string{
		"копейка",
		"копейки",
		"копеек",
	},
	Pieces: [...]string{
		"штука",
		"штуки",
		"штук",
	},
}

const (
	numTrilliard = 1e12
	numBillion   = 1e9
	numMillion   = 1e6
	numThousand  = 1e3
)

func firstDecade(val int, gender int) (string, int) {

	var suffixIdx int
	var ns string

	switch {
	case 0 <= val && val <= 19:
		if 1 <= val && val <= 2 && gender != GenderMale {
			ns = defaultSayWords.FirstDecadeFemale[val-1]
		} else {
			ns = defaultSayWords.FirstDecade[val]
		}
		switch {
		case val == 1:
			suffixIdx = 0
		case 2 <= val && val <= 4:
			suffixIdx = 1
		default:
			suffixIdx = 2
		}
	default:
		panic("TODO error " + strconv.Itoa(val))
	}

	return ns, suffixIdx
}

func lessThousand(val int, gender int) ([]string, int) {
	if val >= 1000 {
		panic("TODO error " + strconv.Itoa(val))
	}

	var suffixIdx int
	var skipZero bool

	ns := make([]string, 0, 3) // TODO optimize size

	if val >= 100 {
		ns = append(ns, defaultSayWords.Hundreds[val/100-1])
		val = val % 100
		skipZero = true
	}

	if val >= 20 {
		ns = append(ns, defaultSayWords.Decades[val/10-2])
		val = val % 10
		skipZero = true
	}

	if val == 0 && skipZero {
		suffixIdx = 2
	} else {
		var firstDecadeName string
		firstDecadeName, suffixIdx = firstDecade(val, gender)
		ns = append(ns, firstDecadeName)
	}

	return ns, suffixIdx
}

// Получает целое число и род числительного.
// Возвращает список строк представляющих числительное и тип суффикса.
// Тип суффикса требуется что бы добавить правильный. Например:
//   0 - один рубль
//   1 - два рубля
//   3 - пять рублей
func SayNumber(val int, gender int) (stringed []string, suffixType int) {

	if val >= numTrilliard {
		panic(errors.New("Unsupported value " + strconv.Itoa(val)))
	}

	var suffixIdx int
	var skipZero bool

	ns := make([]string, 0)

	if val < 0 {
		ns = append(ns, "минус")
		val = -val
	}

	if val >= numBillion {
		var billion []string
		billion, suffixIdx = lessThousand(val/numBillion, GenderMale)
		ns = append(ns, billion...)
		ns = append(ns, defaultSayWords.Billions[suffixIdx])
		val = val % numBillion
		skipZero = true
	}

	if val >= numMillion {
		var million []string
		million, suffixIdx = lessThousand(val/numMillion, GenderMale)
		ns = append(ns, million...)
		ns = append(ns, defaultSayWords.Millions[suffixIdx])
		val = val % numMillion
		skipZero = true
	}

	if val >= numThousand {
		var thousand []string
		thousand, suffixIdx = lessThousand(val/numThousand, GenderFemale)
		ns = append(ns, thousand...)
		ns = append(ns, defaultSayWords.Thousands[suffixIdx])
		val = val % numThousand
		skipZero = true
	}

	if val > 0 {
		var firstThousand []string
		firstThousand, suffixIdx = lessThousand(val, gender)
		ns = append(ns, firstThousand...)
	} else {
		suffixIdx = 2
		if !skipZero {
			ns = append(ns, defaultSayWords.FirstDecade[0])
		}
	}

	return ns, suffixIdx
}

// Сказать, сколько рублей
func SayRub(val int) []string {
	var suffixIdx int
	ns, suffixIdx := SayNumber(val, GenderMale)
	ns = append(ns, defaultSayWords.Rubles[suffixIdx])
	return ns
}

// Сказать, сколько копеек
func SayKopek(val int) []string {
	var suffixIdx int
	ns, suffixIdx := SayNumber(val, GenderFemale)
	ns = append(ns, defaultSayWords.Kopek[suffixIdx])
	return ns
}

// Сказать, сколько штук
func SayPieces(val int) []string {
	var suffixIdx int
	ns, suffixIdx := SayNumber(val, GenderFemale)
	ns = append(ns, defaultSayWords.Pieces[suffixIdx])
	return ns
}
