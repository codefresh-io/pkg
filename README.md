codefresh common go modules

### Using [pre-commit](https://pre-commit.com/) hooks:
When installed correctly, this will run `golangci-lint` before every commit, and `go test` before every push. Any error will cause the commit or push to fail.

1. [Install](https://pre-commit.com/#1-install-pre-commit):  
   `brew install pre-commit`
1. [Install the git hook scripts](https://pre-commit.com/#3-install-the-git-hook-scripts):  
   `pre-commit install -t pre-commit -t pre-push`
