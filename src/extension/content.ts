import axios from "axios"
import type { PlasmoCSConfig } from "plasmo"

export {}

async function sleep(duration: number) {
  await new Promise((res) => setTimeout(res, duration))
}

console.log("Loaded content script")

chrome.runtime.onMessage.addListener(function (request) {
  const key = request.key
  const apiBaseUrl = request.apiBaseUrl

  run(key, apiBaseUrl)

  return true
})

async function run(key: string, apiBaseUrl: string) {
  const out: { promise: Promise<any>; url: string; post: string }[] = []

  for (const elem of Array.from(document.querySelector("[role=feed]").children)) {
    try {
      // Expand the text
      const msg = elem.querySelector("[data-ad-comet-preview=message]")

      if (msg) {
        const more = msg.querySelector("[role=button]")

        if (more) {
          // @ts-ignore
          more.click()
          await sleep(500)
        }
      }

      // Get the post URL
      const share = elem.querySelector('[aria-label="Send this to friends or post it on your Timeline."]')
      // @ts-ignore
      share.click()
      await sleep(1000)

      const shareOptionContainer = document.querySelector("[role=dialog]")
      const shareOptions = shareOptionContainer.querySelectorAll("[role=button]")
      const copyLink = shareOptions[shareOptions.length - 1]

      copyLink.scrollIntoView()

      const url = await new Promise<string>((res) => {
        copyLink.addEventListener("click", () => {
          navigator.clipboard.readText().then((url) => res(url))
        })
      })

      // Post request
      const data = { url, post: msg.textContent }
      const promise = axios.post(`${apiBaseUrl}/rentals`, data, { headers: { "x-api-key": key } })

      out.push({ promise, url, ...data })
    } catch (e) {
      console.error(e)
    }
  }

  console.log("Finishing")

  for (const elem of out) {
    try {
      await elem.promise

      console.log("Added ", { url: elem.url, post: elem.post })
    } catch (e) {
      console.error(e)
    }
  }
}

export const config: PlasmoCSConfig = {
  matches: ["https://www.facebook.com/*"]
}
