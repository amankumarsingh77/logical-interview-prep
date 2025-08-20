package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Metadata   string
}

func main() {
	v1, err := ParseVersion("1.2.3-alpha+build")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Major %d, Minor %d, Patch %d, PreRelease %s, Metadata %s\n", v1.Major, v1.Minor, v1.Patch, v1.PreRelease, v1.Metadata)
	v2, _ := ParseVersion("1.0.0-rc.2")
	log.Println(Compare(v1, v2))
}

func ParseVersion(versionString string) (*Version, error) {
	mainAndMeta := strings.SplitN(versionString, "+", 2)
	mainPart := mainAndMeta[0]
	metadata := ""
	if len(mainAndMeta) > 1 {
		metadata = mainAndMeta[1]
		if metadata == "" {
			return nil, errors.New("build metadata cannot be empty")
		}
	}
	coreAndPre := strings.SplitN(mainPart, "-", 2)
	corePart := coreAndPre[0]
	preRelease := ""
	if len(coreAndPre) > 1 {
		preRelease = coreAndPre[1]
		if preRelease == "" {
			return nil, errors.New("pre-release identifier cannot be empty")
		}
	}
	
	coreParts := strings.Split(corePart, ".")
	if len(coreParts) != 3 {
		return nil, fmt.Errorf("invalid core version format: %s", corePart)
	}

	major, err := strconv.Atoi(coreParts[0])
	if err != nil {
		return nil, fmt.Errorf("major version is not a number: %s", coreParts[0])
	}
	minor, err := strconv.Atoi(coreParts[1])
	if err != nil {
		return nil, fmt.Errorf("minor version is not a number: %s", coreParts[1])
	}
	patch, err := strconv.Atoi(coreParts[2])
	if err != nil {
		return nil, fmt.Errorf("patch version is not a number: %s", coreParts[2])
	}

	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: preRelease,
		Metadata:   metadata,
	}, nil
}

func Compare(v1 *Version, v2 *Version) int {
	if v1.Major != v2.Major {
		if v1.Major > v2.Major {
			return 1
		}
		return -1
	}
	if v1.Minor != v2.Minor {
		if v1.Minor > v2.Minor {
			return 1
		}
		return -1
	}
	if v1.Patch != v2.Patch {
		if v1.Patch > v2.Patch {
			return 1
		}
		return -1
	}
	v1HasPre := v1.PreRelease != ""
	v2HasPre := v2.PreRelease != ""
	if !v1HasPre && v2HasPre {
		return 1
	}
	if v1HasPre && !v2HasPre {
		return -1
	}
	if !v1HasPre && !v2HasPre {
		return 0
	}
	p1Segs := strings.Split(v1.PreRelease, ".")
	p2Segs := strings.Split(v2.PreRelease, ".")
	minLen := min(len(p1Segs), len(p2Segs))
	for i := 0; i < minLen; i++ {
		seg1, seg2 := p1Segs[i], p2Segs[i]
		num1, err1 := strconv.Atoi(seg1)
		num2, err2 := strconv.Atoi(seg2)

		if err1 == nil && err2 == nil {
			if num1 != num2 {
				if num1 > num2 {
					return 1
				}
				return -1
			}
		} else {
			if seg1 != seg2 {
				return strings.Compare(seg1, seg2)
			}
		}
	}

	if len(p1Segs) > len(p2Segs) {
		return 1
	}
	if len(p1Segs) < len(p2Segs) {
		return -1
	}

	return 0
}
