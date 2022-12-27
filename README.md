# Haystack - Workplace Search Engine for Developers

Are you tired of sifting through multiple communication channels and documents in your workplace, trying to find that one piece of information you need. 
Look no further than Haystack! 

![Alternate image text](https://raw.githubusercontent.com/haystackoss/haystack/main/whitefade.svg)

## I'm listening?
Haystack allows you to search all of your workplace apps from a single place.

![Alternate image text](https://raw.githubusercontent.com/haystackoss/haystack/main/asknatural.png)


### Integrations
- [x] Slack
- [x] Confluence
- [x] Notion
- [x] Jira
- [x] Github Projects
- [x] Email

### Natural Langauge
allows searching using natural language.
such as `"How to do X"`, `"how to connect to Y"`, `"Do we support Z"`

### But wait, there's more!
- Go to the relevant matched paragraph directly from the search result.
- Search results are enriched with a summary of matched content and it's relevancy to the query to make it easier for the user to evaluate without entering the page.
- The whole thing happens in the browser, indices are stored locally for added security.


## Under the hood

### Storage
Haystack uses IndexDB for storing result indices and NLP models

### Permissions
Sets up read permissions for workplace apps and stores 3rd party tokens in secure local browser storage.

### Indexing
Indexes each document, message, and email, generates vector embeddings using a fine-tuned TinyBERT based bi-encoder.

### Searching
Each query is converted into a compressed vector embedding (like in the index process) and compared to the most relevant embeddings, with the top 100 results being reranked using a t5-small cross encoder for accuarcy. 
A natural language summary of the top 3 results is then generated based on the original matched paragraph and the users query.

## Next steps
We are currently fine-tuning Haystack for lower end hardware, specifically laptops with no dedicated graphics. 

Meanwhile we are rolling haystack out to developers we know well, or those who show particular interest. 

### Get early access 
Explain why you need this [here](https://m8i3t3b9dp5.typeform.com/to/q2zPGfOU) and we'll send you an invite, otherwise we would be releasing beta invites to a wider audience in a few months anyways.

Cheers!
