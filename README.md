# mycart
> API docs [_here_]().

## Table of Contents
- [mycart](#mycart)
  - [Table of Contents](#table-of-contents)
  - [General Information](#general-information)
    - [Tech used](#tech-used)
    - [Installation](#installation)
      - [Using Git](#using-git)
    - [Running tests](#running-tests)
    - [Setting up environments](#setting-up-environments)
    - [Usage](#usage)
  - [Project Status](#project-status)
  - [Features](#features)
  - [Contact](#contact)

## General Information
- An e-commerce cart-centric REST API the leverages redis for cart.
- This project aims to reduce cart abandonment rate by reducing time to add and remove items from cart.

### Tech used
**Database**
- [x] PostgreSQL

**ORM**
- [x] GORM

**Language**
- [x] Golang
  
**Router**
- [x] Chi
  
**Queue**
- [x] Redis
  
### Installation
#### Using Git
1. Navigate & open CLI into the directory where you want to put this project & clone this project using this command.
   
```bash
git clone https://github.com/Adedunmol/mycart.git
```
2. Run `go get` to install all dependencies

### Running tests
* Run `go test ./... -v` to run unit tests.


### Setting up environments
1. There is a file named `app.env.example` on the root directory of the project
2. Create a new file by copying & pasting the file on the root directory & rename it to just `app.env`
3. The `app.env` file is already ignored, so your credentials inside it won't be committed
4. Change the values of the file. Make changes of comment to the `app.env.example` file while adding new constants to the `app.env` file.

### Usage
* Run `air` to start the application.
* Connect to the API using Postman on port 5000.


## Project Status
Project is: _in progress_.

## Features
- Creating cart in Redis for users.
- Migrate data from Redis to PostgreSQL when TTL for cart expires.
- Migrate data from PostgreSQL to Redis when needed.
- Generation and sending of receipt for users on purchase


## Contact
Created by [Adedunmola](mailto:oyewaleadedunmola@gmail.com) - feel free to contact me!