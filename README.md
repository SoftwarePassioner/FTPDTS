# ftpdts
Ftpdts service do a real-time file generation and exposes them as downloadable FTP files. Files are generated from templates on the fly and can be filled with any user data. Ftpdts has the built-in FTP server based on ftpdt (goftp/server), and the simple webapi to put/get user data.

It was originally designed to implement the trick of escaping to the IOS default browser from Facebook, Instagram, webkit in-app browsers.
See real usage example: https://github.com/starshiptroopers/escfbb

Ftp is used as an intermediate gateway to download the html file with redirect to real web site.
The main reason the ftp is using in this trick is because the ftp protocol is handled with Safary by default. Opening a ftp link in a webkit lead to starting the Safari. At the moment of developing this library it was the only way to escape from Facebook browser to IOS default browser.

Server read the configuration from ftpdts.ini file, templates is stored at ./tmpl folder by default, persistent data storage is at ./data folder

##### Templates:
default.tmpl is the default template file. It used when the ftp client requests the file from the root folder, for example with url: ftp://server-name/UID.html
You can customize your templates and place them into templates folder with a different filename. 
For example, we place a template with name test.tmpl into the templates folder, 
then we can use the url ftp://server/test/UID.html to download the file created from this template and populated with UID data set

##### Datasets:
You can create your own dataset and place the file into the ./data folder. Dataset file name must be in UID format and file must contain a JSON
Also you can post dataset directly to the webAPI endpoint. It will be stored into the memory cache or persistent storage.

##### WebAPI endpoints:
```
POST:
  url: /data?ttl=n
  ttl = time to live in seconds the data will be stored in the memory storage, if ttl = 0 data will be stored into the persistent storage, if ttl is not defined data will be stored into the memory storage with default ttl
  body: data to fill into the template in JSON format
  response:
  	{
	   "code": 0,    		// error code
    	   "message": "OK",		// error message
 	   "uid": "xxxxxxxxxxxxxxxxxxxxxxxx"  // uid the data was stored with
  	}

 GET:
  url: /data?uid=xxxxxxx...xx
  response:
  	{
	    "code": 0,
	    "message": "OK",
	    "data": {
		....
	    },
    	    "createdAt": datetime,
    	    "ttl": 0
	}
```

##### Usage example:

    1. Start the service: docker-compose up
    2. Do the POST request to http://localhost:2000/data with curl
       curl --header "Content-Type: application/json" --request POST --data '{"Title":"Redirect page","Caption":"This is a redirection page","Url": "https://starshiptroopers.dev"}' http://localhost:2000/data
       You will get the data {"code":0,"message":"OK","uid":"Xmnw48xJKpolFYwLn7a0wetEdsTKym1M"}
    3. Use UID to get the file with curl
		   curl ftp://localhost/Xmnw48xJKpolFYwLn7a0wetEdsTKym1M.html
		 You will get the html page with content:
```
      <!DOCTYPE html>
			<html lang="en">
			<head><title>Redirect page</title></head>
			<body><h1>This is a redirection page</h1> <script> window.location.href = "https://starshiptroopers.dev" </script> </body>
			</html>
```

