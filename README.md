# cfgen

This is a codeforces contest scraper that creates a directory for a contest with subdirectories for the problems in the contest.

## Installation

### Via script

The command below installs and sets up the binary automatically. You can read `install.sh` to verify what the script does.
`curl -sSfL https://raw.githubusercontent.com/theAnuragMishra/cfgen/main/install.sh | sh
`

### Manual

Download the binary corresponding to your os from the releases and rename it to whatever you like.
Make it executable: `chmod +x {binary_name}`.

## Usage

Go to the directory where you want all your contests to be and run `cfgen{contest_id}`. Replace `${contest_id}` with the actual contest id which is usually the four digit number in the contest url.

Boom! You are ready to write the solutions!

## Options

- Use `--save-ps` (alias `sps`) if you want to save the problem statements for the problems.
- Use `-u {username}` to include a username comment in the solution template generated.

## Contributing

I'm willing to accept prs introducing new features or improving the ux :D
