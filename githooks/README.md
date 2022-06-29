## Configuration
To configure Git Hooks for your branch / repository do the following:
```sh=
git submodule add git@github.com:rogersdigitalexperience/git-hooks.git git_hooks
git config --global core.hooksPath "./git_hooks"
```

## Prepare-commit-msg
Read details at [prepare-commit-msg](https://git-scm.com/docs/githooks#_prepare_commit_msg).

### Workflow

We handle 3 scenarios according to the Git usage agreements and request:

1) Developer commits message in format `CC-0000: summary` or `REDESIGN: summary`

Then we add initials of the corresponding user which are configured and taken from `git config --get user.name`.

In result, considering user.name="Firstname Lastname", we get a commit message like `[FL] CC-0000: summary` or `[FL] REDESIGN: summary`.

2) Developer commits message in format `[FL] CC-0000: summary` or `[FL] REDESIGN: summary`

We just accept that message.

3) Developer commits message `SOMETHING DIFFERENT ^_^`

We print the following error message:
> **WARNING**  
Aborting commit. 
Acceptable options: 'CC-0000: summary', '[FL] CC-0000: summary', 'REDESIGN: summary', "[FL] REDESIGN: summary"



## Pre-commit
Read details at [pre-commit](https://git-scm.com/docs/githooks#_pre_commit).
  
### Introduction
Pre-commit hook is purposed to verify credentials violation in the staged files.

### Workflow
We simply run checks against each staged file prepared for the commit.

In case of Docker is not found in the system - the hook will print the warning:
> **WARNING**  
DOCKER IS NOT INSTALLED IN THE SYSTEM.
CREDENTIALS VIOLATION SCAN CANNOT BE PERFORMED.
And commit will pass (yet the CI will fail, please keep that in mind).

### Technical details
We are using 2 different tools: 
* cred-alert (has good support of regexp masks)
* detect-secrets (has strong heuristic analysis engine)

Regexp we use for cred-alert:
```regexp=
\.clientId[ ]*=[ ]*[A-Z0-9_+=(){}*#@!%^\-"'~`|.,\[\]].*
\.clientSecret[ ]*=[ ]*[A-Z0-9_+=(){}*#@!%^\-"'~`|.,\[\]].*
\.admin-password[ ]*=[ ]*[A-Z0-9_+=(){}*#@!%^\-"'~`|.,\[\]].*
\.admin-account[ ]*=[ ]*[A-Z0-9_+=(){}*#@!%^\-"'~`|.,\[\]].*
\.jwtSecret[ ]*=[ ]*[A-Z0-9_+=(){}*#@!%^\-"'~`|.,\[\]].*
```
_(cred-alert is converting the text to upper-case register before running verification checks against it)_

### False-positives
In case of you would like to prevent detect-secrets from scanning specific line please add a comment to the end of that line `pragma: allowlist secret`.

### Outputs
#### detect-secrets
```
"application-ctp-client.properties": [
 
"hashed_secret": "d64193a396be117315362d70f1224099f829eec6",
"is_verified": false,
"line_number": 4,
"type": "Secret Keyword"
 
]
```

#### cred-alert

```
[CRED] application-ctp-client.properties:3
[CRED] application-ctp-client.properties:4
```
