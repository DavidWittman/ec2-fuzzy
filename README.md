# 🔎 ec2-fuzzy

fzf-style fuzzy search for AWS EC2 instances. Inspired by [aws-fuzzy-finder](https://github.com/pmazurek/aws-fuzzy-finder)

## Installation

### Homebrew

```
brew tap DavidWittman/ec2-fuzzy https://github.com/DavidWittman/ec2-fuzzy
brew install ec2-fuzzy
```

## Usage

### Fuzzy Search

Just run `ec2-fuzzy`.

### Exact Match

Pass the `--instance` or `-i` flag with the Name or Instance ID of the instance.

```
ec2-fuzzy -i my-instance-name
# or by instance ID
ec2-fuzzy -i i-0675e1acdc61a6cc7
```

If multiple instances match the Name tag provided, ec2-fuzzy will connect to the first match.

## Configuration

All configurations can be passed by flags or environment variables.

 - `--private` or `EC2_FUZZY_PRIVATE=1` - Use private IP address when connecting to instance
 - `--user` or `EC2_FUZZY_USER` - Set SSH username
