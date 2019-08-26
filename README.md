git-switch-trainer
======

## Description

You are forced to use git switch by `git-switch-trainer`.

## Getting Started

### Install
```shell
go get -u github.com/sonatard/git-switch-trainer/
```

### Settings

Add alias to your `.bash_profile`

- .bash_profile

```
alias git=git-switch-trainer
```

```
$ source ~/.bash_profile
```


## Example

```shell
$ git checkout master 
Error: Use git switch or git restore instead of git checkout.
```
