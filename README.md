# seederella
seederella is a powerful database seeding tool that effortlessly brings your
database schema to life with realistic, customizable data. Whether you’re
testing, developing, or demoing, Seederella makes populating your database with
mock data a breeze.

seederella is schema-agnostic, relationship-aware, and driven by simple
configuration files, allowing you to focus on building your application while
it handles the complexities of inserting data into any SQL database.

## Features
- Schema-Agnostic: Compatible with any SQL database, including PostgreSQL, MySQL, SQLite, and more.
- Faker-Driven: Automatically generates realistic, random data using the faker library for various data types like names, emails, dates, and more.
- Relationship-Aware: Handles foreign key constraints and relationships between tables seamlessly.
- Config-Based: Define table/column behaviors using simple JSON or YAML configuration files.
- CLI Tool: Integrate it easily into your CI/CD pipelines or local scripts for quick and hassle-free data bootstrapping.

## Example Config (YAML)
```yaml
driver: postgres
dns: postgresql://postgres:password@127.0.0.1:5432/mydb?sslmode=disable
schema: public

tables:
  users:
    count: 10
    fields:
      name:
        faker: name
      email:
        faker: email
        transform: lower
      role:
        value: "admin"
      created_at:
        faker: date_past

  posts:
    count: 30
    fields:
      title:
        faker: sentence
      content:
        faker: paragraph
      author_id:
        reference: users.id
      modified_by_id:
        same_as: author_id
```

## Why seederella?
Every application needs mock data, but real-world schemas are often messy with
complex relationships, constraints, and varying column needs. seederella makes
managing this complexity effortless.

## Installation
```sh
# Install using Go
go install github.com/QaidVoid/seederella/cmd/seederella@latest
```

## Usage
```sh
seederella - a database seeding tool that brings your schema to life effortlessly and efficiently

Usage:
  seederella [flags]

Flags:
      --clean           Drop and recreate DB schema before seeding
  -c, --config string   Path to config file (default "config.yaml")
      --driver string   Override DB driver
      --dsn string      Override DB DSN
  -h, --help            help for seederella
```

## Contributing
Feel free to fork the repository and open pull requests. Contributions are
always welcome! If you find any bugs or have suggestions, please open an issue
on GitHub.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
