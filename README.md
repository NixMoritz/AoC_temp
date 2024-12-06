# Advent of Code template in Go
Go template project for [advent of code](https://adventofcode.com/) (AOC) is intended for use with gonew see https://go.dev/blog/gonew
This will automatically download input and puzzle instructions. As well as submit answers to the advent of code website.

### DISCLAIMER 
This project is thrown together with late night coding and ductape, so there might be some bugs. If you find any feel free to leave a bug report or do a pull request.  
English is not my first language and I apologize in advance for any spelling errors or grammar mistakes. But I do enjoy learning so feel free to point them out 😀

# Install instructions:
1. [Download and install](https://go.dev/doc/install) a version of Golang >= 1.23.3
2. Do 1 of either:
    * clone repo and optionally manually rename package to you liking   
    ```git clone git@github.com:linusback/aoc-go-template.git```
   * install gonew and initialize it that way, will automatically rename package  
   ```git clone git@github.com:linusback/aoc-go-template.git```  
   ```gonew github.com/linusback/aoc-go-template github.com/your_github/your_project```
3. Generate session Cookie (by logging in) and saving it according to [instructions](https://github.com/linusback/aoc-go-template?tab=readme-ov-file#generate-and-save-session-cookie)
4. run make inside the projects root folder, this will generate files for latest year.  
   ```cd your/project/path && make``` 

### Generate and save Session Cookie
You will need a session Cookie for it to be able to download and submit answer, [detailed instructions](https://www.cookieyes.com/blog/how-to-check-cookies-on-your-website-manually/).
1. Login on advent of code.
2. Verify that you are logged in.
3. Open development tools in your favourite browser (I press F12 on chrome in Ubuntu). 
4. Navigate to Application and copy the value of the session cookie (not my actual cookie feel free to jwt decode it though 😉). 



There are 3 ways to provide the value for the session cookie.

1. In an ADVENT_OF_CODE_SESSION environment variable.
2. In a file called .adventofcode.session (note the dot) in your home directory /home/alice on Linux (not tested for other operating systems).
3. In a file called adventofcode.session (no dot) in your user's config directory /home/alice/.config on Linux (not tested for other operating systems).

inspired by https://github.com/scarvalhojr/aoc-cli

# Usage:
Using generator without any arguments will try to download all available puzzle and inputs for the current year.
Calling the solver will try to solve and submit the current days problem (or latest one).
Just running make will call both without any arguments. 

If you want to generate files or run a solver for specific year day, call the solver/generator like this
```
make run-generator YEAR=2024 DAY=2
make run-solver YEAR=2024 DAY=2
```

The solver (aoc binary) only works if you first generate the needed files.

Generally the filestructure inside of internal should be something like this.
But feel free to add other packages outside of this structure or adding more packages 
([my own version](https://github.com/linusback/aoc) a few of those like 
[pkg/util/input.go](https://github.com/linusback/aoc/blob/main/pkg/util/input.go "handle input files") or [pkg/util/parse.go](https://github.com/linusback/aoc/blob/main/pkg/util/parse.go "parse input rows helpers")).

```
internal/
├─ some-other-package/  <- feel free to create your own packages
├─ year2023/
├─ year2024/
│  ├─ day1/
│  ├─ day2/
│  │  ├─ solve.go   <- write your code here
│  │  ├─ example    <- generated but empty, feel free to edit.
│  │  ├─ input      <- automatically downloaded from AOC
│  │  ├─ puzzle.md  <- automatically downloaded from AOC
│  │  ├─ answer1    <- might be downloaded from AOC or generated once a successfull answer is submitet
│  │  ├─ answer2    <- might be downloaded from AOC or generated once a successfull answer is submitet
│  ├─ solve.go      <- do not edit (might be overwritten by generator)
├─ .gitignore 
├─ solve.go         <- do not edit (might be overwritten by generator)
```

And each days problem will be available in internal/year{year}/day{day}/solve.go

The other solve files (inside of internal/ and internal/year{year}/) are not intended to be edited, since they can be overwritten by generate, they are only meant to be used as switch statements to execute the correct solver by aoc command. 

A typical solve file for a specific day will look like this as a blank template: 
```
package day2

import (
	"log"
	"os"
)

const (
	exampleFile = "./internal/year2024/day2/example"
	inputFile   = "./internal/year2024/day2/input"
)

func Solve() (solution1, solution2 string, err error) {
	log.Println("welcome to day 2 of advent of code")
	b, err := os.ReadFile(exampleFile)
	if err != nil {
		return
	}
	log.Println(string(b))
	return
}

```
Setting either solution1, solution2 or both will submit those to advent of code. 
Unless you already have answered them before (there is an answer1 and/or answer2 file inside the internal/year{year}/day{day} directory).
Previously answered problems should be detected automatically by generator when downloading the puzzle.md file.
Empty strings will not be submitted, hence doing something like this: 
```
package day2
...
func Solve() (solution1, solution2 string, err error) {
	solution1="42"
	return
}
```
Would only try to answer part 1 of day 2, and only if there is not an answer1 file in the folder.

Calling generate again will not change solve.go inside the day package (internal/year{year}/day{day}/solve.go) but might overwrite the 
solve file for the year package internal/year{year}/solve.go as well as internal/solve.go (if needed).  

Happy coding!
