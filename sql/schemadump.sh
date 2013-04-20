pg_dump -h localhost -p 5432 -U askbitcoin -W --schema-only -n askbitcoin projects

# -h for host
# -p for port
# -U username
# -W to prompt password
# --schema-only to avoid dumping data
# -n to specify the schema name
# -C to produce create-database statements, not just create-schema statements
# the final argument is the database name (a database can have many schemas; a schema can have many tables)
