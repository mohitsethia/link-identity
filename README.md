# Link Identity

[Task Details]([url](https://github.com/mohitsethia/link-identity/blob/master/cmd/link-identity-api/README.md))

# How to run this application
`make run-docker`

Endpoints:
1. `localhost:8000/`, `localhost:8000/health/check` <br>
Response:`{"status_code":200,"data":"success"}`

2. `localhost:8000/identify` <br>
Request Payload: 
```
{
   "email": "test1@gmail.com",
   "phone": "+4917612345670"
}
```
Response:
```
{
    "status_code": 200,
    "data": {
        "contact": {
            "PrimaryContactID": 1,
            "emails": [
                "test1@gmail.com",
                "test2@gmail.com"
            ],
            "phoneNumbers": [
                "+4917612345670",
                "+4917612345672",
                "+4917612345671"
            ],
            "secondaryContactIds": [
                2,
                3
            ]
        }
    }
}
```
