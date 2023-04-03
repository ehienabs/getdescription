# How to run this app locally

First clone the repository

```go
git clone https://github.com/ehienabs/getdescription.git
```

To run locally without building  a binary, from the root directory run the following command

```go
go run main.go
```

Additionally, you can build a binary and run it. To build a binary, from the root directory, run the following command

```go
go build -o <name_of_binary>
```

Run your binary using the command below

```bash
./<name_of_binary>
```

## Run API as a container

To run the `getdescription` API as a container, you can build the image and run it using the following commands:

To build the image, run the following command in the root directory of the project.

```docker
docker build -t <name_of_image> .
```

Run the container using your built image

```docker
docker run -p 8080:8080 <name_of_image>
```

You can also run the app as a detached container

```docker
docker run -d -p 8080:8080 <name_of_image>
```

Additionally, pull and run the image from Docker remote registry using the following command

```docker
docker run -p 8080:8080 ehienabs/getdescription:v1
```

# How does it work

The application has two endpoints

- `/home` endpoint which prints the following message.

<aside>
💡 "Add /description?name=<name>&lang=<lang> in your browser, such that name and language are query parameters”

</aside>

- `/description` endpoint which allows you to retrieve a short description of a name from the Wikipedia API using a standard query.
    - The endpoint accepts a query parameter of `?name=<name>&lang=<lang>` where
        - `name` is the name of the person for whom you wish to retrieve a description. First, middle, and last names can be separated by a space.
        - `lang` is an optional parameter that you can use to specify a language.
- The result is a JSON payload with the following schema. Where `name` is the person’s name and `desc` is the short description.

```json
{
   "desc":"\n\n\nKristen Mary Jenner (née  Houghton  HOH-tən, formerly Kardashian; born November 5, 1955) is an American media personality, socialite, and businesswoman.",
   "name":"kris jenner"
}
```