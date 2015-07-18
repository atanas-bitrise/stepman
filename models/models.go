package models

import log "github.com/Sirupsen/logrus"

// -------------------
// --- Models

// InputModel ...
type InputModel struct {
	EnvKey            *string   `json:"env_key,omitempty" yaml:"env_key,omitempty"`
	Title             *string   `json:"title,omitempty" yaml:"title,omitempty"`
	Description       *string   `json:"description,omitempty" yaml:"description,omitempty"`
	Value             *string   `json:"value,omitempty" yaml:"value,omitempty"`
	ValueOptions      *[]string `json:"value_options,omitempty" yaml:"value_options,omitempty"`
	IsRequired        *bool     `json:"is_required,omitempty" yaml:"is_required,omitempty"`
	IsExpand          *bool     `json:"is_expand,omitempty" yaml:"is_expand,omitempty"`
	IsDontChangeValue *bool     `json:"is_dont_change_value,omitempty" yaml:"is_dont_change_value,omitempty"`
}

// OutputModel ...
type OutputModel struct {
	EnvKey      *string `json:"env_key,omitempty" yaml:"env_key,omitempty"`
	Title       *string `json:"title,omitempty" yaml:"title,omitempty"`
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
}

// StepModel ...
type StepModel struct {
	ID                  string            `json:"id"`
	SteplibSource       string            `json:"steplib_source"`
	VersionTag          string            `json:"version_tag"`
	Name                string            `json:"name" yaml:"name"`
	Description         *string           `json:"description,omitempty" yaml:"description,omitempty"`
	Website             string            `json:"website" yaml:"website"`
	ForkURL             *string           `json:"fork_url,omitempty" yaml:"fork_url,omitempty"`
	Source              map[string]string `json:"source" yaml:"source"`
	HostOsTags          *[]string         `json:"host_os_tags,omitempty" yaml:"host_os_tags,omitempty"`
	ProjectTypeTags     *[]string         `json:"project_type_tags,omitempty" yaml:"project_type_tags,omitempty"`
	TypeTags            *[]string         `json:"type_tags,omitempty" yaml:"type_tags,omitempty"`
	IsRequiresAdminUser *bool             `json:"is_requires_admin_user,omitempty" yaml:"is_requires_admin_user,omitempty"`
	Inputs              []*InputModel     `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	Outputs             []*OutputModel    `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

// StepGroupModel ...
type StepGroupModel struct {
	ID       string      `json:"id"`
	Versions []StepModel `json:"versions"`
	Latest   StepModel   `json:"latest"`
}

// StepHash ...
type StepHash map[string]StepGroupModel

// StepCollectionModel ...
type StepCollectionModel struct {
	FormatVersion        string              `json:"format_version" yaml:"format_version"`
	GeneratedAtTimeStamp int64               `json:"generated_at_timestamp" yaml:"generated_at_timestamp"`
	Steps                StepHash            `json:"steps" yaml:"steps"`
	SteplibSource        string              `json:"steplib_source" yaml:"steplib_source"`
	DownloadLocations    []map[string]string `json:"download_locations" yaml:"download_locations"`
}

// WorkFlowModel ...
type WorkFlowModel struct {
	FormatVersion string      `json:"format_version"`
	Environments  []string    `json:"environments"`
	Steps         []StepModel `json:"steps"`
}

// -------------------
// --- Struct methods

// GetStep ...
func (collection StepCollectionModel) GetStep(id, version string) (bool, StepModel) {
	versions := collection.Steps[id].Versions
	if len(versions) > 0 {
		for _, step := range versions {
			if step.VersionTag == version {
				return true, step
			}
		}
	}
	return false, StepModel{}
}

// GetDownloadLocations ...
func (collection StepCollectionModel) GetDownloadLocations(step StepModel) []map[string]string {
	locations := []map[string]string{}
	for _, downloadLocation := range collection.DownloadLocations {
		for key, value := range downloadLocation {
			switch key {
			case "zip":
				url := value + step.ID + "/" + step.VersionTag + "/step.zip"
				locations = append(locations, map[string]string{key: url})
			case "git":
				locations = append(locations, step.Source)
			default:
				log.Error("[STEPMAN] - Invalid download location")
			}
		}
	}
	return locations
}
