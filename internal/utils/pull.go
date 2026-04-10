package utils

import (
	"math/rand"
	"time"
	"zmd-gacha/internal/models"
)

func Pull(config models.GachaPoolConfig, characters []models.Character, user models.UserPool) models.Character {
	sNoUpCharacters := getSNoUpCharacters(characters)
	sUpCharacters := getSUpCharacters(characters)
	aCharacters := getACharacters(characters)
	bCharacters := getBCharacters(characters)
	user.PullCount++
	// 大保底
	if user.PullCount-user.LastUpCount >= config.LimitPity {
		return randomCharacter(sUpCharacters)
	}

	// 小保底
	if user.PullCount-user.LastSCount >= config.SPityEnd {
		if isUp(config.LimitRateWhenS) {
			return randomCharacter(sUpCharacters)
		}
		return randomCharacter(sNoUpCharacters)
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
			if isUp(config.LimitRateWhenS) {
				return randomCharacter(sUpCharacters)
			}
			return randomCharacter(sNoUpCharacters)
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
		if isUp(config.LimitRateWhenS) {
			return randomCharacter(sUpCharacters)
		}
		return randomCharacter(sNoUpCharacters)
	} else if p < config.SRankBaseRate+config.ARankBaseRate {
		return randomCharacter(aCharacters)
	}
	return randomCharacter(bCharacters)
}

func PullTen(config models.GachaPoolConfig, characters []models.Character, user models.UserPool) []models.Character {
	var results []models.Character
	for i := 0; i < 10; i++ {
		result := Pull(config, characters, user)
		if result.Rank == "S" {
			user.LastSCount = user.PullCount + 1
			if result.IsUp {
				user.LastUpCount = user.PullCount + 1
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

// 是否是UP
func isUp(upRateWhenS float64) bool {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	p := r.Float64()
	return p < upRateWhenS
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

func getSNoUpCharacters(characters []models.Character) []models.Character {
	var sNoUpCharacters []models.Character
	for _, char := range characters {
		if char.Rank == "S" && !char.IsUp {
			sNoUpCharacters = append(sNoUpCharacters, char)
		}
	}
	return sNoUpCharacters
}

func getSUpCharacters(characters []models.Character) []models.Character {
	var sUpCharacters []models.Character
	for _, char := range characters {
		if char.IsUp && char.Rank == "S" {
			sUpCharacters = append(sUpCharacters, char)
		}
	}
	return sUpCharacters
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
