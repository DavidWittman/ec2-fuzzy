# 🔎 ec2-fuzzy

fzf-style fuzzy search for AWS EC2 instances. Inspired by [aws-fuzzy-finder](https://github.com/pmazurek/aws-fuzzy-finder)

## Configuration

All configurations can be passed by flags or environment variables.

 - `--private` or `EC2_FUZZY_PRIVATE=1` - Use private IP address when connecting to instance
 - `--user` or `EC2_FUZZY_USER` - Set SSH username
