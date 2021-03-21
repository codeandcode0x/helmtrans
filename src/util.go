package src

import (
	"github.com/ghodss/yaml"
	"k8s.io/client-go/kubernetes"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
	"path"
	"log"
	"strings"
	"fmt"
	"os"
	"os/user"
	"helmtrans/src/k8s"
	"runtime"
	"bytes"
	"errors"
	"os/exec"
	"bufio"
	"io"
)

var clientset *kubernetes.Clientset
var apiClient k8s.K8sClientInf
var configPath string
var namespace string

//init
func init() {
	apiClient = &k8s.K8sClient{
		nil,
		"default",
	}
}

//get resource yaml
func GetResourceYaml(filePath string) []byte {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	return bytes
}

//getResourceType
func GetResourceType(bytes []byte) string{
	typeJson, err := yaml.YAMLToJSON(bytes)
	if err != nil {
		panic("get resource type error !")
	}
	var typeIn interface{}
	json.Unmarshal(typeJson, &typeIn)
	return typeIn.(map[string]interface{})["kind"].(string)
}

//scan data file
func ScanResFile(pathName, outPutPath ,subSourcePath string) {
	// mkdir output path
	if !FFExists(outPutPath) {
		mkErr := os.Mkdir(outPutPath, 0755)
	    if mkErr != nil {
	        log.Fatal(mkErr)
	    }
	}
	// set source path
	sourcePath:= outPutPath
	// read dir
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
    	log.Fatalln("scan data file - read file error!", err)
    }
    //resource delete
	for _, fi := range rd {
		// create dir
		if fi.IsDir() {
			err = os.Mkdir(sourcePath+"/"+fi.Name(), 0755)
		    if err != nil {
		        log.Fatal(err)
		    }
		    // copy templates
		    CopyDir("templates/chart/", sourcePath+"/"+fi.Name()+"/")
		    // subpath
		    subSourcePath = sourcePath+"/"+fi.Name()
		}

	    // scan files
	    if !fi.IsDir() {
	    	dataExt := path.Ext(fi.Name())
	    	if dataExt == ".yaml" {
				ScanFile(pathName, subSourcePath)
	    	}
		}else{
			ScanResFile(pathName+"/"+fi.Name(), outPutPath, subSourcePath)
		}	
	 }
}

// scan res files
func ScanFile(pathName, sourcePath string) {
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
    	log.Fatalln("scan data file - read file error!", err)
    }

	for _, fi := range rd {
	    // scan file
		if !fi.IsDir() {
			dataExt := path.Ext(fi.Name())
			switch(dataExt) {
			case ".yaml":
				bytes := GetResourceYaml(pathName+"/"+fi.Name())
				resourceType := GetResourceType(bytes)
				//write file
				switch(resourceType) {
				case "ConfigMap":
					spec := apiClient.UnmarshalConfigMap(bytes)
					t := readYaml(sourcePath+"/values.yaml")
					t["env"] = spec.Data
					WriteDataFile(JsonToYaml(t), sourcePath+"/values.yaml")
					break
				case "Deployment":
					spec := apiClient.UnmarshalDeployment(bytes)
					containerCount := len(spec.Spec.Template.Spec.Containers)
					// read values yaml
					t := readYaml(sourcePath+"/values.yaml")
					appName := spec.Name
					// just one container
					if containerCount == 1 {
						annotations := spec.Annotations
						env := spec.Spec.Template.Spec.Containers[0].Env
						image := spec.Spec.Template.Spec.Containers[0].Image
						imagePullPolicy := spec.Spec.Template.Spec.Containers[0].ImagePullPolicy
						ports := spec.Spec.Template.Spec.Containers[0].Ports
						readinessProbe := spec.Spec.Template.Spec.Containers[0].ReadinessProbe
						livenessProbe := spec.Spec.Template.Spec.Containers[0].LivenessProbe
						resources := spec.Spec.Template.Spec.Containers[0].Resources
						lifecycle := spec.Spec.Template.Spec.Containers[0].Lifecycle
						// image split
						imageSlice := strings.Split(image, ":")
						imageVersion := "master"
						if len(imageSlice) > 0 {
							imageVersion = imageSlice[1]
						}

						t["applicationName"] = appName
						t["deploymentAnnotations"] = annotations
						t["image"] = map[string]interface{}{"pullPolicy":imagePullPolicy, "repository":imageSlice[0], "tag":imageVersion}
						t["livenessProbe"] = livenessProbe
						t["readinessProbe"] = readinessProbe
						t["resources"] = resources
						t["containerPorts"] = ports
						envVars := make(map[string]string, 0)
						// merge env
						for key, val := range t["env"].(map[string]interface{}) {
							envVars[key] = val.(string)
						}

						for _, val := range env {
							envVars[val.Name] = val.Value
						}

						t["env"] = envVars

						if lifecycle != nil {
							t["shutdownDelay"] =1
							t["lifecycle"] = lifecycle
						}
					}else { // multi containers
						t["containers"] = spec.Spec.Template.Spec.Containers
						t["env"] = make(map[string]interface{}, 0)
					}

					// container count
					t["containerCount"] = containerCount
					// add image pull secrets
					imagePullSecrets := spec.Spec.Template.Spec.ImagePullSecrets
					secrets := []string{}
					for _, val := range imagePullSecrets {
						secrets = append(secrets, val.Name)
					}
					t["imagePullSecrets"] = secrets

					// write data to file
					WriteDataFile(JsonToYaml(t), sourcePath+"/values.yaml")

					// read chart yaml
					c := readYaml(sourcePath+"/Chart.yaml")
					c["name"] = appName
					// write data to file
					WriteDataFile(JsonToYaml(c), sourcePath+"/Chart.yaml")

					break
				case "Service":
					spec := apiClient.UnmarshalService(bytes)
					t := readYaml(sourcePath+"/values.yaml")
					t["serviceEnabled"]=true
					t["service"] = spec.Spec
					WriteDataFile(JsonToYaml(t), sourcePath+"/values.yaml")
					break
				}

				break
			}
		}else{
			ScanFile(pathName+"/"+fi.Name(), sourcePath)
		}
	 }
}

//map to string
func createKeyValuePairs(m map[string]string) string {
    b := new(bytes.Buffer)
    for key, value := range m {
        fmt.Fprintf(b, "%s: \"%s\"\n", key, value)
    }
    return b.String()
}

//josn to 
func JsonToYaml(t map[string]interface {}) string{
	jsonStr, _ := json.Marshal(t)
	yamlStr, err := yaml.JSONToYAML(jsonStr)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return string(yamlStr)
}

//read yaml
func readYaml(path string) map[string]interface{}{
	t := map[string]interface{}{}
    buffer, err := ioutil.ReadFile(path)
    err = yaml.Unmarshal(buffer, &t)
    if err != nil {
        log.Fatalf(err.Error())
    }
    return t
}

//read data file
func ReadDataFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
    if err != nil {
    	log.Fatalln("read "+path+" file failed!", err)
    }
    return data
}

//write data file
func WriteDataFile(yamlStr, filePath string) {
	//create file
	dataFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("create data file failed!", err)
	}
	//close file
	defer dataFile.Close()
	//create bufio write
	dataFileWriter := bufio.NewWriter(dataFile)
	_, errWrite := dataFileWriter.WriteString(string(yamlStr))
	if errWrite != nil {
		log.Fatalln("write data to file err!", errWrite)
	}
	dataFileWriter.Flush()
}

// copy dir
func CopyDir(srcPath string, destPath string) error {
	//检测目录正确性
	if srcInfo, err := os.Stat(srcPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("src 不是一个正确的目录！")
			fmt.Println(e.Error())
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("dest 不是一个正确的目录！")
			fmt.Println(e.Error())
			return e
		}
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			copyFile(path, destNewPath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
	return err
}

// copy files
func copyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	//分割path目录
	destSplitPathDirs := strings.Split(dest, "/")

	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, _ := pathExists(destSplitPath)
			if b == false {
				// fmt.Println("创建目录:" + destSplitPath)
				//创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

//检测文件夹路径时候存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//file exist
func FileExist(path string) bool {
  _, err := os.Lstat(path)
  return !os.IsNotExist(err)
}

//get Home
func GetHome() (string, error) {
    user, err := user.Current()
    if nil == err {
        return user.HomeDir, nil
    }

    // cross compile support
    if "windows" == runtime.GOOS {
        return homeWindows()
    }

    // Unix-like system, so just assume Unix
    return homeUnix()
}

func homeUnix() (string, error) {
    // First prefer the HOME environmental variable
    if home := os.Getenv("HOME"); home != "" {
        return home, nil
    }

    // If that fails, try the shell
    var stdout bytes.Buffer
    cmd := exec.Command("sh", "-c", "eval echo ~$USER")
    cmd.Stdout = &stdout
    if err := cmd.Run(); err != nil {
        return "", err
    }

    result := strings.TrimSpace(stdout.String())
    if result == "" {
        return "", errors.New("blank output when reading home directory")
    }

    return result, nil
}

func homeWindows() (string, error) {
    drive := os.Getenv("HOMEDRIVE")
    path := os.Getenv("HOMEPATH")
    home := drive + path
    if drive == "" || path == "" {
        home = os.Getenv("USERPROFILE")
    }
    if home == "" {
        return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
    }

    return home, nil
}

//file or folder check exists
func FFExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//catch
func  Catch()  {
    if r := recover(); r != nil {
        fmt.Println("ERROR:", r)
        var err error
        switch x := r.(type) {
        case string:
            err = errors.New(x)
        case error:
            err = x
        default:
            err = errors.New("")
        }
        if err != nil {
        }
    }
}

//check error
func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}


