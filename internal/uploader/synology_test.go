package uploader

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/eventials/go-tus"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

func TestAPIs(t *testing.T) {
	go uploadServer(t)

	config := ConfigOptions{
		hostname:  "localhost",
		port:      5000,
		Account:   "admin",
		Passwd:    "SuperD0CK!@#",
		StorePath: t.TempDir(),
	}

	apiInfo, err := GetAPIInfo(config)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Got response: %+v\n", *apiInfo)

	var sid string

	t.Run("auth login api", func(t *testing.T) {
		sid, err = AuthLogin(apiInfo, config)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("Got sid: %s\n", sid)
	})
	t.Run("file download api", func(t *testing.T) {
		file, err := FileDownload(
			apiInfo,
			config,
			[]string{
				"/photo/Superdock/202104/20210401115140/result/mission_finish.jpg",
				"/photo/Superdock/202104/20210401115140/result/ready_to_fly.jpg",
			},
			sid,
		)
		if err != nil {
			t.Fatal(err)
		}
		if err := uploadFile(file); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("auth logout api", func(t *testing.T) {
		if err = AuthLogout(apiInfo, config, sid); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInfo(t *testing.T) {
	requestElem := requestElem{
		schema:  schema,
		host:    "localhost:5000",
		apiName: "SYNO.API.Info",
		version: 1,
		path:    "query.cgi",
		method:  "query",
		params:  "query=all",
	}
	url, err := constructURL(requestElem)
	if err != nil {
		t.Fatal(err)
	}
	apiInfo, err := Info(url.String())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Got response: %+v\n", *apiInfo)
}

func TestConstructURL(t *testing.T) {
	requestElem := requestElem{
		schema:  schema,
		host:    "localhost:5000",
		apiName: "SYNO.API.Info",
		version: 1,
		path:    "query.cgi",
		method:  "query",
	}
	url, err := constructURL(requestElem)
	if err != nil {
		t.Fatal(err)
	}
	if url.String() != "http://localhost:5000/webapi/query.cgi?api=SYNO.API.Info&method=query&version=1" {
		t.Fatalf("got: %s want: %s", url.String(), "http://localhost:5000/webapi/query.cgi?api=SYNO.API.Info&method=query&version=1")
	}
}

func uploadServer(t *testing.T) {
	store := filestore.FileStore{
		Path: t.TempDir(),
	}

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		panic(fmt.Errorf("Unable to create handler: %s", err))
	}

	go func() {
		for {
			event := <-handler.CompleteUploads
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	http.Handle("/files/", http.StripPrefix("/files/", handler))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(fmt.Errorf("Unable to listen: %s", err))
	}
}

func uploadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	client, err := tus.NewClient("http://localhost:8080/files", nil)
	if err != nil {
		return err
	}
	upload, err := tus.NewUploadFromFile(f)
	if err != nil {
		return err
	}
	uploader, err := client.CreateUpload(upload)
	if err != nil {
		return err
	}
	return uploader.Upload()
}
