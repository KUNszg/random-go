
## Installation and building

- clone this repository using
```bash
$ git clone https://github.com/KUNszg/random-go.git
```
- Next you will need to create an account and generate a random.org [API key](https://api.random.org/dashboard) and copy it
- cd into the cloned repository
- paste the API key into config.txt file
- build the docker image
```bash
$ cd .. && sudo docker build random-go
```
- replace <image_id> below with the generated docker image ID 
```bash
$ sudo docker run -p 8080:8080 <image_id>
```
## Usage

- send a GET request to ```$ curl "http://localhost:8080/random/mean?length=<l>&requests=<r>"```, replace \<l> with amount of random values you want to receive and \<r> with number of requests

- example:

```$ curl "http://localhost:8080/random/mean?length=5&requests=4"```

```json
{
   "0":{
      "stddev":3.84708,
      "data":[
         9,
         0,
         3,
         8,
         10
      ]
   },
   "1":{
      "stddev":2.87054,
      "data":[
         9,
         4,
         8,
         5,
         1
      ]
   },
   "2":{
      "stddev":3.18748,
      "data":[
         9,
         9,
         4,
         8,
         1
      ]
   },
   "3":{
      "stddev":2.99333,
      "data":[
         0,
         4,
         1,
         6,
         8
      ]
   },
   "4":{
      "stddev":3.38046,
      "data":[
         0,
         4,
         1,
         6,
         8,
         9,
         9,
         4,
         8,
         1,
         9,
         4,
         8,
         5,
         1,
         9,
         0,
         3,
         8,
         10
      ]
   }
}
```
