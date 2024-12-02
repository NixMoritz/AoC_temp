package generate

import (
	"embed"
	"errors"
	"fmt"
	"github.com/linusback/aoc-go-template/pkg/errorsx"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"slices"
	"text/template"
)

const (
	ErrFailedToGetModuleName errorsx.SimpleError = "failed tog et name of go module"
)

//go:embed templates
var files embed.FS

type YearData struct {
	Imports []string
	Cases   []string
	Year    string
	ModName string
}

type DayData struct {
	Day  string
	Year string
}

func Generate(year string, days []string) error {
	if len(days) == 0 {
		return nil
	}
	moduleName, err := getModuleName()
	if err != nil {
		return err
	}

	err = generateSolver(moduleName, year)
	if err != nil {
		return err
	}

	err = generateYearSolve(moduleName, year, days)
	if err != nil {
		return err
	}
	for _, day := range days {
		err = generatDaySolve(moduleName, year, day)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func generateSolver(moduleName, year string) error {
	years := getYears(moduleName)
	// do not change file if year already implemented
	if slices.Contains(years, year) {
		return nil
	}
	years = append(years, year)
	slices.Sort(years)
	f, err := os.Create("./internal/solve.go")
	if err != nil {
		return err
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			err = errors.Join(err, err2)
		}
	}()

	data := YearData{
		Imports: make([]string, 0, len(years)),
		Cases:   make([]string, 0, len(years)),
		Year:    year,
		ModName: moduleName,
	}
	intendent := ""
	for i, y := range years {
		data.Imports = append(data.Imports, fmt.Sprintf(`%s"%s/internal/year%s"`, intendent, moduleName, y))
		data.Cases = append(data.Cases, fmt.Sprintf(`%scase "%s":
		return year%s.Solve(day)`, intendent, y, y))
		if i == 0 {
			intendent = "\n\t"
		}
	}
	t, err := template.ParseFS(files, "templates/solve.go.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	return nil
}

func generateYearSolve(moduleName, year string, days []string) error {
	var allExists bool
	err := createDirIfNotExists(fmt.Sprintf("./internal/year%s", year))
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("./internal/year%s/solve.go", year)

	allExists, err = checkDays(moduleName, year, days, filePath)
	if err != nil {
		return err
	}
	if allExists {
		return nil
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			err = errors.Join(err, err2)
		}
	}()

	data := YearData{
		Imports: make([]string, 0, len(days)),
		Cases:   make([]string, 0, len(days)),
		Year:    year,
		ModName: moduleName,
	}
	intendent := ""
	for i, day := range days {
		data.Imports = append(data.Imports, fmt.Sprintf(`%s"%s/internal/year%s/day%s"`, intendent, moduleName, year, day))
		data.Cases = append(data.Cases, fmt.Sprintf(`%scase "%s":
		return day%s.Solve()`, intendent, day, day))
		if i == 0 {
			intendent = "\n\t"
		}
	}
	t, err := template.ParseFS(files, "templates/solve_year.go.tmpl")
	if err != nil {
		return err
	}

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}

func checkDays(moduleName, year string, days []string, path string) (allExists bool, err error) {
	exists, err := fileExists(path)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	regexpImport, err := regexp.Compile(fmt.Sprintf(`"%s/internal/year%s/day(\d{1,2})"`, moduleName, year))
	if err != nil {
		return false, err
	}
	byteYears := regexpImport.FindAllSubmatch(fileBytes, -1)

	foundDays := make([]string, 0, len(byteYears))
	for _, yb := range byteYears {
		if len(yb) < 2 {
			continue
		}
		d := string(yb[1])
		if !slices.Contains(foundDays, d) {
			foundDays = append(foundDays, d)
		}
	}
	for _, d := range days {
		if !slices.Contains(foundDays, d) {
			return false, nil
		}
	}

	return true, nil
}

func generatDaySolve(moduleName, year string, day string) error {
	err := createDirIfNotExists(fmt.Sprintf("./internal/year%s/day%s", year, day))
	if err != nil {
		return err
	}

	err = createEmptyFileIfNotExists(fmt.Sprintf("./internal/year%s/day%s/example", year, day))
	if err != nil {
		return err
	}

	solveFilePath := fmt.Sprintf("./internal/year%s/day%s/solve.go", year, day)

	exists, err := fileExists(solveFilePath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	f, err := os.Create(solveFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			err = errors.Join(err, err2)
		}
	}()

	data := DayData{
		Day:  day,
		Year: year,
	}

	t, err := template.ParseFS(files, "templates/solve_day.go.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	return nil
}

func getYears(moduleName string) []string {
	yearRegex, err := regexp.Compile(`"` + moduleName + `/internal/year(\d{4})"`)
	if err != nil {
		return nil
	}
	oldFile, err := os.ReadFile("./internal/solve.go")
	if err != nil {
		return nil
	}
	byteYears := yearRegex.FindAllSubmatch(oldFile, -1)
	result := make([]string, 0, len(byteYears)+1)
	for _, yb := range byteYears {
		if len(yb) < 2 {
			continue
		}
		y := string(yb[1])
		if !slices.Contains(result, y) {
			result = append(result, y)
		}
	}
	slices.Sort(result)
	return result
}

func getModuleName() (string, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "", ErrFailedToGetModuleName
	}

	if len(bi.Deps) == 0 {
		return "", ErrFailedToGetModuleName
	}
	// this is still a bit hacky probably should look at the actual mod file or something
	for _, d := range bi.Deps {
		if d.Version == "(devel)" {
			return d.Path, nil
		}
	}

	return "", ErrFailedToGetModuleName
}

func createDirIfNotExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0775)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func fileExists(filePath string) (exists bool, err error) {
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func createEmptyFileIfNotExists(filePath string) error {
	var f *os.File
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		f, err = os.Create(filePath)
		if err != nil {
			return err
		}
		defer func() {
			err2 := f.Close()
			if err2 != nil {
				err = errors.Join(err, err2)
			}
		}()
	} else if err != nil {
		return err
	}
	return nil
}
