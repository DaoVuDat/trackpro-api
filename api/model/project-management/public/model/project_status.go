//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type ProjectStatus string

const (
	ProjectStatus_Registering ProjectStatus = "registering"
	ProjectStatus_Progressing ProjectStatus = "progressing"
	ProjectStatus_Finished    ProjectStatus = "finished"
)

func (e *ProjectStatus) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "registering":
		*e = ProjectStatus_Registering
	case "progressing":
		*e = ProjectStatus_Progressing
	case "finished":
		*e = ProjectStatus_Finished
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for ProjectStatus enum")
	}

	return nil
}

func (e ProjectStatus) String() string {
	return string(e)
}