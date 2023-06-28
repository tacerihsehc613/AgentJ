# SaramIn Scraper  <img src="https://github.com/tacerihsehc613/AgentJ/assets/127294863/ea206719-2488-4302-95dc-e07783af761e" width="50" height="auto"> 

<img src="https://github.com/tacerihsehc613/AgentJ/assets/127294863/51c45bec-aa22-4e04-af18-79d68d65c94d" width="400" height="auto">

Build a Golang scraping tool to retrieve job data of a specific stack from the SaramIn website.

1. The tool spawns the main goroutines to scrape each page concurrently.
2. Job details of every job card on each page are extracted through additional goroutines.
3.  This second goroutine finds and extracts the job title, company name, location, and desc from the card.
4. Writes the scraped job data to a CSV file named 'jobs.csv'
5. After the file is sent to the client as an attachment, it is discarded from the server.
