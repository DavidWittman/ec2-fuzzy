# 🔎 ec2-fuzzy

fzf-style fuzzy search for AWS EC2 instances. Inspired by [aws-fuzzy-finder](https://github.com/pmazurek/aws-fuzzy-finder)

## Installation

### Homebrew

```
brew tap DavidWittman/ec2-fuzzy https://github.com/DavidWittman/ec2-fuzzy
brew install ec2-fuzzy
```

## Configuration

All configurations can be passed by flags or environment variables.

 - `--private` or `EC2_FUZZY_PRIVATE=1` - Use private IP address when connecting to instance
 - `--user` or `EC2_FUZZY_USER` - Set SSH username
