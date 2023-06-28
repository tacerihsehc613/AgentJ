# SaramIn Scraper  <img src="https://github.com/tacerihsehc613/AgentJ/assets/127294863/fdb1a6d5-7435-44b4-9c6f-b6fa07dff2f7" width="50" height="auto"> 

<img src="https://github.com/tacerihsehc613/AgentJ/assets/127294863/aba74c3f-84a9-405b-9923-d4efd6fdf3d1" width="400" height="auto">

Build a Golang scraping tool to retrieve job data of a specific stack from the SaramIn website.

1. The tool spawns the main goroutines to scrape each page concurrently.
2. Job details of every job card on each page are extracted through additional goroutines.
3.  This second goroutine finds and extracts the job title, company name, location, and desc from the card.
4. Writes the scraped job data to a CSV file named 'jobs.csv'
5. After the file is sent to the client as an attachment, it is discarded from the server.
