package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/ukubenet/metadata-repository/config"
	global "github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func executeAppTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	path, _ := os.Getwd()
	tmplFile := path + "/" + config.Config.Metadata.Path + "/" + config.Config.Metadata.Tmplsubpath + "/" + tmplName + ".tmpl"
	executeTemplate(w, tmplName, []string{tmplFile}, data)
}

func (app *application) editApp(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")

	tmplData := struct {
		AppName string
	}{appName}

	executeAppTemplate(w, "app_edit", tmplData)
}

func (app *application) newApp(w http.ResponseWriter, r *http.Request) {
	executeAppTemplate(w, "app_new", nil)
}

func (app *application) listApps(w http.ResponseWriter, r *http.Request) {
	// todo: should be adapter to get app list from metadata repository
	list := []string{}
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path
	files, err := os.ReadDir(appPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			if file.Name() == ".git" || file.Name() == "metadata" {
				continue
			}
			list = append(list, file.Name())
		}
	}

	tmplData := struct {
		Apps []string
	}{
		list,
	}
	executeAppTemplate(w, "app_list", tmplData)
}

func (app *application) deleteApp(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")

	// todo: should be adapter to delete app
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path + "/" + appName
	err := os.RemoveAll(appPath)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (app *application) postApp(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")

	r.ParseForm()
	newAppName := r.FormValue("AppName")

	// todo: should be adapter to rename app
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path
	if newAppName != "" && appName != newAppName {
		err := os.Rename(appPath+"/"+appName, appPath+"/"+newAppName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, "/v1/", http.StatusSeeOther)
}

func (app *application) duplicateApp(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")

	// todo: should be adapter to rename app
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path

	err := Dir(appPath+"/"+appName, appPath+"/"+appName+"_copy")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/v1/", http.StatusSeeOther)
}

func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				return err
			}
		}
	}
	return nil
}

func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func (app *application) postNewApp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	newAppName := r.FormValue("AppName")

	// todo: should be an adapter to create app
	path, _ := os.Getwd()
	appPath := path + "/" + config.Config.Metadata.Path + "/" + newAppName
	if newAppName != "" {
		subfolders := []string{config.Config.Metadata.Metasubpath, config.Config.Deployer.Entitysubpath}

		// Create the app folder
		err := os.Mkdir(appPath, 0755)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create the metadata and entity subfolders.
		for _, subfolder := range subfolders {
			path := appPath + "/" + subfolder
			err := os.MkdirAll(path, 0755)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	http.Redirect(w, r, "/v1/", http.StatusSeeOther)
}

func (app *application) postAppChatGPT(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var data ChatGPT
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
		return
	}

	// response, err := chatgpt.RenameApp(data.ChatGPT)
	// if err != nil {
	// 	if response != nil {
	// 		// Marshal map to JSON
	// 		jsonString, _ := json.Marshal(response)
	// 		if jsonString != nil {
	// 			http.Error(w, err.Error()+string(jsonString), http.StatusBadRequest)
	// 			return
	// 		}
	// 	}

	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	http.Redirect(w, r, "/v1/", http.StatusSeeOther)
}

func (app *application) runApp(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	appName := params.ByName("app")
	global.SetAppName(appName)

	catalogs, err := metaapi.ReadMetadataList(metadata.Catalog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events, err := metaapi.ReadMetadataList(metadata.Event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmplData := struct {
		Catalogs []string
		Events   []string
	}{
		catalogs,
		events,
	}
	executeAppTemplate(w, "run", tmplData)
}
