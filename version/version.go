package version

import (
	"strconv"
)

type Version struct {
	Major, Minor, Patch int
}

type Semver interface {
	//Returns the string form of the version: eg. 1.2.3
	ToString() string
	// Sets the Major Version (int) from a String
	SetMajor(string) error
	// Sets the Minor Version (int) from a String
	SetMinor(string) error
	// Sets the Patch Version (int) from a String
	SetPatch(string) error
	// Gets the Major Version (int)
	GetMajor() int
	// Gets the Minor Version (int)
	GetMinor() int
	// Gets the Patch Version (int)
	GetPatch() int
	// Increments the Major Version
	IncrementMajor()
	// Increments the Minor Version
	IncrementMinor()
	// Increments the Patch Version
	IncrementPatch()
}

func (v *Version) SetMajor(s string) (err error) {
	v.Major, err = strconv.Atoi(s)
	if err != nil {
		return err
	}
	return nil
}

func (v *Version) SetMinor(s string) (err error) {
	v.Minor, err = strconv.Atoi(s)
	if err != nil {
		return err
	}
	return nil
}

func (v *Version) SetPatch(s string) (err error) {
	v.Patch, err = strconv.Atoi(s)
	if err != nil {
		return err
	}
	return nil
}

func (v *Version) GetMajor() (i int) {
	return v.Major
}

func (v *Version) GetMinor() (i int) {
	return v.Minor
}

func (v *Version) GetPatch() (i int) {
	return v.Patch
}

func (v *Version) IncrementMajor() {
	v.Major++
	v.Minor = 0
	v.Patch = 0
}

func (v *Version) IncrementMinor() {
	v.Minor++
	v.Patch = 0
}

func (v *Version) IncrementPatch() {
	v.Patch++
}

func (v *Version) ToString() (s string) {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
}
