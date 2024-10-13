# Gladnytt.net

Gladnytt.net is a news website that uses sentiment analysis to filter out negative news articles. The articles are fetched from NRK's RSS feed every minute, and sentiment is determined by GPT-4o. The app is built with Go and htmx.

![Screenshot](https://github.com/user-attachments/assets/5d443004-4aa1-4c63-b5ae-3f4c62586f90)

## Running locally

1. Install [Go](https://go.dev/doc/install)
2. Install [Air](https://github.com/air-verse/air) (live reloading for Go)
3. Install [pnpm](https://pnpm.io/installation) (used for Tailwind CSS and Prettier)
4. Clone the repository
5. Copy `.env.example` to `.env` and add your OpenAI API key
6. Run `air` from the project root to start the app

The app will be available at http://localhost:3000 and will automatically reload when changes are made.
