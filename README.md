# cfgen

This is a codeforces contest scraper that creates a directory for a contest with subdirectories for the problems in the contest.

## Usage

- Download the binary corresponding to your operating system.
- Go to the directory where you want all your contests to be and run `cfgen{contest_id}`. Replace `${contest_id}` with the actual contest id which is usually the four digit number in the contest url.

Boom! You are ready to write the solutions!

## Options

- Use `--save-ps` (alias `sps`) if you want to save the problem statements for the problems.
- Use `-u {username}` to include a username comment in the solution template generated.

## Contributing

I'm willing to accept prs introducing new features or improving the ux :D
