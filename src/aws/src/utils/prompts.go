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
	If the post indicates someone is looking for a shared rental property the output should be no.
	Your response should be either 'yes' or 'no' on a new line after "Output:".

	Examples:

	Post:
	Private room in Earlwood, Sydney for $200/week
	• House is only 1km walk to Canterbury train station which takes you to the central in 20 minutes by train
	• Less than 30 minutes from city centre, beautiful beaches, universities, and airport.
	• Plenty of street parking available for free if you have a car.
	• Has a big but low maintenance cemented backyard with privacy on quiet street.
	• Ground floor so you don't have to climb stairs.
	• Well lit, airy and sunny home.
	• Big lounge and spacious kitchen with more than enough shelves to store your pantry items.
	• High ceilings with big wardrobe
	• Double brick home with good insulation.
	• Share bathroom, lounge and kitchen with a mum and daughter only
	• Might suit Females / students / retirees / backpackers / 40+
	Output:
	yes

	Post:
	I'm Jane Doe, 28 from the US. Looking for a place in the areas of Bondi/Bondi Junction/Coogee/Bronte/Woollahra and surrounding. 
	Very easy going, love the beach and into my fitness. Love to cook and have a drink over dinner. Enjoy down time and own space but always happy to socialise and respectful of living space.
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
		price?: number | null;
		bond?: number | null;
		location?: string | null;
		rentalType?: "Apartment" | "House" | null;
		gender?: "Male" | "Female" | null;
		age?: "Young" | "Middle Aged" | "Old" | null;
		duration?: "Short Term" | "Long Term" | null;
		tenant?: "Singles" | "Couples" | null;
		features?: ("Garage" | "WiFi" | "Bills Included" | "Furnished" | "Pets Allowed" | "Garage" | "Mattress" | "Pool" | "Gym")[];
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
	📍Sydney Olympic Park
	🎡🎡🎡MASTER ROOM AVAILABLE 🎡🎡🎡
	Negotiable for single
	👌all bills included 
	🛜Internet✅
	🛁Own bathroom✅
	2 Wardrobe✅ 
	Comfortable new Bed✅
	study desk
	🛒 Close to IGA
	🚏Bus stop 2 min
	🚊Train station 6 min
	🔑Security
	📍Minimum stay 3-6 months
	👍2weeks bond 
	👍🏻2 weeks rent in advance 
	🌲🌲🌲  Beautiful park to walk and do barbecue no need to go far away to have a picnic
	🎓study room, 
	😀 very clean, nice and calm environment 
	🏡 private balcony 
	🚭No smoke 
	🍾No party
	🐶No pet
	📍the room has a great view, quiet, comfy and privacy.
	Please kindly send me chat 480$ per week included bills.💕
	Output:
	{
		"price": 480,
		"bond": null,
		"location": "Sydney Olympic Park",
		"rentalType": "Apartment",
		"gender": null,
		"age": null,
		"duration": "Long Term",
		"tenant": "Singles",
		"features": [
			"Bills Included",
			"WiFi",
			"Mattress"
		]
	}	

	Post:
	{{.post}}
	Output:
	`, []string{"post"})

	return chains.NewLLMChain(llm, prompt)
}
