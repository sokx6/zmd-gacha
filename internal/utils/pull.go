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
	if user.PullCount-user.LastSCount >= config.SPityEnd {
		if isLimited(config.LimitRateWhenS) {
			return randomCharacter(sLimitedCharacters)
		}
		return randomCharacter(sResiCharacters)
	}

	// A保底
	if user.PullCount-user.LastACount >= config.AGuaranteeInterval {
		return randomCharacter(aCharacters)
	}

	// S概率提升
	if user.PullCount-user.LastSCount >= config.SPityStart {
		sRate := config.SRankBaseRate + float64(user.PullCount-user.LastSCount-config.SPityStart)*config.SPityStep
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		p := r.Float64()
		if p < sRate {
			if isLimited(config.LimitRateWhenS) {
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
		if isLimited(config.LimitRateWhenS) {
			return randomCharacter(sLimitedCharacters)
		}
		return randomCharacter(sResiCharacters)
	} else if p < config.SRankBaseRate+config.ARankBaseRate {
		return randomCharacter(aCharacters)
	}
	return randomCharacter(bCharacters)
}

func PullTen(config models.GachaPoolConfig, characters []models.Character, user models.User) []models.Character {
	var results []models.Character
	for i := 0; i < 10; i++ {
		result := Pull(config, characters, user)
		if result.Rank == "S" {
			user.LastSCount = user.PullCount + 1
			if result.IsLimited {
				user.LastLimitedCount = user.PullCount + 1
			}
		}
		if result.Rank == "A" {
			user.LastACount = user.PullCount + 1
		}
		results = append(results, result)
		user.PullCount++
	}
	return results
}

// 是否是限定
func isLimited(limitRateWhenS float64) bool {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	p := r.Float64()
	return p < limitRateWhenS
}

func randomCharacter(characters []models.Character) models.Character {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	if len(characters) == 0 {
		return models.Character{}
	}
	index := r.Intn(len(characters))
	return characters[index]
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
