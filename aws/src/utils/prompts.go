package utils

import (
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

// Create a prompt for a post validation
func NewPostValidation(llm *openai.LLM) *chains.LLMChain {
	prompt := prompts.NewPromptTemplate(`Please return 'yes' if the following post is advertising a shared rental property, and 'no' if otherwise.
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
	1 bed 1 bathroom – Flat
	Coogee, Coogee, NSW, Australia
	Hi everyone, 
	I will be moving to Sydney on the 21st of  October and I’m looking for somewhere to rent. Open to short or long term at the moment and preferably in the Eastern suburbs. 
	Thank you!
	Output:
	no

	Post:
	Furnished Private Rooms (Chippendale)
	Incredible location. 4 mins to Broadway Shopping Center. Walk everywhere.
	- Furnished Private Room $315 pw including bills available 06 July
	- Furnished Ensuite Room (private room with own attached bathroom) $375 pw including bills available 20 July
	The House details: spacious and clean kitchen/dining area/ laundry in good condition:
	- Large balcony
	- Fully furnished kitchen
	- Wardrobes
	- Microwave
	- Washing machine
	- Clothes Dryer Machine
	- Unlimited WIFI internet
	- Fully furnished
	Location by walk:
	- Bus stop (200 meters, 2 mins)
	- Broadway Shopping Centre - Supermarkets Coles & Aldi Broadway Shopping (400 meters, 4 mins)
	- Central Train Station (800 meters, 10 mins) – can go anywhere in Sydney from here.
	No Party
	No Smoker
	Output:
	yes

	Post:
	This is a SHORT-TERM rental from 22 July until 26 August. Perfect place if you want to live in a beautiful and quiet area, close to the water and just 10 walking minutes away from the city.
	The apartment is very spacious and has two floors. We have 3 bedrooms on the upper level and one bathroom that our two flatmates are sharing. Your bedroom has its own bathroom, a big window, a queen bed, two desks incl monitors and large built-in wardrobes with mirrors. 
	Downstairs is the living area with a sofa, a dining table and a fully equipped kitchen. We also have a laundry room and a bathroom downstairs. We have our own entrance and a little outdoor area in front. The main building also has a gym, a sauna and an outdoor pool.
	Coles and Woolworth are less than 10 minutes away. There are also lots of good restaurants, cafés and pubs around.
	A bus stop and the light rail are just 2 minutes away.
	Output:
	Yes

	Post:
	I’m Jane Doe, 28 from the US. Looking for a place in the areas of Bondi/Bondi Junction/Coogee/Bronte/Woollahra and surrounding. 
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
