# Replicate Floor Plan

Experiment with the [Replicate API](https://replicate.com) to send images of real estate floor plans and extract quantifiable data on them using AI.

This was just a prototype using LLaVA. It is not effective as the model needs training to work correctly.

## Running

- `cp cmd/app/.env.example cmd/app/.env`
- Fill in env variables
- `cd cmd/app && go run .`
