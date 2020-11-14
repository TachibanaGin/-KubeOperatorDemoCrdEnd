package utils

func LabelsForFront(name string) map[string]string {
	return map[string]string{"app": name}
}

func LabelsForRs(name string) string {
	return "app="+ name
}
