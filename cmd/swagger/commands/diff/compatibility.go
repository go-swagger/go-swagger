package diff

import (
	"fmt"
)

func getCompatibilityForChange(forDiff SpecChangeCode, where DataDirection) Compatibility {
	compatibility := Breaking
	switch forDiff {
	case DeletedProperty:
		nonBreakingIf(where == Request, &compatibility)
	case AddedProperty:
		nonBreakingIf(where == Response, &compatibility)
	case AddedRequiredProperty:
	case ChangedOptionalToRequiredParam:
	case DeletedOptionalParam:
		compatibility = NonBreaking
	case DeletedDeprecatedEndpoint:
		compatibility = NonBreaking
	case DeletedEndpoint:
	case AddedRequiredParam:
	case DeletedRequiredParam:
		compatibility = NonBreaking
	case ChangedRequiredToOptional, AddedEndpoint:
		compatibility = NonBreaking
	case WidenedType:
		nonBreakingIf(where == Request, &compatibility)
	case NarrowedType:
		nonBreakingIf(where == Response, &compatibility)
	case AddedEnumValue:
		nonBreakingIf(where == Request, &compatibility)
	case DeletedEnumValue:
		nonBreakingIf(where == Response, &compatibility)
	case AddedOptionalParam:
		compatibility = NonBreaking
	case ChangedRequiredToOptionalParam:
		compatibility = NonBreaking
	case AddedResponse:
		compatibility = NonBreaking
	case DeletedResponse:
		compatibility = Breaking
	case ChangedType:
		compatibility = Breaking
	case AddedResponseHeader:
		compatibility = NonBreaking
	case ChangedResponseHeader:
	case DeletedResponseHeader:
	default:
		fmt.Printf("ERROR: Unknown diff type")
	}
	return compatibility
}

func nonBreakingIf(cond bool, compatibility *Compatibility) {
	if cond {
		*compatibility = NonBreaking
	}
}
