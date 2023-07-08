go build retrieve_rentals/main.go && zip retrieve_rentals.zip main && rm main
go build process_rental/main.go && zip process_rental.zip main && rm main

terraform apply

rm *.zip