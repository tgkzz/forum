
# Forum

This project consists in creating a web forum that allows :

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

## Run Locally

run the project using next command

```bash
  go run .
```

or use make file

```bash
  make run
```

or use docker

```bash
  make build
  make docker-run
```

in order to stop container use next command

```bash
  make stop
```

in both ways server will run on the next route

```
http://localhost:4000/
```




## ERD Diagram

![ERD](https://tinypic.host/images/2023/12/14/Untitled-Diagram.drawio1.png)
## API Reference

#### Get all items

```http
  GET /
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| none | none | displays home page and all availale posts |

#### signup

```http
  GET /signup
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `username` | string | username you want on cinema forum (must be unique) |
| `email` | string | email you want (required to be normal) |
| `password` | string | password you will be logging |


#### signin

```http
  GET /signup
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `username` | string | username you have taken  |
| `password` | string | password you have put while signining up |



#### Get post information

```http
  GET /posts/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | the id`s post |

#### Get your posts

```http
  GET /myposts
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| none      | none | nothing |


#### Get posts you have liked

```http
  GET /liked
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| none      | none | nothing |



#### leave the comment under the post

```http
  POST /posts/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | displays information about post |
| `session`      | `string` | session you are currently using |
| `text`      | `string` | text under the comment |

and there`re many other routes i wanted to describe, but there are too many of them, so it's better to use swagger documentation, so routes documentation would be automated


## Run Locally

run the project using next command

```bash
  go run ./cmd
```

or use make file

```bash
  make simple
```

or use docker

```bash
  make build
  make run
```

in order to stop container use next command

```bash
  make stop
```

in both ways server will run on the next route

```
http://localhost:4000/
```
