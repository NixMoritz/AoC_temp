# Advent of Code template in Go
Go template project for [advent of code](https://adventofcode.com/) (AOC) is intended for use with gonew see https://go.dev/blog/gonew
This will automatically download input and puzzle instructions. As well as submit answers to the advent of code website.

Feel free to submit an issue report if there is anything you would like to see here.

### Todo

- [ ] Add more comprehensive tests for packages.
- [ ] Add godoc documentation for packages in pkg.
- [ ] Create template for submitting issue reports.
- [ ] Redo structure as to minimize the size of final aoc binary (will only grow over time).
  - [ ] Add separate run binary for individual days, add those to Makefile 
- [ ] Change from args to using flags (as in options --date for aoc and generator).
    - [ ] Add --help flag as well as switching puzzle input file via flags.
- [ ] Make generated files have proper file ending.
- [ ] Switch from log to slog package to allow for more granular control of logging in solver and generator.
- [ ] Rename "aoc" binary to "solver" (it is reference as such in other places of the project).

### In Progress

- [ ] Simplify coding structure and update README.md.

### Done ✓

- [x] Create TODO.md for project.
- [x] Add check for already submitted answers.
- [x] Publish and tag first version.

