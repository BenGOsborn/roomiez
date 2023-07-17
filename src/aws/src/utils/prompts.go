package utils

import (
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

// Create a prompt for a post validation
func NewPostValidation(llm *openai.Chat) *chains.LLMChain {
	prompt := prompts.NewPromptTemplate(`Please return 'yes' if the following post is advertising a shared rental property, and 'no' if otherwise. Some examples have been provided below to give you context.
	For example, if the post is advertising a shared rental property the output should be yes.
	If a post is advertising a rental but it is not a shared property the output should be no.
	If the post indicates someone is looking for a shared rental property the output should be no.
	Your response should be either 'yes' or 'no' on a new line after "Output:".

	Examples:

	Post:
	Private room in Earlwood, Sydney for $200/week
	• House is only 1km walk to Canterbury train station which takes you to the central in 20 minutes by train
	• Less than 30 minutes from city centre, beautiful beaches, universities, and airport.
	• Plenty of street parking available for free if you have a car.
	• Share bathroom, lounge and kitchen with a mum and daughter only
	• Might suit Females / students / retirees / backpackers / 40+
	Output:
	yes

	Post:
	I'm Jane Doe, 28 from the US. Looking for a place in the areas of Bondi/Bondi Junction/Coogee/Bronte/Woollahra and surrounding. 
	Very easy going, love the beach and into my fitness. Love to cook and have a drink over dinner.
	Preferably a furnished place but willing to furnish the bedroom and have a weekly budget of $400incl bills - willing to pay more depending on property etc.
	Available to move in from 19th July🙂
	Output:
	no

	Post:
	{{.post}}
	Output:
	`, []string{"post"})

	return chains.NewLLMChain(llm, prompt)
}

// Create a prompt for post data extraction
func NewPostExtraction(llm *openai.Chat) *chains.LLMChain {
	prompt := prompts.NewPromptTemplate(`Given the following post, please extract the information and format it as a JSON object that strictly follows the schema described below. Ensure that each extracted field corresponds to its respective field in the schema.
	For example, if the post contains information about the price, assign the extracted price to the "price" field in the JSON object.
	It is very important that the JSON generated matches the schema exactly.
	Note that 'Young' is classified as anyone between 18 and 35 years old, 'Middle Aged' is anyone between 36 and 55 years old, and 'Old' is anyone older than 56.
	In addition, the price should ALWAYS be the weekly price.
	If there is any ambiguity and the field supports null, it is best to put null.
	Your response should be a JSON string on a new line after "Output:".

	JSON Schema:
	type RentalSchema = {
		price: number | null;
		bond: number | null;
		location: string | null;
		rentalType: "Apartment" | "House" | null;
		gender: "Male" | "Female" | null;
		age: "Young" | "Middle Aged" | "Old" | null;
		duration: "Short Term" | "Long Term" | null;
		tenant: "Singles" | "Couples" | null;
		features: ("Garage" | "WiFi" | "Bills Included" | "Furnished" | "Pets Allowed" | "Mattress" | "Pool" | "Gym")[];
	}

	Some examples have been provided below to give you context.

	Examples:

	Post:
	LIVE BY THE BEACH, Maroubra
	Short term room for rent, $300/week includ bills, bond $1080
	Looking for a girl to live with us from 
	Saturday July 15th - August 26th 
	$300/week including bills (electricity, gas and wifi) 
	Furnished bedroom with queen bed 
	House fully furnished
	You'll be living with 2 working girls in their early 30s
	Bus to city outside of the door - 396, 396x, 397 
	10 min walk to Maroubra Beach 
	5 min walk to Maroubra Junction - Coles, Aldi, Chemist Warehouse other buses to Bondi Junction, Airport 
	Free street parking, can always find parking
	DM for more details!	
	Output:
	{
		"price": 300,
		"bond": 1080,
		"location": "Maroubra",
		"rentalType": "House",
		"gender": "Female",
		"age": "Young",
		"duration": "Short Term",
		"tenant": "Singles",
		"features": [
			"WiFi",
			"Bills Included",
			"Furnished"
		]
	}

	Post:
	{{.post}}
	Output:
	`, []string{"post"})

	return chains.NewLLMChain(llm, prompt)
}

// Create a prompt for post summarization
func NewPostDescription(llm *openai.Chat) *chains.LLMChain {
	prompt := prompts.NewPromptTemplate(`Given the following post, please summarize it into a very short paragraph (less than 60 words). You will return
	Make sure to include the price, bond, and location if they exist, as well as any other important information mentioned in the post (e.g. age preference, gender preference, whether it is a house or apartment).
	Only include the relevant information about the rental, don't include anything along the lines of "DM for details" or "Message owner".
	Your response should be a summary string on a new line after "Summary:".

	Some examples have been provided below to give you context.

	Examples:

	Post:
	LIVE BY THE BEACH, Maroubra
	Short term room for rent, $300/week includ bills, bond $1080
	Looking for a girl to live with us from 
	Saturday July 15th - August 26th 
	$300/week including bills (electricity, gas and wifi) 
	Furnished bedroom with queen bed 
	House fully furnished
	You'll be living with 2 working girls in their early 30s
	Bus to city outside of the door - 396, 396x, 397 
	10 min walk to Maroubra Beach 
	5 min walk to Maroubra Junction - Coles, Aldi, Chemist Warehouse other buses to Bondi Junction, Airport 
	Free street parking, can always find parking
	DM for more details!	
	Summary:
	Short term room for rent in Maroubra. Price is $300/week and bond is $1080. Looking for a girl to live with two working girls in their early 30s.
	Bus to city outside, 10 min walk to Maroubra Beach. 5 min walk to Maroubra Junction. Free street parking.

	Post:
	{{.post}}
	Summary:
	`, []string{"post"})

	return chains.NewLLMChain(llm, prompt)
}
