package helpers

import (
	"fmt"
	"strconv"
	"strings"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

func InvalidVersionFormatMessage(field string) *core.IValidMessage {
	return &core.IValidMessage{
		Name:    field,
		Code:    "INVALID_VERSION_FORMAT",
		Message: fmt.Sprintf("The field %s must be in format n.n.n where n is integer", field),
	}
}

func IsValidVersionUpdate(updatingVersion *string, currentVersion string, field string) (bool, *core.IValidMessage) {
	split := strings.Split(currentVersion, ".")
	major := split[0]
	minor := split[1]
	patch := split[2]
	majorInt, _ := strconv.Atoi(major)
	minorInt, _ := strconv.Atoi(minor)
	patchInt, _ := strconv.Atoi(patch)

	allowMajor := fmt.Sprintf("%s.0.0", strconv.Itoa(majorInt+1))
	allowMinor := fmt.Sprintf("%s.%s.0", major, strconv.Itoa(minorInt+1))
	allowPatch := fmt.Sprintf("%s.%s.%s", major, minor, strconv.Itoa(patchInt+1))

	if utils.GetString(updatingVersion) != allowMajor && *updatingVersion != allowMinor && *updatingVersion != allowPatch {
		return false, &core.IValidMessage{
			Name:    field,
			Code:    "VERSION_DISALLOWED",
			Message: fmt.Sprintf("The current version (%s) on field %s can only be updated to %s, %s or %s", currentVersion, field, allowMajor, allowMinor, allowPatch),
			Data:    []string{allowMajor, allowMinor, allowPatch},
		}
	}

	return true, nil
}

func IsValidVersionFormat(version *string, field string) (bool, *core.IValidMessage) {
	if version == nil {
		return true, nil
	}

	split := strings.Split(*version, ".")

	if len(split) != 3 {
		return false, InvalidVersionFormatMessage(field)
	}

	for _, i := range split {
		_, err := strconv.Atoi(i)
		if err != nil {
			return false, InvalidVersionFormatMessage(field)
		}
	}

	return true, nil
}
