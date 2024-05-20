package goodbaddev

import "fmt"

type PurchaseHistory struct{
	TotalAmount int
	PurChaseFrequencyPerMonth int
	ReturnRate float64
}

// Policyパターン
type ExcellentCustomerRule interface {
	ok(ph PurchaseHistory) bool
}

type GoldPerchaseAmountRule struct {}
func (gpar GoldPerchaseAmountRule) ok(ph PurchaseHistory) bool {
	return ph.TotalAmount >= 100000
}

type PurChaseFrequencyRule struct{}
func (gpar PurChaseFrequencyRule) ok(ph PurchaseHistory) bool {
	return ph.PurChaseFrequencyPerMonth >= 10
}

type ReturnRateRule struct{}
func (gpar ReturnRateRule) ok(ph PurchaseHistory) bool {
	return ph.ReturnRate <= 0.001
}

type ExcellentCustomerPolicy struct{
	CustomerRules []ExcellentCustomerRule
}
func NewExcellentCustomerPolicy() ExcellentCustomerPolicy {
	return ExcellentCustomerPolicy {
		CustomerRules: make([]ExcellentCustomerRule, 0),
	}
}
func (ecp *ExcellentCustomerPolicy) Add(ecr ExcellentCustomerRule) {
	ecp.CustomerRules = append(ecp.CustomerRules, ecr)
}
func (ecp *ExcellentCustomerPolicy) ComplyWithAll(ph PurchaseHistory) bool {
	for _, cr := range ecp.CustomerRules {
		if !cr.ok(ph) {
			return false
		}
	}
	return true
}

type GoldCustomerPolicy struct {
	GoldCustomerRules ExcellentCustomerPolicy
}
func (gcr *GoldCustomerPolicy) NewGoldCustomerPolicy() {
	gcr.GoldCustomerRules = NewExcellentCustomerPolicy()
	gcr.GoldCustomerRules.Add(&GoldPerchaseAmountRule{})
	gcr.GoldCustomerRules.Add(&PurChaseFrequencyRule{})
	gcr.GoldCustomerRules.Add(&ReturnRateRule{})
}
func (gcr *GoldCustomerPolicy) ComplyWithAll(ph PurchaseHistory) bool {
	return gcr.GoldCustomerRules.ComplyWithAll(ph)
}

type SilverCustomerPolicy struct {
	SilverCustomerRules ExcellentCustomerPolicy
}
func (gcr *SilverCustomerPolicy) NewSilverCustomerPolicy() {
	gcr.SilverCustomerRules = NewExcellentCustomerPolicy()
	gcr.SilverCustomerRules.Add(&PurChaseFrequencyRule{})
	gcr.SilverCustomerRules.Add(&ReturnRateRule{})
}
func (gcr *SilverCustomerPolicy) ComplyWithAll(ph PurchaseHistory) bool {
	return gcr.SilverCustomerRules.ComplyWithAll(ph)
}

func Policy() {
	user1 := PurchaseHistory{
		TotalAmount: 100000,
		PurChaseFrequencyPerMonth: 20,
		ReturnRate: 0.0,
	}
	user2 := PurchaseHistory{
		TotalAmount: 1000,
		PurChaseFrequencyPerMonth: 20,
		ReturnRate: 0.0,
	}

	goldCustomerPolicy := GoldCustomerPolicy{}
	goldCustomerPolicy.NewGoldCustomerPolicy()
	silverCustomerPolicy := SilverCustomerPolicy{}
	silverCustomerPolicy.NewSilverCustomerPolicy()

	if goldCustomerPolicy.ComplyWithAll(user1) {
		fmt.Println("user1はゴールド会員です")
	} else if silverCustomerPolicy.ComplyWithAll(user1) {
		fmt.Println("user1はシルバー会員です")
	} else {
		fmt.Println("user1は通常会員です")
	}
	fmt.Println("=====")
	if goldCustomerPolicy.ComplyWithAll(user2) {
		fmt.Println("user2はゴールド会員です")
	} else if silverCustomerPolicy.ComplyWithAll(user2) {
		fmt.Println("user2はシルバー会員です")
	} else {
		fmt.Println("user2は通常会員です")
	}

	
}


// 通常の会員判定
// func isGoldCustomer(ph PurchaseHistory) bool {
// 	if ph.TotalAmount >= 100000 {
// 		if ph.PurChaseFrequencyPerMonth >= 10 {
// 			if ph.ReturnRate <= 0.001 {
// 				return true
// 			}
// 		}
// 	}

// 	return false
// }

// func isSilverCustomer(ph PurchaseHistory) bool {
// 	if ph.PurChaseFrequencyPerMonth >= 10 {
// 		if ph.ReturnRate <= 0.001 {
// 			return true
// 		}
// 	}

// 	return false
// }