package model

const noOpAction = "SELECT 1"

func ExecuteWithNoOpActions(
	resourceName string,
) *ExecuteModel {
	return Execute(resourceName, noOpAction, noOpAction)
}
