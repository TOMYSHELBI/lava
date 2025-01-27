package rpcconsumer

import (
	"fmt"
	"strconv"

	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	"github.com/lavanet/lava/protocol/common"
	"github.com/lavanet/lava/utils"
)

type RelayErrors struct {
	relayErrors       []RelayError
	onFailureMergeAll bool
}

// checking the errors that appeared the most and returning the number of errors that were the same and the index of one of them
func (r *RelayErrors) findMaxAppearances(input map[string][]int) (maxVal int, indexToReturn int) {
	var maxValIndexArray []int // one of the indexes
	for _, val := range input {
		if len(val) > maxVal {
			maxVal = len(val)
			maxValIndexArray = val
		}
	}
	if len(maxValIndexArray) > 0 {
		indexToReturn = maxValIndexArray[0]
	} else {
		indexToReturn = -1
	}
	return
}

func (r *RelayErrors) GetBestErrorMessageForUser() RelayError {
	bestIndex := -1
	bestResult := github_com_cosmos_cosmos_sdk_types.ZeroDec()
	errorMap := make(map[string][]int)
	for idx, relayError := range r.relayErrors {
		errorMessage := relayError.err.Error()
		errorMap[errorMessage] = append(errorMap[errorMessage], idx)
		if relayError.ProviderInfo.ProviderQoSExcellenceSummery.IsNil() || relayError.ProviderInfo.ProviderStake.Amount.IsNil() {
			continue
		}
		currentResult := relayError.ProviderInfo.ProviderQoSExcellenceSummery.MulInt(relayError.ProviderInfo.ProviderStake.Amount)
		if currentResult.GTE(bestResult) { // 0 or 1 here are valid replacements, so even 0 scores will return the error value
			bestResult.Set(currentResult)
			bestIndex = idx
		}
	}

	errorCount, index := r.findMaxAppearances(errorMap)
	if index >= 0 && errorCount >= (len(r.relayErrors)/2) {
		// we have majority of errors we can return this error.
		return r.relayErrors[index]
	}

	if bestIndex != -1 {
		// Return the chosen error.
		// Print info for the consumer to know which errors happened
		utils.LavaFormatInfo("Failed all relays", utils.LogAttr("error_map", errorMap))
		return r.relayErrors[bestIndex]
	}
	// if we didn't manage to find any index return all.
	utils.LavaFormatError("Failed finding the best error index in GetErrorMessageForUser", nil, utils.LogAttr("relayErrors", r.relayErrors))
	if r.onFailureMergeAll {
		return RelayError{err: r.mergeAllErrors()}
	}
	// otherwise return the first element of the RelayErrors
	return r.relayErrors[0]
}

func (r *RelayErrors) getAllUniqueErrors() []error {
	allErrors := make([]error, len(r.relayErrors))
	repeatingErrors := make(map[string]struct{})
	for idx, relayError := range r.relayErrors {
		errString := relayError.err.Error() // using strings to filter repeating errors
		_, ok := repeatingErrors[errString]
		if ok {
			continue
		}
		repeatingErrors[errString] = struct{}{}
		allErrors[idx] = relayError.err
	}
	return allErrors
}

func (r *RelayErrors) mergeAllErrors() error {
	mergedMessage := ""
	allErrors := r.getAllUniqueErrors()
	allErrorsLength := len(allErrors)
	for idx, message := range allErrors {
		mergedMessage += strconv.Itoa(idx) + ". " + message.Error()
		if idx < allErrorsLength {
			mergedMessage += ", "
		}
	}
	return fmt.Errorf(mergedMessage)
}

type RelayError struct {
	err          error
	ProviderInfo common.ProviderInfo
	response     *relayResponse
}
