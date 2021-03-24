// Package rule is to control the audit rule behaviors
package rule

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/XrayR-project/XrayR/api"
)

type RuleManager struct {
	InboundRule         *sync.Map // Key: Tag, Value: []api.DetectRule
	InboundDetectResult *sync.Map // key: Tag, Value: []api.DetectResult
}

func New() *RuleManager {
	return &RuleManager{
		InboundRule:         new(sync.Map),
		InboundDetectResult: new(sync.Map),
	}
}

func (r *RuleManager) UpdateRule(tag string, newRuleList []api.DetectRule) error {
	if value, ok := r.InboundRule.LoadOrStore(tag, newRuleList); ok {
		oldRuleList := value.([]api.DetectRule)
		if !reflect.DeepEqual(oldRuleList, newRuleList) {
			r.InboundRule.Store(tag, newRuleList)
		}
	}
	return nil
}

func (r *RuleManager) GetDetectResult(tag string) (*[]api.DetectResult, error) {
	emptyDetectResult := new([]api.DetectResult)
	if value, ok := r.InboundDetectResult.LoadAndDelete(tag); ok {
		result := value.([]api.DetectResult)
		return &result, nil
	} else {
		return emptyDetectResult, nil
	}
}

func (r *RuleManager) Detect(tag string, destination string, email string) (reject bool) {
	reject = false
	var hitRuleID int = -1
	// If we have some rule for this inbound
	if value, ok := r.InboundRule.Load(tag); ok {
		ruleList := value.([]api.DetectRule)
		for _, r := range ruleList {
			if matchRule(r.Pattern, destination) {
				hitRuleID = r.ID
				reject = true
				break
			}
		}
		// If we hit some rule
		if reject && hitRuleID != -1 {
			l := strings.Split(email, "|")
			uid, err := strconv.Atoi(l[len(l)-1])
			if err != nil {
				newError(fmt.Sprintf("Record illegal behavior failed! Cannot find user's uid: %s", email)).AtDebug().WriteToLog()
				return reject
			}
			detectResult := []api.DetectResult{api.DetectResult{UID: uid, RuleID: hitRuleID}}
			// If there are any hit history
			if v, ok := r.InboundDetectResult.LoadOrStore(tag, detectResult); ok {
				resultList := v.([]api.DetectResult)
				resultList = append(resultList, detectResult...)
				r.InboundDetectResult.Store(tag, resultList)
			}
		}
	}
	return reject
}

func matchRule(rule string, destination string) (hit bool) {
	hit = false
	// Check Regex
	if regexp.MustCompile(rule).Match([]byte(destination)) {
		return true
	}
	return hit
}
