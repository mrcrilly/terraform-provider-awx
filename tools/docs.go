package tools

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrcrilly/terraform-provider-awx/awx"
)

const (
	cloudMark      = "awx"
	cloudTitle     = "AWX"
	cloudMarkShort = "awx"
	docRoot        = "../docs"
)

func GenerateProviderCoumentation() error {
	log.Printf("some test")

	provider := awx.Provider()
	vProvider := runtime.FuncForPC(reflect.ValueOf(awx.Provider).Pointer())
	// document for DataSources

	fname, _ := vProvider.FileLine(0)
	fpath := filepath.Dir(fname)
	log.Printf("generating doc from: %s\n", fpath)

	for k, v := range provider.DataSourcesMap {
		genDoc("data_source", fpath, k, v)
	}

	// document for Resources
	for k, v := range provider.ResourcesMap {
		genDoc("resource", fpath, k, v)
	}

	return nil
}

// genDoc generating doc for resource
func genDoc(dtype, fpath, name string, resource *schema.Resource) {
	log.Printf("some test")
	data := map[string]string{
		"name":              name,
		"dtype":             strings.Replace(dtype, "_", "", -1),
		"resource":          name[len(cloudMark)+1:],
		"cloud_mark":        cloudMark,
		"cloud_title":       cloudTitle,
		"example":           "",
		"description":       "",
		"description_short": "",
		"import":            "",
	}

	fname := fmt.Sprintf("%s_%s.go", dtype, data["resource"])
	log.Printf("[START]get description from file: %s\n", fname)

	description, err := getFileDescription(fmt.Sprintf("%s/%s", fpath, fname))
	if err != nil {
		log.Printf("[SKIP!]get description failed, skip: %s", err)
		return
	}

	description = strings.TrimSpace(description)
	if description == "" {
		log.Printf("[SKIP!]description empty, skip: %s\n", fname)
		return
	}

	importPos := strings.Index(description, "\nImport\n")
	if importPos != -1 {
		data["import"] = strings.TrimSpace(description[importPos+8:])
		description = strings.TrimSpace(description[:importPos])
	}

	pos := strings.Index(description, "\nExample Usage\n")
	if pos != -1 {
		data["example"] = strings.TrimSpace(description[pos+15:])
		description = strings.TrimSpace(description[:pos])
	} else {
		log.Printf("[SKIP!]example usage missing, skip: %s\n", fname)
		return
	}

	data["description"] = description
	pos = strings.Index(description, "\n\n")
	if pos != -1 {
		data["description_short"] = strings.TrimSpace(description[:pos])
	} else {
		data["description_short"] = description
	}

	requiredArgs := []string{}
	optionalArgs := []string{}
	attributes := []string{}
	subStruct := []string{}

	var keys []string
	for k := range resource.Schema {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := resource.Schema[k]
		//if v.Description == "" {
		//	continue
		//}
		if v.Required {
			opt := "Required"
			if v.ForceNew {
				opt += ", ForceNew"
			}
			requiredArgs = append(requiredArgs, fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description))
			subStruct = append(subStruct, getSubStruct(0, k, v)...)
		} else if v.Optional {
			opt := "Optional"
			if v.ForceNew {
				opt += ", ForceNew"
			}
			optionalArgs = append(optionalArgs, fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description))
			subStruct = append(subStruct, getSubStruct(0, k, v)...)
		} else {
			attrs := getAttributes(0, k, v)
			if len(attrs) > 0 {
				attributes = append(attributes, attrs...)
			}
		}
	}

	sort.Strings(requiredArgs)
	sort.Strings(optionalArgs)
	sort.Strings(attributes)

	requiredArgs = append(requiredArgs, optionalArgs...)
	data["arguments"] = strings.Join(requiredArgs, "\n")
	if len(subStruct) > 0 {
		data["arguments"] += "\n" + strings.Join(subStruct, "\n")
	}
	data["attributes"] = strings.Join(attributes, "\n")
	dtypePath := ""
	if dtype == "data_source" {
		dtypePath = "data-sources"
	} else if dtype == "resource" {
		dtypePath = "resources"
	}
	fname = fmt.Sprintf("%s/%s/%s.md", docRoot, dtypePath, data["resource"])
	fd, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("[FAIL!]open file %s failed: %s", fname, err)
		return
	}

	defer func() {
		if e := fd.Close(); e != nil {
			log.Printf("[FAIL!]close file %s failed: %s", fname, e)
		}
	}()
	t := template.Must(template.New("t").Parse(docTPL))
	err = t.Execute(fd, data)
	if err != nil {
		log.Printf("[FAIL!]write file %s failed: %s", fname, err)
		return
	}

	log.Printf("[SUCC.]write doc to file success: %s", fname)
}

// getAttributes get attributes from schema
func getAttributes(step int, k string, v *schema.Schema) []string {
	attributes := []string{}
	ident := strings.Repeat(" ", step*2)

	//if v.Description == "" {
	//	return attributes
	//}

	if v.Computed {
		if _, ok := v.Elem.(*schema.Resource); ok {
			listAttributes := []string{}
			for kk, vv := range v.Elem.(*schema.Resource).Schema {
				attrs := getAttributes(step+1, kk, vv)
				if len(attrs) > 0 {
					listAttributes = append(listAttributes, attrs...)
				}
			}
			slistAttributes := ""
			sort.Strings(listAttributes)
			if len(listAttributes) > 0 {
				slistAttributes = "\n" + strings.Join(listAttributes, "\n")
			}
			attributes = append(attributes, fmt.Sprintf("%s* `%s` - %s%s", ident, k, v.Description, slistAttributes))
		} else {
			attributes = append(attributes, fmt.Sprintf("%s* `%s` - %s", ident, k, v.Description))
		}
	}

	return attributes
}

// getFileDescription get description from go file
func getFileDescription(fname string) (string, error) {
	fset := token.NewFileSet()

	parsedAst, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	return parsedAst.Doc.Text(), nil
}

// getSubStruct get sub structure from go file
func getSubStruct(step int, k string, v *schema.Schema) []string {
	subStructs := []string{}

	if v.Description == "" {
		return subStructs
	}

	if v.Type == schema.TypeMap || v.Type == schema.TypeList || v.Type == schema.TypeSet {
		if _, ok := v.Elem.(*schema.Resource); ok {
			subStructs = append(subStructs, fmt.Sprintf("\nThe `%s` object supports the following:\n", k))
			requiredArgs := []string{}
			optionalArgs := []string{}
			attributes := []string{}

			var keys []string
			for kk := range v.Elem.(*schema.Resource).Schema {
				keys = append(keys, kk)
			}
			sort.Strings(keys)
			for _, kk := range keys {
				vv := v.Elem.(*schema.Resource).Schema[kk]
				if vv.Description == "" {
					vv.Description = "************************* Please input Description for Schema ************************* "
				}
				if vv.Required {
					opt := "Required"
					if vv.ForceNew {
						opt += ", ForceNew"
					}
					requiredArgs = append(requiredArgs, fmt.Sprintf("* `%s` - (%s) %s", kk, opt, vv.Description))
				} else if vv.Optional {
					opt := "Optional"
					if vv.ForceNew {
						opt += ", ForceNew"
					}
					optionalArgs = append(optionalArgs, fmt.Sprintf("* `%s` - (%s) %s", kk, opt, vv.Description))
				} else {
					attrs := getAttributes(0, kk, vv)
					if len(attrs) > 0 {
						attributes = append(attributes, attrs...)
					}
				}
			}
			sort.Strings(requiredArgs)
			subStructs = append(subStructs, requiredArgs...)
			sort.Strings(optionalArgs)
			subStructs = append(subStructs, optionalArgs...)
			sort.Strings(attributes)
			subStructs = append(subStructs, attributes...)

			for _, kk := range keys {
				vv := v.Elem.(*schema.Resource).Schema[kk]
				subStructs = append(subStructs, getSubStruct(step+1, kk, vv)...)
			}
		}
	}
	return subStructs
}
