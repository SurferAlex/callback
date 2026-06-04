package botapp

import (
	"os"
	"strconv"
	"strings"
)

type PlanOffer struct {
	Code       string
	Label      string
	Duration   string
	PriceRub   int
	Devices    int
	PlanName   string
}

func LoadPlans() []PlanOffer {
	return []PlanOffer{
		{Code: "1m", Label: "1 месяц", Duration: "1 месяц", PriceRub: priceEnv("PLAN_PRICE_1M", 299), Devices: 2, PlanName: "Premium"},
		{Code: "3m", Label: "3 месяца", Duration: "3 месяца", PriceRub: priceEnv("PLAN_PRICE_3M", 799), Devices: 2, PlanName: "Premium"},
		{Code: "6m", Label: "6 месяцев", Duration: "6 месяцев", PriceRub: priceEnv("PLAN_PRICE_6M", 1490), Devices: 2, PlanName: "Premium"},
		{Code: "12m", Label: "12 месяцев", Duration: "12 месяцев", PriceRub: priceEnv("PLAN_PRICE_12M", 2490), Devices: 2, PlanName: "Premium"},
	}
}

func GetPlan(code string) (PlanOffer, bool) {
	for _, p := range LoadPlans() {
		if p.Code == code {
			return p, true
		}
	}
	return PlanOffer{}, false
}

func priceEnv(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}
