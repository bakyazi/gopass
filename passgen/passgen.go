package passgen

import (
	"crypto/rand"
	"github.com/pkg/errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

var (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
)

var setMap = map[rune]string{
	'l': lowerCharSet,
	'L': lowerCharSet,
	'u': upperCharSet,
	'U': upperCharSet,
	's': specialCharSet,
	'S': specialCharSet,
	'n': numberSet,
	'N': numberSet,
}

var optionPattern = `^(\d+)\/([lsunLSUN]{1,4})$`

var (
	ErrLengthShouldBeGreater = errors.New("length of password should be greater than number of set selected")
	ErrFailedOptionParse     = errors.New("failed to parse given option")
)

func Generate(options string) (string, error) {
	pg := &passwordGenerator{}
	err := pg.parseOption(options)
	if err != nil {
		return "", err
	}
	return pg.generate()
}

type passwordGenerator struct {
	sets   []string
	length int
}

func (p *passwordGenerator) parseOption(opt string) error {
	re := regexp.MustCompile(optionPattern)
	ss := re.FindStringSubmatch(opt)
	if len(ss) != 3 {
		return errors.Wrap(ErrFailedOptionParse, opt)
	}
	length, err := strconv.Atoi(ss[1])
	if err != nil {
		return err
	}
	p.length = length
	charSets := map[rune]struct{}{}
	for _, o := range []rune(ss[2]) {
		charSets[o] = struct{}{}
	}
	for cs := range charSets {
		p.sets = append(p.sets, setMap[cs])
	}

	return nil
}

func (p passwordGenerator) generate() (string, error) {
	if p.length < len(p.sets) {
		return "", ErrLengthShouldBeGreater
	}
	var sb strings.Builder
	for _, s := range p.sets {
		char, err := p.nextChar([]rune(s))
		if err != nil {
			return "", err
		}
		_, err = sb.WriteRune(char)
		if err != nil {
			return "", err
		}
	}

	for i := 0; i < p.length-len(p.sets); i++ {
		set, err := p.nextSet()
		if err != nil {
			return "", err
		}
		char, err := p.nextChar(set)
		if err != nil {
			return "", err
		}
		_, err = sb.WriteRune(char)
		if err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}

func (p passwordGenerator) nextSet() ([]rune, error) {
	setIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(p.sets))))
	if err != nil {
		return nil, err
	}
	return []rune(p.sets[setIndex.Int64()]), nil
}

func (p passwordGenerator) nextChar(set []rune) (rune, error) {
	runeIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
	if err != nil {
		return 0, err
	}
	return set[runeIndex.Int64()], nil
}
