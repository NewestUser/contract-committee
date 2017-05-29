# Contract Committee

A service aiming to detect backwards incompatible changes in API's.  
Heavily rely on request/response assertions.


## Example 1
```json
{
	"strict":"true",
	"request":{
		"endpoint":"http://srvice-id/accounts",
		"method":"post",
		"headers":[
			{"key":"content-type", "value":"application-json"}
		],
		"body":{
			"name" : "{{ $name:=rndString | save }}",
			"status" : "ENABLED",
			"credentials": {
				"customerNumber":"{{ $custNum := rndNum | save }}",
				"projectId":"{{ $rndNum(2,5) }}",
				"apiKey":"{{ $rndAlphaNum(6,7) }}"
			}	
		}
	},
	"expectation":{
		"status":201,
		"headers":[""],
		"body":{
			"id":"$isReturned",
			"name":"$isSame",
			"status":"$isSame",
			"credentials":{
				"customerNumber":"$isSame",
				"projectId":"$isSame",
				"apiKey":"$isSame"
			},
			"createdBy":"$isReturned"
		}
	}
}

```


## Example 2
```json

{
	"strict":"true",
	"request":{
		"endpoint":"http://srvice-id/accounts/$pop(body.id)",
		"metohd":"put",
		"method":[
			{"key":"content-type", "value":"application-json"}
		],
		"body":{
			"name" : "{{$pop(body.name)}}",
			"status" : "DISABLED",
			"credentials": {
				"customerNumber":"{{$pop(body.customerNumber)}}",
				"projectId":"{{$pop(body.projectId)}}",
				"apiKey":"{{$pop(body.apiKey)}}"
			}	
		}
	},
	"expectation":{
		"status":201,
		"headers":[""],
		"body":{
			"id":"$isReturned",
			"name":"$isSame",
			"status":"$isSame",
			"credentials":{
				"customerNumber":"$isSame",
				"projectId":"$isSame",
				"apiKey":"$isSame"
			},
			"createdBy":"$isReturned"
		}
	}
}

```