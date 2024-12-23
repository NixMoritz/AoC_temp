package aoc

import (
	"bytes"
	"errors"
	"fmt"
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/linusback/aoc-go-template/pkg/errorsx"
	"github.com/linusback/aoc-go-template/pkg/filenames"
	"github.com/linusback/aoc-go-template/pkg/util"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

type Part string

const (
	Part1 Part = "1"
	Part2 Part = "2"
)

const (
	ErrCouldNotFindSessionCookie errorsx.SimpleError = "could not find session cookie"
	ErrWrongAnswer               errorsx.SimpleError = "that's not the right answer"
	ErrCouldNotFinsExample       errorsx.SimpleError = "could not extract example.txt from puzzle.md"
)

var (
	sessionClient       *http.Client
	sessionErr          error
	answerWrongMsg      = []byte(`That's not the right answer`)
	answerToRecentlyMsg = []byte(`You gave an answer too recently`)
	answerCorrectMsg    = []byte(`That's the right answer!`)
	exampleRegExp       *regexp.Regexp
)

func Send(part Part, year, day, answer string) error {
	if len(answer) == 0 {
		log.Printf("empty answer part %s, skipping...\n", part)
		return nil
	}
	answerFile := fmt.Sprintf("./internal/year%s/day%s/%s", year, day, answerFile(part))
	exists, err := util.FileExists(answerFile)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("already answer part %s, skipping...\n", part)
		return checkExistingAnswer(answerFile, answer)
	}
	client, err := getSessionClient()
	if err != nil {
		return err
	}

	respBody, err := sendAnswer(client, part, year, day, answer)
	if err != nil {
		return err
	}

	err = checkResponseAndCache(respBody, answer, answerFile)
	if err != nil {
		return err
	}
	log.Printf("part %s: That's the right answer!", part)
	if part != Part1 {
		return nil
	}

	err = reDownloadPuzzle(client, year, day)
	if err != nil {
		return err
	}

	return nil
}

func checkExistingAnswer(filename, answer string) error {
	existingAnswer, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	existingAnswer = bytes.TrimSpace(existingAnswer)
	if !util.BytesEqualString(existingAnswer, answer) {
		return fmt.Errorf("existing answer: \"%s\" doesn't match new \"%s\"", existingAnswer, answer)
	}
	return nil
}

func checkResponseAndCache(body []byte, answer, filePath string) error {
	reg, err := regexp.Compile("(?i)(?s)<main>(?P<main>.*)</main>")
	if err != nil {
		return err
	}
	mainBody := reg.Find(body)
	if bytes.Contains(mainBody, answerWrongMsg) {
		return ErrWrongAnswer
	}
	if bytes.Contains(mainBody, answerToRecentlyMsg) {
		return errors.New(string(mainBody))
	}
	if !bytes.Contains(mainBody, answerCorrectMsg) {
		return errors.New(string(mainBody))
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		err2 := file.Close()
		if err2 != nil {
			err = errors.Join(err, err2)
		}
	}()
	_, err = file.WriteString(answer)
	if err != nil {
		return err
	}

	return nil
}

func sendAnswer(client *http.Client, part Part, year, day, answer string) (respBody []byte, err error) {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://adventofcode.com/%s/day/%s/answer", year, day), strings.NewReader(fmt.Sprintf("level=%s&answer=%s", part, answer)))
	if err != nil {
		return nil, err
	}

	request.Header.Set("content-type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		err2 := resp.Body.Close()
		if err2 != nil {
			err = errors.Join(err, err2)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed request status: %d", resp.StatusCode)
	}
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return
}

func Download(year string, days []string) error {
	client, err := getSessionClient()
	if err != nil {
		return err
	}
	for _, day := range days {
		err = download(client, year, day)
		if err != nil {
			log.Printf("error while downloading files for year %s day %s: %v", year, day, err)
		}
	}
	return nil
}

func download(client *http.Client, year, day string) error {
	location := fmt.Sprintf("./internal/year%s/day%s/", year, day)
	inputFound, puzzleFound, err := filesAlreadyExists(location)
	if err != nil {
		return err
	}
	if inputFound && puzzleFound {
		log.Printf("internal/year%s/day%s/%s: file already exists skipping...\n", year, day, filenames.PuzzleFile)
		log.Printf("internal/year%s/day%s/%s: file already exists skipping...\n", year, day, filenames.InputFile)
		return nil
	}

	errCh := make(chan error, 2)
	wg := new(sync.WaitGroup)
	if !puzzleFound {
		log.Printf("generating internal/year%s/day%s/%s\n", year, day, filenames.PuzzleFile)
		exampleFile := location + filenames.ExampleFile
		downloadFileAsync(wg, errCh, client, year, day, location, filenames.PuzzleFile, "", parseHtmlToMarkdown(exampleFile))
	} else {
		log.Printf("internal/year%s/day%s/%s: file already exists skipping...\n", year, day, filenames.PuzzleFile)
	}
	if !inputFound {
		log.Printf("generating internal/year%s/day%s/%s\n", year, day, filenames.InputFile)
		downloadFileAsync(wg, errCh, client, year, day, location, filenames.InputFile, "/input", io.Copy)
	} else {
		log.Printf("internal/year%s/day%s/%s: file already exists skipping...\n", year, day, filenames.InputFile)
	}

	wg.Wait()
	close(errCh)
	for e := range errCh {
		err = errors.Join(err, e)
	}
	return err
}

func reDownloadPuzzle(client *http.Client, year, day string) error {
	location := fmt.Sprintf("./internal/year%s/day%s/", year, day)
	err := downloadFile(client, year, day, location, filenames.PuzzleFile, "", parseHtmlToMarkdown(""))
	if err != nil {
		return err
	}
	return nil
}

func getSessionClient() (*http.Client, error) {
	sync.OnceFunc(func() {
		jar, err := cookiejar.New(nil)
		if err != nil {
			sessionErr = err
			return
		}

		u, err := url.Parse("https://adventofcode.com/")
		if err != nil {
			sessionErr = err
			return
		}

		val, err := getSessionCookie()
		if err != nil {
			sessionErr = err
			return
		}

		jar.SetCookies(u, []*http.Cookie{
			{
				Name:  "session",
				Value: val,
			},
		})

		sessionClient = &http.Client{
			Jar: jar,
		}
	})()

	return sessionClient, sessionErr
}

func downloadFileAsync(wg *sync.WaitGroup, errCh chan<- error, client *http.Client, year, day, location, file, endpoint string, handleResp func(io.Writer, io.Reader) (int64, error)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := downloadFile(client, year, day, location, file, endpoint, handleResp)
		if err != nil {
			errCh <- err
			return
		}
	}()
}

func downloadFile(client *http.Client, year, day, location, file, endpoint string, handleResp func(io.Writer, io.Reader) (int64, error)) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://adventofcode.com/%s/day/%s%s", year, day, endpoint), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed request status: %d", resp.StatusCode)
	}
	f, err := os.OpenFile(location+file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = handleResp(f, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func parseHtmlToMarkdown(exampleFile string) func(io.Writer, io.Reader) (int64, error) {
	return func(w io.Writer, r io.Reader) (int64, error) {
		reg, err := regexp.Compile("(?i)(?s)<main>(?P<main>.*)</main>")
		if err != nil {
			return 0, err
		}
		by, err := io.ReadAll(r)
		if err != nil {
			return 0, err
		}
		by, err = htmltomarkdown.ConvertReader(bytes.NewReader(reg.Find(by)),
			converter.WithDomain("https://adventofcode.com/"))
		if err != nil {
			return 0, err
		}
		// maybe try parsing looking for example use this regex "For example:\s*```[\s]?([^`]*)```"
		if len(exampleFile) > 0 {
			err = writeExampleFile(exampleFile, by)
			if err != nil {
				log.Printf("warning: unable to extract example from puzzle.md please fill out manualy")
				err = nil
			}
		}

		return io.Copy(w, bytes.NewReader(by))
	}
}

func getSessionCookie() (string, error) {
	var (
		ok   bool
		s    string
		err  error
		home string
	)
	s, ok = os.LookupEnv("ADVENT_OF_CODE_SESSION")
	if ok {
		return s, nil
	}
	home, err = os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	s, err = getSessionFromFile(home + "/.config/adventofcode.session")
	if err == nil {
		return s, nil
	}
	sessionCookieError := err
	s, err = getSessionFromFile(home + "/.adventofcode.session")
	if err == nil {
		return s, nil
	}
	return "", errors.Join(ErrCouldNotFindSessionCookie, sessionCookieError, err)
}

func getSessionFromFile(s string) (string, error) {
	f, err := os.Open(s)
	if err != nil {
		return "", err
	}
	sb := new(strings.Builder)
	_, err = io.Copy(sb, f)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(sb.String()), nil
}

func filesAlreadyExists(dir string) (inputFound, puzzleFound bool, err error) {
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range dirEntry {
		if e.Name() == filenames.InputFile && e.Type().IsRegular() {
			inputFound = true
		}
		if e.Name() == filenames.PuzzleFile && e.Type().IsRegular() {
			puzzleFound = true
		}
		if inputFound && puzzleFound {
			return
		}
	}

	return
}

func writeExampleFile(filename string, b []byte) (err error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err = getExampleFromPuzzleFile(b)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(b))
	if err != nil {
		return err
	}
	return nil
}

func getExampleFromPuzzleFile(b []byte) (res []byte, err error) {
	if exampleRegExp == nil {
		exampleRegExp, err = regexp.Compile("(?i)for example:\\s*```\\s?([^`]*)```")
		if err != nil {
			return nil, err
		}
	}
	sub := exampleRegExp.FindSubmatch(b)
	if len(sub) < 2 {
		return nil, ErrCouldNotFinsExample
	}
	return sub[1], nil
}

func answerFile(part Part) string {
	switch part {
	case Part1:
		return filenames.Answer1
	case Part2:
		return filenames.Answer2
	default:
		panic(fmt.Sprintf("unknown aoc.Part type %s", part))
	}
}
