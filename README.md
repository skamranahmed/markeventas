# MarkEventAs

![header-image](https://user-images.githubusercontent.com/43776315/163416945-3bf9d38e-110c-42ef-836c-a9626a629f70.png)

🗓 Create Google Calendar Events from Tweet.

You just need to reply `@markeventas <event_name> | <date> | <time> | <timezone>` to any tweet.

## Demo

## Why I built this?
### 😩 The problem
- There were a lot of events happening on twitter specially twitter spaces and I was missing most of them.

- Twitter does provide an option of setting a reminder for Twitter Spaces but at the same time it doesn't have an option to collectively go through the list of all the Twitter Spaces for which one has set a reminder.

- A solution that I could think of was to why not link all the Twitter Spaces that I wish to attend, directly with my Google Calendar.

- Most of us already use Google Calendar in our daily lives and this would be an easy way to track all the twitter events as well.

### 💡 The Solution

- You just need to reply `@markeventas <event_name> | <date> | <time> | <timezone>` to any tweet.

- This is not only limited to Twitter Spaces. You can literally create an event for any tweet. Example:
    - Suppose someone tweeted about organizing a meetup on Sunday at 9PM Indian Standart Time.
        - You can create an event for it by replying `@markeventas Sunday Meetup | 25th Aug, 2022 | 9 PM | IST` to the tweet.

## Tech Stack

**Backend:** Golang

**Database:** PostgreSQL

**Infra:** Docker, AWS ECR, AWS EKS

## DB Schema

![DB Schema](https://user-images.githubusercontent.com/43776315/161421459-c6d881d5-b40f-4361-9629-c0a533115a00.png)

## Developer

- [@skamranahmed](https://github.com/skamranahmed)
