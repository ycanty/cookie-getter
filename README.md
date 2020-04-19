# cookie-getter - Get cookies saved in web browsers

This small CLI utility allows to extract cookies from web browser cookie storage,
decrypt them and print their value.

## Installation

```bash
go get -u github.com/getmelisted/cookie-getter
```

This will download, build and install the latest released version
of `cookie-getter` in ${GOPATH}/bin.  Make sure this location is in your PATH:
```
export PATH=$(go env GOPATH)/bin:${PATH}
```

## Usage

```
cookie-getter --help
```

## Main use case

This utility was created with the intent to use it from automation scripts,
as an easy means to get authentication cookies like `JSESSIONID` from the web browser,
for the scripts to be able to authenticate to REST APIs that are behind OKTA
authentication.

The overall process is as follows:

1. Login normally to OKTA in the web browser, then open the application
   from its OKTA chicklet, for example JIRA.
2. In the automation script that needs to access the JIRA REST API, get the
   `JSESSIONID` cookie from the web browser, with:
    ```bash  
    cookie-getter --short --domain jira.domain.com --name JSESSIONID
    ```
3. Use that `JSESSIONID` as a cookie with every request sent to the JIRA
   API.  For example, with curl:
   ```bash
   ID=$(cookie-getter --short --domain jira.domain.com --name JSESSIONID)
   curl -L -v \
   -b "JSESSIONID=${ID};" \
   -D- \
   -X GET \
   -H "Content-Type: application/json" \
   https://jira.domain.com/rest/api/2/search?jql=assignee="someone"
   ```
