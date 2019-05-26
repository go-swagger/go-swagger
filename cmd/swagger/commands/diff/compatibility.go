package diff

import (
	"fmt"
)

func getCompatabilityForChange(forDiff SpecChangeCode, where DataDirection) Compatability {
	compatability := Breaking
	switch forDiff {
	case DeletedProperty:
		nonBreakingIf(where == Request, &compatability)
	case AddedProperty:
		nonBreakingIf(where == Response, &compatability)
	case AddedRequiredProperty:
	case ChangedOptionalToRequiredParam:
	case DeletedOptionalParam:
		compatability = NonBreaking
	case DeletedEndpoint:
	case AddedRequiredParam:
	case DeletedRequiredParam:
		compatability = NonBreaking
	case ChangedRequiredToOptional, AddedEndpoint:
		compatability = NonBreaking
	case WidenedType:
		nonBreakingIf(where == Request, &compatability)
	case NarrowedType:
		nonBreakingIf(where == Response, &compatability)
	case AddedEnumValue:
		nonBreakingIf(where == Request, &compatability)
	case DeletedEnumValue:
		nonBreakingIf(where == Response, &compatability)
	case AddedOptionalParam:
		compatability = NonBreaking
	case ChangedRequiredToOptionalParam:
		compatability = NonBreaking
	case AddedResponse:
		compatability = NonBreaking
	case DeletedResponse:
		compatability = Breaking
	case ChangedType:
		compatability = Breaking
	case AddedResponseHeader:
		compatability = NonBreaking
	case ChangedResponseHeader:
	case DeletedResponseHeader:
	default:
		fmt.Printf("ERROR: Unknown diff type")
	}
	return compatability
}

func nonBreakingIf(cond bool, compatability *Compatability) {
	if cond {
		*compatability = NonBreaking
	}
}
