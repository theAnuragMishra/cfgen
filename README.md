# cfgen

This is a codeforces contest scraper that creates a directory for a contest with subdirectories for the problems in the contest.

## Setup

Clone this repo- `git clone https://github.com/theAnuragMishra/cfgen.git`.
Go to `main.go` and replace the username value in line 46 with your username (doesn't really matter, the username is used only to add author to the solutions).
Build the binary- `go build -o cfgen`.
Move this binary by cutting and pasting to a folder where you want all your contests to be.
Now you're ready to use cfgen!

## Usage

Go to the directory where you want all your contests to be and run `cfgen -c {contest_id}`. Replace `${contest_id}` with the actual contest id which is usually the four digit number in the contest url.
Boom! You are ready to write the solutions!
