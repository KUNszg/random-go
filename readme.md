
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
$ cd .. && docker build random 
```
- replace <image_id> below with the generated docker image ID 
```bash
$ docker run -p 8080:8080 <image_id>
```
## Usage

- send a GET request ```$ curl localhost:8080/random/mean?length=<l>&requests=<r>``` in your browser, replace \<l> with amount of random values you want to receive and \<r> with number of requests
