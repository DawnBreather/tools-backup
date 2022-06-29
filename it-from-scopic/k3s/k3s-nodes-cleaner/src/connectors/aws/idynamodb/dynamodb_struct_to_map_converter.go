package idynamodb

import (
	"github.com/fatih/structs"
)

func convertStructToMapInterface(object interface{}) map[string]interface{}{
	return structs.Map(object)
}