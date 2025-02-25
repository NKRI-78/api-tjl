# Api Spec

## LOGIN 

Request : 
- Method : POST
- Endpoint : `/api/v1/login`
- Body :

```json
{
    "val": "089670558381",
    "password": "12345678",
}
```

Response : 

```json
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJpZCI6IjUyNzU3OWI5LTBmMDItNDY2NS05MWIyLTY3MzVlNzc1YTg1ZSJ9.3ELrZBcwmgUy_mtNL5xiWoq36v7_rnkrfLcz9bnRmwc"
    }
}
```

## REGISTER

Request : 
- Method : POST
- Endpoint : `/api/v1/register`
- Body :

```json
{
    "email": "reihanagam7@gmail.com",
    "phone": "089670558381",
    "password": "12345678",
    "app_name": "myapp"
}
```

Response : 

```json
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJpZCI6IjAzMWZjMWVmLWRmYzEtNDA3Ny05NTM3LTFmNmFjYWI3ZDJjNSJ9.BsjXvoHoEkoLVTY2-XLXJgYu8jfU5AG-aKxN5r-iN34"
    }
}
```

## VERIFY OTP 

Request : 
- Method : POST
- Endpoint : `/api/v1/verify-otp`
- Body :

```json
{
    "val": "reihanagam8@gmail.com",
    "otp": "SEEA"
}
```

Response : 

```json
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJpZCI6IjAzMWZjMWVmLWRmYzEtNDA3Ny05NTM3LTFmNmFjYWI3ZDJjNSJ9.BsjXvoHoEkoLVTY2-XLXJgYu8jfU5AG-aKxN5r-iN34"
    }
}
```

## RESEND OTP

Request : 
- Method : POST
- Endpoint : `/api/v1/resend-otp`
- Body :

```json
{
    "val": "reihanagam8@gmail.com",
}
```

Response : 

```json
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {
        "otp": "A8F1"
    }
}
```

## CONTENT ALL

Request : 
- Method : GET
- Endpoint : `/api/v1/content?page=1&limit=10&search&app_name=saka`

Response : 

```json
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "total": 2,
    "per_page": 1,
    "prev_page": 1,
    "next_page": 2,
    "current_page": 1,
    "next_url": "http://localhost:5004?page=2",
    "prev_url": "http://localhost:5004?page=1",
    "data": [
        {
            "id": "b4413807-b5f6-48a5-abda-15673622f2f1",
            "title": "Title 1",
            "desc": "Description 1",
            "medias": [],
            "likes": [
                {
                    "id": "0eca071f-4dd7-49a4-bf39-c774dae5f04e",
                    "user": {
                        "id": "8697faf2-c76e-4713-9514-e862fc0505f0",
                        "name": ""
                    }
                }
            ],
            "unlikes": [
                {
                    "id": "87d7322d-05e8-4e14-8ac6-2cc199562ebc",
                    "user": {
                        "id": "8697faf2-c76e-4713-9514-e862fc0505f0",
                        "name": ""
                    }
                }
            ],
            "comments": [
                {
                    "id": "4ce6dce2-00ee-4dd0-9562-7fa2bf933e6b",
                    "comment": "hello world",
                    "user": {
                        "user_id": "8697faf2-c76e-4713-9514-e862fc0505f0",
                        "name": ""
                    }
                }
            ],
            "app": {
                "id": "dcba8c60-62f4-416b-953a-f6da0a607862",
                "name": "saka"
            },
            "user": {
                "fullname": "",
                "email": "reihanagam7@gmail.com",
                "phone": "089670558381"
            },
            "type": "NEWS",
            "created_at": "2024-09-21 13:06"
        },
        {
            "id": "b4413807-b5f6-48a5-abda-15673622f2f2",
            "title": "Title 2",
            "desc": "Description 2",
            "medias": [],
            "likes": [],
            "unlikes": [],
            "comments": [],
            "app": {
                "id": "dcba8c60-62f4-416b-953a-f6da0a607862",
                "name": "saka"
            },
            "user": {
                "fullname": "",
                "email": "reihanagam7@gmail.com",
                "phone": "089670558381"
            },
            "type": "EVENT",
            "created_at": "2024-09-21 16:23"
        }
    ]
}
```

## CREATE CONTENT

Request : 
- Method : POST
- Endpoint : `/api/v1/content`
- Body :

```json
{
    "id": "b4413807-b5f6-48a5-abda-15673622f2f2",
    "title": "Title 2",
    "desc": "Description 2",
    "type_id": 2
}
```

Response : 

```json 
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {}
}
```

## LIKE CONTENT

Request : 
- Method : POST
- Endpoint : `/api/v1/content/like`
- Body :

```json
{
    "content_id": "b4413807-b5f6-48a5-abda-15673622f2f1"
}
```

Response : 

```json 
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {}
}
```

## UNLIKE CONTENT

Request : 
- Method : POST
- Endpoint : `/api/v1/content/unlike`
- Body :

```json
{
    "content_id": "b4413807-b5f6-48a5-abda-15673622f2f1"
}
```

Response : 

```json 
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {}
}
```

## CREATE COMMENT CONTENT

Request : 
- Method : POST
- Endpoint : `/api/v1/content/comment`
- Body :

```json
{
    "content_id": "b4413807-b5f6-48a5-abda-15673622f2f1",
    "comment": "hello world"
}
```

Response : 

```json 
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {}
}
```

## DELETE CONTENT

Request : 
- Method : DELETE
- Endpoint : `/api/v1/content/delete`
- Body :

```json
{
    "id": "b4413807-b5f6-48a5-abda-15673622f2f2",
}
```

Response : 

```json 
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {}
}
```

## DELETE COMMENT CONTENT

Request : 
- Method : DELETE
- Endpoint : `/api/v1/content/comment/delete`
- Body :

```json
{
    "id": "b4413807-b5f6-48a5-abda-15673622f2f2",
}
```

Response : 

```json 
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {}
}
```

## MEDIA UPLOAD

Request Form-Data :  
- Method : POST
- Endpoint : `/api/v1/media/upload`
- Body :

```json
{
    "file": "[type file]",
    "folder": "[type text]",
}
```

Response :

```json
{
    "status": 200,
    "error": false,
    "message": "Successfully",
    "data": {
        "path": "http://localhost:5004/test/fspmi.png",
        "filename": "fspmi.png",
        "size": 57466,
        "mime": "image/png"
    }
}
```

## MEMBERNEAR

Request Form-Data : 
- Method : GET
- Headers : APP_NAME 
- Endpoint : `/api/v1/membernear/all?lat=-6.176132&lng=106.822864`