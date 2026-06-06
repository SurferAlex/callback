package botapp

// Callback data constants (max 64 bytes each).
const (
	CBTrialMenu     = "trial_menu"
	CBTrialActivate = "trial_activate"
	CBBuyVPN        = "buy_vpn"
	CBPlan1m        = "plan_1m"
	CBPlan3m        = "plan_3m"
	CBPlan6m        = "plan_6m"
	CBPlan12m       = "plan_12m"
	CBPaymentSoon   = "payment_soon"
	CBChangePlan    = "change_plan"
	CBGetConfig     = "get_config"
	CBRefreshConfig = "refresh_config"
	CBInstructions  = "instructions"
	CBGuideIOS      = "guide_ios"
	CBGuideAndroid  = "guide_android"
	CBGuideWin      = "guide_win"
	CBGuideMac      = "guide_mac"
	CBSupport       = "support"
	CBCabinet       = "cabinet"
	CBWebCabinet    = "web_cabinet"
	CBBackMain      = "back_main"
)

func PlanFromCallback(data string) string {
	switch data {
	case CBPlan1m:
		return "1m"
	case CBPlan3m:
		return "3m"
	case CBPlan6m:
		return "6m"
	case CBPlan12m:
		return "12m"
	default:
		return ""
	}
}
