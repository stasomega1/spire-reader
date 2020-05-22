package services

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"spire-reader/app/model"
	"spire-reader/app/store"
	"strings"
)

type SpireService struct {
	store  *store.Store
	logger *logrus.Logger
}

func NewSpireService(store *store.Store, logger *logrus.Logger) *SpireService {
	return &SpireService{
		store:  store,
		logger: logger,
	}
}

func (spireService *SpireService) Version() (string, error) {
	return spireService.store.PgdbRepository().Version()
}

func (spireService *SpireService) GetExampleRunData() (*model.SpireRuns, error) {
	//files, err := getFileNames("./exampledata/IRONCLAD")

	return nil, nil
}

func (spireService *SpireService) TestFunc(path string) (*model.SpireRun, error) {
	fileNames, err := getFileNames(path)
	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s/%s", path, fileNames[0])
	spireService.logger.Info(filePath)
	testFile, err := getFile(filePath)
	if err != nil {
		return nil, err
	}

	runData, err := getRunData(testFile)
	if err != nil {
		return nil, err
	}

	parsedRunData, err := parseJsonRunData(runData)
	if err != nil {
		return nil, err
	}

	return parsedRunData, nil
}

//Returns all .run file names in directory
func getFileNames(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(info.Name(), ".run") {
			return nil
		}
		files = append(files, info.Name())
		return err
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func getFile(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func getRunData(file []byte) (*model.JsonSpireRun, error) {
	jsonSpireRun := &model.JsonSpireRun{}

	if err := json.Unmarshal(file, &jsonSpireRun); err != nil {
		panic(err)
	}

	return jsonSpireRun, nil
}

func parseJsonRunData(jsonRunData *model.JsonSpireRun) (*model.SpireRun, error) {
	//Init pathPerFloorNormal
	pathPerFloorNormal := make([]string, len(jsonRunData.PathPerFloor))
	for key, val := range jsonRunData.PathPerFloor {
		valStr := ""
		if val != nil {
			valStr = *val
		} else {
			valStr = "N"
		}
		pathPerFloorNormal[key] = valStr
	}
	jsonRunData.PathPerFloorNormal = pathPerFloorNormal

	//Start parsing
	spireRun := &model.SpireRun{
		CharacterChosen: model.NewJsonNullString(jsonRunData.CharacterChosen),
		ItemsPurchased:  model.NewJsonNullStringSlice(jsonRunData.ItemsPurchased),
		PathPerFloor:    model.NewJsonNullStringSlice(jsonRunData.PathPerFloorNormal),
		FloorReached:    model.NewJsonNullInt32(int32(jsonRunData.FloorReached)),
		Playtime:        model.NewJsonNullInt32(int32(jsonRunData.Playtime)),
		Victory:         model.NewJsonNullBool(jsonRunData.Victory),
	}

	return spireRun, nil
}
