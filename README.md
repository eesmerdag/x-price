This is a simple application to show how to show one product's prices fetched from different sources.

As a solution:

1. i just followed the algorithm keeping the prices in memory and updating them after every minute with lazy loading.
2. And using concurrency is a good to be faster not to make one service to wait for others. That also optimize CPU
   usage.

This repo has two different projects to run:

* price-api to provide sample test data for prices
* web-app to have simple web app page with prices info

Please read each projects readme.md files to see how to run.