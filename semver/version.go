package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version ...
type Version struct {
	major int64
	minor int64
	patch int64
}

var versionRegex = regexp.MustCompile(`^([0-9]+)(\.[0-9]+)?(\.[0-9]+)?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

// LessThan ...
const LessThan = -1

// Equal ...
const Equal = 0

// GreaterThan ...
const GreaterThan = 1

// Compare compares this version to another one. It returns -1, 0, or 1 if
// the version smaller, equal, or larger than the other version.
// Compare compares the given version to the current version. It returns -1 if
// it's smaller, 0 if it's equal or 1 if it's greater
func (v *Version) Compare(version *Version) int {
	if d := compareSegment(v.major, version.major); d != 0 {
		return d
	}
	if d := compareSegment(v.minor, version.minor); d != 0 {
		return d
	}
	if d := compareSegment(v.patch, version.patch); d != 0 {
		return d
	}

	// at this point, it's equal which we return 0
	return 0
}

func NewVersion(s string) (*Version, error) {
	m := versionRegex.FindStringSubmatch(s)
	if m == nil {
		return nil, errors.New("Invalid SemVer")
	}

	v := Version{}
	i, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing version segment: %s", err)
	}
	v.major = i

	v.minor = 0
	if m[2] != "" {
		i, err = strconv.ParseInt(strings.TrimPrefix(m[2], "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version segment: %s", err)
		}
		v.minor = i
	}

	v.patch = 0
	if m[3] != "" {
		i, err = strconv.ParseInt(strings.TrimPrefix(m[3], "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version segment: %s", err)
		}
		v.patch = i
	}

	return &v, nil
}

func compareSegment(v, o int64) int {
	if v < o {
		return -1
	}
	if v > o {
		return 1
	}

	return 0
}
