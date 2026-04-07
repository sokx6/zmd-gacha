package utils

import (
	"math/rand"
	"time"
	"zmd-gacha/internal/models"
)

func Pull(config models.GachaPoolConfig, characters []models.Character, user models.User) models.Character {
	sResiCharacters := getSResiCharacters(characters)
	sLimitedCharacters := getSLimitedCharacters(characters)
	aCharacters := getACharacters(characters)
	bCharacters := getBCharacters(characters)
	user.PullCount++
	// 大保底
	if user.PullCount-user.LastLimitedCount >= config.LimitPity {
		return randomCharacter(sLimitedCharacters)
	}

	// 小保底
	if user.PullCount-user.LastSCount >= 80 {
		if isLimited() {
			return randomCharacter(sLimitedCharacters)
		}
		return randomCharacter(sResiCharacters)
	}

	// A保底
	if user.PullCount-user.LastACount >= 10 {
		return randomCharacter(aCharacters)
	}

	// S概率提升
	if user.PullCount-user.LastSCount >= config.SPityStart {
		sRate := config.SRankBaseRate + float64(user.PullCount-user.LastSCount-config.SPityStart)*config.SPityStep
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		p := r.Float64()
		if p < sRate {
			if isLimited() {
				return randomCharacter(sLimitedCharacters)
			}
			return randomCharacter(sResiCharacters)
		} else if p < sRate+config.ARankBaseRate {
			return randomCharacter(aCharacters)
		}
		return randomCharacter(bCharacters)
	}

	// 正常抽卡
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	p := r.Float64()
	if p < config.SRankBaseRate {
		return randomCharacter(sResiCharacters)
	} else if p < config.SRankBaseRate+config.ARankBaseRate {
		return randomCharacter(aCharacters)
	}
	return randomCharacter(bCharacters)
}

// 是否是限定
func isLimited() bool {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	switch r.Intn(2) {
	case 0:
		return false
	case 1:
		return true
	}
	return false
}

func randomCharacter(characters []models.Character) models.Character {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return characters[r.Intn(len(characters))]
}

func getSResiCharacters(characters []models.Character) []models.Character {
	var sResiCharacters []models.Character
	for _, char := range characters {
		if char.Rank == "S" && !char.IsLimited {
			sResiCharacters = append(sResiCharacters, char)
		}
	}
	return sResiCharacters
}

func getSLimitedCharacters(characters []models.Character) []models.Character {
	var sLimitedCharacters []models.Character
	for _, char := range characters {
		if char.IsLimited && char.Rank == "S" {
			sLimitedCharacters = append(sLimitedCharacters, char)
		}
	}
	return sLimitedCharacters
}

func getACharacters(characters []models.Character) []models.Character {
	var aCharacters []models.Character
	for _, char := range characters {
		if char.Rank == "A" {
			aCharacters = append(aCharacters, char)
		}
	}
	return aCharacters
}

func getBCharacters(characters []models.Character) []models.Character {
	var bCharacters []models.Character
	for _, char := range characters {
		if char.Rank == "B" {
			bCharacters = append(bCharacters, char)
		}
	}
	return bCharacters
}
