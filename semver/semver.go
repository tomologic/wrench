package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Semver struct {
	Major    int
	Minor    int
	Patch    int
	Snapshot string
}

type SemverList []Semver

func Parse(s string) (Semver, error) {
	sv := Semver{}

	// Remove leading v
	s = strings.TrimLeft(s, "v")

	// Split semver into maximum 3 parts
	parts := strings.SplitN(s, ".", 3)

	// 3 parts should be found
	if len(parts) != 3 {
		return sv, errors.New("Semver requires 3 parts Major.Minor.Patch[-snapshot]")
	}

	// Make sure no part empty
	for i, p := range parts {
		if strings.TrimSpace(p) == "" {
			return sv, errors.New(fmt.Sprintf("Part %d empty", i+1))
		}
	}

	var err error
	if sv.Major, err = strconv.Atoi(parts[0]); err != nil {
		return sv, errors.New(fmt.Sprintf("Unable to convert major '%s' to int", parts[0]))
	}

	if sv.Minor, err = strconv.Atoi(parts[1]); err != nil {
		return sv, errors.New(fmt.Sprintf("Unable to convert minor '%s' to int", parts[1]))
	}

	// Split Patch and Snapshot
	parts = strings.SplitN(parts[2], "-", 2)

	if sv.Patch, err = strconv.Atoi(parts[0]); err != nil {
		return sv, errors.New(fmt.Sprintf("Unable to convert patch '%s' to int", parts[0]))
	}

	if len(parts) == 2 {
		sv.Snapshot = strings.TrimSpace(parts[1])
	}

	return sv, nil
}

func (s Semver) Equals(v Semver) bool {
	return ((s.Major == v.Major) &&
		(s.Minor == v.Minor) &&
		(s.Patch == v.Patch) &&
		(s.Snapshot == v.Snapshot))
}

func (s *Semver) Bump(level string) error {
	level = strings.ToLower(level)
	if level == "major" {
		s.Major += 1
		s.Minor = 0
		s.Patch = 0
	} else if level == "minor" {
		s.Minor += 1
		s.Patch = 0
	} else if level == "patch" {
		s.Patch += 1
	} else {
		return errors.New(fmt.Sprintf("Unknown level '%s'", level))
	}
	s.Snapshot = ""
	return nil
}

func (s Semver) String() string {
	if s.Snapshot == "" {
		return fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
	} else {
		return fmt.Sprintf("v%d.%d.%d-%s", s.Major, s.Minor, s.Patch, s.Snapshot)
	}
}

func (s Semver) IsReleaseVersion() bool {
	if s.Snapshot == "" {
		return true
	}
	return false
}

func (s SemverList) Len() int {
	return len(s)
}

func (s SemverList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SemverList) Less(i, j int) bool {
	if s[i].Major < s[j].Major {
		return true
	} else if s[i].Major == s[j].Major {
		if s[i].Minor < s[j].Minor {
			return true
		} else if s[i].Minor == s[j].Minor {
			if s[i].Patch < s[j].Patch {
				return true
			}
		}
	}
	return false
}
