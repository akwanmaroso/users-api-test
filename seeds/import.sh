#! /bin/bash

mongoimport --host mongodb --db users --collection users --type json --file /seeds/users.json --jsonArray