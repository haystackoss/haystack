# Haystack - Workplace Search Engine for Developers

Are you tired of sifting through multiple communication channels and documents in your workplace, trying to find that one piece of information you need. 
Look no further than Haystack! 

Haystack is a search engine that allows you to search all of your workplace applications from one place.

![Alternate image text](https://raw.githubusercontent.com/haystackoss/haystack/main/asknatural.png)

### Integrations
- [x] Slack
- [x] Confluence
- [x] Notion
- [x] Jira
- [x] Github Projects
- [x] Email

### Natural Langauge
it also allows you to search using natural language.
such as `"How to do X"`, `"how to connect to Y"`, `"Do we support Z"`

### Features
- Lets you go directly to the relevant paragraph in the search result.
- Adds additional information to search results to make them easier to evaluate.
- It's all done in the browser, with the option to store results locally for added security.


## Under the hood

### Storage
Haystack uses IndexDB for storing result indices and NLP models

### Permissions
Sets up read permissions for workplace apps with secure token storage.

### Indexing
Indexes each document, message, and email, generating vector embeddings with a fine-tuned mini BERT based bi-encoder for quick searches.

### Searching
When a query is entered, it is converted into a vector embedding and compared to the most relevant embeddings, with the top results being reranked using a cross encoder. A natural language summary of the top 3 results is then generated based on the original matched paragraph and user query.

## Next steps
We are currently fine-tuning Haystack for a 9th gen i5 with no dGPU, and rolling it out slowly to developers who can get the most value from it. If you want to get early access, click the button on our landing page to fill out a quick 3-question form. 
Don't miss out on the chance to revolutionize your workplace search experience with Haystack!
