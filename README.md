# pragmatic.cooking

A web app to help you keep track of your favorite recipes.

<img width="1910" height="1033" alt="image" src="https://github.com/user-attachments/assets/f19146f7-c231-4d90-b81d-bde486498999" />
<img width="1910" height="615" alt="image" src="https://github.com/user-attachments/assets/ca62d9c4-1511-417a-85c1-589660294aec" />
<img width="1910" height="1327" alt="image" src="https://github.com/user-attachments/assets/82270e12-54fb-4fe5-a40f-7bc3d3085a32" />

## Why

While I always enjoy cooking, it can be a chore to figure out what to cook.  I created this project for the Boot.Dev 2025 Hackathon, but I am hoping that it will help me keep track of what I enjoy cooking. 

## View the Project

I am hosting the project at https://app.pragmatic.cooking/

It is fully functional there, and is the easiest way to experience the project.

## About the Project

### Key Features

- Knowing what ingredients are needed for a recipe
- Knowing what you can make with each ingredient
- Links to the original recipe

### Technical Features:

- User authentication using magic links
- Sending magic links to users using the Mailgun API
- Storing JWT tokens in cookies for later use
- Storing content in a Postgres database
- Listing ingredients and recipes from the database using templ
- Adding ingredients and recipes to the database and refreshing the list using HTMX
- Deleting ingredients and recipes from the database and refreshing the list using HTMX
- Preventing users from deleting other users' content
- Preventing ingredients in use by a recipe from being deleted
- A list of recipes owned by the current user displayed on their account page

### Future Features

- Options for each ingredient, including price and retailer
- User ingredient option preferences
- Using the Kroger API to get current ingredient prices
- Cost estimates for each recipe, using user ingredient option preferences
- Bookmarked Recipes list on User's Account Page
- Shopping list
- Weekly meal plan

# Manual Setup

## Prerequisites

To run this project, you will need:
1. A PostgreSQL database (this can be on the same server as the backend, if you would like)
2. A backend server with Go and templ installed.

## Get the repo

1. Clone the repo `git clone https://github.com/Rodabaugh/pragmatic-cooking/`
2. Navigate to the program dir `pragmatic-cooking`

## Configuration

Environment variables are used for configuration. Create a .env file in the root of the project dir. You need to specify your DB_URL, PLATFORM, MG_API_KEY, and JWT_SECRET. DB_URL is the url for your database. Platform can either be "prod" or "dev". You can generate your JWT_SECRET with `openssl rand -base64 64`. Your `.env` file should look something like the one below. Please be sure to create your own JWT_SECRET and use your own DB_URL.
```
PLATFORM=prod
DB_URL="postgres://postgresUser:postgresPass@localhost:5432/pragmatic_cooking?sslmode=disable"
JWT_SECRET="ALItvAPa64TLZ4wjqWsaiVW3ZrQ7ZT209sAkIsos8K3p6ldeMb+K5Ji5j90kI4cQ
k0I6WY6KgXALHP7EjeLXOw=="
MG_API_KEY="KEY HERE"
```

A port may also be specified using ```PORT=1234```. If a port is not specified, it will default to 8080.

## Setting up the database

Goose is used to manage the database migrations. Install goose with `go install github.com/pressly/goose/v3/cmd/goose@latest`

Navigate to the sql/schema dir `cd sql/schema`

Setup the database using goose `goose postgres <connection_string> up` e.g `goose postgres postgres://postgresUser:postgresPass@localhost:5432/pragmatic_cooking?sslmode=disable up`

## Compile and run the backend

Once your .env has been configured, and your database is setup, it is time to build and run the backend.

Build the application with `make build`

Run the backend application with `./pragmatic-cooking`

Once the backend server is running, you can setup your server to run the application as a service. 

## Success

At this point, the pragmatic.cooking application should be running on your server.

