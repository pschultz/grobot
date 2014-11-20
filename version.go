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
	if v.raw == NoVersion.raw {
		v.Major = 0
		v.Minor = 0
		return nil
	}

	if strings.HasPrefix(v.raw, "branch:") {
		v.Major = 0
		v.Minor = 0
		v.Branch = v.raw[7:]
		return nil
	}

	versionParts := strings.SplitN(v.raw, ".", 2)
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

	return nil
}

func (v *Version) String() string {
	return v.raw
}

func (v *Version) GreaterThen(other *Version) bool {
	if v.Major == other.Major {
		return v.Minor > other.Minor
	}
	return v.Major > other.Major
}

func (v *Version) LowerThen(other *Version) bool {
	return other.GreaterThen(v)
}
