package aoc

import (
	"bytes"
	"errors"
	"fmt"
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/linusback/aoc-go/pkg/errorsx"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type Part string

const (
	aocCli          = "aoc"
	Part1      Part = "1"
	Part2      Part = "2"
	InputFile       = "input"
	PuzzleFile      = "puzzle.md"
)

const (
	ErrCmdIsNil                  errorsx.SimpleError = "cmd was nil for some reason"
	ErrCouldNotFindSessionCookie errorsx.SimpleError = "could not find session cookie"
)

var (
	sessionClient *http.Client
	sessionErr    error
)

func Send(part Part, day, answer string) error {
	cmd := exec.Command(aocCli, "submit", string(part), answer, "-d", day)
	if cmd == nil {
		return ErrCmdIsNil
	}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func Download(year string, day string) error {
	location := fmt.Sprintf("./internal/year%s/day%s/", year, day)
	inputFound, puzzleFound, err := filesAlreadyExists(location)
	if err != nil {
		return err
	}

	client, err := getSessionClient()
	if err != nil {
		return err
	}

	errCh := make(chan error, 2)
	wg := new(sync.WaitGroup)
	if !puzzleFound {
		downloadAsync(wg, errCh, client, year, day, location, PuzzleFile, "", parseHtmlToMarkdown)
	}
	if !inputFound {
		downloadAsync(wg, errCh, client, year, day, location, InputFile, "/input", io.Copy)
	}

	wg.Wait()
	close(errCh)
	for e := range errCh {
		err = errors.Join(err, e)
	}
	return err
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

func downloadAsync(wg *sync.WaitGroup, errCh chan<- error, client *http.Client, year, day, location, file, endpoint string, handleResp func(io.Writer, io.Reader) (int64, error)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := download(client, year, day, location, file, endpoint, handleResp)
		if err != nil {
			errCh <- err
			return
		}
	}()
}

func download(client *http.Client, year, day, location, file, endpoint string, handleResp func(io.Writer, io.Reader) (int64, error)) error {
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

func parseHtmlToMarkdown(w io.Writer, r io.Reader) (int64, error) {
	reg, err := regexp.Compile("(?i)(?s)<main>(?P<main>.*)</main>")
	if err != nil {
		return 0, err
	}
	by, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	by, err = htmltomarkdown.ConvertReader(bytes.NewReader(reg.Find(by)))
	if err != nil {
		return 0, err
	}
	return io.Copy(w, bytes.NewReader(by))
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
		if e.Name() == InputFile && e.Type().IsRegular() {
			inputFound = true
		}
		if e.Name() == PuzzleFile && e.Type().IsRegular() {
			puzzleFound = true
		}
		if inputFound && puzzleFound {
			return
		}
	}

	return
}
