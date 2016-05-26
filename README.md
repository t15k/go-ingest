# go-ingest
A high performant HTTP endpoint for data ingestion.


    cd mixpaneld
    nc -l localhost 4343 &
    go run mixpaneld.go -socketout,localhost:4343
    # then from another terminal
    curl "localhost:8080/track/?data=ew0KICAgICJldmVudCI6ICJTaWduZWQgVXAiLA0KICAgICJwcm9wZXJ0aWVzIjogew0KICAgICAgICAiZGlzdGluY3RfaWQiOiAiMTM3OTMiLA0KICAgICAgICAidG9rZW4iOiAiZTNiYzQxMDAzMzBjMzU3MjI3NDBmYjhjNmY1YWJkZGMiLA0KICAgICAgICAiUmVmZXJyZWQgQnkiOiAiRnJpZW5kIg0KICAgIH0NCn0="
