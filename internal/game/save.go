package game

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"maps"

	"github.com/joho/godotenv"
)

const SAVE_FILE_FOLDER = "go-game-space-shooter"

func NewSave() *Save {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		HandleError(err)
	}

	data := make(map[string]any)

	saveFileName := Configs["SAVE_FILE_NAME"]
	if saveFileName == "" {
		HandleError(errors.New("save file name cannot be empty"))
	}
	if filepath.Ext(saveFileName) != ".save" {
		HandleError(errors.New("save file extension must be \".save\""))
	}

	return &Save{
		path:     userConfigDir,
		filename: saveFileName,
		data:     data,
	}
}

func (s *Save) Save(game *Game) (bool, error) {
	err := s.createFileIfNotExists()
	if err != nil {
		return false, err
	}

	origData := s.LoadSave(game, true)

	data := make(map[string]string, len(origData))
	maps.Copy(data, origData)

	// Save: Highscore
	data["HIGHSCORE"] = strconv.FormatInt(game.score.GetHighScore(), 10)

	if !maps.Equal(data, origData) {
		godotenv.Write(data, filepath.Join(s.path, SAVE_FILE_FOLDER, s.filename))
	}

	return true, nil
}

func (s *Save) LoadSave(game *Game, return_only bool) (data map[string]string) {
	saveFile := filepath.Join(s.path, SAVE_FILE_FOLDER, s.filename)

	data, err := godotenv.Read(saveFile)

	if return_only {
		return data
	}

	if err == nil {
		// Load: Highscore
		highscore, err := strconv.ParseInt(data["HIGHSCORE"], 10, 64)
		if err == nil {
			game.score.SetHighScore(highscore)
		}
	}

	return data
}

func (s *Save) createFileIfNotExists() error {
	err := os.Mkdir(filepath.Join(s.path, SAVE_FILE_FOLDER), 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}

	_, err = os.ReadFile(filepath.Join(s.path, SAVE_FILE_FOLDER, s.filename))
	if err != nil {
		err = os.WriteFile(filepath.Join(s.path, SAVE_FILE_FOLDER, s.filename), []byte(""), 0660)
		if err != nil {
			return err
		}
	}

	return nil
}
