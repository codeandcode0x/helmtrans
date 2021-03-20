package src

// yaml to helm 
func YamltoHelm(filePath, outPutPath string) {
	defer Catch()
	ScanResFile(filePath, outPutPath, "")
}


// helm to yaml
func HelmtoYaml(yamlStr string) {
	defer Catch()
}


// helm check
func HelmCheck(filePath string) {
	defer Catch()
}














