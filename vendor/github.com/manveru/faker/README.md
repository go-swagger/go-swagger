# Faker for Go

## Usage

    package main

    import (
      "github.com/manveru/faker"
    )

    func main() {
      fake := faker.New("en")
      println(fake.Name())  //> "Adriana Crona"
      println(fake.Email()) //> charity.brown@fritschbotsford.biz
    }

Inspired by the ruby faker gem, which is a port of the Perl Data::Faker library.
