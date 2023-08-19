# k8f 
## Release Notes
### Changes
- added AWS Credentials validation
  - by default AWS credentials and STS is beeing validated, in case of issues the creds are removed from list.
  - added force stop flag `--validate` if credentials are not validated.
  - improved Unit tests
- added `-q` `--quit` flag to show only error message in the cli, this is useful for scripting and automation.
<!-- ### Known Issues -->
### Known Issues
<!-- ## Contributors -->
<!-- ## Bugfix -->
<!-- ## Braking changes -->     
