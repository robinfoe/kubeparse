package task

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/allegro/bigcache"
	"github.com/robinfoe/kuberparse/internal/app"
	"gopkg.in/yaml.v3"
)

type Parser struct {
	Opts  *ParseOptions
	Cache *bigcache.BigCache
}

func (p *Parser) Parse() error {

	log.Println("start parsing...")

	apps, _ := p.parseDeploymentConfig()
	cms, _ := p.parseConfigmap()

	keyRef := make(map[string]int)

	// sorting map keys
	keys := make([]string, 0, len(apps))
	for k := range apps {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	appList := make([]*app.AppConfig, len(keys))

	// appending services call for each application
	for i, k := range keys {
		a := apps[k]
		keyRef[k] = i
		keyRef[strings.Split(k, "|")[1]] = i

		for _, cn := range a.ConfigMaps {
			fmt.Printf("%s|%s\n", a.Namespace, cn)
			cm := cms[(fmt.Sprintf("%s|%s", a.Namespace, cn))]
			if cm != nil {
				a.ServiceCall = append(a.ServiceCall, cm.ServiceCall...)
				a.DBProp = append(a.DBProp, cm.DBProperties)
			}

		}
		appList[i] = a
	}

	// call coordinate
	cc := make([][]int, len(keys))

	for i, a := range appList {
		cc[i] = make([]int, len(keys))
		fmt.Printf("%s -- %s\n", a.Namespace, a.Name)
		for _, s := range a.ServiceCall {

			fmt.Println(s.Url)
			if s.IsValidService() {
				var validPtr = -1

				if val, ok := keyRef[s.GetKey()]; ok {
					validPtr = val
					// cc[i][val] = cc[i][val] + 1
					// append service name details ?

				} else {
					if val, ok := keyRef[s.Name]; ok {
						validPtr = val
						// cc[i][val] = cc[i][val] + 1
					}
				}

				if validPtr >= 0 {
					cc[i][validPtr] = cc[i][validPtr] + 1
					pa := appList[validPtr]
					if pa.RefferedBy == nil {
						pa.RefferedBy = make(map[string]*app.AppConfig)
					}

					if _, ok := pa.RefferedBy[a.GetKey()]; !ok {
						pa.RefferedBy[a.GetKey()] = a
					}

				}

			}
		}
	}

	//	generate CSV
	//  start generate coordinate
	csvRaw := make([][]string, len(keys)+3) // 3 row for empty header
	csvRaw[0] = make([]string, len(keys)+2)
	csvRaw[1] = make([]string, len(keys)+2)
	csvRaw[2] = make([]string, len(keys)+2)

	for i, a := range appList {
		csvRaw[i+3] = make([]string, len(keys)+2)

		// Row HEADER
		csvRaw[0][i+2] = a.Namespace
		csvRaw[1][i+2] = a.Name

		// Column Header
		csvRaw[i+3][0] = a.Namespace
		csvRaw[i+3][1] = a.Name

		cv := cc[i]
		for ci, val := range cv {
			csvRaw[i+3][ci+2] = fmt.Sprintf("%v", val)
		}

	}

	// write svc call detail
	file, err := os.Create("svc-call-detail.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range csvRaw {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}

	filens, err := os.Create("svc-call-namespace.txt")
	checkError("Cannot create file", err)
	defer file.Close()

	// writing as namespace
	ns := ""
	nsMap := map[string]string{}

	for _, a := range appList {
		if ns != a.Namespace {
			ns = a.Namespace
			_, err := filens.WriteString(fmt.Sprintf("\n\n\npackage [%s]\n\n", ns))
			checkError("Cannot write to file", err)
		}

		for _, s := range a.ServiceCall {
			if s.IsValidService() {
				kloc := -1
				if val, ok := keyRef[s.GetKey()]; ok {
					kloc = val
				} else {
					if val, ok := keyRef[s.Name]; ok {
						kloc = val
					}
				}

				if kloc > -1 {
					svc := appList[kloc]
					nsMapKey := fmt.Sprintf("%s|%s", ns, svc.Namespace)
					if _, ok := nsMap[nsMapKey]; !ok {
						// not found
						nsMap[nsMapKey] = nsMapKey
						_, err := filens.WriteString(fmt.Sprintf("[%s] --> [%s]\n", ns, svc.Namespace))
						checkError("Cannot write to file", err)
					}

				}
			}
		}
	}

	filep, err := os.Create("parent-child.txt")
	checkError("Cannot create file", err)
	defer filep.Close()
	// write parent child dependency

	ns = ""
	for _, k := range keys {
		a := apps[k]

		if ns != a.Namespace {
			ns = a.Namespace
			_, err := filep.WriteString(fmt.Sprintf("Namespace ::  [%s]\n", ns))
			checkError("Cannot write to file", err)
		}

		_, era := filep.WriteString(fmt.Sprintf("\n\tApp : [%s] Replica - [%s] ====================\n", a.Name, a.Replicas))
		checkError("Cannot write to file", era)

		_, err := filep.WriteString(fmt.Sprintf("\t\tDB Properties ::  \n"))
		checkError("Cannot write to file", err)

		for _, d := range a.DBProp {

			for k, v := range d.(map[string]string) {
				_, erk := filep.WriteString(fmt.Sprintf("\t\t%s : %s\n", k, v))
				checkError("Cannot write to file", erk)
			}

		}

		for _, r := range a.RefferedBy {
			_, err := filep.WriteString(fmt.Sprintf("\t\t R : [%s] - [%s]\n", r.Namespace, r.Name))
			checkError("Cannot write to file", err)
		}

		_, err = filep.WriteString("\n\n")
		checkError("Cannot write to file", err)

	}

	return nil
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func (p *Parser) parseConfigmap() (map[string]*app.CmConfig, error) {

	log.Printf("Configmap Location : %s", p.Opts.Path.ConfigMap)

	cms := make(map[string]*app.CmConfig)

	yamlFile, err := ioutil.ReadFile(p.Opts.Path.ConfigMap)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	cmList := &ConfigMapList{}
	err = yaml.Unmarshal(yamlFile, cmList)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, i := range cmList.Items {

		cm := &app.CmConfig{
			Namespace:    i.Metadata.Namespace,
			Name:         i.Metadata.Name,
			RawData:      i.Data,
			DBProperties: make(map[string]string),
		}

		cm.GrabLinks()

		// store into map
		cms[cm.GetKey()] = cm
	}

	return cms, nil
}

func (p *Parser) parseDeploymentConfig() (map[string]*app.AppConfig, error) {

	log.Printf("Deployment Config Location : %s", p.Opts.Path.DeploymentConfig)

	apps := make(map[string]*app.AppConfig)

	yamlFile, err := ioutil.ReadFile(p.Opts.Path.DeploymentConfig)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	dcList := &DeploymentConfigList{}

	err = yaml.Unmarshal(yamlFile, dcList)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	// log.Printf("Kind %s", dcList.APIVersion)
	for _, i := range dcList.Items {
		// log.Printf("namespace: %s , name : %s", i.Metadata.Namespace, i.Metadata.Name)

		app := &app.AppConfig{
			Namespace: i.Metadata.Namespace,
			Name:      i.Metadata.Name,
			Replicas:  i.Spec.Replicas,
		}

		// for each containers
		for _, c := range i.Spec.Template.Spec.Containers {
			// env
			for _, e := range c.Env {
				if !(e.ValueFrom.ConfigMapKeyRef.Key == "") {
					app.ConfigMaps = append(app.ConfigMaps, e.ValueFrom.ConfigMapKeyRef.Name)
				}
			}

			// env from
			for _, e := range c.EnvFrom {
				app.ConfigMaps = append(app.ConfigMaps, e.ConfigMapRef.Name)
			}
		}

		//for each volumes
		for _, v := range i.Spec.Template.Spec.Volumes {
			if !(v.ConfigMap.Name == "") {
				app.ConfigMaps = append(app.ConfigMaps, v.ConfigMap.Name)
			}
		}

		// store into map
		apps[app.GetKey()] = app
	}

	return apps, nil

}
