//gomin:pkg:default:file:regex:profiles.go:exclude
package profiles

import "golang.org/x/tools/cover"

type ProfilesByName struct {
	profiles map[string]*cover.Profile
}

func New(profiles []*cover.Profile) *ProfilesByName {
	profilesByName := &ProfilesByName{
		profiles: make(map[string]*cover.Profile),
	}

	for _, profile := range profiles {
		profilesByName.profiles[profile.FileName] = profile
	}

	return profilesByName
}

func (profiles *ProfilesByName) Get(fileName string) *cover.Profile {
	return profiles.profiles[fileName]
}
