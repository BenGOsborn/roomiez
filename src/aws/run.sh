cd src

go build retrieve_rentals/main.go && zip ../retrieve_rentals.zip main && rm main
go build process_rental/main.go && zip ../process_rental.zip main && rm main
go build get_fields/main.go && zip ../get_fields.zip main && rm main
go build subscribe/main.go && zip ../subscribe.zip main && rm main
go build unsubscribe/main.go && zip ../unsubscribe.zip main && rm main

cd ..

terraform apply

rm *.zip