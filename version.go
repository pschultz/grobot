package grobot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var NoVersion = &Version{raw: "none"}

type Version struct {
	raw    string
	Major  int
	Minor  int
	Patch  int
	Branch string
}

func NewVersion(raw string) *Version {
	version := Version{raw: raw}
	version.parse()
	return &version
}

func (v *Version) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &v.raw)
	if err != nil {
		return err
	}

	return v.parse()
}

func (v *Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.raw)
}

func (v *Version) parse() (err error) {
	v.Major = 0
	v.Minor = 0
	v.Patch = 0

	if v.raw == NoVersion.raw {
		return nil
	}

	if strings.HasPrefix(v.raw, "branch:") {
		v.Branch = v.raw[7:]
		return nil
	}

	versionParts := strings.SplitN(v.raw, ".", 3)
	v.Major, err = strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("Could not parse major version from %s : %s", v.raw, err.Error())
	}

	if len(versionParts) > 1 {
		v.Minor, err = strconv.Atoi(versionParts[1])
		if err != nil {
			return fmt.Errorf("Could not parse minor version from %s : %s", v.raw, err.Error())
		}
	}

	if len(versionParts) > 2 {
		v.Patch, err = strconv.Atoi(versionParts[2])
		if err != nil {
			return fmt.Errorf("Could not parse patch version from %s : %s", v.raw, err.Error())
		}
	}

	return nil
}

func (v *Version) String() string {
	if v.Branch != "" || (v.Major == 0 && v.Minor == 0 && v.Patch == 0) {
		return v.raw
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) GreaterThen(other *Version) (bool, error) {
	if v.Branch != "" {
		if other.Branch == v.Branch {
			return false, nil
		}
		return false, fmt.Errorf("Can not compare two branch versions")
	}

	if v.Major == other.Major {
		if v.Minor == other.Minor {
			return v.Patch > other.Patch, nil
		}
		return v.Minor > other.Minor, nil
	}
	return v.Major > other.Major, nil
}

func (v *Version) LowerThen(other *Version) (bool, error) {
	return other.GreaterThen(v)
}
