# kbbi-api

[![godoc](https://godoc.org/github.com/raf555/kbbi-api/pkg/kbbi?status.svg)][godoc]

[godoc]: https://godoc.org/github.com/raf555/kbbi-api/pkg/kbbi

Probably the most complete public KBBI API you will ever find. 


## Documentation

- Swagger docs: [https://kbbi.raf555.dev/swagger/index.html](https://kbbi.raf555.dev/swagger/index.html)

- Sample API endpoint: [https://kbbi.raf555.dev/api/v1/entry/apel](https://kbbi.raf555.dev/api/v1/entry/apel)

- Sample response:

```json
{
  "lemma": "bermalas-malasan",
  "entries": [
    {
      "entry": "ber.ma.las-ma.las.an",
      "baseWord": "malas",
      "entryVariants": [],
      "pronunciation": "",
      "definitions": [
        {
          "definition": "bermalas-malas",
          "referencedLemma": "",
          "labels": [
            {
              "code": "v",
              "name": "Verba",
              "kind": "Kelas Kata"
            }
          ],
          "usageExamples": [
            "dia selalu ~ saat hari libur"
          ]
        }
      ],
      "nonStandardWords": [],
      "variants": [],
      "compoundWords": [],
      "derivedWords": [],
      "proverbs": [],
      "metaphors": []
    }
  ]
}
```


## Data Source

Latest edition: `Oktober 2023`

The dictionary is mirrored from [Official KBBI Application][] `v1.0.0` with some hand-edited data and customly decoded for author's requirement.

> [!NOTE]  
> You might encounter some lemmas that are having different information from the official KBBI website (or even missing). I tend to use the official site as the source of truth, so in case you find some, please make a new issue, I'll try to get it fixed ASAP.

The dictionary used by the server will be updated as soon as new version of the application is released.

## Issues

For any issues, be it from the API server or from the dictionary, or questions or inquiries or suggestions, feel free to raise a new issue. I'll try to look into it as soon as possible.

If you are willing to open PR to fix any open issues, I'm more than welcome to review it.

## Background and Motivation

**TL;DR**. Official KBBI website sucks, I build my own API.

<details>
  <summary>Expand</summary>

Due to recent [Official KBBI Website][] introducing Cloudflare firewall to their site and limiting user's request to only a couple of lemmas for each day, my personal chatbot which scraps the website for the lemma information became unusable. Even as an actual user, it is kind of frustrating, really.

<img width="942" alt="Image" src="https://github.com/user-attachments/assets/7dc09b77-cde6-4140-ab84-f129823c7816" />

I did a bit of research to find a free public KBBI API on the internet, but most of them don't really give the information that I need that I have used on the chatbot (e.g. they does not fully cover the KBBI lemma response cases). Most of the APIs I found are also doing scraping to the KBBI website, which makes them unusable anyway. I found some that uses offline data though, but most of them are outdated already, and they don't really fit into my chatbot.

Since I'm too lazy to make a Cloudflare bypasser, I decided to make this API server. Since I want to make the information provided by this API to be as complete as possible and as fast as possible (for my chatbot), I opted for looking into the [Official KBBI Application][] since it is offline and it should have all the information I need.

Long story short, I was able to scrap all lemmas from there *(I won't tell how I did this though (yet?))*. All data used in this API is completely from the application (with some additional hand-edited data). They are then decoded and parsed to fit my requirement. It ends up perfectly as I wanted. The final product is the one you see on the API response.

Feel free to use the API as much as you want, there is no rate limiting as of now (**not yet, at least**). As long as the server can handle the traffic and does not exceed the free resources usage on the cloud provider I used, lol.
</details>

[Official KBBI Website]: https://kbbi.kemdikbud.go.id/

[Official KBBI Application]: https://play.google.com/store/apps/details?id=yuku.kbbi5
